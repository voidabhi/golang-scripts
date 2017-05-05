
package main
import "fmt"

func ingest(out <-chan []string) {
	out <- []string{"aaaa", "bbb"}
	out <- []string{"cccccc", "dddddd"}
	out <- []string{"e", "fffff", "g"}
	close(out)
}

func process(in <-chan []string, out <-chan int) {
	for data := range in {
		for _, word := range data {
			out <- len(word)
		}
	}
}

func store(in <-chan int) {
	for data := range in {
		fmt.Println(data)
	}
}

func main() {
	concurrency := 4
	
	// stage 1 ingest data from source
	in := make(chan []string)
	go ingest(in)

	// stage 2 - process data
	reduced := make(chan int)
	
	var wg sync.WaitGroup
	wg.Add(concurrency)
	
	for i := 0; i < concurrency; i++ {
		go func() {
			process(in, reduced)
			wg.Done()
		}()
	}
	
	go func() {
		wg.Wait()
		close(reduced)
	}

	// stage 3 - store
	store(reduced)
}
