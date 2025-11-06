package main

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var regexTarget *regexp.Regexp

// startWebServer initializes and starts the web server.
func startWebServer() error {
	fmt.Println("Starting web server...")

	// Compile the regex
	regexTarget = regexp.MustCompile(`^` + config.BaseURLPath + `/([^/]+)/ical$`)

	router = gin.Default()
	setupRoutes()

	err := router.Run(config.BindAddress)
	return err
}

// setupRoutes configures the HTTP routes for the web server.
func setupRoutes() {
	basePath := config.BaseURLPath
	
	router.GET(basePath+"/health", handleHealthCheck)
	
	// Setup configured target routes
	for name := range config.Targets {
		routePath := fmt.Sprintf("%s/%s/%s", basePath, name, "ical")
		router.GET(routePath, handleCalendarRequest)
	}
}

// handleHealthCheck responds with a simple health status.
func handleHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// handleCalendarRequest processes calendar requests based on the URL and query parameters.
func handleCalendarRequest(c *gin.Context) {
	// Get the target from the URL
	matches := regexTarget.FindStringSubmatch(c.Request.URL.Path)
	if len(matches) < 2 {
		c.String(400, "Invalid URL path")
		return
	}
	// Extract the captured group
	calTarget := matches[1]

	// Extract the calender type
	calType := c.Query("calType")
	if calType == "" {
		c.String(400, "Missing required parameter: calType")
		return
	}
	
	// Check if target exists
	if _, exists := config.Targets[calTarget]; !exists {
		c.String(404, "Unknown calendar target: %s", calTarget)
		return
	}
	// Check if calType is valid
	if _, exists := config.Groupings[calType]; !exists {
		c.String(400, "Unknown calendar type: %s", calType)
		return
	}

	// Get all the request params without calType
	params := c.Request.URL.Query()
	delete(params, "calType")

	// Request and filter the calendar
	calendar, err := requestCalendar(calType, calTarget, params)
	if err != nil {
		c.String(500, "Failed to fetch calendar: %s", err.Error())
		return
	}

	// Serialize the calendar to iCal format
	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s-%s.ics\"", calTarget, calType))
	c.Header("Cache-Control", "no-cache, must-revalidate")
	c.String(200, calendar.Serialize())
}