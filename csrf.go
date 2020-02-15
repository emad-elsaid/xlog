package xlog

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

const (
	CSRFFile    = "security-token"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CSRFToken() []byte {
	if !csrfExists() {
		writeCSRF()
	}

	return readCSRF()
}

func csrfExists() bool {
	_, err := os.Stat(CSRFFile)
	return err == nil
}

func readCSRF() []byte {
	dat, err := ioutil.ReadFile(CSRFFile)
	if err != nil {
		log.Fatal("Can't read security-token file")
		return []byte{}
	}

	return dat
}

func writeCSRF() {
	err := ioutil.WriteFile(CSRFFile, generateCSRF(), 0644)
	if err != nil {
		log.Fatal("Can't write " + CSRFFile + " file")
	}
}

func generateCSRF() []byte {
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}
