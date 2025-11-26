package http

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Subscription *SubscriptionHandler
}

func NewRouter(sub *SubscriptionHandler) *Router {
	return &Router{
		Subscription: sub,
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/subscriptions")
	{
		api.POST("", r.Subscription.Create)
		api.GET("/:id", r.Subscription.GetByID)
		api.PUT("/:id", r.Subscription.Update)
		api.DELETE("/:id", r.Subscription.Delete)
		api.GET("", r.Subscription.List)
		api.GET("/total", r.Subscription.GetTotalSum)
	}

	return router
}
