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
	"strconv"
	"time"
)

func GetFileScanner(fileName string) (*bufio.Scanner, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(file), nil
}


func SaveFile(file multipart.File) (string, error) {
	now := strconv.Itoa(int(time.Now().Unix()))

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	err := os.MkdirAll(now, 0755)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", now, now), buf.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return now, nil
}

func PartitionFile(filename string) error {
	return exec.Command("sh", "partition.sh", filename).Run()
}