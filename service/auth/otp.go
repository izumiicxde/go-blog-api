package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/izumii.cxde/blog-api/types"
)

func GenerateOTP() string {
	// Generate a random number between 100000 and 999999
	randomNumber := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", randomNumber)
}

func ValidateOTP(otp string, u types.User, userId int64) bool {
	if u.Otp == otp && u.OtpExpiration.After(time.Now()) && !u.Verified && u.ID == uint(userId) {
		return true
	}
	return false
}
