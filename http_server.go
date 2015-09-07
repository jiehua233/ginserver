package main

import (
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func httpServer() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if cfg.DevMode {
		router.Use(gin.Logger())
	}
	if cfg.Sentry != "" && Raven != nil {
		router.Use(recovery(Raven))
	} else {
		router.Use(gin.Recovery())
	}

	router.GET("/", handler)
	log.Info("Start Gin Http Server at: ", cfg.Server)
	log.Error("Gin Http Server Error: ", router.Run(cfg.Server))
}

// 参考: http://golang.org/pkg/net/http/#Request
// HTTP defines that header names are case-insensitive.
// The request parser implements this by canonicalizing the
// name, making the first character and any characters
// following a hyphen uppercase and the rest lowercase.
func handler(c *gin.Context) {
	r := c.Request // *http.Request
	h := r.Header  // http.Header  map[string][]string
	h.Add("RemoteAddr", r.RemoteAddr)
	c.JSON(http.StatusOK, h)
}

func recovery(client *raven.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {

		}()
		c.Next()
	}
}
