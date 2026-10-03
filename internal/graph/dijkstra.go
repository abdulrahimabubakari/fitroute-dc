package graph

import "math"

type PathResult struct {
	Path       []int64
	Distance   float64 // meters
	NodesVisited int
	Found      bool
}

// Dijkstra finds the shortest path from start to end using a binary min-heap
// with decrease-key support, giving O((V + E) log V) time complexity.
func Dijkstra(g *Graph, start, end int64) PathResult {
	dist := make(map[int64]float64)
	prev := make(map[int64]int64)
	visited := make(map[int64]bool)

	for id := range g.Nodes {
		dist[id] = math.Inf(1)
	}
	dist[start] = 0

	pq := NewMinHeap()
	pq.Push(start, 0)

	nodesVisited := 0

	for pq.Len() > 0 {
		current, currentDist := pq.Pop()

		if visited[current] {
			continue
		}
		visited[current] = true
		nodesVisited++

		if current == end {
			break
		}

		for _, edge := range g.Adjacency[current] {
			if visited[edge.To] {
				continue
			}
			newDist := currentDist + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = current
				if pq.Contains(edge.To) {
					pq.DecreaseKey(edge.To, newDist)
				} else {
					pq.Push(edge.To, newDist)
				}
			}
		}
	}

	if math.IsInf(dist[end], 1) {
		return PathResult{Found: false, NodesVisited: nodesVisited}
	}

	// Reconstruct path by walking prev[] backwards from end to start.
	path := []int64{}
	for at := end; ; {
		path = append([]int64{at}, path...)
		if at == start {
			break
		}
		at = prev[at]
	}

	return PathResult{
		Path:         path,
		Distance:     dist[end],
		NodesVisited: nodesVisited,
		Found:        true,
	}
}