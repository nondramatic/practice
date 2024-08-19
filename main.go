package main

import (
	"math/rand"
	"strings"
)

func main() {

	wordGen(100, 10)

}

func wordGen(nDistinct, wordLen int) func() string {
	vocab := make([]string, nDistinct)
	for i := range nDistinct {
		word := randomString(wordLen)
		vocab[i] = word
	}
	return func() string {
		word := vocab[rand.Intn(nDistinct)]
		return strings.Clone(word)
	}
}

func randomString(n int) string {
	// 脑子进煎鱼了
	const letters = "eddycjyabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		ret[i] = letters[int(b[0])%len(letters)]
		i++
	}
	return string(ret)
}
