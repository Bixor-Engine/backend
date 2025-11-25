package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

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

func connectDB() *sql.DB {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logError("DATABASE_URL environment variable is required. Please create a .env file (copy from .env.example)")
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

	return db
}

func listUsers(db *sql.DB, filter string) {
	query := `
		SELECT id, first_name, last_name, username, email, role, status, 
			   email_status, phone_number, created_at, last_login_at
		FROM users 
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if filter != "" {
		query += ` AND (LOWER(username) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1) OR LOWER(first_name) LIKE LOWER($1) OR LOWER(last_name) LIKE LOWER($1))`
		args = append(args, "%"+filter+"%")
	}

	query += ` ORDER BY created_at DESC LIMIT 20`

	rows, err := db.Query(query, args...)
	if err != nil {
		logError(fmt.Sprintf("Failed to query users: %v", err))
		return
	}
	defer rows.Close()

	fmt.Printf("\n%s%-36s %-15s %-15s %-20s %-10s %-10s %-6s %-15s%s\n",
		ColorBlue, "ID", "First Name", "Last Name", "Email", "Role", "Status", "Email", "Username", ColorReset)
	fmt.Println(strings.Repeat("-", 150))

	count := 0
	for rows.Next() {
		var id, firstName, lastName, username, email, role, status string
		var emailStatus bool
		var phoneNumber, lastLoginAt *string
		var createdAt string

		err := rows.Scan(&id, &firstName, &lastName, &username, &email, &role, &status,
			&emailStatus, &phoneNumber, &createdAt, &lastLoginAt)
		if err != nil {
			logError(fmt.Sprintf("Error scanning row: %v", err))
			continue
		}

		emailStatusStr := "No"
		if emailStatus {
			emailStatusStr = "Yes"
		}

		fmt.Printf("%-36s %-15s %-15s %-20s %-10s %-10s %-6s %-15s\n",
			id[:36], firstName, lastName, email, role, status, emailStatusStr, username)
		count++
	}

	if count == 0 {
		logWarn("No users found matching the criteria")
	} else {
		logInfo(fmt.Sprintf("Found %d users", count))
	}
}

func activateUser(db *sql.DB, identifier string) {
	// Try to find user by username or email
	query := `
		UPDATE users 
		SET status = 'active', updated_at = NOW()
		WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1)) 
		AND deleted_at IS NULL
		AND status != 'active'
		RETURNING id, username, email, status
	`

	var id, username, email, status string
	err := db.QueryRow(query, identifier).Scan(&id, &username, &email, &status)

	if err == sql.ErrNoRows {
		logWarn(fmt.Sprintf("No inactive user found with identifier: %s", identifier))
		return
	} else if err != nil {
		logError(fmt.Sprintf("Failed to activate user: %v", err))
		return
	}

	logSuccess(fmt.Sprintf("User activated successfully!"))
	fmt.Printf("  ID: %s\n", id)
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  New Status: %s\n", status)
}

func deactivateUser(db *sql.DB, identifier string) {
	query := `
		UPDATE users 
		SET status = 'pending', updated_at = NOW()
		WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1)) 
		AND deleted_at IS NULL
		AND status = 'active'
		RETURNING id, username, email, status
	`

	var id, username, email, status string
	err := db.QueryRow(query, identifier).Scan(&id, &username, &email, &status)

	if err == sql.ErrNoRows {
		logWarn(fmt.Sprintf("No active user found with identifier: %s", identifier))
		return
	} else if err != nil {
		logError(fmt.Sprintf("Failed to deactivate user: %v", err))
		return
	}

	logSuccess(fmt.Sprintf("User deactivated successfully!"))
	fmt.Printf("  ID: %s\n", id)
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  New Status: %s\n", status)
}

func deleteTestUsers(db *sql.DB) {
	query := `
		DELETE FROM users 
		WHERE username LIKE '%test_%' 
		OR email LIKE '%test_%@%'
	`

	result, err := db.Exec(query)
	if err != nil {
		logError(fmt.Sprintf("Failed to delete test users: %v", err))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	logSuccess(fmt.Sprintf("Deleted %d test users", rowsAffected))
}

func showUserDetails(db *sql.DB, identifier string) {
	query := `
		SELECT id, first_name, last_name, username, email, password, 
			   email_status, phone_number, phone_status, referred_by,
			   address, city, country, role, status, kyc_status,
			   twofa_enabled, last_login_at, last_login_ip, language,
			   timezone, created_at, updated_at
		FROM users 
		WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1) OR id::text = $1)
		AND deleted_at IS NULL
	`

	var user struct {
		ID, FirstName, LastName, Username, Email, Password string
		EmailStatus, PhoneStatus, TwoFAEnabled             bool
		PhoneNumber, ReferredBy, Address, City, Country    *string
		Role, Status, KYCStatus, LastLoginIP, Language     string
		Timezone, CreatedAt, UpdatedAt                     string
		LastLoginAt                                        *string
	}

	err := db.QueryRow(query, identifier).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password,
		&user.EmailStatus, &user.PhoneNumber, &user.PhoneStatus, &user.ReferredBy,
		&user.Address, &user.City, &user.Country, &user.Role, &user.Status, &user.KYCStatus,
		&user.TwoFAEnabled, &user.LastLoginAt, &user.LastLoginIP, &user.Language,
		&user.Timezone, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logWarn(fmt.Sprintf("No user found with identifier: %s", identifier))
		return
	} else if err != nil {
		logError(fmt.Sprintf("Failed to query user: %v", err))
		return
	}

	fmt.Printf("\n%sUser Details:%s\n", ColorCyan, ColorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Name: %s %s\n", user.FirstName, user.LastName)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s (Verified: %t)\n", user.Email, user.EmailStatus)

	if user.PhoneNumber != nil {
		fmt.Printf("Phone: %s (Verified: %t)\n", *user.PhoneNumber, user.PhoneStatus)
	} else {
		fmt.Printf("Phone: Not provided\n")
	}

	fmt.Printf("Role: %s\n", user.Role)
	fmt.Printf("Status: %s\n", user.Status)
	fmt.Printf("KYC Status: %s\n", user.KYCStatus)
	fmt.Printf("2FA Enabled: %t\n", user.TwoFAEnabled)

	if user.Address != nil || user.City != nil || user.Country != nil {
		fmt.Printf("Address: %s, %s, %s\n",
			getStringValue(user.Address), getStringValue(user.City), getStringValue(user.Country))
	}

	fmt.Printf("Language: %s\n", user.Language)
	fmt.Printf("Timezone: %s\n", user.Timezone)
	fmt.Printf("Created: %s\n", user.CreatedAt[:19])
	fmt.Printf("Updated: %s\n", user.UpdatedAt[:19])

	if user.LastLoginAt != nil {
		fmt.Printf("Last Login: %s\n", (*user.LastLoginAt)[:19])
		if user.LastLoginIP != "" {
			fmt.Printf("Last Login IP: %s\n", user.LastLoginIP)
		}
	} else {
		fmt.Printf("Last Login: Never\n")
	}

	fmt.Printf("Password Hash: %s...\n", user.Password[:50])
}

func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return "N/A"
}

func showHelp() {
	fmt.Printf("\n%sBixor User Manager%s\n", ColorCyan, ColorReset)
	fmt.Println("==================")
	fmt.Println()
	fmt.Println("Usage: go run cmd/user-manager/main.go <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list [filter]           - List users (optional filter by name/email/username)")
	fmt.Println("  show <identifier>       - Show detailed user information")
	fmt.Println("  activate <identifier>   - Activate a user (set status to 'active')")
	fmt.Println("  deactivate <identifier> - Deactivate a user (set status to 'pending')")
	fmt.Println("  cleanup                 - Delete all test users")
	fmt.Println("  help                    - Show this help message")
	fmt.Println()
	fmt.Println("Identifier can be: username, email, or user ID")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/user-manager/main.go list")
	fmt.Println("  go run cmd/user-manager/main.go list test")
	fmt.Println("  go run cmd/user-manager/main.go show alice_test_20241201_123456")
	fmt.Println("  go run cmd/user-manager/main.go activate alice_test_20241201_123456")
	fmt.Println("  go run cmd/user-manager/main.go cleanup")
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := strings.ToLower(os.Args[1])

	if command == "help" {
		showHelp()
		return
	}

	// Connect to database for all other commands
	db := connectDB()
	defer db.Close()

	switch command {
	case "list":
		filter := ""
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		listUsers(db, filter)

	case "show":
		if len(os.Args) < 3 {
			logError("Missing user identifier. Usage: show <username|email|id>")
			return
		}
		showUserDetails(db, os.Args[2])

	case "activate":
		if len(os.Args) < 3 {
			logError("Missing user identifier. Usage: activate <username|email|id>")
			return
		}
		activateUser(db, os.Args[2])

	case "deactivate":
		if len(os.Args) < 3 {
			logError("Missing user identifier. Usage: deactivate <username|email|id>")
			return
		}
		deactivateUser(db, os.Args[2])

	case "cleanup":
		fmt.Print("Are you sure you want to delete all test users? (y/N): ")
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) == "y" || strings.ToLower(confirmation) == "yes" {
			deleteTestUsers(db)
		} else {
			logInfo("Cleanup cancelled.")
		}

	default:
		logError(fmt.Sprintf("Unknown command: %s", command))
		showHelp()
	}
}
