package main

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"

	emailmodule "github.com/sujithsss01/email_server/app_server/internal/modules/email"
)

type Module interface {
	RegisterRoutes(*gin.RouterGroup)
}

type App struct {
	modules []Module
}

func (a *App) registerRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	appRouter := router.Group("/api/v1")
	for _, module := range a.modules {
		module.RegisterRoutes(appRouter)
	}
}

func NewApp(session *gocql.Session) *App {
	modules := []Module{
		emailmodule.NewEmailModule(session),
	}
	return &App{
		modules: modules,
	}
}
