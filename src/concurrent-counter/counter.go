package counter

import (
	"fmt"
	"sync"
)

func inc(i *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	*i++
	mu.Unlock()
}

func Count() {
	var val int = 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go inc(&val, &mu, &wg)
	}

	wg.Wait()

	fmt.Println("result: ", val)
}
