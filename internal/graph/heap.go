package graph

// pqItem is one entry in the priority queue: a node and its current best-known distance.
type pqItem struct {
	nodeID int64
	dist   float64
	index  int // position in the heap slice, maintained for decrease-key support
}

// MinHeap is a binary min-heap ordered by dist, with O(log n) push/pop
// and O(log n) decrease-key via index tracking (so Dijkstra doesn't need
// to push duplicate entries and filter them out later).
type MinHeap struct {
	items   []*pqItem
	indexOf map[int64]*pqItem // nodeID -> its item, for fast decrease-key lookup
}

func NewMinHeap() *MinHeap {
	return &MinHeap{
		items:   make([]*pqItem, 0),
		indexOf: make(map[int64]*pqItem),
	}
}

func (h *MinHeap) Len() int { return len(h.items) }

func (h *MinHeap) Push(nodeID int64, dist float64) {
	item := &pqItem{nodeID: nodeID, dist: dist, index: len(h.items)}
	h.items = append(h.items, item)
	h.indexOf[nodeID] = item
	h.siftUp(item.index)
}

// Pop removes and returns the node with the smallest dist.
func (h *MinHeap) Pop() (int64, float64) {
	top := h.items[0]
	last := len(h.items) - 1
	h.swap(0, last)
	h.items = h.items[:last]
	delete(h.indexOf, top.nodeID)
	if len(h.items) > 0 {
		h.siftDown(0)
	}
	return top.nodeID, top.dist
}

// DecreaseKey lowers a node's distance if newDist is better than what's recorded,
// and re-heapifies. This is the key efficiency trick over a naive "push every update".
func (h *MinHeap) DecreaseKey(nodeID int64, newDist float64) {
	item, exists := h.indexOf[nodeID]
	if !exists || newDist >= item.dist {
		return
	}
	item.dist = newDist
	h.siftUp(item.index)
}

func (h *MinHeap) Contains(nodeID int64) bool {
	_, exists := h.indexOf[nodeID]
	return exists
}

func (h *MinHeap) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].index = i
	h.items[j].index = j
}

func (h *MinHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.items[i].dist >= h.items[parent].dist {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

func (h *MinHeap) siftDown(i int) {
	n := len(h.items)
	for {
		left, right := 2*i+1, 2*i+2
		smallest := i
		if left < n && h.items[left].dist < h.items[smallest].dist {
			smallest = left
		}
		if right < n && h.items[right].dist < h.items[smallest].dist {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.swap(i, smallest)
		i = smallest
	}
}