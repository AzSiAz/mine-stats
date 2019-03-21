//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

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
	"mine-stats/public"
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

	if *prod {
		r.GET("/index.html", h.ServeIndex)
		r.GET("/", h.ServeIndex)
		r.GET("/*", echo.WrapHandler(public.Handler))
	} else {
		r.Static("/", "public/dist")
	}

	if *metrics {
		r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	api := r.Group("/api")
	{
		authApi := api.Group("/auth")
		{
			authApi.POST("/signup", h.SignUpHandler)
			authApi.POST("/login", h.LoginHandler)
			authApi.GET("/logout", h.LogoutHandler, middleware.CheckAuth)
			authApi.GET("/me", h.MeHandler, middleware.CheckAuth)
		}
		srvApi := api.Group("/server", middleware.CheckAuth)
		{
			srvApi.GET("", h.ListOwnServer)
			srvApi.GET("/:id", h.OneOwnServer)
			srvApi.POST("", h.AddServer)
			srvApi.PUT("", h.UpdateServer)
			srvApi.DELETE("/:id", h.DeleteServer)

			//https://github.com/asdine/storm/issues/212
			//statsApi := srvApi.Group("/stats")
			//{
			//
			//}
		}
		admApi := api.Group("/admin", middleware.CheckAuth, middleware.CheckAdmin)
		{
			admApi.GET("/server", h.AdminListServer)
			admApi.GET("/server/:id", h.AdminOneServer)
			admApi.DELETE("/server/:id", h.AdminDeleteServer)
			admApi.GET("/user", h.AdminListUser)
			admApi.GET("/user/:id", h.AdminOneUser)
			admApi.DELETE("/user/:id", h.AdminDeleteUser)
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
