package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FindHash(w http.ResponseWriter, r *http.Request) () {
	desired := strings.TrimPrefix(r.URL.Path, "/sync/")
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

	for {
		b, err := buf.ReadBytes('\n')
		b = []byte(strings.TrimSuffix(string(b), "\n"))

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprint(w, "Erro ao ler arquivo.")
			break
		}

		hash := fmt.Sprintf("%x", md5.Sum(b))
		if hash == desired {
			fmt.Fprintf(w, "Hash encontrada!")
			return
		}

	}

	fmt.Fprintf(w, "Falhei miseravelmente")
}
