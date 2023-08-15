package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Producer struct {
	lock *sync.Mutex
	wg   *sync.WaitGroup
}

type Consumer struct {
	num  int
	lock *sync.Mutex
	wg   *sync.WaitGroup
}

func (p *Producer) Produce(ctx context.Context, buffer chan int, offset time.Duration) {
	defer p.wg.Done()
	for {
		select {
		case <-time.After(offset):
			p.lock.Lock()
			buffer <- rand.Intn(100)
			p.lock.Unlock()
		case <-ctx.Done():
			return
		}
	}

}

func (c *Consumer) Consume(ctx context.Context, buffer chan int) {
	defer c.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case target := <-buffer:
			c.lock.Lock()
			fmt.Printf("consumer %d: consumed %d\n", c.num, target)
			c.lock.Unlock()
		}
	}
}

func timeBomb(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-time.After(time.Second * 30):
		cancel()
	case <-ctx.Done():
		return
	}
}

func main() {
	consumers := 100
	producers := 50
	lock := sync.Mutex{}
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	buffer := make(chan int, 1024)
	wg := sync.WaitGroup{}
	for i := 1; i <= producers; i++ {
		p := &Producer{
			lock: &lock,
			wg:   &wg,
		}
		wg.Add(1)
		go p.Produce(ctx, buffer, time.Second*time.Duration(i))
	}
	for i := 1; i <= consumers; i++ {
		c := &Consumer{
			num:  i,
			lock: &lock,
			wg:   &wg,
		}
		wg.Add(1)
		go c.Consume(ctx, buffer)
	}
	wg.Add(1)
	go timeBomb(ctx, cancel, &wg)
	wg.Wait()
}
