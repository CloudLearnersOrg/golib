package sentinel

import "strings"

// ParseStringToHosts converts a comma-separated string of hosts into a slice of hosts.
// This is useful for loading host configurations from environment variables.
//
// Example:
//
//	hosts := sentinel.ParseStringToHosts("api.example.com,localhost:3000")
//	// hosts = []string{"api.example.com", "localhost:3000"}
//
// Empty or whitespace-only entries are automatically filtered out.
func ParseStringToHosts(hostsString string) []string {
	if hostsString == "" {
		return []string{}
	}

	// Split on commas and trim whitespace
	hosts := []string{}
	for _, host := range strings.Split(hostsString, ",") {
		trimmedHost := strings.TrimSpace(host)
		if trimmedHost != "" {
			hosts = append(hosts, trimmedHost)
		}
	}

	return hosts
}
