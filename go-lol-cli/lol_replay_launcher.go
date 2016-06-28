package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charlywazzap/go-lol"
)

type LolReplayLauncher interface {
	Launch(address string, region *lol.Region, id lol.GameID, encryptionKey string) error
}

type version struct {
	numbers []int64
}

var versionRx = regexp.MustCompile(`\A[0-9]+(\.[0-9]+)*\z`)

func newVersion(name string) (version, error) {
	res := version{}
	if versionRx.MatchString(name) == false {
		return res, fmt.Errorf("Invalid version name %s", name)
	}
	nberStr := strings.Split(name, ".")
	res.numbers = make([]int64, 0, len(nberStr))

	for _, nStr := range nberStr {
		v, err := strconv.ParseInt(nStr, 10, 64)
		if err != nil {
			return version{}, fmt.Errorf("Invalid version %s: %s", name, err)
		}
		res.numbers = append(res.numbers, v)
	}
	return res, nil
}

func (v version) String() string {
	asStr := make([]string, 0, len(v.numbers))
	for _, n := range v.numbers {
		asStr = append(asStr, fmt.Sprintf("%d", n))
	}
	return strings.Join(asStr, ".")
}

type versionList []version

func (l versionList) Len() int {
	return len(l)
}

func (v version) less(o version) bool {
	size := len(v.numbers)
	if len(o.numbers) < size {
		size = len(o.numbers)
	}
	for i := 0; i < size; i++ {
		if v.numbers[i] < o.numbers[i] {
			return true
		}
		if v.numbers[i] > o.numbers[i] {
			return false
		}
	}

	// equal up to that point
	if len(v.numbers) < len(o.numbers) {
		return true
	}
	return false
}

func (l versionList) Less(i, j int) bool {
	return l[i].less(l[j])
}

func (l versionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

const (
	MaestroParam1 = "8394"
	MaestroParam2 = "LoLLauncher"
)
