package generateotp

import (
	"math/rand"
	"time"
)

// init initializes the random number generator seed.
// It uses the current Unix timestamp in nanoseconds to ensure that the random numbers 
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateOTP generates a one-time password (OTP) of the specified length.
// The OTP consists of numeric characters only (0-9).
func GenerateOTP(length int) string {
	// Set of characters to choose from for OTP generation
	characters := "0123456789"
	// Initialize an empty string to hold the generated OTP.
	otp := ""

	// Loop through 'length' times, randomly selecting a character for each position.
	for i := 0; i < length; i++ {
		// Generate a random index to select a character from the set.
		randomIndex := rand.Intn(len(characters))

		// Append the selected character to the OTP string.
		otp += string(characters[randomIndex])
	}

	// Return the final generated OTP as a string.
	return otp
}
