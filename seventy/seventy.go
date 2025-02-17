package main

import "fmt"

func main() {
	ten1()
}

func one() {
	var k = 1
	var s = []int{1, 2}
	// 多重赋值
	// 1. 计算等号左边的索引表达式，接着计算等号右边的表达式
	// 2. 赋值
	k, s[k] = 0, 3
	fmt.Println(s[0] + s[1])
}

func one1() {
	var k = 9
	for _, k = range []int{1} {
		fmt.Println("nihao", k)
	}
	fmt.Println(k)
	for k = 0; k < 3; k++ {
	}
	fmt.Println(k)
	for k = range (*[3]int)(nil) {
	}
	fmt.Println(k)
}

func two() {
	nil := 123
	fmt.Println(nil)
	// var _ map[string]int=nil
	// 引发 panic 这个时候 nil 是 整型，不能赋值给 map[string]int,而不是空
}

func two1() {
	var x int8 = -128
	var y = x / -1
	fmt.Println(y) // 输出 -128,因为溢出了
}

func F(n int) func() int {
	return func() int {
		n++
		return n
	}
}

// defer() 后面的函数如果带有参数，会优先计算参数，并将结果存储在栈中，等到真正执行 defer 的时候取出来
func three() {
	f := F(5)
	defer func(i int) {
		fmt.Println("C", f())
	}(10)
	defer fmt.Println("B", f()) // defer 后的 fmt.Println() 带参数
	i := f()
	fmt.Println("C", i)
}

// recover 必须在 defer() 函数中直接调用才有效
func four() {
	defer func() {
		recover()
	}()
	panic(2)
}

// 触发 panic(2) 后，执行 defer() 函数，先执行倒数第一个 defer 函数，这个 defer 函数里面还有一个 defer 函数，由于这个 defer 函数带参数
// 因此先计算这个函数的参数并将结果压入栈中，因此这个 defer 函数捕获到 panic(2)；然后继续执行剩下的代码，又引发 panic，执行倒数第二个 deger
// 函数，捕获到 panic(1)
func four1() {
	defer func() {
		fmt.Println(recover())
	}()

	defer func() {
		defer fmt.Println(recover())
		panic(1)
	}()

	// defer recover()
	panic(2)
}

// todo 疑惑: 为什么 ”Q“ 这一个 recover 无法捕获 panic(2)
func five() {
	defer func() {
		fmt.Println("C", recover())
	}()
	defer func() {
		defer func() {
			fmt.Println("Q", recover())
		}()
		fmt.Println("nihao")
		//panic(1)
	}()
	panic(2)
}

func five1() {
	panic(1)
}

type T struct {
	n int
}

func six() {
	ts := [2]T{}
	for i, t := range ts { // 使用的是数组 ts 的副本
		switch i {
		case 0:
			t.n = 3 // 由于 t 使用的是 ts 的副本，因此对 t.n 的修改不起作用
			ts[1].n = 9
		case 1:
			fmt.Print(t.n, " ") // 输出 ts 原来的副本值，ts[1].n=9 的修改不起作用，因此输出 0
		}
	}
	fmt.Print(ts) // 输出 [{0} {9}]
}

func six1() {
	ts := [2]T{}
	for i, t := range &ts { // 循环变量 t 是原数组元素的副本
		switch i {
		case 0:
			t.n = 3 // 由于 t 使用的是 ts 的副本，因此对 t.n 的修改不起作用
			ts[1].n = 9
		case 1:
			fmt.Print(t.n, " ") // 输出 ts 更新后的值，ts[1].n=9 ，输出 9
		}
	}
	fmt.Print(ts) // 输出 [{0} {9}]
}

func seven() {
	ts := [2]T{}
	for i := range ts[:] { // 使用的是切片的副本，没有复制底层数组，此副本切片与原数组共享底层数组
		switch i {
		case 0:
			ts[1].n = 9
		case 1:
			fmt.Print(ts[i].n, " ") // 输出 ts 更新后的值，ts[1].n=9 ，输出 9
		}
	}
	fmt.Print(ts) // 输出 [{0} {9}]
}

func seven1() {
	ts := [2]T{}
	for i := range ts[:] { // 循环变量 t 是原数组元素的副本
		switch t := &ts[i]; i {
		case 0:
			t.n = 3 // t 为指针，因此修改有效
			ts[1].n = 9
		case 1:
			fmt.Print(t.n, " ") // 输出 ts 更新后的值，ts[1].n=9 ，输出 9
		}
	}
	fmt.Print(ts) // 输出 [{0} {9}]
}

func eight() {
	// goto 不能跳转到其他函数或者内存代码
	//for i:=0;i<10;i++{
	//	loop:
	//		println(i)
	//}
	//goto loop
}

// 在 go 1.18 中，v 不会重新声明，go 1.22 之后，for-range 循环中，每次循环都会创建一个新的变量
func eight1() {
	x := []int{0, 1, 2}
	y := [3]*int{}
	for i, v := range x {
		defer func() {
			print(v)
		}()
		y[i] = &v
	}
	print(*y[0], *y[1], *y[2])
}

func nine() {
	var t []int
	t = append(t, 1)

	var s []int
	fmt.Println(s == nil)
	s = make([]int, 0) // 相当于分配了空间
	fmt.Println(s == nil)
	s = append(s, 1)

	var m map[string]int
	fmt.Println(m == nil)
	m = make(map[string]int) // map 一定要分配空间后才能往其中插入键值对
	fmt.Println(m == nil)
	m["one"] = 1
}

func test(x int) (func(), func()) {
	return func() {
			println(x)
			x += 10
		}, func() {
			println(x)
		}
}

// 闭包引用相同变量
func nine2() {
	a, b := test(100)
	a()
	b()
}

func ten() {
	// 单引号表示 rune，本质是一个 byte 类型数组
	str := 'c' + 'c'
	fmt.Println(str)
	// 双引号表示 string
	str1 := "abc" + "abc"
	fmt.Println(str1)
	fmt.Sprintf("abc%d", 123)
}

func ten1() {
	println(DeferTest1(1))
	println(DeferTest2(1))
	println(DeferTest3(1))
}

// 函数返回值名字会在函数起始处被初始化为对应类型的零值，其作用域为整个函数。
// 在 return 之前 defer 会被执行

func DeferTest1(i int) (r int) {
	r = i
	defer func() {
		r += 3
	}()
	return r
}

// return 先赋值给 r，然后再执行 defer 中的代码，最后再返回
func DeferTest2(i int) (r int) {
	defer func() {
		fmt.Println("Defer before", r)
		r += i
		fmt.Println("Defer", r)
	}()
	return 2
}

func DeferTest3(i int) int {
	t := i
	defer func() {
		t += 3 // 这个时候，对 t 的修改不影响返回值
	}()
	return t
}
