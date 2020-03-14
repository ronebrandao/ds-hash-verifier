package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"sync"
	"time"
)

func runSync() {
	start := time.Now()
	file, err := os.Open("rockyou.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo")
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// time.Sleep(time.Microsecond * 300)
		hash := fmt.Sprintf("%x", md5.Sum([]byte(scanner.Text())))

		if hash == "fd9cabd4def5137a73d682f4dd963e57" {
			fmt.Printf("\nEncontrar Sincrono: %2fs\n", time.Since(start).Seconds())
		}
	}

	fmt.Printf("\nTerminei Sincrono: %2fs\n", time.Since(start).Seconds())
}

func runAsync() {
	start := time.Now()
	file, err := os.Open("rockyou.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo")
	}

	scanner := bufio.NewScanner(file)

	lines := make(chan string, 1000)
	// found := make(chan bool)

	wg := sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, i int) {
			fmt.Println("Comecando a rotina ", i)
			defer wg.Done()

			for line := range lines {
				// time.Sleep(time.Microsecond * 300)
				hash := fmt.Sprintf("%x", md5.Sum([]byte(line)))
				if hash == "fd9cabd4def5137a73d682f4dd963e57" {
					fmt.Printf("\nEncontrar Assincrono: %2fs\n", time.Since(start).Seconds())
					break
					// found <- true
				}
			}

			fmt.Println("Acabando a rotina", i)
		}(&wg, i)

	}

	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	wg.Wait()

	// close(found)
	// wasFound := <-found

	// if !wasFound {
	// 	fmt.Println("Falhei miseravelmente")
	// }

	fmt.Printf("\nTerminei Assincrono: %2fs\n", time.Since(start).Seconds())
}

func main() {
	runSync()
	runAsync()
}
