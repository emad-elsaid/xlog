package xlog

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

func CSRFToken() []byte {
	if !csrfExists() {
		writeCSRF()
	}

	return readCSRF()
}

func csrfExists() bool {
	_, err := os.Stat("security-token")
	return err == nil
}

func readCSRF() []byte {
	dat, err := ioutil.ReadFile("security-token")
	if err != nil {
		log.Fatal("Can't read security-token file")
		return []byte{}
	}

	return dat
}

func writeCSRF() {
	err := ioutil.WriteFile("security-token", generateCSRF(), 0644)
	if err != nil {
		log.Fatal("Can't write security-token file")
	}
}

func generateCSRF() []byte {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}
