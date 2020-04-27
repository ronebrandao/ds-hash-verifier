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

var lines chan string
var found chan bool
var desiredHash string
var root = "server/files/"

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

	fileName, err := SaveFile(file)

	if err != nil {
		fmt.Fprint(w, "Ocorreu um erro ao salvar o arquivo em disco.", err)
		return
	}

	err = PartitionFile("files/" + *fileName)

	if err != nil {
		fmt.Fprint(w, "Erro ao dividir arquivo.", err)
		return
	}

	lines = make(chan string, 1000)
	found = make(chan bool, 1)

	wg := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go worker(ctx, cancel, &wg)
	}

	wg.Add(4)

	go collector(root+"/xaa.txt", lines, &wg)
	go collector(root+"/xab.txt", lines, &wg)
	go collector(root+"/xac.txt", lines, &wg)
	go collector(root+"/xad.txt", lines, &wg)

	wg.Wait()

	close(found)
	wasFound := <-found

	if wasFound {
		fmt.Fprintf(w, "Encontrei a hash")
	} else {
		fmt.Fprintf(w, "Falhei miseravelmente")
	}
}

func collector(fileName string, lines chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, _ := os.Open(fileName)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines <- scanner.Text()
	}
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
