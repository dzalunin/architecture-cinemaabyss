package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()
	router.GET("/health", health)
	api := router.Group("/api")
	{
		api.Any("/events/*proxyPath", reverseProxy(cfg.EventsURL))
		api.Any("/movies/*proxyPath", moviesProxy(cfg))
		api.Any("/users/*proxyPath", reverseProxy(cfg.MonolithURL))
		api.Any("/payments/*proxyPath", reverseProxy(cfg.MonolithURL))
		api.Any("/subscriptions/*proxyPath", reverseProxy(cfg.MonolithURL))
	}

	router.Run(":" + cfg.Port)
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": true})
}

func moviesProxy(cfg *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		var target = cfg.MonolithURL

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Float64()*100 < float64(cfg.MoviesPercent) {
			target = cfg.MoviesURL
		}

		log.Printf("Percent: %d, target: %v", cfg.MoviesPercent, target.Host)
		forwardRequest(target, context)
	}
}

func reverseProxy(target *url.URL) gin.HandlerFunc {
	return func(context *gin.Context) {
		forwardRequest(target, context)
	}
}

func forwardRequest(target *url.URL, c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(target)

	c.Request.URL.Scheme = target.Scheme
	c.Request.URL.Host = target.Host
	c.Request.Host = target.Host
	c.Request.URL.Path = strings.TrimSuffix(c.Request.URL.Path, "/")
	c.Request.URL.RawPath = strings.TrimSuffix(c.Request.URL.RawPath, "/")

	proxy.ServeHTTP(c.Writer, c.Request)
}
