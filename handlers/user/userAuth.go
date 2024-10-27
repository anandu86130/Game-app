package user

import (
	"fmt"
	database "gaming/database"
	"gaming/generateotp"
	"gaming/jwt"
	"gaming/model"
	"gaming/utility"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Define user model variable
var user model.User

// Define user role constant
const RoleUser = "user"

// Signup function handles user registration
func Signup(c *gin.Context) {
	// Bind the JSON input to the User model
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("error when binding json: %v", err) // Log binding error
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Error when binding JSON",
		})
		return
	}

	// Check if the user already exists in the database
	var existingUsers []model.User
	result := database.DB.Where("email = ?", user.Email).Find(&existingUsers)

	// If user already exists, return a conflict status
	if result.RowsAffected > 0 {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "This user already exists",
		})
		return
	}

	// Hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "error when hashing password",
		})
		return
	}
	user.Password = string(hashedPassword) // Save the hashed password

	// Generate a new OTP (One-Time Password)
	otp := generateotp.GenerateOTP(6)
	newOTP := model.OTP{
		Email: user.Email,
		Otp:   otp,
		Exp:   time.Now().Add(5 * time.Minute), // Set expiration time for the OTP
	}
	fmt.Println("Generated OTP:", otp)

	// Check if an OTP already exists for this user
	var existingOTP model.OTP
	if err := database.DB.Where("email = ?", user.Email).First(&existingOTP).Error; err == nil {
		// If an existing OTP is found, update it
		existingOTP.Otp = otp
		existingOTP.Exp = time.Now().Add(5 * time.Minute) // Reset expiration time
		if err := database.DB.Save(&existingOTP).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "failed to update OTP in the database",
			})
			return
		}
	} else {
		// Create a new OTP entry in the database
		if err := database.DB.Create(&newOTP).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "Failed to store OTP to the database",
			})
			return
		}
	}

	// Send the OTP to the user's email
	utility.SendOTPByEmail(newOTP.Email, newOTP.Otp)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"Message": "OTP sent successfully, please verify OTP",
	})
}

// VerifyOTP function handles OTP verification for user registration
func VerifyOTP(c *gin.Context) {
	var otp model.OTP
	err := c.BindJSON(&otp) // Bind the JSON input to the OTP model
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind json",
		})
		return
	}

	// Retrieve the existing OTP for the user from the database
	var existingotp model.OTP
	result := database.DB.Where("email=?", user.Email).First(&existingotp)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch otp",
		})
		return
	}

	// Check if the OTP is expired
	currentTime := time.Now()
	if currentTime.After(existingotp.Exp) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "otp expired",
		})
		return
	}

	// Check if the provided OTP matches the stored OTP
	if existingotp.Otp != otp.Otp {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "otp expired",
		})
		return
	}

	// Create the new user in the database
	create := database.DB.Create(&user)
	fmt.Println(user)
	if create.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to create user",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"Message": "User created successfully",
		})
	}
}

// Login function handles user authentication
func Login(c *gin.Context) {
	var userlogin model.User
	err := c.ShouldBindJSON(&userlogin) // Bind the JSON input for login
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to bind json",
		})
		return
	}

	// Retrieve the existing user from the database
	var existinguser model.User
	result := database.DB.Where("email=?", userlogin.Email).First(&existinguser)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "incorrect email or password",
		})
		return
	}

	// Compare the provided password with the hashed password
	password := bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte(userlogin.Password))
	if password != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "incorrect email or password",
		})
		return
	} else {
		// Generate a JWT token for the authenticated user
		jwt.JwtToken(c, existinguser.UserID, userlogin.Email, RoleUser)
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"Message": "Login successfully",
		})
	}
}