package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AddStandartGin - adds standart gin values to specialized
func (r *Router) addStandartGin(c *gin.Context, temp gin.H) gin.H {
	for k, v := range r.getStandartGin(c) {
		temp[k] = v
	}
	return temp
}

// GetStandartGin - used for adding standart gin values to all pages
func (r *Router) getStandartGin(c *gin.Context) gin.H {
	standard := gin.H{
		"version": r.generateSessionToken(),
	}

	return standard
}

func (r *Router) generateSessionToken() string {
	return fmt.Sprint(time.Now().UnixNano())
}

func (r *Router) render(c *gin.Context, templateName string, data gin.H) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, r.addStandartGin(c, data))
	}
}

// Shows start page
func (r *Router) showIndexPage(c *gin.Context) {
	r.render(c, "index.html", gin.H{
		"title": "Monopoly",
	})
}

func (r *Router) showGamePage(c *gin.Context) {
	r.render(c, "game.html", gin.H{
		"title":   "Monopoly",
		"playing": true,
	})
}
