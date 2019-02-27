package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/mholt/certmagic"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	prod       = flag.Bool("prod", false, "use it to be launch in production mode (https,...)")
	emailSSL   = flag.String("email_ssl", "stef@azsiaz.tech", "use it to change default email address to use for ssl certificate")
	stagingSSL = flag.Bool("staging_ssl", false, "use a staging env on for certmagic lib on Let's Encrypt")
)

func init() {
	flag.Parse()

	if *stagingSSL {
		certmagic.CA = certmagic.LetsEncryptStagingCA
	}
	if *prod {
		certmagic.Agreed = true
		certmagic.Email = *emailSSL

		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	store := openStore()
	defer store.Close()

	rtr := setupRouter()
	go launchWebServer(rtr)

	t := make(chan struct{})
	<-t
}

func launchWebServer(rtr *gin.Engine) {
	if *prod {
		log.Fatal(certmagic.HTTPS([]string{"azsiaz.cloud"}, rtr))
	} else {
		log.Fatal(rtr.Run(":8080"))
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, world")
	})

	return r
}

func openStore() *Store {
	log.Info("Opening database")
	store, err := NewStore("db.storm")
	if err != nil {
		log.
			WithError(err).
			Fatalln("Error opening DB")
	}
	log.Info("Database opened")

	return store
}
