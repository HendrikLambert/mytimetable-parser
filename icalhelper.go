package main

import (
	"fmt"
	"slices"
	"strings"

	ics "github.com/arran4/golang-ical"
)

// detectICalEntryType extracts the event type from the description field.
// The type is expected to be in the format "Type: <type>" on any line.
// Returns an empty string if no type is found.
func detectICalEntryType(e *ics.VEvent) string {
	descProp := e.GetProperty(ics.ComponentPropertyDescription)
	if descProp == nil {
		return ""
	}

	desc := descProp.Value

	// Search through all lines of the description
	lines := strings.Split(desc, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this line starts with "Type: "
		if strings.HasPrefix(line, "Type: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Type: "))
		}
	}

	// No type found
	return ""
}

// matchEventToGrouping matches an event type to a main grouping category.
// It searches through the configured groupings to find which category the event type belongs to.
// Returns the grouping name (e.g., "tentamen", "hoorcollege") or the default group if no match is found.
func matchEventToGrouping(eventType string) string {
	// Iterate through all groupings
	for groupName, types := range config.Groupings {
		// Check if the event type matches any type in this grouping
		if slices.Contains(types, eventType) {
			return groupName
		}
	}

	fmt.Printf("[iCal] Unknown event type '%s', using default group '%s'\n", eventType, config.DefaultGroup)
	return config.DefaultGroup
}

// getICalEntryGrouping extracts the event type and matches it to a grouping.
// Returns the grouping name or the default group if no match is found.
func getICalEntryGrouping(e *ics.VEvent) string {
	eventType := detectICalEntryType(e)
	if eventType == "" {
		return config.DefaultGroup
	}

	return matchEventToGrouping(eventType)
}