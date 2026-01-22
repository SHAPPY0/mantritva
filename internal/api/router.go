package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/config"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRouter(db *gorm.DB, mongoDB *config.MongoDB) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(CORSMiddleware())
	r.Static("/static", "./statics")

	repos := repository.NewRepos(db, mongoDB)
	services := service.NewServices(repos)

	v1 := r.Group("/api/v1")
	{
		countryRouter(v1, db)
		stateRouter(v1, db)
		districtRouter(v1, db)
		loksabhaRouter(v1, db)
		vidhanSabhaRouter(v1, db)
		partyRouter(v1, repos)
		allianceRouter(v1, db)
		memberRouter(v1, repos)
		rulersRouter(v1, repos)
		leadershipRouter(v1, repos)
		newsPressRouter(v1, services)
		feedRouter(v1, services)
		userRouter(v1, repos)
		consoleRouter(v1, services)
		entityRouter(v1, repos)
		grievanceRouter(v1, services)
		reportRouter(v1, services)
		eventRouter(v1, services)
	}
	return r
}
