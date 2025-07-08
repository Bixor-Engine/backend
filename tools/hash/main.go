package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Bixor-Engine/backend/internal/models"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔐 Argon2 Hash Tool for Bixor Engine 🔐")
	fmt.Println("=====================================")
	fmt.Println("")

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Hash a password (Encrypt)")
		fmt.Println("2. Verify a hash (Validate)")
		fmt.Println("3. Exit")
		fmt.Print("\nEnter your choice (1-3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			hashPassword(reader)
		case "2":
			verifyHash(reader)
		case "3":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please enter 1, 2, or 3.")
		}

		fmt.Println("")
		fmt.Println("─────────────────────────────────────")
		fmt.Println("")
	}
}

func hashPassword(reader *bufio.Reader) {
	fmt.Println("")
	fmt.Println("🔒 PASSWORD HASHING")
	fmt.Println("===================")
	fmt.Print("Enter password to hash: ")

	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("❌ Password cannot be empty!")
		return
	}

	fmt.Printf("🔄 Hashing password: '%s'\n", password)

	// Hash the password using our Argon2 implementation
	hash, err := models.HashPassword(password, nil)
	if err != nil {
		fmt.Printf("❌ Failed to hash password: %v\n", err)
		return
	}

	fmt.Println("")
	fmt.Println("✅ PASSWORD HASHED SUCCESSFULLY!")
	fmt.Println("================================")
	fmt.Printf("📝 Plain Text: %s\n", password)
	fmt.Printf("🔐 Argon2 Hash: %s\n", hash)
	fmt.Println("")
	fmt.Printf("📋 Hash Length: %d characters\n", len(hash))

	// Show hash components
	parts := strings.Split(hash, "$")
	if len(parts) >= 4 {
		fmt.Printf("🔧 Algorithm: %s\n", parts[1])
		fmt.Printf("🔧 Version: %s\n", parts[2])
		fmt.Printf("🔧 Parameters: %s\n", parts[3])
	}

	fmt.Println("")
	fmt.Println("💡 You can now copy this hash and use option 2 to verify it!")
}

func verifyHash(reader *bufio.Reader) {
	fmt.Println("")
	fmt.Println("🔍 HASH VERIFICATION")
	fmt.Println("====================")
	fmt.Print("Enter the plain text password: ")

	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("❌ Password cannot be empty!")
		return
	}

	fmt.Print("Enter the Argon2 hash to verify against: ")
	hash, _ := reader.ReadString('\n')
	hash = strings.TrimSpace(hash)

	if hash == "" {
		fmt.Println("❌ Hash cannot be empty!")
		return
	}

	fmt.Printf("🔄 Verifying password: '%s'\n", password)
	fmt.Printf("🔄 Against hash: %s\n", hash)
	fmt.Println("")

	// Verify the password
	isValid, err := models.VerifyPassword(password, hash)
	if err != nil {
		fmt.Printf("❌ Failed to verify hash: %v\n", err)
		fmt.Println("")
		fmt.Println("💡 Common issues:")
		fmt.Println("   - Hash format is incorrect")
		fmt.Println("   - Hash was generated with different algorithm")
		fmt.Println("   - Hash is corrupted or incomplete")
		return
	}

	fmt.Println("🔍 VERIFICATION RESULT")
	fmt.Println("======================")
	if isValid {
		fmt.Println("✅ MATCH! The password is CORRECT!")
		fmt.Printf("✅ '%s' matches the provided hash\n", password)
		fmt.Println("✅ User would be able to login with this password")
	} else {
		fmt.Println("❌ NO MATCH! The password is INCORRECT!")
		fmt.Printf("❌ '%s' does NOT match the provided hash\n", password)
		fmt.Println("❌ User would NOT be able to login with this password")
	}

	fmt.Println("")
	fmt.Println("🔐 Security Note: This is exactly how login authentication works!")
}
