package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)
import "time"
// Bounded box around Howard University: Shaw, Le Droit Park, part of Columbia Heights/U Street
const (
	south = 38.912
	west  = -77.028
	north = 38.927
	east  = -77.010
)
const (
	streetSouth = 38.916
	streetWest  = -77.024
	streetNorth = 38.923
	streetEast  = -77.015
)

const overpassURL = "https://overpass.kumi.systems/api/interpreter"


func fetchOverpass(query string) ([]byte, error) {
	data := url.Values{}
	data.Set("data", query)

	var lastErr error
	for attempt := 1; attempt <= 4; attempt++ {
		req, err := http.NewRequest("POST", overpassURL, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("User-Agent", "FitRouteDC/1.0 (personal portfolio project)")

		client := &http.Client{Timeout: 180 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			fmt.Printf("Attempt %d failed (network error): %v. Retrying...\n", attempt, err)
			time.Sleep(time.Duration(attempt) * 5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return io.ReadAll(resp.Body)
		}

		body, _ := io.ReadAll(resp.Body)
		lastErr = fmt.Errorf("overpass returned %d: %s", resp.StatusCode, string(body))

		if resp.StatusCode == 504 || resp.StatusCode == 503 {
			fmt.Printf("Attempt %d failed (server busy, %d). Retrying in %ds...\n", attempt, resp.StatusCode, attempt*5)
			time.Sleep(time.Duration(attempt) * 5 * time.Second)
			continue
		}

		return nil, lastErr
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}

func main() {
	bbox := fmt.Sprintf("%f,%f,%f,%f", south, west, north, east)
	
		streetBbox := fmt.Sprintf("%f,%f,%f,%f", streetSouth, streetWest, streetNorth, streetEast)

	fmt.Println("Fetching street network...")
	streetQuery := fmt.Sprintf(`
		[out:json][timeout:90];
		way["highway"~"^(primary|secondary|tertiary|residential|unclassified|living_street)$"](%s);
		(._;>;);
		out body;
	`, streetBbox)

	streetData, err := fetchOverpass(streetQuery)
	if err != nil {
		fmt.Println("Error fetching streets:", err)
		os.Exit(1)
	}
	if err := os.WriteFile("data_streets.json", streetData, 0644); err != nil {
		fmt.Println("Error saving streets:", err)
		os.Exit(1)
	}
	fmt.Println("Saved data_streets.json")

	fmt.Println("Fetching gyms...")
	gymQuery := fmt.Sprintf(`
    [out:json][timeout:60];
    (
      node["leisure"="fitness_centre"](%s);
      node["leisure"="sports_centre"](%s);
      way["leisure"="fitness_centre"](%s);
      way["leisure"="sports_centre"](%s);
    );
    out center;
`, bbox, bbox, bbox, bbox)

	gymData, err := fetchOverpass(gymQuery)
	if err != nil {
		fmt.Println("Error fetching gyms:", err)
		os.Exit(1)
	}
	if err := os.WriteFile("data_gyms.json", gymData, 0644); err != nil {
		fmt.Println("Error saving gyms:", err)
		os.Exit(1)
	}
	fmt.Println("Saved data_gyms.json")

	var gymParsed struct {
		Elements []json.RawMessage `json:"elements"`
	}
	json.Unmarshal(gymData, &gymParsed)
	fmt.Printf("Found %d gyms in the bounded area.\n", len(gymParsed.Elements))
}