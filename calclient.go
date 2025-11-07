package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/go-resty/resty/v2"
)

// cacheEntry stores a cached parsed calendar and its timestamp
type cacheEntry struct {
	calendar  *ics.Calendar
	timestamp time.Time
}

// calendarCache stores cached calendar responses
var calendarCache = make(map[string]cacheEntry)
var cacheMutex sync.RWMutex

// httpClient is a shared resty client for all calendar requests (safe for concurrent use)
var httpClient = resty.New()

const cacheDuration = 5 * time.Minute

// generateCacheKey creates a unique cache key from the request parameters
func generateCacheKey(calTarget string, params url.Values) string {
	return fmt.Sprintf("%s:%s", calTarget, params.Encode())
}

// fetchCalendarData retrieves calendar data from the remote server, using cache if available.
//
// Parameters:
//   - calTarget: The target calendar to request (e.g., "leidencal", "delftcal")
//   - params: Additional query parameters to include in the request
// Returns the parsed calendar and an error if the request fails
//
func fetchCalendarData(calTarget string, params url.Values) (*ics.Calendar, error) {
	cacheKey := generateCacheKey(calTarget, params)

	// Check cache first
	cacheMutex.RLock()
	if entry, exists := calendarCache[cacheKey]; exists {
		if time.Since(entry.timestamp) < cacheDuration {
			cacheMutex.RUnlock()
			fmt.Printf("[Calendar] Cache hit for %s (age: %s)\n", calTarget, time.Since(entry.timestamp).Round(time.Second))
			return entry.calendar, nil
		}
	}
	cacheMutex.RUnlock()

	// Get the target URL from config
	targetURL := config.Targets[calTarget].URL

	// Send GET request with query parameters using shared client
	resp, err := httpClient.R().
		SetQueryParamsFromValues(params).
		Get(targetURL)

	if err != nil {
		return nil, fmt.Errorf("failed to send calendar request: %w", err)
	}

	// Log the response
	fmt.Printf("[Calendar] Fetched %s | %s | %d bytes\n", calTarget, resp.Status(), len(resp.Body()))

	// Parse the iCal data
	calendar, err := ics.ParseCalendar(strings.NewReader(string(resp.Body())))
	if err != nil {
		return nil, fmt.Errorf("failed to parse iCal data: %w", err)
	}

	// Store parsed calendar in cache
	cacheMutex.Lock()
	calendarCache[cacheKey] = cacheEntry{
		calendar:  calendar,
		timestamp: time.Now(),
	}
	cacheMutex.Unlock()

	return calendar, nil
}

// requestCalendar sends a request to the calendar server to fetch the specified calendar.
//
// Parameters:
//   - calType: The type of calendar to request (e.g., "tentamen", "hoorcollege")
//   - calTarget: The target calendar to request (e.g., "leidencal", "delftcal")
//   - params: Additional query parameters to include in the request
// Returns the filtered calendar and an error if the request fails
//
func requestCalendar(calType string, calTarget string, params url.Values) (*ics.Calendar, error) {
	// Fetch calendar data (from cache or remote server)
	calendar, err := fetchCalendarData(calTarget, params)
	if err != nil {
		return nil, err
	}

	// Create a new calendar for the filtered events
	filteredCalendar := ics.NewCalendar()

	// Copy all properties from original calendar
	filteredCalendar.CalendarProperties = append(filteredCalendar.CalendarProperties, calendar.CalendarProperties...)

	// Update cal name
	filteredCalendar.SetName(fmt.Sprintf("Filtered calendar - %s", calType))

	// Filter events by calType
	for _, event := range calendar.Events() {
		grouping := getICalEntryGrouping(event)
		if grouping == calType {
			// Remove course code from summary if present
			summary := event.GetProperty(ics.ComponentPropertySummary).Value
			if idx := strings.Index(summary, " - "); idx != -1 {
				// Found " - ", remove everything before it (course code)
				event.SetSummary(summary[idx+3:])
			}
			filteredCalendar.AddVEvent(event)
		}
	}

	return filteredCalendar, nil
}