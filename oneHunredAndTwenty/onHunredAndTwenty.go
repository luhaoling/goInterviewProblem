package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	seven()
}

// 使用 println() 打印结果不同。每个 goroutine 都会分配相应的栈内存。随着程序的运行，栈内容发生增长或者缩小，协程会重新申请栈内存块。
// 如这题，循环调用 f() ，发生深度递归，栈内存不断增大，当超过范围时，重新申请栈内存，然后 val 的地址会发生变化。
// 如果将 println() 换成 fmt.Println() 后会发现打印结果相同。这是因为 fmt.Println() 使变量 val 发生了逃逸，逃逸到堆内存，
// 即使协程重新申请资源，val 变量在堆内存的地址也不会改变

// println() 和 fmt.Println() 的底层差异
// println() 是 go 的内置函数，属于 builtin 包，直接输出到标准错误，不涉及接口或反射操作。它的参数处理更简单，
// 编译器可以将其参数分配到栈上，且不会强制变量逃逸到堆上
// fmt.Println() 是标准库函数，接收 interface{} 类型的参数。由于需要支持任意类型的动态解析，底层会触发反射和接口类型转换，
// 导致编译器无法确定变量的声明周期，从而强制变量逃逸到堆。
func one() {
	var val int
	println(&val)
	//fmt.Println(&val)
	f(10000)
	println(&val)
	//fmt.Println(&val)
}

func f(i int) {
	if i--; i == 0 {
		return
	}
	f(i)
}

// a 是指向变量 val 的指针，val 变量的地址发生了改变，a 指向 val 新的地址是由内存管理自动实现的
func one1() {
	var val int
	a := &val
	println(a)
	f(10000)
	b := &val
	println(b)
	println(a == b)
}

func two() {
	x := []int{100, 200, 300, 400, 500, 600, 700}
	// 在声明时直接初始化的切片的长度和容量都一样
	println("x len and cap", len(x), cap(x))
	twohundred := &x[1]
	x = append(x, 800) // 容量不足，重新申请内存
	for i := range x {
		x[i]++
	}
	fmt.Println(*twohundred)

	x = make([]int, 0, 8)
	x = append(x, 100, 200, 300, 400, 500, 600, 700)
	twohundred = &x[1]
	x = append(x, 800) // 容量足够，不会重新申请内存
	for i := range x {
		x[i]++
	}
	fmt.Println(*twohundred)
}

// 错误声明方法
var Err error = errors.New("ele")

// 常量声明
const Pi float64 = 3.1415926
const zero = 0.0
const (
	size int64 = 1024
	eof        = -1
)

const u, vfloat32 = 0, 3
const a, b, c = 3, 4, "ze"

func link(p ...interface{}) {
	fmt.Println(p)
}

func three() {
	link("seek", 1, 2, 3, 4)
	a := []int{1, 2, 3, 4}
	link("seek", a)

	tmplink := make([]interface{}, 0, len(a)+1)
	tmplink = append(tmplink, "seek")
	for _, li := range a {
		tmplink = append(tmplink, li)
	}
	link(tmplink...)
}

func four() {
	// 0 开头表示八进制；0x 开头表示16进制
	ns := []int{010: 200, 005: 100}
	print(len(ns))
}

func four1() {
	i := 0
	f := func() int {
		i++
		return i
	}
	c := make(chan int, 1)
	c <- f() // 执行 f() 后往 c 中发送 i，这时 i=1
	select {
	case c <- f(): // 执行这一个分支，执行 f() ,然后往 c 中发送 i，这时 i=2，由于 c 中已有一个值，这里会阻塞，然后执行 default 分支
	default:
		fmt.Println(i)
	}
}

// 拆解 four1
func four2() {
	i := 0
	f := func() int {
		fmt.Println("incr")
		i++
		return i
	}
	c := make(chan int)
	for j := 0; j < 2; j++ {
		select {
		case c <- f():
		default:
			fmt.Println("nihao")
		}
	}
	fmt.Println(i)
}

var y int

func f1(i int) int {
	return 7
}

// switch 的两种用法
func five() {
	//用法1： 无表达式直接判断 case 条件（布尔表达式匹配）
	switch y = f1(2); {
	case y == 7:
		//return
	}
	// 用法1相当于下面的代码
	y = f1(2)
	switch {
	case y == 7:
	}

	// 基于 y 的值匹配（等值比较）
	switch y {
	case 0:
		fmt.Println("0")
	case 7:
		fmt.Println("7")
	case 8:
		fmt.Println("8")
	default:
		fmt.Println("no ")
	}
}

// switch 详细用法：
func five1() {
	// 1. type Switch(类型判断)
	var i interface{} = "hello"
	switch v := i.(type) {
	case int:
		fmt.Println("整型", v)
	case string:
		fmt.Println("字符串", v)
	default:
		fmt.Println("未知", v)
	}

	// 2. 无表达式的布尔条件判断
	a, b := 1, 2
	switch {
	case a == 1 && b == 2:
		fmt.Println("条件1成立")
	case a > 5:
		fmt.Println("条件2成立")
	}

	// 3. fallthrough 强制执行后续分支
	switch a := 1; a {
	case 1:
		fmt.Println("a=1")
		fallthrough
	case 2:
		fmt.Println("a=2")
	}

	// 4. 多值匹配
	num := 3
	switch num {
	case 1, 3, 5:
		fmt.Println("奇数小于6")
	case 2, 4:
		fmt.Println("偶数小于6")
	}

	// 5.灵活处理 default 分支
	switch a := 3; a {
	default:
		fmt.Println("默认")
		fallthrough
	case 1:
		fmt.Println("a==1")
	}

	// 6. 初始化语句与局部作用域
	calculate := func() int {
		return 1
	}
	switch a := calculate(); a {
	case 1:
		fmt.Println("a==1")
	}
}

func five3() {
	a := []int{1, 2, 3, 4}
	b := variadic(a...)
	b[0], b[1] = b[1], b[0]
	fmt.Println(a)
}

// 切片作为参数传入可变函数时不会创建新的切片
func variadic(ints ...int) []int {
	return ints
}

const (
	o = 1 << iota
	two1
)

const (
	greeting = "hello"
	on       = 1 << iota
	tw
)

func six() {
	fmt.Println(o)
	fmt.Println(two1)

	fmt.Println(on)
	fmt.Println(tw)
}

// go 中大多数类型都可以转化为有效的 json 文本，除了 channel、complex、函数外

const N = 10

// 在 1.22.0 之前，执行下面代码 m 的数量不为10
// 循环变量行为变化
// go 1.21 及之前，循环变量 i 在每次迭代的时候复用同一个内存地址，导致竞态
// go 1.22 及之后，每次迭代都会创建新的 i 副本，每个 goroutine 捕获的时当前迭代的独立 i 值，避免竞态
func seven() {
	m := make(map[int]int)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			m[i] = i
			mu.Unlock()
		}()
	}
	wg.Wait()
	println(len(m))
}

func seven1() {
	m := make(map[int]int)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			m[i] = i
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	println(len(m))
}

func eight() {
	// 从一个已经关闭的 channel 中接收数据，如果缓冲区为空，返回一个零值
	ch := make(chan int)
	close(ch)
	fmt.Println(<-ch)
}

const (
	_      = iota
	c1 int = (10 * iota)
	c2
	d = iota
)

func eight1() {
	fmt.Printf("%d-%d-%d", c1, c2, d) // 10-20-3
}

var ErrDidNotWork = errors.New("did not work")

func DoThenThing(reallyDoIt bool) (err error) {
	if reallyDoIt {
		// 这里面的 err 的作用域在 if {} 里面
		result, err := tryTheThing()
		if err != nil || result != "it worked" {
			err = ErrDidNotWork
		}
	}
	return err
}

func tryTheThing() (string, error) {
	return "", ErrDidNotWork
}

func nine() {
	fmt.Println(DoThenThing(true))
	fmt.Println(DoThenThing(false))
}

func Ten() {
	fmt.Println(len("你好bj!"))
}

func Ten1() {
	intmap := map[int]string{
		1: "a",
		2: "bb",
		3: "ccc",
	}
	v, err := GetValue(intmap, 3)
	fmt.Println(v, err)
}

func GetValue(m map[int]string, id int) (string, bool) {
	if _, exist := m[id]; exist {
		return "存在数据", true
	}
	//return nil,false
	return "", false
}
