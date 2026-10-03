package graph

import "math"

type Point struct {
	Lat float64
	Lon float64
}

type Edge struct {
	To     int64
	Weight float64 // meters
}

type Graph struct {
	Nodes     map[int64]Point
	Adjacency map[int64][]Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:     make(map[int64]Point),
		Adjacency: make(map[int64][]Edge),
	}
}

// Haversine computes the great-circle distance between two points in meters.
func Haversine(a, b Point) float64 {
	const earthRadiusMeters = 6371000.0

	lat1 := a.Lat * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180
	dLat := (b.Lat - a.Lat) * math.Pi / 180
	dLon := (b.Lon - a.Lon) * math.Pi / 180

	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))

	return earthRadiusMeters * c
}

func (g *Graph) AddEdge(from, to int64, weight float64) {
	g.Adjacency[from] = append(g.Adjacency[from], Edge{To: to, Weight: weight})
	g.Adjacency[to] = append(g.Adjacency[to], Edge{To: from, Weight: weight})
}

// NearestNode does a brute-force search for the graph node closest to the given point.
// Fine for a graph this size; would need a spatial index (k-d tree, grid) at larger scale.
func (g *Graph) NearestNode(p Point) (int64, float64) {
	var bestID int64
	bestDist := math.MaxFloat64
	for id, point := range g.Nodes {
		d := Haversine(p, point)
		if d < bestDist {
			bestDist = d
			bestID = id
		}
	}
	return bestID, bestDist
}