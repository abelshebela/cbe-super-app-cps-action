package queue

import "container/heap"

// RetryHeap implements heap.Interface for RetryJob items,
// ordered by executeAt (min-heap). The earliest job to execute is at index 0.
type RetryHeap []*RetryJob

func (h RetryHeap) Len() int           { return len(h) }
func (h RetryHeap) Less(i, j int) bool { return h[i].executeAt.Before(h[j].executeAt) }
func (h RetryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *RetryHeap) Push(x interface{}) {
	item := x.(*RetryJob)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *RetryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*h = old[0 : n-1]
	return item
}

func (h *RetryHeap) Peek() *RetryJob {
	if h.Len() == 0 {
		return nil
	}
	return (*h)[0]
}

func InitHeap(h *RetryHeap) {
	heap.Init(h)
}
