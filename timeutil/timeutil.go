package timeutil

import (
	"time"
)

//millisecond base
func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
