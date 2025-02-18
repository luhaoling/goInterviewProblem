package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

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

// 一、Named（命名类型） 和 Unamed Types（未命令类型/无名类型）
// named types 有两类
// 1. 内置的类型，如 int,int64,string,bool
// 2. 用 type 关键字声明的类型，如 type Foo string
// unamed types: 基于已有的 named types 声明出来的组合类型，如 struct{},[]string, interface{}, map[int]bool
// named types 可以作为方法的接受者，unamed type 不能
type Map map[string]string

func (m Map) Set(key, value string) {
	m[key] = value
}

// invalid receiver type map[string]string(map[string]string is an unnamed type)
//func (m map[string]string)Set1(key string,values string){
//m[key]=value
//}

// 二、底层类型
// 每种类型 T 都有一个底层类型，如果 T 是预声明类型或者类型字面量，它的底层类型就是 T 本身；
// 否则，T 的底层类型就是其类型声明中引用类型的底层类型
type (
	B1 string
	B2 B1
	B3 []B1
	B4 B3
)

// string, B1 和 B2 的底层类型是 string；
// []B1,B3 和 B4 的底层类型是 []B1
// （判断类型是否相同）
// 因此：所有基于相同 unmaned types 声明的变量类型都相同；
// 而对于 named types 变量而言，即使它们的底层类型相同，它们也是不同类型

var x struct{ I int } // unamed type

var x2 struct{ I int } // unamed type

type Fo struct{ I int }  // name type
var y Fo                 // named type
type Bar struct{ I int } // named type
var z Bar                // named type
func ten2() {
	x = y  // x 和 y 类型相同
	y = x  // x 和 y 类型相同
	x = x2 // x 和 x2 类型相同
	x = z  // x 和 z 类型相同
	// y= z // y 和 z 的类型不同
}

// 三、不同类型之间是不可以赋值的
type MyInt int

var i int = 2
var i2 MyInt

// i=12 （虽然它们的底层类型相同，但是它们是不同的类型）
// 对于拥有相同底层类型的变量而言，还有一个概念是：可复制性。即：底层类型相同的两个变量可以赋值的条件是：至少有一个不是 unamed type
// 即：在Go中，两个类型是否可以直接赋值，取决于它们的底层类型是否相同，以及它们是否都是命名类型
// 如果两个类型具有相同的底层类型，并且至少其中一个是未命名类型，那么它们之间是可以相互赋值的。
// 而如果两个都是命名类型，即使它们的底层类型相同，也不能直接赋值，除非它们是通过类型别名（type alias）定义的。

// 四、关于类型继承
// 当你使用 type 声明一个新类型后，它不会集成原有类型的方法集
type User struct {
	Name string
}

func (u *User) SetName(name string) {
	u.Name = name
}

type Employee User

func test1() {
	//employee:=new(Employee)
	//employee.SetName("jack")
	//error employee.SetName undefined
}

// 一个编程小技巧：可以将原有类型作为一个匿名字段内嵌到 struct 当中继承它的方法

type Employee1 struct {
	User  // annonymous field
	title string
}

func test2() {
	employee := new(Employee1)
	employee.SetName("jack")
}

type o = int // 类型别名
type u int

// 类型别名例外：类型别名与原类型视为同一类型，不受此限制。
func test3() {
	var o1 o
	var u1 u
	o1 = o(u1)
	fmt.Println(o1)
}

// *******************************************
// 总结如下：至少有一个是未命名类型
// 在底层类型相同的情况下，只有参与赋值的类型中有一个是未命名的，则允许赋值。类型别名与原类型视为同一类型，不受此限制,允许赋值。
// 具体示例：
// 1. 命名类型 vs 未命名类型（允许赋值）
type My []int

var a1 My = []int{1}
var b1 []int = a1

// 2. 两个命名类型（禁止赋值）
type c1 int
type c2 int

var c11 c1 = 1

//c2=c1(不能赋值)

// 3. 两个未命名类型（允许赋值）
var a111 []int
var b111 []int = a111

// 4. 类型别名（允许赋值）
type Myq = int

var myq Myq = 10
var myc = myq
