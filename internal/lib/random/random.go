package random

import (
	"math/rand"
	"time"
)

func GetRandomByLength(aliasLangth int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var a = make([]rune, aliasLangth)

	for i := range a {
		a[i] = letterRunes[r.Intn(len(letterRunes))]
	}

	return string(a)
}
