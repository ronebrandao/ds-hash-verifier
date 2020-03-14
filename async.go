package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var lines chan string
var found chan bool
var desiredHash string

func FindHashAsync(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())

	desiredHash = strings.TrimPrefix(r.URL.Path, "/async/")

	scanner, err := GetFileScanner("rockyou.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo.")
	}

	lines = make(chan string, 1000)
	found = make(chan bool, 1)

	wg := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go worker(ctx, cancel, &wg)
	}

	go collector(scanner, lines)

	wg.Wait()

	close(found)
	wasFound := <-found

	if wasFound {
		fmt.Fprintf(w, "Encontrei a hash")
	} else {
		fmt.Fprintf(w, "Falhei miseravelmente")
	}
}

func collector(scanner *bufio.Scanner, lines chan string) {
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
}

func worker(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case line, ok := <-lines:
			if !ok {
				return
			}
			hash := fmt.Sprintf("%x", md5.Sum([]byte(line)))
			if hash == desiredHash {
				fmt.Println("Achei")
				found <- true
				cancel()
			}
		case <-ctx.Done():
			fmt.Println("O cara achou, to parando.")
			return
		}
	}
}
