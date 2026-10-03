package main

import (
	"fmt"
	"time"
	"fitroute-dc/internal/graph"
)

func main() {
	g, err := graph.BuildFromOSMFile("data_streets.json")
	if err != nil {
		fmt.Println("Error building graph:", err)
		return
	}

		startPoint := graph.Point{Lat: 38.9218, Lon: -77.0195} // near Howard's main campus
	endPoint := graph.Point{Lat: 38.9227058, Lon: -77.0225837} // Banneker Recreation Center

	start, startDist := g.NearestNode(startPoint)
	end, endDist := g.NearestNode(endPoint)

	fmt.Printf("Start: nearest graph node is %d (%.0fm from target point)\n", start, startDist)
	fmt.Printf("End: nearest graph node is %d (%.0fm from target point)\n\n", end, endDist)

	fmt.Printf("Benchmarking route from node %d to node %d (averaged over 1000 runs)\n\n", start, end)

	const iterations = 1000

	var dijkstraResult, astarResult graph.PathResult

	dijkstraStart := time.Now()
	for i := 0; i < iterations; i++ {
		dijkstraResult = graph.Dijkstra(g, start, end)
	}
	dijkstraAvg := time.Since(dijkstraStart) / iterations

	astarStart := time.Now()
	for i := 0; i < iterations; i++ {
		astarResult = graph.AStar(g, start, end)
	}
	astarAvg := time.Since(astarStart) / iterations

	fmt.Println("=== Dijkstra ===")
	fmt.Printf("Distance: %.1f meters | Nodes explored: %d | Avg time: %v\n",
		dijkstraResult.Distance, dijkstraResult.NodesVisited, dijkstraAvg)

	fmt.Println("\n=== A* ===")
	fmt.Printf("Distance: %.1f meters | Nodes explored: %d | Avg time: %v\n",
		astarResult.Distance, astarResult.NodesVisited, astarAvg)

	nodeReduction := 100 * (1 - float64(astarResult.NodesVisited)/float64(dijkstraResult.NodesVisited))
	timeReduction := 100 * (1 - float64(astarAvg)/float64(dijkstraAvg))

	fmt.Printf("\nA* explored %.1f%% fewer nodes than Dijkstra\n", nodeReduction)
	fmt.Printf("A* ran %.1f%% faster than Dijkstra\n", timeReduction)

	if dijkstraResult.Distance != astarResult.Distance {
		fmt.Println("\nWARNING: distances differ - something is wrong, A* should match Dijkstra's optimal distance")
	} else {
		fmt.Println("\nBoth algorithms agree on the optimal distance, confirming A* with this admissible heuristic finds the same shortest path as Dijkstra, just with fewer nodes explored.")
	}
}