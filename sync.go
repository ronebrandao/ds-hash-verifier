package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)

func FindHash(w http.ResponseWriter, r *http.Request) () {
	desired := strings.TrimPrefix(r.URL.Path, "/sync/")

	scanner, err := GetFileScanner("rockyou.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo.")
	}

	for scanner.Scan() {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(scanner.Text())))
		if hash == desired {
			fmt.Fprintf(w, "Encontrei a hash")
			return
		}
	}

	fmt.Fprintf(w, "Falhei miseravelmente")
}
