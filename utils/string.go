package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

//RandomString get random string with length
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//WriteContentToFile write to file
func WriteContentToFile(content []string, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return fmt.Errorf("Failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range content {
		_, err = datawriter.WriteString(data + "\n")
		if err != nil {
			return fmt.Errorf("Failed to write to file %s", err.Error())
		}
	}

	err = datawriter.Flush()
	if err != nil {
		return fmt.Errorf("Failed to flush file %s", err.Error())
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("Failed to close file %s", err.Error())
	}
	return nil
}
