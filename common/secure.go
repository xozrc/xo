package common

import (
	"math/rand"
	"strconv"
)

func RandUserCode() string {
	return strconv.Itoa(rand.Int())
}
