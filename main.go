package main

import (
	"math/rand"
)

func main() {

	m := wordGen(100, 10)
	println(m[0], m[1])

}

func wordGen(nDistinct, wordLen int) []string {
	vocab := make([]string, nDistinct)
	for i := range nDistinct {
		word := randomString(wordLen)
		vocab[i] = word
	}
	return vocab
}

func randomString(n int) string {
	// 脑子进煎鱼了
	const letters = "eddycjyabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, n)
	r := rand.New(rand.NewSource(rand.Int63()))
	for i := 0; i < n; {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			panic(err)
		}
		ret[i] = letters[int(b[0])%len(letters)]
		i++
	}
	return string(ret)
}
