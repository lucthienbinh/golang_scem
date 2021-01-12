package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/api/router"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// RunServer will start 2 server for app and web
func RunServer() {
	// gin.SetMode(gin.ReleaseMode)
	// export GIN_MODE=debug

	webServer := &http.Server{
		Addr:         os.Getenv("WEB_PORT"),
		Handler:      webRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// appServer := &http.Server{
	// 	Addr:         os.Getenv("APP_PORT"),
	// 	Handler:      appRouter(),
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	g.Go(func() error {
		err := webServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	// g.Go(func() error {
	// 	err := appServer.ListenAndServe()
	// 	if err != nil && err != http.ErrServerClosed {
	// 		log.Fatal(err)
	// 	}
	// 	return err
	// })

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func webRouter() http.Handler {
	e := gin.Default()
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"http://127.0.0.1:3000"}
		config.AllowCredentials = true
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-CSRF-Token", "Accept"}
		e.Use(cors.New(config))
	}

	e.Static("/api/images", os.Getenv("IMAGE_FILE_PATH"))
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	e.MaxMultipartMemory = 8 << 20 // 8 MiB

	webAuth := e.Group("/web-auth")
	router.WebAuthRoutes(webAuth)

	api := e.Group("/api")
	// Active web auth
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		api.Use(middleware.ValidateWebSession())
	}
	router.UserRoutes(api)
	router.WebOrderRoutes(api)
	router.OrderShipRoutes(api)
	return e
}

func appRouter() http.Handler {
	e := gin.Default()
	e.Static("/api/image", "./public/upload/images")

	appAuth := e.Group("/app-auth")
	router.AppAuthRoutes(appAuth)

	fcmAuth := e.Group("/fcm-auth")
	router.AppFMCToken(fcmAuth)

	api := e.Group("/api")

	// Active app auth
	if os.Getenv("RUN_APP_AUTH") == "yes" {
		api.Use(middleware.ValidateAppToken())
	}
	router.UserRoutes(api)

	return e
}
