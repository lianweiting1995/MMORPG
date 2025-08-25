package helper_test

import (
	"MMORPG/helper"
	"fmt"
	"testing"
	"time"
)

func TestSupervisor(t *testing.T) {
	helper.Push("demo", func() {
		i := 1
		for {
			fmt.Println("demo")
			time.Sleep(1 * time.Second)
			i += 1
			if i == 4 {
				panic("demo panic")
			}
		}
	})

	time.Sleep(time.Hour)
}
