package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func generateKeyPair() (solana.PublicKey, solana.PrivateKey) {
	privKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		os.Exit(1)
	}
	return privKey.PublicKey(), privKey
}

func isMatch(address, start, end string, caseSensitive bool) bool {
	if !caseSensitive {
		address = strings.ToLower(address)
		start = strings.ToLower(start)
		end = strings.ToLower(end)
	}
	return strings.HasPrefix(address, start) && strings.HasSuffix(address, end)
}

func worker(start, end string, caseSensitive bool, results chan<- [2]string, wg *sync.WaitGroup, attempts *uint64, done <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			pubKey, privKey := generateKeyPair()
			address := pubKey.String()
			atomic.AddUint64(attempts, 1)
			if isMatch(address, start, end, caseSensitive) {
				select {
				case results <- [2]string{address, privKey.String()}:
					return
				case <-done:
					return
				}
			}
		}
	}
}

func main() {
	start := flag.String("s", "", "The starting string of the address")
	end := flag.String("e", "", "The ending string of the address")
	threads := flag.Int("t", int(float64(runtime.NumCPU())*float64(0.5)), "The number of threads to use")
	caseSensitive := flag.Bool("c", false, "Case sensitive match")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	if len(os.Args) > 1 && os.Args[1] == "-h" {
		flag.Usage()
		os.Exit(0)
	}

	wsClient, err := ws.Connect(context.Background(), rpc.MainNetBeta_WS)
	if err != nil {
		fmt.Printf("Error connecting to websocket: %v\n", err)
		os.Exit(1)
	}
	defer wsClient.Close()

	// Determina il numero di worker
	numWorkers := 1
	if *start != "" || *end != "" || *caseSensitive {
		numWorkers = *threads
	}

	results := make(chan [2]string)
	done := make(chan struct{})
	var wg sync.WaitGroup
	var attempts uint64
	var once sync.Once

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Number of attempts: %d\n", atomic.LoadUint64(&attempts))
			case <-done:
				return
			}
		}
	}()

	fmt.Printf("Starting with %d threads\n", numWorkers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(*start, *end, *caseSensitive, results, &wg, &attempts, done)
	}

	foundResult := <-results
	once.Do(func() { close(done) })
	fmt.Printf("Found matching address: %s\n", foundResult[0])
	fmt.Printf("Private key: %s\n", foundResult[1])

	close(results)
	wg.Wait()
}
