package players

import (
	database "gaming/database"
	"gaming/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTeam(c *gin.Context) {
	var team model.Team
	// Bind the incoming JSON to the team struct
	if err := c.ShouldBindJSON(&team); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind JSON",
		})
		return
	}

	// Check if the user exists
	var user model.User
	if err := database.DB.First(&user, team.PlayerID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	var league model.League
	if err := database.DB.First(&league, team.LeagueID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "league not found",
		})
		return
	}
	// Check if a team with the same name already exists
	var existingTeam model.Team
	if err := database.DB.Where("name = ?", team.Name).First(&existingTeam).Error; err == nil {
		// If no error, it means the team already exists
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "this team already exists, please create a different team",
		})
		return
	}
	// Proceed to create the team since it doesn't already exist
	if err := database.DB.Create(&team).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to create team",
		})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "team created successfully",
		"data":    team,
	})
}

func CreateTeamA(c *gin.Context) {
	var team model.TeamA

	// Bind the incoming JSON to the team struct
	if err := c.ShouldBindJSON(&team); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind JSON",
		})
		return
	}

	var tournament model.Tournament
	if err := database.DB.First(&tournament, team.TournamentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "tournament not found",
		})
		return
	}
	// Check if the user exists
	var user model.User
	if err := database.DB.First(&user, team.PlayerID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	// Check if a team with the same name already exists
	var existingTeam model.TeamA
	if err := database.DB.Where("name = ?", team.Name).First(&existingTeam).Error; err == nil {
		// If no error, it means the team already exists
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "this team already exists, please create a different team",
		})
		return
	}

	// Proceed to create the team since it doesn't already exist
	if err := database.DB.Create(&team).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to create team",
		})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "team created successfully",
		"data":    team,
	})
}

func CreateTeamB(c *gin.Context) {
	var team model.TeamB
	// Bind the incoming JSON to the team struct
	if err := c.ShouldBindJSON(&team); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind JSON",
		})
		return
	}

	var tournament model.Tournament
	if err := database.DB.First(&tournament, team.TournamentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "tournament not found",
		})
		return
	}
	// Check if the user exists
	var user model.User
	if err := database.DB.First(&user, team.PlayerID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	// Check if a team with the same name already exists
	var existingTeam model.TeamB
	if err := database.DB.Where("name = ?", team.Name).First(&existingTeam).Error; err == nil {
		// If no error, it means the team already exists
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "this team already exists, please create a different team",
		})
		return
	}

	// Proceed to create the team since it doesn't already exist
	if err := database.DB.Create(&team).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to create team",
		})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "team created successfully",
		"data":    team,
	})
}
