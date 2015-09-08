package main

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func httpServer() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if cfg.DevMode {
		log.Info("Use Gin Logger")
		router.Use(gin.Logger())
	}
	if cfg.Sentry != "" && Raven != nil {
		log.Info("Use Custom Recovery")
		router.Use(recovery(Raven))
	} else {
		log.Info("Use Gin Recovery")
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
			for _, item := range c.Errors {
				packet := raven.NewPacket(item.Err.Error(), &raven.Message{Message: item.Err.Error(), Params: []interface{}{item.Meta}}, raven.NewHttp(c.Request))
				_, ch := client.Capture(packet, nil)
				if err := <-ch; err != nil {
					log.Error("Gin Error: ", err)
				}
			}
			if rval := recover(); rval != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)

				message := fmt.Sprint(rval)
				trace := raven.NewStacktrace(0, 2, nil)
				packet := raven.NewPacket(message, raven.NewException(errors.New(message), trace), raven.NewHttp(c.Request))
				client.Capture(packet, nil)
				log.Error("Gin Error: ", message)
			}
		}()
		c.Next()
	}
}
