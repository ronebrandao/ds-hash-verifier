package main

import (
	"bufio"
	"os"
)

func GetFileScanner(fileName string) (*bufio.Scanner, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(file), nil
}
