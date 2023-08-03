package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func helloFromGoroutine(ctx context.Context, num int, execute chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	prefix := fmt.Sprintf("goroutine %d:", num)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s context cancelled, cleaning up...\n", prefix)
			return
		case <-execute:
			fmt.Printf("%s hello!\n", prefix)
		}
	}
}

func callExec(ctx context.Context, duration time.Duration, execChan chan struct{}, wg *sync.WaitGroup, threads int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("context canceled, callExec shutting down!\n")
			return
		case <-time.After(duration):
			for i := 0; i < threads; i++ {
				execChan <- struct{}{}
			}
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("not enough arguments, please specify the number of goroutines\n")
	}
	threads, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("using, %d goroutines\n", threads)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	exec := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go callExec(ctx, time.Second*1, exec, &wg, threads)
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go helloFromGoroutine(ctx, i, exec, &wg)
	}
	wg.Wait()
}
