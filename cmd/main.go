package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/fadygamilm/go-leak-detector/internal/parser"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		go func() {
			time.Sleep(200 * time.Millisecond) // <-- keep nested goroutine alive
			wg.Done()
		}()
		time.Sleep(200 * time.Millisecond) // <-- keep first goroutine alive
		wg.Done()
	}()
	time.Sleep(50 * time.Millisecond) // let goroutines start

	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	wg.Wait()
	parser := parser.New(context.Background())
	result := parser.Parse(buf, n)
	log.Println(result[0])
	log.Println("=================")
	log.Println(result[1])
	log.Println("=================")
	log.Println(result[2])
}
