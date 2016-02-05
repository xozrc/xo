package common

import (
	"math/rand"

	"strconv"
)

func RandUserCode() string {
	return strconv.Itoa(rand.Int())
}

func TokenGenerate(name string, password string, now int64) (token string) {

	return
}
