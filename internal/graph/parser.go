package graph

import (
	"encoding/json"
	"os"
)

type osmElement struct {
	Type  string            `json:"type"`
	ID    int64             `json:"id"`
	Lat   float64           `json:"lat,omitempty"`
	Lon   float64           `json:"lon,omitempty"`
	Nodes []int64           `json:"nodes,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type osmResponse struct {
	Elements []osmElement `json:"elements"`
}

// BuildFromOSMFile reads an Overpass-format JSON file (streets) and builds a Graph.
func BuildFromOSMFile(path string) (*Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parsed osmResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	g := NewGraph()

	// First pass: register every node's coordinates.
	for _, el := range parsed.Elements {
		if el.Type == "node" {
			g.Nodes[el.ID] = Point{Lat: el.Lat, Lon: el.Lon}
		}
	}

	// Second pass: connect consecutive nodes within each way (street segment).
	edgeCount := 0
	skipped := 0
	for _, el := range parsed.Elements {
		if el.Type != "way" || len(el.Nodes) < 2 {
			continue
		}
		for i := 0; i < len(el.Nodes)-1; i++ {
			a, b := el.Nodes[i], el.Nodes[i+1]
			pa, okA := g.Nodes[a]
			pb, okB := g.Nodes[b]
			if !okA || !okB {
				skipped++
				continue
			}
			weight := Haversine(pa, pb)
			g.AddEdge(a, b, weight)
			edgeCount++
		}
	}

	_ = edgeCount
	_ = skipped
	return g, nil
}