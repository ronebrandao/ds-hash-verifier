package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
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

	err := r.ParseMultipartForm(200 << 20)

	if err != nil {
		fmt.Fprintf(w, "O arquivo deve ter no máximo 200MB")
	}

	file, _, _ := r.FormFile("wordlist")

	if file == nil {
		fmt.Fprint(w, "Erro, nenhuma wordlist foi encontrada!")
		return
	}

	var buf bytes.Buffer

	io.Copy(&buf, file)

	lines = make(chan string, 1000)
	found = make(chan bool, 1)

	wg := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go worker(ctx, cancel, &wg)
	}

	go collector(buf, lines)

	wg.Wait()

	close(found)
	wasFound := <-found

	if wasFound {
		fmt.Fprintf(w, "Encontrei a hash")
	} else {
		fmt.Fprintf(w, "Falhei miseravelmente")
	}
}

func collector(buf bytes.Buffer, lines chan string) {
	for {
		b, err := buf.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		lines <- strings.TrimSuffix(string(b), "\n")
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
