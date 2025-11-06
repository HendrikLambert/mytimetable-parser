package main

import (
	"fmt"
	"slices"
	"strings"

	ics "github.com/arran4/golang-ical"
)

// detectICalEntryType extracts the event type from the description field.
// The type is expected to be on the first line in the format "Type: <type>".
// Returns an empty string if no type is found.
func detectICalEntryType(e *ics.VEvent) string {
	descProp := e.GetProperty(ics.ComponentPropertyDescription)
	if descProp == nil {
		return ""
	}

	desc := descProp.Value

	// Get the first line of the description
	firstLine := strings.Split(desc, "\n")[0]
	firstLine = strings.TrimSpace(firstLine)

	// Check if it starts with "Type: "
	if strings.HasPrefix(firstLine, "Type: ") {
		return strings.TrimSpace(strings.TrimPrefix(firstLine, "Type: "))
	} else {
		fmt.Printf("[iCal] Unknown event type '%s'\n", firstLine)
		return ""
	}

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