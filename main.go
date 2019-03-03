//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mholt/certmagic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"mine-stats/handler"
	"mine-stats/handler/middleware"
	"mine-stats/store"
	"net/http"
	"os"
)

var (
	prod       = flag.Bool("prod", false, "use it to be launch in production mode")
	https       = flag.Bool("https", false, "use it to use https in production mode")
	emailSSL   = flag.String("email_ssl", "stef@azsiaz.tech", "use it to change default email address to use for ssl certificate")
	stagingSSL = flag.Bool("staging_ssl", false, "use a staging env on for certmagic lib on Let's Encrypt")
	metrics    = flag.Bool("metrics", false, "expose a prometheus metrics endpoint")
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
	}
}

func main() {
	st := openStore()
	defer st.Close()

	go setupJobs(st)

	rtr := setupRouter(st)
	go launchWebServer(rtr)

	t := make(chan struct{})
	<-t
}

func setupJobs(st *store.Store) {
	srvs, err := st.GetMinecraftServerList()
	if err != nil {
		os.Exit(1)
	}

	for _, srv := range srvs {
		println(srv.Name)
	}
}

func setupRouter(st *store.Store) *echo.Echo{
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
		r.GET("/index.html", ServeIndex)
		r.GET("/", ServeIndex)
		r.GET("/*", echo.WrapHandler(Handler))
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
			srvApi.GET("", h.ListServer)
			srvApi.POST("", h.AddServer)
			srvApi.PUT("", h.UpdateServer)
			srvApi.DELETE("", h.DeleteServer)

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
	st, err := store.NewStore("db.storm")
	if err != nil {
		log.
			WithError(err).
			Fatalln("Error opening DB")
	}
	log.Info("Database opened")

	return st
}

func ServeIndex(c echo.Context) error {
	htmlb, err := ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	// convert to string
	html := string(htmlb)
	return c.HTML(http.StatusOK, html)
}