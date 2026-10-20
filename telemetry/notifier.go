package main

import (
	"sync"
)

// PodNotifier Go channels are like pipes in C/UNIX
type PodNotifier struct {
	mu          sync.Mutex
	subscribers map[string][]chan string
}

func NewPodNotifier() *PodNotifier {
	return &PodNotifier{
		subscribers: make(map[string][]chan string),
	}
}

func (pn *PodNotifier) Subscribe(podName string) chan string {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	ch := make(chan string, 1)
	pn.subscribers[podName] = append(pn.subscribers[podName], ch)
	return ch
}

func (pn *PodNotifier) Unsubscribe(podName string, ch chan string) {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	subs := pn.subscribers[podName]
	for i, sub := range subs {
		if sub == ch {
			pn.subscribers[podName] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (pn *PodNotifier) Publish(podName string, ip string) {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	if subs, exists := pn.subscribers[podName]; exists {
		for _, ch := range subs {
			select {
			case ch <- ip:
			default:
			}
		}
		delete(pn.subscribers, podName)
	}
}
