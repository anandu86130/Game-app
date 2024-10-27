package result

import (
	database "gaming/database"
	"gaming/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TournamentResult handles the request to fetch tournament results and winning teams
func TournamentResult(c *gin.Context) {
	var tournaments []model.Tournament // Slice to hold tournaments

	// Fetch tournaments from the database, preloading associated TeamA and TeamB relationships
	result := database.DB.Preload("TeamA").Preload("TeamB").Find(&tournaments)
	if result.Error != nil { // Check for errors during fetching
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to fetch tournaments", 
		})
		return // Exit if an error occurs
	}

	var tournamentView []gin.H // Slice to hold formatted tournament results
	for _, tournament := range tournaments { // Iterate over each tournament
		var winningTeam gin.H    // Map to hold the winning team's details
		var winningScore float64  // Variable to track the highest score

		// Compare scores between TeamA and find the winning team
		for _, teamA := range tournament.TeamA {
			if teamA.Score > winningScore { // If TeamA's score is higher than the current winning score
				winningTeam = gin.H{ // Update winning team details
					"team_id":   teamA.ID,    
					"team_name": teamA.Name,   
					"score":     teamA.Score,   
				}
				winningScore = teamA.Score // Update winning score
			}
		}

		// Compare scores between TeamB and find the winning team
		for _, teamB := range tournament.TeamB {
			if teamB.Score > winningScore { // If TeamB's score is higher than the current winning score
				winningTeam = gin.H{ // Update winning team details
					"team_id":   teamB.ID,     
					"team_name": teamB.Name,  
					"score":     teamB.Score,  
				}
				winningScore = teamB.Score // Update winning score
			}
		}

		// Create a map for tournament details, including the winning team's information
		details := gin.H{
			"tournament_id":   tournament.ID,          
			"tournament_name": tournament.Name,        
			"start_time":      tournament.StartTime,  
			"prize_pool":      tournament.PrizePool,  
			"winning_team":    winningTeam,            
		}
		tournamentView = append(tournamentView, details) // Append tournament details to the tournament view
	}

	// Respond with the tournament results
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "This team has won in the tournament",
		"data":    tournamentView,
	})
}

// TournamentPriceDistribution determines the winning team and distributes the prize
func TournamentPriceDistribution(c *gin.Context) {
	var tournaments []model.Tournament // Slice to hold tournaments

	// Fetch tournaments from the database, preloading associated TeamA and TeamB relationships
	result := database.DB.Preload("TeamA").Preload("TeamB").Find(&tournaments)
	if result.Error != nil { // Check for errors during fetching
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to fetch tournaments", 
		})
		return // Exit if an error occurs
	}

	var tournamentView []gin.H // Slice to hold formatted tournament results
	for _, tournament := range tournaments { // Iterate over each tournament
		// Create a map for tournament details
		details := gin.H{
			"tournament_id":   tournament.ID,         
			"tournament_name": tournament.Name,     
			"start_time":      tournament.StartTime,  
			"prize_pool":      tournament.PrizePool, 
		}
		tournamentView = append(tournamentView, details) // Append tournament details to the tournament view
	}

	// Respond with the tournament prize distribution results
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Congratulations, you have won in this tournament",
		"data":    tournamentView,
	})
}