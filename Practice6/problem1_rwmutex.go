package main

import (
	"fmt"
	"sync"
)

func main() {
	safeMap := make(map[string]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(i)

		go func(key int) {
			mu.Lock()
			safeMap["key"] = key
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.Unlock()
	value := safeMap["key"]
	mu.RUnlock()
	fmt.Printf("Value: %d\n", value)
}
