package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Metrics struct {
	processed int64
	failed    int64
}

func main() {
	dir := flag.String("dir", ".", "Directory to scan")
	workers := flag.Int("workers", 4, "Initial number of worker goroutines")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal...")
		cancel()
	}()

	jobs := make(chan string, 100)
	var wg sync.WaitGroup
	var metrics Metrics

	var errMu sync.Mutex
	var errors []error

	// Start metrics reporter
	go metricsReporter(ctx, jobs, &metrics)

	// Start initial worker pool
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg, &metrics, &errMu, &errors)
	}

	// Start worker autoscaler
	go workerAutoscaler(ctx, jobs, &wg, &metrics, &errMu, &errors, *workers)

	// Walk directory
	go func() {
		defer close(jobs)
		err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}

			select {
			case jobs <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})

		if err != nil && err != context.Canceled {
			fmt.Println("Walk error:", err)
		}
	}()

	wg.Wait()

	fmt.Println("\nProcessing complete")
	fmt.Println("Files processed:", atomic.LoadInt64(&metrics.processed))
	fmt.Println("Files failed:", atomic.LoadInt64(&metrics.failed))

	if len(errors) > 0 {
		fmt.Println("Some errors occurred:")
		for _, err := range errors {
			fmt.Println("-", err)
		}
	}
}

func worker(
	ctx context.Context,
	id int,
	jobs <-chan string,
	wg *sync.WaitGroup,
	metrics *Metrics,
	errMu *sync.Mutex,
	errors *[]error,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down...\n", id)
			return
		case path, ok := <-jobs:
			if !ok {
				return
			}

			err := processFile(path)
			if err != nil {
				atomic.AddInt64(&metrics.failed, 1)

				errMu.Lock()
				*errors = append(*errors, err)
				errMu.Unlock()

				continue
			}

			atomic.AddInt64(&metrics.processed, 1)
		}
	}
}

func processFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("hash %s: %w", path, err)
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	time.Sleep(50 * time.Millisecond)

	fmt.Printf("Processed: %s | SHA256: %s\n", path, hash)
	return nil
}

// Live metrics reporter
func metricsReporter(ctx context.Context, jobs chan string, metrics *Metrics) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Metrics reporter shutting down...")
			return
		case <-ticker.C:
			processed := atomic.LoadInt64(&metrics.processed)
			failed := atomic.LoadInt64(&metrics.failed)
			queueLength := len(jobs)
			goroutines := runtime.NumGoroutine()

			fmt.Printf("\n[METRICS] Processed: %d | Failed: %d | Queue: %d | Goroutines: %d\n",
				processed, failed, queueLength, goroutines)
		}
	}
}

// Worker Autoscaler
func workerAutoscaler(ctx context.Context, jobs chan string, wg *sync.WaitGroup, metrics *Metrics, errMu *sync.Mutex, errors *[]error, initialWorkers int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	workerID := initialWorkers
	maxWorkers := 20
	minWorkers := 2
	activeWorkers := initialWorkers

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queueLength := len(jobs)

			// Scale up
			if queueLength > 50 && activeWorkers < maxWorkers {
				add := 2
				for i := 0; i < add && activeWorkers < maxWorkers; i++ {
					wg.Add(1)
					workerID++
					go worker(ctx, workerID, jobs, wg, metrics, errMu, errors)
					activeWorkers++
					fmt.Printf("Autoscaler: Spawned extra worker %d (total workers: %d)\n", workerID, activeWorkers)
				}
			}

			// Scale down (conceptual, we can't forcibly stop workers without context)
			if queueLength < 10 && activeWorkers > minWorkers {
				activeWorkers-- // track logical reduction; idle workers will naturally exit when queue is empty
				fmt.Printf("Autoscaler: Reducing worker count (logical total: %d)\n", activeWorkers)
			}
		}
	}
}
