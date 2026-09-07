package network

import (
	"container/heap"
	"fmt"
)

// PathOptimizer implements routing algorithms (e.g., Dijkstra's) for optimal backup data pathing.
type PathOptimizer struct{}

// NewPathOptimizer creates a new PathOptimizer.
func NewPathOptimizer() *PathOptimizer {
	return &PathOptimizer{}
}

// Node represents a node in the network graph.
type Node struct {
	ID      string
	Edges   []*Edge
	Distance float64 // Used during Dijkstra's algorithm
	Visited bool    // Used during Dijkstra's algorithm
	Previous *Node  // Used during Dijkstra's algorithm
}

// Edge represents a connection between two nodes in the network graph.
type Edge struct {
	Target  *Node
	Weight  float64 // e.g., latency, cost, or distance
}

// Path represents a route through the network.
type Path struct {
	Nodes   []string
	TotalWeight float64
}

// FindShortestPath uses Dijkstra's algorithm to find the shortest path between two nodes.
func (po *PathOptimizer) FindShortestPath(nodes map[string]*Node, startID, endID string) (*Path, error) {
	start := nodes[startID]
	end := nodes[endID]

	if start == nil || end == nil {
		return nil, ErrNodeNotFound
	}

	// Initialize nodes for Dijkstra's algorithm
	for _, node := range nodes {
		node.Distance = float64(^uint(0) >> 1) // Max int value
		node.Visited = false
		node.Previous = nil
	}

	start.Distance = 0

	// Priority queue for Dijkstra's algorithm
	pq := make(NodePriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &NodeItem{node: start, priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*NodeItem)
		current := item.node

		if current.Visited {
			continue
		}

		current.Visited = true

		if current == end {
			// Found the shortest path, reconstruct it
			path := []string{}
			totalWeight := 0.0
			for at := end; at != nil; at = at.Previous {
				path = append([]string{at.ID}, path...)
				if at.Previous != nil {
					// Find the edge weight to add to total
					for _, edge := range at.Previous.Edges {
						if edge.Target == at {
							totalWeight += edge.Weight
							break
						}
					}
				}
			}
			return &Path{
				Nodes:   path,
				TotalWeight: totalWeight,
			}, nil
		}

		// Update distances to neighbors
		for _, edge := range current.Edges {
			neighbor := edge.Target
			if neighbor.Visited {
				continue
			}

			newDistance := current.Distance + edge.Weight
			if newDistance < neighbor.Distance {
				neighbor.Distance = newDistance
				neighbor.Previous = current
				heap.Push(&pq, &NodeItem{node: neighbor, priority: newDistance})
			}
		}
	}

	return nil, ErrNoPathFound
}

// ErrNodeNotFound is returned when start or end node is not found.
var ErrNodeNotFound = fmt.Errorf("start or end node not found")

// ErrNoPathFound is returned when no path exists between nodes.
var ErrNoPathFound = fmt.Errorf("no path found between nodes")

// NodeItem is an item in the priority queue.
type NodeItem struct {
	node     *Node
	priority float64
	index    int // The index in the heap
}

// NodePriorityQueue implements heap.Interface for Nodes.
type NodePriorityQueue []*NodeItem

func (pq NodePriorityQueue) Len() int { return len(pq) }

func (pq NodePriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq NodePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *NodePriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*NodeItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *NodePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.index = -1 // For safety
	*pq = old[0 : n-1]
	return item
}