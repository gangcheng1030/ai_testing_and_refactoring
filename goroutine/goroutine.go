package goroutine

import "time"

var counter int

func InitCounter(initialValue int) {
	counter = initialValue
}

func StartCounter() {
	// 启动一个goroutine，每隔1秒增加counter的值
	go func() {
		for {
			time.Sleep(1 * time.Second)
			counter++
		}
	}()
}
