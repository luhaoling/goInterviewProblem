package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func main() {
	ten1()
}

// 除 init() 函数外，一个包内不允许有其他同名函数
func init() {
	fmt.Println("test1")
}
func init() {
	fmt.Println("test2")
}

func f() {}

// func f(){}
func one() {}

var x int

//https://cloud.tencent.com/developer/article/2138066
// init 函数特点：
// 1. init 函数是可选的，可以没有
// 2. 与 main 函数一样，不能有入参与返回值
// 3. 与 main 函数一样，自动执行，并且不能被其他函数调用
// 4. 一个包内可以有多个 init 函数。即可以在包的多个源文件中定义多个 init 函数。一般建议在与包同名源文件中写一个。这样可读性好且便于维护
// 5. 一个源文件可以有多个 init 函数

// init 函数执行顺序
//  1. 单个文件的 init() 函数调用顺序与其定义顺序一致，从上到下
//  2. 同一个包中不同源文件的 init 函数执行顺序是根据文件名的字典序来确定的
//  3. main 包导入多个包时 init 执行顺序
//     对于不同的包，如果不相互依赖的话，按照 main 包中导入顺序调用包的 init 函数，最后再调用 main 包中的 init 函数
func init() {
	x++
	fmt.Println("x")
}
func two() {
	//init() //init() 函数不能被其他函数（包括 main 函数调用）
	fmt.Println(x)
}

// 比较大小的一个新方法：copy 函数。返回长度较小的那一个
func min(a int, b int) {
	var min = 0
	min = copy(make([]struct{}, a), make([]struct{}, b))
	fmt.Printf("The min of %d and %d is %d\n", a, b, min)
}

func three() {
	min(123, 256)
}

// four
// C 函数不能通过编译。default 属于关键字。string 和 len 是预定义表示符，可以在局部使用。
func A(string string) string {
	return string + string
}

func B(len int) int {
	return len + len
}

// 不能通过编译
//func C(val,default string)string{
//	if val==""{
//		return default
//	}
//	return val
//}

var nil = new(int) // 值为 new(int) 返回的地址，动态类型为 *int
// 注意：指针判断是否相等时，有两个条件：值和动态类型
func four1() {
	var p *int // 声明时，动态类型是 *int,值为 nil
	fmt.Printf("%#v,%v\n", nil, nil)
	fmt.Printf("%#v,%v\n", p, p)

	if p == nil {
		fmt.Println("p is nil")
	} else {
		fmt.Println("p is not nil")
	}
}

type foo struct{ Val int }
type bar struct{ Val int }

func five() {
	a := &foo{Val: 5}
	b := &foo{Val: 5}
	fmt.Println(a == b)
	fmt.Printf("%p\n", a)
	fmt.Printf("%p\n", b)
	c := foo{Val: 5}
	d := bar{Val: 5}
	fmt.Println(c == foo(d)) // 强制转换
	e := bar{Val: 5}
	f := bar{Val: 5}
	fmt.Println(e == f)
}

func A1() int {
	time.Sleep(100 * time.Millisecond)
	return 1
}

func B1() int {
	time.Sleep(1000 * time.Millisecond)
	return 2
}

// 注意不要被 A1、B1 函数中的 time.Sleep() 误导。
// 在 select 中，如果有多个 case 同时就绪（通道操作可以立即完成）时，select 会随机选择一个执行
func five1() {
	ch := make(chan int, 1)
	go func() {
		select {
		case ch <- A1():
		case ch <- B1():
		default:
			ch <- 3
		}
	}()
	fmt.Println(<-ch)
}

type Point struct{ x, y int }

func six() {
	s := []Point{
		{1, 2},
		{3, 4},
	}
	for _, p := range s {
		p.x, p.y = p.y, p.x
	}
	fmt.Println(s)
}

//修复代码

func six1() {
	s := []*Point{
		&Point{1, 2},
		&Point{3, 4},
	}
	for _, p := range s {
		p.x, p.y = p.y, p.x
	}
	fmt.Println(*s[0])
	fmt.Println(*s[1])
}

// 隐患：get() 函数返回的切片与原切片共用底层数组，如果在调用函数里面修改返回的切片，将会影响原切片
func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0])
	return raw[:3]
}

func seven() {
	data := get()
	fmt.Println(len(data), cap(data), &data[0])
}

// 修改
func getUpdate() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0])
	res := make([]byte, 3)
	copy(res, raw[:3])
	return res
}

func seven1() {
	data := get()
	fmt.Println(len(data), cap(data), &data[0])
}

func modifyMap(m map[string]interface{}) {
	m["age"] = 20
	m["sex"] = "man"
}

func modifySlice(m []int) {
	m[0] = 10
	m[1] = 20
}

// 在Go语言中，map作为引用类型(参数传入的是指针[指向 hmap 结构的指针])，在函数间传递时发生的元素级修改（增删改）都会直接影响原始map，只有对map变量本身的重新赋值（改变指针指向）才不会影响原始map
// go 中，严格意义的引用类型有三种：slice、map、channel，它们可以直接操作共享的底层数据
// 非严格引用类型：接口、函数。
// 接口：存储指针，行为取决于具体类型
// 函数：传递函数指针，不拷贝代码逻辑
// 一般默认接口和函数也属于引用类型
// todo：根据训练营的内容，准确划分 go 数据类型
func eight() {
	// 在 map 反序列化时 json.unmarshal() 的入参必须是 map 的地址（json.Unmarshal() 要求第二个参数必须为指针类型）
	data := []byte(`{"name":"Alice","age":25}`)
	var m map[string]interface{}
	err := json.Unmarshal(data, &m) // 注意：如果 第二个参数没有传入指针，goland 并不会爆红
	fmt.Println("err:", err)
	fmt.Println(m)
	modifyMap(m)
	fmt.Println(m)
	arr := []int{1, 2}
	fmt.Println(arr)
	modifySlice(arr)
	fmt.Println(arr)
}

// 使用值类型接收者定义的方法，调用的时候，使用的是值的副本，对副本操作不会影响原来的值。如果想要在调用函数中修改原值，可以使用指针接受者的方法
type Foo struct {
	val int
}

func (f *Foo) Inc(inc int) {
	f.val += inc
}

func (f Foo) Inc1(inc int) {
	f.val += inc
}

// 总结：值类型和指针类型都可以调用值接收者方法和指针接收者方法
// 在 go 中，值类型可以调用指针接收者方法（编译器隐式取地址）
// 在 go 中，指针类型可以调用值接收者方法（编译器隐式解引用）
// 只有可寻址的值类型才能调用指针接受者方法。
// 空指针调用值接收者方法会导致 panic
func eight1() {
	f := &Foo{}
	var a Foo
	a.Inc1(100)
	fmt.Println("a", a.val)
	f.Inc(100)
	fmt.Println(f.val)

	Foo{}.Inc1(11)
	//只有可寻址的值类型才能调用指针接受者方法。
	//Foo{}.Inc(11)

	// 空指针调用值接收者(或指针接受者)方法会导致 panic
	var f1 *Foo
	f1.Inc1(100)

}

func testq(i int) (ret int) {
	ret = i * 2
	if ret < 10 {
		ret = 10
		return ret
	}
	return
}

func nine() {
	result := testq(10)
	fmt.Println(result)
}

// 奇诡代码：不建议这么做
func nine1() {
	true := false
	fmt.Println(true)
}

func waShadwDefer(i int) (ret int) {
	ret = i * 2
	// 局部变量
	if ret > 10 {
		ret := 10
		defer func() {
			ret = ret + 1
		}()
	}
	return
}

func nine2() {
	result := waShadwDefer(50)
	fmt.Println(result)
}

func ten() {
	m := map[string]int{
		"G": 7, "A": 1,
		"C": 3, "E": 5,
		"D": 4, "B": 2,
		"F": 6, "I": 9,
		"H": 8,
	}
	var order []string
	var order1 []string
	// 注意，k 是 key
	for k, _ := range m {
		order = append(order, k)
	}
	fmt.Println(order)
	// 注意，k 是 key
	for q := range m {
		order1 = append(order, q)
	}
	fmt.Println(order1)
}

type UserAges struct {
	ages       map[string]int
	sync.Mutex // 相当于 UserAge 继承了 sync.Mutex。所以 UserAges 可以调用 Mutex 的方法
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

// map 并发读写不安全，存在读写资源竞争的情况
func (ua *UserAges) Get(name string) int {
	ua.Lock()
	defer ua.Unlock()
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

func ten1() {
	count := 1000
	gw := sync.WaitGroup{}
	gw.Add(count * 3)
	u := UserAges{ages: map[string]int{}}
	add := func(i int) {
		u.Add(fmt.Sprintf("user_%d", i), i)
		gw.Done()
	}
	for i := 0; i < count; i++ {
		go add(i)
		go add(i)
	}
	for i := 0; i < count; i++ {
		go func(i int) {
			defer gw.Done()
			u.Get(fmt.Sprintf("user_%d", i))
		}(i)
	}
	gw.Wait()
	fmt.Println(len(u.ages))
	fmt.Println("Done")
}
