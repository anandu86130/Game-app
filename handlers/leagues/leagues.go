package leagues

import (
    "gaming/model"
    "net/http"

    database "gaming/database"
    "github.com/gin-gonic/gin"
)

// CreateLeagues handles the creation of a new league.
func CreateLeagues(c *gin.Context) {
    var league model.League
    // Bind JSON request to the league model
    if err := c.ShouldBindJSON(&league); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "status":  http.StatusBadRequest,
            "message": "failed to bind json",
        })
        return
    }

    // Check if a league with the same name already exists
    var existingleague model.League
    if err := database.DB.Where("name = ?", league.Name).First(&existingleague).Error; err == nil {
        // If no error, the league already exists
        c.AbortWithStatusJSON(http.StatusConflict, gin.H{
            "status":  http.StatusConflict,
            "message": "this league already exists, please create a different league",
        })
        return
    }

    // Create a new league in the database
    if err := database.DB.Create(&league).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "status":  http.StatusInternalServerError,
            "message": "failed to create league in the database",
        })
        return
    }

    // Respond with success and return the created league details
    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "league created successfully",
        "data":    league,
    })
}


// ViewLeagues handles fetching the list of all leagues and their associated teams.
func ViewLeagues(c *gin.Context) {
    var leagues []model.League
    result := database.DB.Preload("Teams").Find(&leagues) // Preload the 'Teams' relationship
    if result.Error != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "status":  http.StatusBadRequest,
            "message": "failed to fetch leagues",
        })
        return
    }

    var leagueview []gin.H
    for _, fetchleagues := range leagues {
        // Create a slice to hold team details for each league
        var teamview []gin.H
        for _, team := range fetchleagues.Teams {
            teamDetails := gin.H{
                "team_id":   team.ID,
                "team_name": team.Name,
                "score":     team.Score,
            }
            teamview = append(teamview, teamDetails)
        }

        // Populate league details with teams
        details := gin.H{
            "league_id":   fetchleagues.ID,
            "league_name": fetchleagues.Name,
            "start_time":  fetchleagues.StartTime,
            "prize_pool":  fetchleagues.PrizePool,
            "teams":       teamview,
        }
        leagueview = append(leagueview, details)
    }

    // Respond with the list of leagues and their teams
    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "fetched leagues successfully",
        "data":    leagueview,
    })
}


// JoinLeague handles the request for a user to join a league.
func JoinLeague(c *gin.Context) {
    id := c.GetUint("userid") // Retrieve user ID from context (usually set during authentication)
    
    // Struct for the incoming request containing the league ID
    var Req struct {
        LeagueID uint `json:"league_id"`
    }

    // Bind the JSON request to the Req struct
    if err := c.ShouldBindJSON(&Req); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "status":  http.StatusBadRequest,
            "message": "failed to bind json",
        })
        return
    }

    // Find the league and preload its teams
    var league model.League
    if err := database.DB.Preload("Teams").First(&league, Req.LeagueID).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "status":  http.StatusInternalServerError,
            "message": "failed to find league",
        })
        return
    }

    // Check if the user exists
    var user model.User
    if err := database.DB.First(&user, id).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "status":  http.StatusNotFound,
            "message": "user not found",
        })
        return
    }

    // Create a slice to hold team details
    var teamview []gin.H
    for _, team := range league.Teams {
        teamDetails := gin.H{
            "team_id":   team.ID,
            "team_name": team.Name,
            "score":     team.Score,
        }
        teamview = append(teamview, teamDetails)
    }

    // Prepare the response with league and team details
    leagueDetails := gin.H{
        "league_id":   league.ID,
        "league_name": league.Name,
        "start_time":  league.StartTime,
        "prize_pool":  league.PrizePool,
        "teams":       teamview,
    }

    // Respond with the league details after joining
    c.JSON(http.StatusOK, gin.H{
        "status":         http.StatusOK,
        "message":        "You joined this league",
        "league_details": leagueDetails,
    })
}
