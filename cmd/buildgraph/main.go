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

	totalEdges := 0
	for _, edges := range g.Adjacency {
		totalEdges += len(edges)
	}

	fmt.Printf("Graph built: %d nodes, %d directed edge entries (%d undirected connections)\n",
		len(g.Nodes), totalEdges, totalEdges/2)
}