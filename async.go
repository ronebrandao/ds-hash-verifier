package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

//var lines chan string
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

	filePath, err := SaveFile(file)

	if err != nil {
		fmt.Fprint(w, "Ocorreu um erro ao salvar o arquivo em disco.", err)
		return
	}

	err = PartitionFile(filePath)

	if err != nil {
		fmt.Fprint(w, "Erro ao dividir arquivo.", err)
		return
	}

	found = make(chan bool, 1)

	collectorWg := sync.WaitGroup{}
	collectorWg.Add(4)

	go collector(ctx, cancel, filePath+"/xaa", &collectorWg)
	go collector(ctx, cancel, filePath+"/xab", &collectorWg)
	go collector(ctx, cancel, filePath+"/xac", &collectorWg)
	go collector(ctx, cancel, filePath+"/xad", &collectorWg)

	collectorWg.Wait()

	close(found)
	wasFound := <-found

	if wasFound {
		fmt.Fprintf(w, "Encontrei a hash")
	} else {
		fmt.Fprintf(w, "Falhei miseravelmente")
	}
}

func collector(ctx context.Context, cancel context.CancelFunc, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	lines := make(chan string, 1000)

	file, _ := os.Open(fileName)

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(lines)

		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	wg2 := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg2.Add(1)

		go worker(ctx, cancel, lines, &wg2)
	}

	wg2.Wait()
}

func worker(ctx context.Context, cancel context.CancelFunc, lines chan string, wg *sync.WaitGroup) {
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
