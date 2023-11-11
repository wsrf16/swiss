package timercache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	i := 1
	build := Build[int](3*time.Second, func() int {
		i++
		return i
	})
	build.Start()
	for {
		fmt.Println(build.GetData())
		// fmt.Println(*build.TryGetData(100 * time.Millisecond))
		time.Sleep(1 * time.Second)
	}
	select {}
}
