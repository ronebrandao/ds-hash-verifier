package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"time"
)

func GetFileScanner(fileName string) (*bufio.Scanner, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(file), nil
}


func SaveFile(file multipart.File) (*string, error) {
	now := time.Now().Format(time.RFC3339Nano)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	err := os.MkdirAll("server/files/", 0755)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fmt.Sprintf("server/files/%s.txt", now), buf.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	return &now, nil
}

func PartitionFile(filename string) error {
	return exec.Command("sh", "server/partition.sh", filename).Run()
}