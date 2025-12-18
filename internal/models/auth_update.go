package models

// UpdateProfileRequest represents the request payload for updating user profile
type UpdateProfileRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=50"`
	LastName    string  `json:"last_name" binding:"required,min=2,max=50"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
}

// UpdateSettingsRequest represents the request payload for updating user settings
type UpdateSettingsRequest struct {
	Language string `json:"language" binding:"required,len=2"` // e.g., "en"
	Timezone string `json:"timezone" binding:"required"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
}

// ToggleTwoFARequest represents the request payload for toggling 2FA
type ToggleTwoFARequest struct {
	Enable bool   `json:"enable"`
	Code   string `json:"code" binding:"required,len=6,numeric"` // OTP to verify action
}
