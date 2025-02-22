package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

func main() {
	Ten3()
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

func two() {
	// 声明无缓冲通道
	// 无缓冲需要一直有接收者接收数据，写操作才会继续，不然会一直阻塞
	ch := make(chan interface{})
	// 声明缓冲为 1 的通道
	// 缓冲为 1 的通道，即使没有接收者也不会阻塞，因为缓冲大小是1.
	// 只有当放第二个值的时候，第一个还没有被取走的时候才会阻塞
	ch1 := make(chan interface{}, 1)
	fmt.Println(ch)
	fmt.Println(ch1)
}

var mu sync.Mutex

var chain string

// fatal error
// 使用 mu 加锁后，不能继续对其加锁，否则会导致死锁
func two1() {
	chain = "main"
	A()
	fmt.Println(chain)
}
func A() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + "  -->A"
	B()
}

func B() {
	chain = chain + " -->B"
	C()
}

func C() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + "  -->C"
}

func three() {
	fmt.Println(doubleScoure(0))
	fmt.Println(doubleScoure(20))
	fmt.Println(doubleScoure(50))
}

// return value 不是原子操作，它在编译器中分为两部分，返回赋值和 return。defer 在 return 返回赋值后执行
func doubleScoure(source float32) (score float32) {
	defer func() {
		if score < 1 || score >= 100 {
			score = source
		}
	}()
	return source * 2
}

var mu1 sync.RWMutex
var count int

// fatal error
// 当写锁阻塞时，新的读锁是无法申请的，从而导致死锁
func three1() {
	go A1()
	time.Sleep(2 * time.Second)
	mu1.Lock()
	fmt.Println("写锁")
	defer mu1.Unlock()
	count++
	fmt.Println(count)
}
func A1() {
	mu1.RLock()
	defer mu1.RUnlock()
	B1()
}

func B1() {
	time.Sleep(5 * time.Second)
	C1()
}

func C1() {
	mu1.RLock()
	defer mu1.RUnlock()
}

// panic
// waitgroup 在调用 Wait() 之后不能再调用 Add() 方法
func four() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
		wg.Add(1)
	}()
	wg.Wait()
}

var c = make(chan int)
var a int

func f() {
	a = 1
	c <- 0
}

func five() {
	go f()
	<-c
	print(a)
}

type MyMutex struct {
	count int
	sync.Mutex
}

// fatal error
// 加锁后复制变量，锁的状态也复制，mu1 处于加锁状态，再加锁会死锁
func five1() {
	var mu MyMutex
	mu.Lock()
	var mu1 = mu
	mu.count++
	mu.Unlock()
	mu1.Lock()
	mu1.Lock()
	mu1.count++
	mu1.Unlock()
	fmt.Println(mu.count, mu1.count)
}

// ch 未初始化，引发 fatal error
func six() {
	var ch chan int
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch)
	}()
	<-ch
	fmt.Println(count)
}

// 程序执行到第二个 groutine 时，ch 还未初始化，导致第二个 goroutine 阻塞。第一个 goroutine 不会阻塞
// 对nil通道的读写操作都会导致goroutine永久阻塞。这是因为nil通道没有具体的底层数据结构（如hchan结构体），无法进行正常的发送或接收操作。
// 当执行接收操作时，运行时系统会检查通道是否为nil，如果是，则会将当前的goroutine挂起，等待永远不会到来的数据，从而导致永久阻塞。
func six1() {
	var ch chan int
	fmt.Println(ch)
	go func() {
		ch = make(chan int, 1)
		fmt.Println("成功初始化:", ch)
		ch <- 1
	}()
	go func(ch chan int) {
		time.Sleep(time.Second)
		<-ch
	}(ch) // 值传递，传入的 ch 为 nil
	c := time.Tick(1 * time.Second)
	for range c {
		fmt.Printf("#goroutines:%d\n", runtime.NumGoroutine())
	}
}

func six2() {
	ch := make(chan int)
	//var ch chan int
	go func() {
		ch <- 0
	}()
	<-ch
}

func seven() {
	var m sync.Map
	m.LoadOrStore("a", 1)
	m.Delete("a")
	//fmt.Println(m.Len())// 没有 m.Len() 方法
}

// append() 不是并发安全的
func seven1() {
	var wg sync.WaitGroup
	wg.Add(2)
	var ints = make([]int, 0, 1000)
	go func() {
		for i := 0; i < 1000; i++ {
			ints = append(ints, i)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			ints = append(ints, i)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println(len(ints))
}

// go 语法规定，小写开头的写法、属性或 struct 是私有的。在 json 解码或转码时无法实现私有属性的转换
type People struct {
	name string `json:"name"`
}

func eight() {
	js := `{
"name":"11"
}`
	var p People
	err := json.Unmarshal([]byte(js), &p)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("people:", p)
}

var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "0.0.0.0", "ip address")
	flag.IntVar(&port, "port", 8000, "port number")
}

// 注意 flag 包的用法
func eight1() {
	flag.Parse()
	fmt.Printf("%s:%d", ip, port)
}

// panic 协程还没来得及执行，chan 就 close() 了，往已经关闭的 chan 写数据会引发 panic
func nine() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println("a:", a)
		}
	}()
	close(ch)
	fmt.Println("ok")
	time.Sleep(time.Second * 20)
}

func nine1() {
	ch := make(chan int)
	close(ch)
	ch <- 1
}

type S struct {
	v int
}

// 升序排序
func nine2() {
	s := []S{{1}, {2}, {3}}
	sort.Slice(s, func(i, j int) bool { return s[i].v < s[j].v }) // 需要记住用法
	fmt.Printf("%#v", s)
}

type T struct {
	V int
}

func (t *T) Incr(wg *sync.WaitGroup) {
	t.V++
	wg.Done()
}

func (t *T) Print() {
	time.Sleep(1)
	fmt.Print(t.V)
}

// 随机输出大小写字母
func Ten() {
	var wg sync.WaitGroup
	wg.Add(10)
	var ts = make([]*T, 10)
	for i := 0; i < 10; i++ {
		ts[i] = &T{i}
	}
	for _, t := range ts {
		go t.Incr(&wg)
	}
	wg.Wait()
	for _, t := range ts {
		t.Print()
	}
	time.Sleep(5 * time.Second)

}

// 随机输出大小写字母
func Ten1() {
	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(2 * N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			//runtime.Gosched()
			fmt.Printf("%c", 'a'+i)
		}(i)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%c", 'A'+i)
		}(i)
	}
	wg.Wait()
}

const N = 26

// 随机输出大写字母，再输出小写字母
func Ten2() {
	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(2 * N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			runtime.Gosched() // runtime.Gosched() 的核心作用是通过强制让出 CPU，确保大写字母 goroutine 优先执行。结合 GOMAXPROCS(1) 的单线程限制
			fmt.Printf("%c", 'a'+i)
		}(i)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%c", 'A'+i)
		}(i)
	}
	wg.Wait()
}

// fatal error: concurrent map read and map write
// map 并发读写引发 fatal
func Ten3() {
	m := make(map[int]int) // 创建 map

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动写 goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			m[i] = i // 并发写操作
		}
	}()

	// 启动读 goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_ = m[i] // 并发读操作
		}
	}()

	wg.Wait()
}
