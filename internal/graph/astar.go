package graph

import "math"

// AStar finds the shortest path using Dijkstra's algorithm augmented with a
// haversine-distance heuristic, biasing exploration toward the goal rather
// than expanding uniformly in all directions like plain Dijkstra does.
func AStar(g *Graph, start, end int64) PathResult {
	gScore := make(map[int64]float64) // known cost from start to this node
	prev := make(map[int64]int64)
	visited := make(map[int64]bool)

	for id := range g.Nodes {
		gScore[id] = math.Inf(1)
	}
	gScore[start] = 0

	endPoint := g.Nodes[end]
	heuristic := func(nodeID int64) float64 {
		return Haversine(g.Nodes[nodeID], endPoint)
	}

	// The heap is ordered by fScore (gScore + heuristic), not raw distance -
	// this is the entire mechanism that lets A* "look toward" the goal.
	pq := NewMinHeap()
	pq.Push(start, heuristic(start))

	nodesVisited := 0

	for pq.Len() > 0 {
		current, _ := pq.Pop()

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
			tentativeG := gScore[current] + edge.Weight
			if tentativeG < gScore[edge.To] {
				gScore[edge.To] = tentativeG
				prev[edge.To] = current
				fScore := tentativeG + heuristic(edge.To)
				if pq.Contains(edge.To) {
					pq.DecreaseKey(edge.To, fScore)
				} else {
					pq.Push(edge.To, fScore)
				}
			}
		}
	}

	if math.IsInf(gScore[end], 1) {
		return PathResult{Found: false, NodesVisited: nodesVisited}
	}

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
		Distance:     gScore[end],
		NodesVisited: nodesVisited,
		Found:        true,
	}
}