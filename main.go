package main

import (
	"gaming/config"
	"gaming/handlers/leagues"
	"gaming/handlers/result"
	team "gaming/handlers/team"
	"gaming/handlers/tournament"
	"gaming/handlers/user"
	"gaming/jwt"

	database "gaming/database"

	"github.com/gin-gonic/gin"
)

func main() {
	//connect env
	config.LoadEnv()
	//connect database
	database.DBconnect()

	//initialize gin
	r := gin.Default()
	//user authentication
	r.POST("/user/signup", user.Signup)
	r.POST("/user/verification", user.VerifyOTP)
	r.POST("/user/login", user.Login)
	//user middleware
	r.GET("/user/profile", jwt.AuthMiddleware("user"), user.UserProfile)
	r.PATCH("/user/profile", jwt.AuthMiddleware("user"), user.EditUser)
	r.POST("/user/leagues", jwt.AuthMiddleware("user"), leagues.CreateLeagues)
	r.GET("/user/leagues", jwt.AuthMiddleware("user"), leagues.ViewLeagues)
	r.POST("/user/league/team", jwt.AuthMiddleware("user"), team.CreateTeam)
	r.POST("/user/leagues/join", jwt.AuthMiddleware("user"), leagues.JoinLeague)
	r.POST("/user/tournament", jwt.AuthMiddleware("user"), tournament.CreateTournament)
	r.GET("/user/tournament", jwt.AuthMiddleware("user"), tournament.ViewTournaments)
	r.POST("/user/tournament/teamA", jwt.AuthMiddleware("user"), team.CreateTeamA)
	r.POST("/user/tournament/teamB", jwt.AuthMiddleware("user"), team.CreateTeamB)
	r.POST("user/tournament/join", jwt.AuthMiddleware("user"), tournament.JoinTournament)
	r.GET("/user/leagues/result", jwt.AuthMiddleware("user"), result.LeagueResult)
	r.GET("/user/leagues/price", jwt.AuthMiddleware("user"), result.PriceDistribution)
	r.GET("/user/tournament/result", jwt.AuthMiddleware("user"), result.TournamentResult)
	r.GET("/user/tournament/price", jwt.AuthMiddleware("user"), result.TournamentPriceDistribution)

	r.Run(":8080")
}
