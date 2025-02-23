package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"testing"
)

func main() {
	eight()
}

func one() {
	// 取反
	fmt.Println(^1)
}

type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}

func (s *Slice) Add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Print(elem)
	return s
}

func one1() {
	s := NewSlice()
	// defer 函数中的参数（包括接收者）是在 defer 语句出现的位置做计算的，而不是函数正在执行的时候做计算的。
	// 因此， s.Add(1) 会先于 s.Add(3) 执行
	defer s.Add(1).Add(2)
	s.Add(3)
}

func two() {
	s := NewSlice()
	defer func() {
		// s.Add(1).Add(2) 作为一个整体包在一个匿名函数中，会延迟执行
		s.Add(1).Add(2)
	}()
	s.Add(3)
}

type Orange struct {
	Quantity int
}

func (o *Orange) Increase(n int) {
	o.Quantity += n
}

func (o *Orange) Decrease(n int) {
	o.Quantity -= n
}

func (o *Orange) String() string {
	return fmt.Sprintf("%#v", o.Quantity)
}

func two1() {
	// String() 是指针方法，而不是值方法，因此使用 Println() 输出时不会调用到 String() 方法
	var orange Orange // 输出 {5}
	//orange := &Orange{} // 输出 5
	orange.Increase(10)
	orange.Decrease(5)
	fmt.Println(orange)
}

func test() []func() {
	var funs []func()
	for i := 0; i < 2; i++ {
		funs = append(funs, func() {
			println(&i, i)
		})
	}
	return funs
}

func three() {
	funs := test()
	for _, f := range funs {
		f()
	}
}

var f = func(i int) {
	print("x")
}

// 注意不是递归输出 109877654321
func three2() {
	f := func(i int) {
		print(i)
		if i > 0 {
			f(i - 1)
		}
	}
	f(10)
}

// go 的调度器在 1.14 版本之前使用了协作式调度（需要主动让出 CPU），这使得 for{} 独占 cpu 资源，导致其他 goroutine 被饿死。
// 然后 go 在1.14之后使用了抢占式调度(基于信号的抢占式调度)，因此子 goroutine 有机会执行。具体来说就是，子 goroutine 的 fmt.Print()
// 触发了抢占，因此得以执行。
// 拓展内容：抢占式调度依赖检查点：函数调用或循环中的特定条件会触发抢占，空循环有可能会绕开这些检查点
func four() {
	// 程序只会使用一个操作系统线程，即单核运行
	runtime.GOMAXPROCS(1)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Print(i)
		}
	}()
	for {
	}
}

// 修改
func four1() {
	runtime.GOMAXPROCS(1)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		os.Exit(0)
	}()

	// select{} 属于阻塞操作，允许调度器将 CPU 时间让给其他 goroutine。for{} 是忙等待，不触发调度。
	// 另外使用 select{} 作为阻塞方法不会占用 CPU，而 for{} 会占用 100% CPU
	select {}
}

func five() {
	f, err := os.Open("file")
	if err != nil {
		return
	}
	defer f.Close() // 该语句放在 err 判断之后
	b, err := ioutil.ReadAll(f)
	println(string(b))
}

type S1 struct{}

func (s1 S1) f() {
	fmt.Println("S1.f()")
}

func (s1 S1) g() {
	fmt.Println("S1.g()")
}

type S2 struct {
	S1
}

func (s2 S2) f() {
	fmt.Println("S2.f()")
}

type I interface {
	f()
}

func printType(i I) {
	fmt.Printf("%T\n", i)

	if s1, ok := i.(S1); ok {
		s1.f()
		s1.g()
	}
	if s2, ok := i.(S2); ok {
		s2.f()
		s2.g()
	}
}

// 类型断言，S2 嵌套了 S1。S2 自己没有实现 g(),因此调用的是 S1 的 g()
func six() {
	printType(S1{})
	printType(S2{})
}

func six2() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("1")
		wg.Done()
		//wg.Add(1) 使用 wg.Add() 没有使用 wg.Done() ,导致 panic

	}()
	wg.Wait()
}

// cap 函数的用法
func seven() {
	arr := [2]int{1, 2}
	// 返回数组元素个数
	fmt.Println(cap(arr))

	slice := []int{1, 2}
	// 返回 slice 的最大容量
	fmt.Println(cap(slice))

	ch := make(chan int, 2)
	// 返回 channel 的容量
	fmt.Println(cap(ch))
}

// 可变函数是指针传递
func hello(num ...int) {
	num[0] = 18
}

func Test(t *testing.T) {
	i := []int{5, 6, 7}
	hello(i...)
	fmt.Println(i[0])
}
func seven2() {
	t := &testing.T{}
	Test(t)
}

func alwaysFalse() bool {
	return false
}

// 输出 true。go 代码断行规则导致的
func eight() {
	switch alwaysFalse(); // go 编译器会自动在这后面加上 ';' , 这使得 switch alwaysFalse() 等价于 switch alwaysFalse(); true 。因此输出 true
	{
	case true:
		fmt.Println(true)
	case false:
		fmt.Println(false)
	}
}

type ConfigOne struct {
	Daemon string
}

// 拓展: 如果类型重新定义 String() 方法，使用 Printf()、Print()、Println()、Sprintf() 等格式化输出是会自动使用 String() 方法。
// 因此：在 String() 方法中应避免直接或间接调用自身。解决办法有：直接访问字段值或使用 %#v 等绕过 String() 的格式化动词。
// eg1: return fmt.Printf("print:"%v",c.Daemon)
// eg2: fmt.Sprintf("print: %#v", c)
func (c *ConfigOne) String() string {
	return fmt.Sprintf("print: %#v\n", c)
}

// 触发无限递归循环，导致栈溢出。
// 递归逻辑：fmt.Sprintf → 调用 c.String() → 内部再次调用 fmt.Sprintf → 再次调用 c.String() → 循环直至栈溢出
func nine() {
	c := &ConfigOne{}
	c.String()
	fmt.Println(c) // 重新定义了 String() 方法后，使用 fmt.Println() 后，会自动调用自定义后的 String() 方法
}

// 下面代码的问题：
// 1.在协程中使用了 wg.Add()
// 2.使用了 sync.WaitGroup 副本
func ten() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		go func(wg sync.WaitGroup, i int) {
			wg.Add(1)
			fmt.Printf("i:%d\n", i)
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
	fmt.Println("exit")
}

// 修正1：
// 不要使用 wg 的副本，并且将 wg.Add() 写在协程外
func ten1() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("i:%d\n", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("exit")
}

// 修正2：
// 不适用 wg 的副本（wg 是一个指针），并且将 wg.Add() 写在协程外
func ten2() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			fmt.Printf("i:%d\n", i)
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
	fmt.Println("exit")
}
