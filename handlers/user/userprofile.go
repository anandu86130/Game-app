package user

import (
	database "gaming/database"
	"gaming/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Define a struct to represent a Person with UserID, Name, and Email fields
type Person struct {
	UserID uint   // Unique identifier for the user
	Name   string // User's name
	Email  string // User's email
}

// UserProfile handles retrieving a user's profile information
func UserProfile(c *gin.Context) {
	var user model.User
	userid := c.GetUint("userid") // Retrieve the user ID from the context

	// Query the database for the user with the given user ID
	result := database.DB.Where("user_id=?", userid).First(&user)
	if result.Error != nil {
		// If there is an error finding the user, return an internal server error response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to find user",
		})
		return
	}

	// Create a Person instance with the user's details
	users := Person{
		UserID: userid,
		Name:   user.Name,
		Email:  user.Email,
	}

	// Return the user's profile information as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"user": users, // Include the user profile in the response
	})
}

// EditUser handles updating a user's profile information
func EditUser(c *gin.Context) {
	var edit model.User
	id := c.GetUint("userid") // Retrieve the user ID from the context

	// Bind the incoming JSON payload to the edit variable
	result := c.BindJSON(&edit)
	if result != nil {
		// If binding fails, return an internal server error response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to bind json",
		})
		return
	}

	var editUser model.User
	// Fetch the user from the database using the user ID
	fetch := database.DB.First(&editUser, id)
	if fetch.Error != nil {
		// If there is an error fetching the user, return an internal server error response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch user",
		})
		return
	}

	if fetch.RowsAffected > 0 {
		// If the user is found, update the user's details in the database
		database.DB.Model(&editUser).Updates(edit)
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK, // Corrected "stauts" to "status"
			"Message": "User updated successfully",
		})
	} else {
		// If no rows were affected, return an internal server error response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to update user",
		})
		return
	}
}
