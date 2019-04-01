package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mine-stats/handler"
	"mine-stats/handler/middleware"
	"mine-stats/jobs"
	"mine-stats/store"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mholt/certmagic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	prod       = flag.Bool("prod", false, "use it to be launch in production mode")
	https      = flag.Bool("https", false, "use it to use https in production mode")
	emailSSL   = flag.String("email_ssl", "stef@azsiaz.tech", "use it to change default email address to use for ssl certificate")
	stagingSSL = flag.Bool("staging_ssl", false, "use a staging env on for certmagic lib on Let's Encrypt")
	metrics    = flag.Bool("metrics", false, "expose a prometheus metrics endpoint")
	firstAdmin = flag.Bool("first_admin", false, "First signup is an admin user")
)

func init() {
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	if *stagingSSL {
		certmagic.CA = certmagic.LetsEncryptStagingCA
	}
	if *https {
		certmagic.Agreed = true
		certmagic.Email = *emailSSL
	}
	if *prod {
		//	For later
		log.Info("Launching mine-stats in production mode")
	}
}

func main() {
	st := openStore()
	defer st.Close()

	go setupJobs(st)

	rtr := setupRouter(st)
	go launchWebServer(rtr)

	<-exit()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobs.ShutDownJob()

	if err := rtr.Shutdown(ctx); err != nil {
		log.WithError(err).Fatalln()
	}
}

// maybe https://github.com/jiansoft/robin later
func setupJobs(st *store.Store) {
	srvs, err := st.GetMinecraftServerList()
	if err != nil {
		os.Exit(1)
	}

	for _, srv := range srvs {
		j := jobs.NewJob(&srv)
		jobs.AddJob(j)

		<-time.After(500 * time.Millisecond)
	}
	log.Infoln("Done loading jobs")
}

func setupRouter(st *store.Store) *echo.Echo {
	r := echo.New()
	h := handler.NewHandler(st, *prod)

	r.Pre(echoMiddleware.RemoveTrailingSlash())
	r.Use(
		echoMiddleware.RequestID(),
		echoMiddleware.Logger(),
		echoMiddleware.Recover(),
		echoMiddleware.Secure(),
	)

	if *metrics {
		r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	api := r.Group("/api")
	{
		authAPI := api.Group("/auth")
		{
			authAPI.POST("/signup", h.SignUpHandler)
			authAPI.POST("/login", h.LoginHandler)
			authAPI.GET("/logout", h.LogoutHandler, middleware.CheckAuth)
			authAPI.GET("/me", h.MeHandler, middleware.CheckAuth)
		}
		srvAPI := api.Group("/server", middleware.CheckAuth)
		{
			srvAPI.GET("", h.ListOwnServer)
			srvAPI.GET("/:id", h.OneOwnServer)
			srvAPI.POST("", h.AddServer)
			srvAPI.PUT("", h.UpdateServer)
			srvAPI.DELETE("/:id", h.DeleteServer)

			//https://github.com/asdine/storm/issues/212
			//statsApi := srvApi.Group("/stats")
			//{
			//
			//}
		}
		admAPI := api.Group("/admin", middleware.CheckAuth, middleware.CheckAdmin)
		{
			admAPI.GET("/server", h.AdminListServer)
			admAPI.GET("/server/:id", h.AdminOneServer)
			admAPI.DELETE("/server/:id", h.AdminDeleteServer)
			admAPI.GET("/user", h.AdminListUser)
			admAPI.GET("/user/:id", h.AdminOneUser)
			admAPI.DELETE("/user/:id", h.AdminDeleteUser)
		}
	}

	return r
}

func launchWebServer(rtr *echo.Echo) {
	if *prod && *https {
		log.Fatal(certmagic.HTTPS([]string{"azsiaz.cloud"}, rtr))
	} else {
		log.Fatal(http.ListenAndServe(":8080", rtr))
	}
}

func openStore() *store.Store {
	log.Info("Opening database")
	st, err := store.NewStore("db.storm", *firstAdmin)
	if err != nil {
		log.
			WithError(err).
			Fatalln("Error opening DB")
	}
	log.Info("Database opened")

	return st
}

func exit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	return ch
}
