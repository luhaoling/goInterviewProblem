package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	ten1()
}

func one() {
	var a []int = nil
	a, a[0] = []int{1, 2}, 9 // 先计算左边，再计算右边，然后赋值。计算左边时，a[0] 没有值
	fmt.Println(a)
}

// 数值溢出，当 i 的值为 0、128 时会发生相等的情况。byte 是 int8 的别名
// 溢出规则：超出数值范围的部分会通过模运算被截断
// 快捷计算实际值的规则(eg:int8)：
// 实际值=i-256(i>127)
// eg（129）:129-256=-127

func two() {
	count := 0
	for i := range [256]struct{}{} {
		m, n := byte(i), int8(i)
		if n == -n {
			fmt.Println("n", n)
			fmt.Println("n", -n)
			count++
		}
		if m == -m {
			fmt.Println("m", m)
			fmt.Println("m", -m)
			count++
		}
	}
	fmt.Println(count)
}

const (
	azero = iota
	aone  = iota
)

const (
	c = "ma"
	a = iota
	b = iota
)

// 在一个常量声明代码块中，如果 iota 没有出现在第一行，则常量的值是非0值
func two1() {
	fmt.Println(azero)
	fmt.Println(aone)
	fmt.Println(a)
	fmt.Println(b)
}

type data struct {
	name string
}

func (p data) print() {
	fmt.Println("name:", p.name)
}

type printer interface {
	print()
}

func three() {
	d1 := data{"one"}
	d1.print()

	// 注意点：是否实现了 printer 接口，值方法和指针方法是不同的
	// 值方法不可以调用指针方法；指针方法可以调用值方法和指针方法
	//var in printer = data{"two"}
	var in printer = &data{"two"}
	in.print()
}

func four() {
	a := []int{0, 1, 2}
	s := a[1:2]
	s[0] = 11         // 对 s 的修改同时影响 a（直接在对应的位置上修改）
	s = append(s, 12) // 对 s 的修改同时影响 a（直接在对应的位置上修改）
	s = append(s, 13) // 再次插入的时候，s 发生了扩容，此后对 s 的修改将不会影响 a
	s[0] = 21
	fmt.Println(a)
	fmt.Println(s)
}

func five() {
	// TrimRight() 方法，将第二个参数里面所有字符拿出来处理，只要与其中任何一个字符相等，便会将其删除
	// TrimRight 返回字符串 s 的一个切片，
	//该切片会移除所有末尾（trailing）存在于 cutset 中的 Unicode 码点。
	//(这里的关键点是“trailing”指的是字符串末尾的连续字符，且这些字符在cutset中存在就会被移除)
	//若要移除特定后缀字符串，请改用 [TrimSuffix]。
	fmt.Println(strings.TrimRight("ABBA", "BA"))
	fmt.Println(strings.TrimRight("¡¡¡Hello, Gophers!!!", "!¡"))
	// 正确截取字符串的方法(移除特定后缀)
	fmt.Println(strings.TrimSuffix("ABBA", "BA"))
}

func five1() {
	var src, dst []int
	src = []int{1, 2, 3}
	copy(dst, src) // copy 函数返回 len(dst)、len(src) 之间的最小值, 并且依据这个值把 src 的内容赋值到 dst 里面
	fmt.Println(dst)
}

// 修改
// 预先分配合适的空间
func five1_update1() {
	var src, dst []int
	src = []int{1, 2, 3}
	dst = make([]int, len(src))
	n := copy(dst, src)
	fmt.Println(n, dst)
}

// 使用 append()
func five1_update2() {
	var src, dst []int
	src = []int{1, 2, 3}
	dst = append(dst, src...)
	fmt.Println("dst:", dst)
}

func six() {
	n := 43210
	// * / % 的优先级相同，从左向右结合
	fmt.Println(n/60*60, " hours and ", n%60*60, " seconds")
	// 修复
	fmt.Println(n/(60*60), " hours and ", n%(60*60), " seconds")
}

// 八进制以 0 开头
const (
	Century = 100 // 十进制
	Decade  = 010 // 八进制, 十进制为 8
	Year    = 001 // 八进制, 十进制为 2
)

func six1() {
	fmt.Println(Century + 2*Decade + 2*Year)
}

func seven() {
	// ^ 作为二元运算符时表示按位异或
	fmt.Println(3^2+4^2 == 5^2)  // true
	fmt.Println(6^2+8^2 == 10^2) // false
	// 0011 ^ 0010=0001 =1
	// 0100 ^ 0010=0110 =6
	// 0101 ^ 0010=0111 =7
	// 所以为 true
}

func Foo() error {
	var err *os.PathError = nil // 声明并定义
	return err
}

// 只有在值和动态类型都为 nil 的情况下，接口值才为 nil
// 在 eight 中
// err 表示值为 nil，动态类型为 *os.PathError
// nil 表示值为 nil，动态类型也为 nil

func eight() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)  // false
	fmt.Printf("%#v\n", err) // "%#v\n" 可以打印某一个变量的详情 (*os.PathError)(nil)
}

func Foo1() (err *os.PathError) { // 这里的 err *os.PathError 只是声明，并没有定义
	return
}

func eight1() {
	err := Foo1()
	fmt.Println(err)
	fmt.Println(err == nil)
}

func nine() {
	v := []int{1, 2, 3}
	for i, n := 0, len(v); i < n; i++ { // 终止条件只计算一次，编译可以通过
		v = append(v, i)
	}
	fmt.Println(v)
}

type P *int
type Q *int

func nine1() {
	var p P = new(int)
	*p += 8
	var x *int = p
	var q Q = x
	*q++
	fmt.Println(*p, *q) // 9,9  p 指针和 q 指针指向同一个地址
}

// 不能通过编译；不同类型的值时不能相互赋值的，即使底层类型一样
type T int

func F(t T) {}
func ten() {
	//var q int
	//F(q)
}

// 可以通过编译
// 底层类型相同的变量可以相互赋值的重要条件，至少有一个不是有名类型（unamed type）,该准则也叫：可复制性
// named type: 主要有两类：
// 1. 内置的类型，如 int, int64, float
// 2. 用 type 关键字声明的类型，如 type Foo string
// unamed type: 基于已有的 named types 声明出的组合类型，如 struct{},[]string,interface{},map[string]bool
type T1 []int

func F1(t T1) {}
func ten1() {
	var q []int
	F1(q)
}

// Named 和 Unamed Types

var x struct{ I int } // unamed type

var x2 struct{ I int } // unamed type

type Fo struct{ I int }  // name type
var y Fo                 // named type
type Bar struct{ I int } // named type
var z Bar                // named type
func ten2() {
	x = y
	y = x
	x = x2
	x = z
}
