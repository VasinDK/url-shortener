package random_test

import (
	"mod_shortener/internal/lib/random"
	"testing"
	"time"
)

var cases = []struct {
	name  string
	lenth int
}{
	{"1", 1},
	{"2", 2},
	{"3", 3},
	{"8", 8},
	{"8", 8},
	{"8", 8},
	{"8", 8},
}

func TestGetRandomByLength(t *testing.T) {
	var arRes = make(map[string]struct{})
	for _, pair := range cases {
		time.Sleep(time.Microsecond)
		var res = random.GetRandomByLength(pair.lenth)

		if len(res) != pair.lenth {
			t.Error("expected:", pair.name)
		}

		if _, ok := arRes[res]; ok {
			t.Error("double", arRes)
		}

		arRes[res] = struct{}{}
	}
}
