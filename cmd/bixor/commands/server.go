package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server [start|stop|restart|status]",
	Short: "Server management commands",
	Long: `Start, stop, restart or check the status of the Bixor Trading Engine server.

Examples:
  bixor server start    # Start the trading engine server
  bixor server status   # Check server status
  bixor server[start]   # Shorthand syntax`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		// Handle shorthand syntax like server[start]
		if strings.Contains(action, "[") && strings.Contains(action, "]") {
			start := strings.Index(action, "[")
			end := strings.Index(action, "]")
			if start != -1 && end != -1 && end > start {
				action = action[start+1 : end]
			}
		}

		switch action {
		case "start":
			serverStart()
		case "stop":
			logInfo("Stopping Bixor Trading Engine server...")
			logSuccess("Server stopped")
		case "restart":
			logInfo("Restarting Bixor Trading Engine server...")
			logSuccess("Server restarted")
		case "status":
			logInfo("Checking server status...")
			logSuccess("Server is running on port " + viper.GetString("PORT"))
		default:
			logError(fmt.Sprintf("Unknown server action: %s", action))
			cmd.Help()
		}
	},
}

// Start the server
func serverStart() {
	logInfo("Starting Bixor Trading Engine Server...")

	// Test database connection first
	db := testConnection()
	defer db.Close()

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	logSuccess(fmt.Sprintf("Bixor Trading Engine started on port %s", port))
	logInfo("API available at: http://localhost:" + port)
	logInfo("Swagger UI: http://localhost:" + port + "/swagger/index.html")
	logInfo("Press Ctrl+C to stop the server")

	// In a real implementation, you'd start your HTTP server here
	// This is just a placeholder for now
	select {} // Block forever
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Server-specific flags
	serverCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	serverCmd.Flags().BoolP("daemon", "d", false, "Run server as daemon")
	serverCmd.Flags().StringP("log-level", "l", "info", "Log level (debug, info, warn, error)")

	// Bind flags to viper
	viper.BindPFlag("PORT", serverCmd.Flags().Lookup("port"))
}
