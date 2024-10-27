package result

import (
	database "gaming/database"
	"gaming/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LeagueResult handles the request to fetch league results along with team details
func LeagueResult(c *gin.Context) {
	var league []model.League // Slice to hold the leagues

	// Fetch leagues from the database, preloading the associated Teams
	result := database.DB.Preload("Teams").Find(&league)
	if result.Error != nil { // Check for errors during fetching
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to fetch leagues", // Error message if fetching fails
		})
		return // Exit if an error occurs
	}

	var leagueview []gin.H                // Slice to hold formatted league results
	for _, fetchleagues := range league { // Iterate over each league
		// Create a slice to hold team details for each league
		var teamview []gin.H
		for _, team := range fetchleagues.Teams { // Iterate over teams in the league
			// Create a map to hold team details
			teamDetails := gin.H{
				"team_id":   team.ID,
				"team_name": team.Name,
				"score":     team.Score,
			}
			teamview = append(teamview, teamDetails) // Append team details to the slice
		}

		// Create a map for league details including teams
		details := gin.H{
			"league_name": fetchleagues.Name,
			"prize_pool":  fetchleagues.PrizePool,
			"teams":       teamview,
		}
		leagueview = append(leagueview, details) // Append league details to the league view
	}

	// Respond with the league results
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "league result fetched successfully",
		"data":    leagueview, // Data containing leagues and teams
	})
}

// PriceDistribution handles the request to fetch prize distribution based on team scores
func PriceDistribution(c *gin.Context) {
	var league []model.League // Slice to hold the leagues

	// Fetch leagues from the database, preloading the associated Teams
	result := database.DB.Preload("Teams").Find(&league)
	if result.Error != nil { // Check for errors during fetching
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to fetch leagues",
		})
		return // Exit if an error occurs
	}

	var leagueview []gin.H                // Slice to hold formatted league results
	for _, fetchleagues := range league { // Iterate over each league
		// Create a slice to hold team details for each league
		var teamview []gin.H
		var winningTeam *model.Team               // Pointer to hold the winning team
		for _, team := range fetchleagues.Teams { // Iterate over teams in the league
			// Create a map to hold team details
			teamDetails := gin.H{
				"team_id":   team.ID,
				"team_name": team.Name,
				"score":     team.Score,
			}
			teamview = append(teamview, teamDetails) // Append team details to the slice

			// Check if this team has the highest score
			if winningTeam == nil || team.Score > winningTeam.Score {
				winningTeam = &team // Set winning team to the current team
			}
		}

		// Determine the prize distribution based on the winning team
		var prizeDistribution gin.H
		if winningTeam != nil { // If a winning team is found
			prizeDistribution = gin.H{
				"winning_team_id":   winningTeam.ID,
				"winning_team_name": winningTeam.Name,
				"prize":             fetchleagues.PrizePool,
			}
		}

		// Create a map for league details including teams and prize distribution
		details := gin.H{
			"league_name":        fetchleagues.Name,
			"prize_pool":         fetchleagues.PrizePool,
			"teams":              teamview,
			"prize_distribution": prizeDistribution,
		}
		leagueview = append(leagueview, details) // Append league details to the league view
	}

	// Respond with the prize distribution results
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Congratulations, this team has won in this league",
		"data":    leagueview, // Data containing leagues, teams, and prize distribution
	})
}
