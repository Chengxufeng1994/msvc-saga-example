package product

import (
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/middleware"
	v1 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/product/controller/v1"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	app              *application.ProductApplication
	jwtAuthenticator *middleware.JwtAuthenticator
	engine           *gin.Engine
}

func NewRouter(engine *gin.Engine, app *application.ProductApplication, jwtAuthenticator *middleware.JwtAuthenticator) *Router {
	return &Router{
		app:              app,
		jwtAuthenticator: jwtAuthenticator,
		engine:           engine,
	}
}

func (r *Router) RegisterRoutes() {
	// K8s probe for kubernetes health checks.
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The server is up and running.")
	})

	// prometheus probe for prometheus pull;
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Handling a page not found endpoint -.
	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "The requested page is not found. Please try later!"})
	})

	productController := v1.NewProductController(r.app.ProductService)
	v1 := r.engine.Group("/api/v1")
	productGroup := v1.Group("/product")
	productGroup.Use(r.jwtAuthenticator.Auth())
	{
		productGroup.POST("/", productController.CreateProduct)
		productGroup.GET("/", productController.ListProducts)
		productGroup.GET("/:product_id", productController.ListProducts)
	}
}
