package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func addBackendSecret(req *http.Request) {
	// Load environment variables
	godotenv.Load()
	backendSecret := os.Getenv("BACKEND_SECRET")
	if backendSecret != "" {
		req.Header.Set("X-Backend-Secret", backendSecret)
	}
}

const BaseURL = "http://localhost:8080/api/v1"

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	User    interface{} `json:"user,omitempty"`
	Tokens  interface{} `json:"tokens,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	fmt.Println("Bixor Authentication Tool")
	fmt.Println("=========================")

	for {
		showMenu()
		choice := getInput("Choose an option: ")

		switch choice {
		case "1":
			registerUser()
		case "2":
			loginUser()
		case "3":
			verifyUser()
		case "4":
			refreshToken()
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println()
	fmt.Println("1. Register account")
	fmt.Println("2. Login")
	fmt.Println("3. Verify user (activate account)")
	fmt.Println("4. Refresh token")
	fmt.Println("5. Exit")
	fmt.Println()
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func registerUser() {
	fmt.Println("\n--- Register New Account ---")

	firstName := getInput("First Name: ")
	lastName := getInput("Last Name: ")
	username := getInput("Username: ")
	email := getInput("Email: ")
	password := getInput("Password: ")

	if firstName == "" || lastName == "" || username == "" || email == "" || password == "" {
		fmt.Println("All fields are required")
		return
	}

	req := RegisterRequest{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  password,
	}

	jsonData, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", BaseURL+"/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	addBackendSecret(httpReq)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result APIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 201 {
		fmt.Println("Registration successful!")
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Email: %s\n", email)
		fmt.Println("Note: Account status is 'pending' - needs activation for login")
	} else {
		fmt.Println("Registration failed:")
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		if result.Message != "" {
			fmt.Printf("Message: %s\n", result.Message)
		}
	}
}

func loginUser() {
	fmt.Println("\n--- Login ---")

	email := getInput("Email: ")
	password := getInput("Password: ")

	if email == "" || password == "" {
		fmt.Println("Email and password are required")
		return
	}

	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", BaseURL+"/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	addBackendSecret(httpReq)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result APIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 200 {
		fmt.Println("Login successful!")
		if result.Message != "" {
			fmt.Printf("Message: %s\n", result.Message)
		}
		fmt.Println("JWT tokens generated successfully")

		// Display the actual JWT tokens
		if result.Tokens != nil {
			tokensMap, ok := result.Tokens.(map[string]interface{})
			if ok {
				fmt.Println("\n--- JWT TOKENS ---")
				if accessToken, exists := tokensMap["access_token"]; exists {
					fmt.Printf("Access Token: %s\n", accessToken)
				}
				if refreshToken, exists := tokensMap["refresh_token"]; exists {
					fmt.Printf("Refresh Token: %s\n", refreshToken)
				}
				if tokenType, exists := tokensMap["token_type"]; exists {
					fmt.Printf("Token Type: %s\n", tokenType)
				}
				if expiresIn, exists := tokensMap["expires_in"]; exists {
					fmt.Printf("Expires In: %v seconds\n", expiresIn)
				}
				fmt.Println("--- END TOKENS ---")
			}
		}
	} else {
		fmt.Println("Login failed:")
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		if result.Message != "" {
			fmt.Printf("Message: %s\n", result.Message)
		}
	}
}

func verifyUser() {
	fmt.Println("\n--- Verify User Account ---")

	email := getInput("Enter user email to verify: ")

	if email == "" {
		fmt.Println("Email is required")
		return
	}

	fmt.Printf("Activating user with email: %s\n", email)

	// Connect to database and activate user
	db := connectDB()
	if db == nil {
		return
	}
	defer db.Close()

	// Update user status to active
	query := `UPDATE users SET status = 'active', updated_at = NOW() WHERE LOWER(email) = LOWER($1) AND status != 'active' RETURNING username, email`

	var username, userEmail string
	err := db.QueryRow(query, email).Scan(&username, &userEmail)

	if err == sql.ErrNoRows {
		fmt.Printf("No inactive user found with email: %s\n", email)
		fmt.Println("User may already be active or email doesn't exist")
		return
	} else if err != nil {
		fmt.Printf("Failed to activate user: %v\n", err)
		return
	}

	fmt.Println("User activated successfully!")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Email: %s\n", userEmail)
	fmt.Println("User can now login with their credentials")
}

func connectDB() *sql.DB {
	// Load environment variables
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL environment variable is required. Please create a .env file (copy from .env.example)")
		return nil
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("Cannot connect to database: %v\n", err)
		db.Close()
		return nil
	}

	return db
}

func refreshToken() {
	fmt.Println("\n--- Refresh Token ---")

	refreshToken := getInput("Enter refresh token: ")

	if refreshToken == "" {
		fmt.Println("Refresh token is required")
		return
	}

	reqData := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, _ := json.Marshal(reqData)
	httpReq, err := http.NewRequest("POST", BaseURL+"/auth/refresh", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	addBackendSecret(httpReq)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result APIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == 200 {
		fmt.Println("Token refresh successful!")
		if result.Message != "" {
			fmt.Printf("Message: %s\n", result.Message)
		}

		// Display the new JWT tokens
		if result.Tokens != nil {
			tokensMap, ok := result.Tokens.(map[string]interface{})
			if ok {
				fmt.Println("\n--- NEW JWT TOKENS ---")
				if accessToken, exists := tokensMap["access_token"]; exists {
					fmt.Printf("New Access Token: %s\n", accessToken)
				}
				if newRefreshToken, exists := tokensMap["refresh_token"]; exists {
					fmt.Printf("New Refresh Token: %s\n", newRefreshToken)
				}
				if tokenType, exists := tokensMap["token_type"]; exists {
					fmt.Printf("Token Type: %s\n", tokenType)
				}
				if expiresIn, exists := tokensMap["expires_in"]; exists {
					fmt.Printf("Expires In: %v seconds\n", expiresIn)
				}
				fmt.Println("--- END NEW TOKENS ---")
			}
		}
	} else {
		fmt.Println("Token refresh failed:")
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		if result.Message != "" {
			fmt.Printf("Message: %s\n", result.Message)
		}
	}
}
