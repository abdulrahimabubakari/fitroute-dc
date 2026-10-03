package main

import (
	"fmt"
	"fitroute-dc/internal/graph"
)

func main() {
	g, err := graph.BuildFromOSMFile("data_streets.json")
	if err != nil {
		fmt.Println("Error building graph:", err)
		return
	}

	// Pick two real node IDs from the graph to route between.
	// grab arbitrary-but-real ones: the first node as start,
	// and a node reasonably far into the map as the destination.
	var start, end int64
	count := 0
	for id := range g.Nodes {
		if count == 0 {
			start = id
		}
		if count == len(g.Nodes)/2 {
			end = id
		}
		count++
	}

	fmt.Printf("Routing from node %d to node %d...\n", start, end)
	result := graph.Dijkstra(g, start, end)

	if !result.Found {
		fmt.Println("No path found.")
		return
	}

	fmt.Printf("Path found: %d nodes, %.1f meters, explored %d nodes\n",
		len(result.Path), result.Distance, result.NodesVisited)
}