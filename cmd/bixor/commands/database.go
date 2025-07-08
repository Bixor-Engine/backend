package commands

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Colors for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

// Logger functions
func logInfo(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", ColorGreen, ColorReset, msg)
}

func logWarn(msg string) {
	fmt.Printf("%s[WARN]%s %s\n", ColorYellow, ColorReset, msg)
}

func logError(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, msg)
}

func logSuccess(msg string) {
	fmt.Printf("%s[SUCCESS]%s %s\n", ColorCyan, ColorReset, msg)
}

// Test database connection
func testConnection() *sql.DB {
	logInfo("Testing database connection...")

	databaseURL := viper.GetString("DATABASE_URL")
	if databaseURL == "" {
		logError("DATABASE_URL not found. Please set it in .env file or environment.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		logError(fmt.Sprintf("Failed to open database connection: %v", err))
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		logError("Cannot connect to database. Please check your DATABASE_URL.")
		db.Close()
		os.Exit(1)
	}

	logInfo("Database connection successful")
	return db
}

// Execute SQL file
func executeSQLFile(db *sql.DB, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filepath, err)
	}

	_, err = db.Exec(string(content))
	return err
}

// Create database if it doesn't exist
func createDatabaseIfNotExists() {
	logInfo("Checking if database exists...")

	databaseURL := viper.GetString("DATABASE_URL")
	if databaseURL == "" {
		logError("DATABASE_URL not found. Please set it in .env file or environment.")
		os.Exit(1)
	}

	// First, try to connect directly to the target database
	// If it works, the database exists and we're good to go
	testDB, err := sql.Open("postgres", databaseURL)
	if err == nil {
		if pingErr := testDB.Ping(); pingErr == nil {
			testDB.Close()
			logInfo("Database connection successful - database exists")
			return
		}
		testDB.Close()
	}

	// If direct connection failed, try to create the database
	logInfo("Direct connection failed, attempting to create database...")

	// Parse database URL to get database name and connection without database
	// postgres://user:pass@host:port/dbname -> postgres://user:pass@host:port/postgres
	parts := strings.Split(databaseURL, "/")
	if len(parts) < 4 {
		logError("Invalid DATABASE_URL format")
		os.Exit(1)
	}

	dbName := strings.Split(parts[len(parts)-1], "?")[0] // Remove query params

	// Try different administrative databases in order of preference
	adminDBs := []string{"postgres", "template1", dbName}
	var adminDB *sql.DB
	var adminURL string

	for _, adminDBName := range adminDBs {
		adminURL = strings.Join(parts[:len(parts)-1], "/") + "/" + adminDBName
		if adminURLWithParams := adminURL; strings.Contains(databaseURL, "?") {
			queryParams := strings.Split(databaseURL, "?")[1]
			adminURL = adminURLWithParams + "?" + queryParams
		}

		adminDB, err = sql.Open("postgres", adminURL)
		if err == nil {
			if pingErr := adminDB.Ping(); pingErr == nil {
				break // Successfully connected to admin database
			}
			adminDB.Close()
			adminDB = nil
		}
	}

	if adminDB == nil {
		logError("Cannot connect to any administrative database (postgres, template1). Please ensure PostgreSQL is running and accessible.")
		os.Exit(1)
	}
	defer adminDB.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = adminDB.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		logError(fmt.Sprintf("Failed to check if database exists: %v", err))
		os.Exit(1)
	}

	if !exists {
		logWarn(fmt.Sprintf("Database '%s' does not exist. Create it? (y/N)", dbName))
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your choice: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			logInfo(fmt.Sprintf("Creating database '%s'...", dbName))
			_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
			if err != nil {
				logError(fmt.Sprintf("Failed to create database: %v", err))
				os.Exit(1)
			}
			logSuccess(fmt.Sprintf("Database '%s' created successfully", dbName))
		} else {
			logInfo("Database creation cancelled.")
			os.Exit(0)
		}
	} else {
		logInfo(fmt.Sprintf("Database '%s' already exists", dbName))
	}
}

// Database setup - create database if needed, run init, migrations, and seeders
func databaseSetup(db *sql.DB) {
	logInfo("Setting up Bixor Trading Engine database...")

	// Create database if it doesn't exist
	createDatabaseIfNotExists()

	// Connect to the target database (it should exist now)
	if db != nil {
		db.Close()
	}
	db = testConnection()
	defer db.Close()

	// Run init.sql first
	initFile := "database/init.sql"
	if _, err := os.Stat(initFile); err == nil {
		logInfo("Running database initialization...")
		if err := executeSQLFile(db, initFile); err != nil {
			logError(fmt.Sprintf("Failed to run init.sql: %v", err))
			os.Exit(1)
		}
		logSuccess("Database initialization completed")
	}

	// Then run migrations
	databaseMigrate(db)

	// Finally run seeders
	databaseSeed(db)

	logSuccess("Database setup completed successfully!")
}

// Check if migrations have been applied
func getMigrationsStatus(db *sql.DB) map[string]bool {
	appliedMigrations := make(map[string]bool)

	// Check if schema_migrations table exists
	var tableExists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'schema_migrations'
	)`
	err := db.QueryRow(query).Scan(&tableExists)
	if err != nil || !tableExists {
		return appliedMigrations // Empty map if table doesn't exist
	}

	// Get applied migrations
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return appliedMigrations
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err == nil {
			appliedMigrations[version] = true
		}
	}

	return appliedMigrations
}

// Run migrations
func databaseMigrate(db *sql.DB) {
	logInfo("Running migrations...")

	migrationDir := "database/migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		logWarn(fmt.Sprintf("Migration directory not found: %s", migrationDir))
		return
	}

	var migrations []string
	err := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			migrations = append(migrations, path)
		}
		return nil
	})

	if err != nil {
		logError(fmt.Sprintf("Failed to read migration directory: %v", err))
		return
	}

	if len(migrations) == 0 {
		logWarn("No migration files found")
		return
	}

	sort.Strings(migrations)

	// Get already applied migrations
	appliedMigrations := getMigrationsStatus(db)

	// Filter out already applied migrations
	var pendingMigrations []string
	for _, migration := range migrations {
		migrationName := filepath.Base(migration)

		// Extract just the numeric version (001, 002, etc.) from filename like "001_create_schema_migrations.sql"
		parts := strings.Split(migrationName, "_")
		if len(parts) == 0 {
			continue
		}
		version := parts[0] // Get "001" from "001_create_schema_migrations.sql"

		if !appliedMigrations[version] {
			pendingMigrations = append(pendingMigrations, migration)
		} else {
			logInfo(fmt.Sprintf("Skipping already applied migration: %s (version %s)", migrationName, version))
		}
	}

	// Debug: Show what we found
	if len(appliedMigrations) > 0 {
		logInfo(fmt.Sprintf("Found %d already applied migrations in database", len(appliedMigrations)))
	}

	if len(pendingMigrations) == 0 {
		logSuccess("All migrations are already applied! Database is up to date.")
		return
	}

	// Show confirmation for pending migrations
	logWarn(fmt.Sprintf("Found %d pending migration(s). Apply them? (y/N)", len(pendingMigrations)))
	for _, migration := range pendingMigrations {
		fmt.Printf("  - %s\n", filepath.Base(migration))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		logInfo("Migration cancelled.")
		return
	}

	migrationCount := 0
	for _, migration := range pendingMigrations {
		migrationName := filepath.Base(migration)
		logInfo(fmt.Sprintf("Applying migration: %s", migrationName))

		if err := executeSQLFile(db, migration); err != nil {
			logError(fmt.Sprintf("✗ Failed to apply %s: %v", migrationName, err))
			os.Exit(1)
		}

		logSuccess(fmt.Sprintf("✓ %s applied successfully", migrationName))
		migrationCount++
	}

	logSuccess(fmt.Sprintf("All %d new migrations applied successfully!", migrationCount))
}

// Run seeders
func databaseSeed(db *sql.DB) {
	logInfo("Running seeders...")

	seederDir := "database/seeders"
	if _, err := os.Stat(seederDir); os.IsNotExist(err) {
		logWarn(fmt.Sprintf("Seeder directory not found: %s", seederDir))
		return
	}

	var seeders []string
	err := filepath.WalkDir(seederDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			seeders = append(seeders, path)
		}
		return nil
	})

	if err != nil {
		logError(fmt.Sprintf("Failed to read seeder directory: %v", err))
		return
	}

	if len(seeders) == 0 {
		logWarn("No seeder files found")
		return
	}

	sort.Strings(seeders)

	seederCount := 0
	for _, seeder := range seeders {
		seederName := filepath.Base(seeder)
		logInfo(fmt.Sprintf("Running seeder: %s", seederName))

		if err := executeSQLFile(db, seeder); err != nil {
			logWarn(fmt.Sprintf("⚠ %s had issues (this may be normal for conflict handling): %v", seederName, err))
		} else {
			logSuccess(fmt.Sprintf("✓ %s completed successfully", seederName))
		}
		seederCount++
	}

	logSuccess(fmt.Sprintf("All %d seeders completed!", seederCount))
}

// Reset database
func databaseReset(db *sql.DB) {
	logWarn("This will drop all tables and data. Are you sure? (y/N)")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		logInfo("Resetting database...")

		// Drop all tables
		dropSQL := `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;`

		if _, err := db.Exec(dropSQL); err != nil {
			logError(fmt.Sprintf("Failed to drop tables: %v", err))
			os.Exit(1)
		}

		logInfo("All tables dropped. Running fresh setup...")

		// Run setup without database existence check (we know it exists since we just dropped tables from it)
		databaseSetupWithoutDBCheck(db)
	} else {
		logInfo("Reset cancelled.")
	}
}

// Database setup without database existence check (used for reset operations)
func databaseSetupWithoutDBCheck(db *sql.DB) {
	logInfo("Setting up Bixor Trading Engine database...")

	// Run init.sql first
	initFile := "database/init.sql"
	if _, err := os.Stat(initFile); err == nil {
		logInfo("Running database initialization...")
		if err := executeSQLFile(db, initFile); err != nil {
			logError(fmt.Sprintf("Failed to run init.sql: %v", err))
			os.Exit(1)
		}
		logSuccess("Database initialization completed")
	}

	// Then run migrations
	databaseMigrate(db)

	// Finally run seeders
	databaseSeed(db)

	logSuccess("Database setup completed successfully!")
}

// Show migration status
func databaseStatus(db *sql.DB) {
	logInfo("Migration status:")
	fmt.Println()

	query := `
	SELECT 
		version as "Version",
		description as "Description",
		applied_at as "Applied At",
		applied_by as "Applied By"
	FROM schema_migrations 
	ORDER BY version;`

	rows, err := db.Query(query)
	if err != nil {
		logWarn("schema_migrations table not found. Run 'bixor database setup' first.")
		return
	}
	defer rows.Close()

	// Print header
	fmt.Printf("%-10s %-30s %-20s %-15s\n", "Version", "Description", "Applied At", "Applied By")
	fmt.Println(strings.Repeat("-", 80))

	// Print rows
	for rows.Next() {
		var version, description, appliedAt, appliedBy string
		if err := rows.Scan(&version, &description, &appliedAt, &appliedBy); err != nil {
			logError(fmt.Sprintf("Error scanning row: %v", err))
			continue
		}
		fmt.Printf("%-10s %-30s %-20s %-15s\n", version, description, appliedAt, appliedBy)
	}

	fmt.Println()
	logInfo(`Use "SELECT * FROM database_info;" to see database version info`)
}

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database [setup|migrate|seed|reset|status]",
	Short: "Database management commands",
	Long: `Manage your Bixor Trading Engine database with various operations:

- setup: Create database if needed + run init + migrations + seeders (first-time setup)
- migrate: Run pending migrations only (with safety checks and confirmation)
- seed: Run database seeders  
- reset: Reset database (drop all + full setup)
- status: Show migration status

Note: Use quotes for shorthand syntax:
  bixor "database[setup]"    # Same as: bixor database setup
  bixor "database[migrate]"  # Same as: bixor database migrate`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Handle shorthand syntax like database[setup]
		action := args[0]
		if strings.Contains(action, "[") && strings.Contains(action, "]") {
			// Extract action from database[setup] format
			start := strings.Index(action, "[")
			end := strings.Index(action, "]")
			if start != -1 && end != -1 && end > start {
				action = action[start+1 : end]
			}
		}

		switch action {
		case "setup":
			// For setup, we handle database creation, so don't test connection first
			databaseSetup(nil)
		case "migrate", "seed", "reset", "status":
			// These commands need existing database connection
			db := testConnection()
			defer db.Close()

			switch action {
			case "migrate":
				databaseMigrate(db)
			case "seed":
				databaseSeed(db)
			case "reset":
				databaseReset(db)
			case "status":
				databaseStatus(db)
			}
		default:
			logError(fmt.Sprintf("Unknown database action: %s", action))
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	// Here you will define your flags and configuration settings.
	databaseCmd.Flags().BoolP("force", "f", false, "Force operation without confirmation")
	databaseCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
