package main

import (
	"fmt"
	"time"
)

func main() {
	one()
}

var ch chan int = make(chan int)

// 从小到大找出 17 和 38 的 3 个公倍数
func generate() {
	for i := 17; i < 5000; i += 17 {
		ch <- i
		time.Sleep(1 * time.Millisecond)
	}
	close(ch)
}

func one() {
	timeout := time.After(800 * time.Millisecond)
	go generate()
	found := 0
	for {
		select {
		case i, ok := <-ch:
			if ok {
				if i%38 == 0 {
					fmt.Println(i, "is a multiple of 17 and 38")
					found++
					if found == 3 {
						break // 不能跳出外层循环
					}
				}
			} else {
				break // 不能跳出外层循环
			}
		case <-timeout:
			fmt.Println("timed out")
			break
		}
	}
	fmt.Println("The end")
}

func oneUpdate() {
	timeout := time.After(800 * time.Millisecond)
	go generate()
	found := 0
MainLoop:
	for {
		select {
		case i, ok := <-ch:
			if ok {
				if i%38 == 0 {
					fmt.Println(i, "is a multiple of 17 and 38")
					found++
					if found == 3 {
						break MainLoop
					}
				}
			} else {
				break MainLoop
			}
		case <-timeout:
			fmt.Println("timed out")
			break
		}
	}
	fmt.Println("The end")
}
