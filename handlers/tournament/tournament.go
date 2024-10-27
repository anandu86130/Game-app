package tournament

import (
	"gaming/model"
	"net/http"

	database "gaming/database"

	"github.com/gin-gonic/gin"
)

// CreateTournament handles the creation of a new tournament.
func CreateTournament(c *gin.Context) {
	var tournament model.Tournament

	// Bind the incoming JSON request to the tournament model
	if err := c.ShouldBindJSON(&tournament); err != nil {
		// If binding fails, respond with a Bad Request status
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind json",
		})
		return
	}

	// Check if a tournament with the same name already exists
	var existingTournament model.Tournament
	if err := database.DB.Where("name = ?", tournament.Name).First(&existingTournament).Error; err == nil {
		// If no error, it means the tournament already exists
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "this tournament already exists, please create a different tournament",
		})
		return
	}

	// Create the new tournament in the database
	if err := database.DB.Create(&tournament).Error; err != nil {
		// If creation fails, respond with an Internal Server Error status
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to create tournament in the database",
		})
		return
	}

	// Respond with success and return the created tournament details
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "tournament created successfully",
		"data":    tournament,
	})
}

// ViewTournaments handles fetching the list of all tournaments along with their teams.
func ViewTournaments(c *gin.Context) {
	var tournaments []model.Tournament

	// Fetch all tournaments and preload their associated teams
	result := database.DB.Preload("TeamA").Preload("TeamB").Find(&tournaments)
	if result.Error != nil {
		// If fetching fails, respond with a Bad Request status
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to fetch tournaments",
		})
		return
	}

	var tournamentView []gin.H
	// Iterate through each tournament to build the response structure
	for _, tournament := range tournaments {
		// Prepare Team A details
		var teamAView []gin.H
		for _, teamA := range tournament.TeamA {
			teamADetails := gin.H{
				"team_id":   teamA.ID,
				"team_name": teamA.Name,
				"score":     teamA.Score,
			}
			teamAView = append(teamAView, teamADetails)
		}

		// Prepare Team B details
		var teamBView []gin.H
		for _, teamB := range tournament.TeamB {
			teamBDetails := gin.H{
				"team_id":   teamB.ID,
				"team_name": teamB.Name,
				"score":     teamB.Score,
			}
			teamBView = append(teamBView, teamBDetails)
		}

		// Construct tournament details with teams
		details := gin.H{
			"tournament_id":   tournament.ID,
			"tournament_name": tournament.Name,
			"start_time":      tournament.StartTime,
			"prize_pool":      tournament.PrizePool,
			"team_a":          teamAView,
			"team_b":          teamBView,
		}
		tournamentView = append(tournamentView, details)
	}

	// Respond with the list of tournaments and their details
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "fetched tournaments successfully",
		"data":    tournamentView,
	})
}

// JoinTournament handles a user's request to join a tournament.
func JoinTournament(c *gin.Context) {
	id := c.GetUint("userid") // Retrieve user ID from context

	// Struct to capture the incoming request for joining a tournament
	var Req struct {
		TournamentID uint `json:"tournament_id"` // tournament ID from request
	}

	// Bind the JSON request to the Req struct
	if err := c.ShouldBindJSON(&Req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to bind json",
		})
		return
	}

	// Find the tournament and preload TeamA and TeamB
	var tournament model.Tournament
	if err := database.DB.Preload("TeamA").Preload("TeamB").First(&tournament, Req.TournamentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to find tournament",
		})
		return
	}

	// Check if the user exists in the database
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	// Prepare the response with tournament details and team details
	var teamA []gin.H
	var teamB []gin.H

	// Collect details for Team A
	for _, team := range tournament.TeamA {
		teamA = append(teamA, gin.H{
			"team_id":   team.ID,
			"team_name": team.Name,
			"score":     team.Score,
		})
	}

	// Collect details for Team B
	for _, team := range tournament.TeamB {
		teamB = append(teamB, gin.H{
			"team_id":   team.ID,
			"team_name": team.Name,
			"score":     team.Score,
		})
	}

	// Construct tournament details for response
	tournamentDetails := gin.H{
		"tournament_id":   tournament.ID,
		"tournament_name": tournament.Name,
		"start_time":      tournament.StartTime,
		"prize_pool":      tournament.PrizePool,
		"team_a":          teamA,
		"team_b":          teamB,
	}

	// Respond with the tournament details and confirmation of joining
	c.JSON(http.StatusOK, gin.H{
		"status":          http.StatusOK,
		"message":         "You joined this tournament",
		"tournament_info": tournamentDetails,
	})
}