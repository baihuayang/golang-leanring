package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reader"
	"runtime"
	"strings"
	"time"
	"tree"
	"wc"

	_ "github.com/go-sql-driver/mysql"
	// "code.google.com/p/go-tour/wc"
	// "pic"
	// "kanon/stringutil"
	// "unicode"
)

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func nextInt(b []byte, pos int) (value, nextPos int) {
	for ; pos < len(b) && !isDigit(b[pos]); pos++ {
	}
	x := 0
	for ; pos < len(b) && isDigit(b[pos]); pos++ {
		x = x*10 + int(b[pos]) - '0'
	}
	return x, pos
}

func compare(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		switch {
		case a[i] > b[i]:
			return 1
		case a[i] < b[i]:
			return -1
		}
	}
	switch {
	case len(a) > len(b):
		return 1
	case len(a) < len(b):
		return -1
	}
	return 0
}

func sliceAppend(slice *[]int) {
	*slice = append(*slice, 1)
	fmt.Println(*slice, "func")
}

type Student struct {
	name  string
	age   int
	score int
	next  *Student
}

func (s Student) String() string {
	return fmt.Sprintf("name:%s, age:%d, score:%d", s.name, s.age, s.score)
}

func insertHead(p *Student) *Student {
	for i := 0; i < 10; i++ {
		stu := Student{
			name:  fmt.Sprintf("stu%d", i),
			age:   rand.Intn(100),
			score: rand.Intn(100),
		}
		stu.next = p
		p = &stu
	}
	return p
}

var c, java, python, lua bool

const (
	Big   = 1 << 100
	Samll = Big >> 99
)

func needInt(x int) int {
	return x*10 + 1
}

func needFloat(x float64) float64 {
	return x * 0.1
}

func Println(a ...interface{}) {
	fmt.Println(a...)
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %g", e)
}

func sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	z := 1.0
	old := z
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
		// fmt.Println(z, i)
		if old == z {
			return z, nil
		}
		old = z
	}
	return z, nil
}

type (
	Vertex struct {
		x float64
		y float64
	}
)

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (this *Vertex) setVertex(x, y float64) {
	this.x = x
	this.y = y
}

func Pic(dx, dy int) [][]uint8 {
	ret := make([][]uint8, dy)
	for i := range ret {
		ret[i] = make([]uint8, dx)
	}

	for y, row := range ret {
		for x := range row {
			ret[x][y] = uint8(x * y)
		}
	}

	return ret
}

var m = map[string]int{
	"kanon":   100,
	"scarlet": 120,
}

func WordCount(s string) map[string]int {
	ret := make(map[string]int)
	for _, w := range strings.Fields(s) {
		ret[w]++
	}
	return ret
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	pre, cur := 1, 0
	return func() int {
		// count += 1
		// if count == 1 {
		// 	return 1
		// }
		cur, pre = cur+pre, cur
		// cur += pre
		// pre = cur - pre
		return cur
	}
}

type IPAddr [4]byte

func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func run() error {
	rand.Seed(time.Now().Unix())
	r := rand.Intn(10)
	Println("rand int = ", r)
	if r < 5 {
		return &MyError{
			time.Now(),
			"it didn't work",
		}
	}
	return nil
}

type MyReader struct{}

func (r MyReader) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		b[i] = 'A'
	}
	return len(b), nil
}

type rot13Reader struct {
	r io.Reader
}

func (r *rot13Reader) Read(b []byte) (n int, err error) {
	n, err = r.r.Read(b)
	if err != nil {
		return n, err
	}
	for i := 0; i < n; i++ {
		// fmt.Printf("before = %c, ", b[i])
		switch {
		case b[i] >= 'A' && b[i] <= 'Z':
			b[i] = 'A' + ((b[i]+13)-'A')%26
		case b[i] >= 'a' && b[i] <= 'z':
			b[i] = 'a' + ((b[i]+13)-'a')%26
		default:
		}
		// fmt.Printf("after = %c\n", b[i])
	}
	// fmt.Printf("%q\n", b[:n])
	return n, err
}

type hello struct{}

func (h *hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello!")
}

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}

type Struct struct {
	Greeting string
	Punct    string
	Who      string
}

func (s *Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprint(s.Greeting, s.Punct, s.Who))
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	Println("111111")
	c <- sum
	Println("done!!!")
}

func gofibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y
		c <- x
	}
	close(c)
}

func walkImpl(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walkImpl(t.Left, ch)
	ch <- t.Value
	walkImpl(t.Right, ch)
}

func Walk(t *tree.Tree, ch chan int) {
	walkImpl(t, ch)
	close(ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		x, ok1 := <-ch1
		y, ok2 := <-ch2
		if ok1 != ok2 || x != y {
			return false
		}
		if !ok1 {
			break
		}
	}

	// for i := 0; i < 10; i++ {
	// 	x, y := <-ch1, <-ch2
	// 	if x != y {
	// 		return false
	// 	}
	// }
	return true
}

func quickSort(values []int, start int, end int) {
	if start >= end {
		return
	}
	i, j := start, end
	key := values[start]
	for i < j {
		for i < j && values[j] >= key {
			j--
		}
		values[i], values[j] = values[j], values[i]
		for i < j && values[i] <= key {
			i++
		}
		values[j], values[i] = values[i], values[j]
	}
	quickSort(values, 0, i)
	quickSort(values, i+1, end)
}

func QuickSort(values []int) {
	if len(values) <= 1 {
		return
	}
	quickSort(values, 0, len(values)-1)
}

func Quick2Sort(values []int) {
	if len(values) <= 1 {
		return
	}
	mid, i := values[0], 1
	head, tail := 0, len(values)-1
	for head < tail {
		if values[i] > mid {
			values[i], values[tail] = values[tail], values[i]
			tail--
		} else {
			values[i], values[head] = values[head], values[i]
			i++
			head++
		}
	}
	QuickSort(values[:head])
	Quick2Sort(values[head+1:])
}

func testSlice(sl []int) {
	sl[0] = 5
}

type Pointer *float32

func Float32ToPoint(f float32) Pointer {
	return Pointer(&f)
}

func Kbcs(params ...interface{}) {
	// println(reflect.TypeOf(params))
	println(len(params))
}

func main() {
	rand.Seed(time.Now().Unix())
	// fmt.Printf(stringutil.Revers e(stringutil.Reverse("hello kanonlee!\n")))
	// b := []byte{'1','2','3','a','b','c','4','5','6'}
	// fmt.Println(unhex('f'))
	// fmt.Println(nextInt(b, 0))
	// var stu *Student = new(Student)
	// stu := new(Student)
	// stu.name = "Kanon"
	// stu.age = 26
	// stu.score = 59
	// fmt.Printf("%p\n", stu)
	// stu = insertHead(stu)
	// var i int
	// fmt.Println(i, c, java, python, lua)

	// fmt.Println(needInt(Samll))
	// fmt.Println(needFloat(Samll))
	// fmt.Println(needFloat(Big))
	// Println(sqrt(2))

	switch os := runtime.GOOS; os {
	case "darwin":
		Println("OS X")
	case "linux":
		Println("Linux")
	default:
		fmt.Printf("%s", os)
	}

	// vertex := Vertex{x:0, y:0}
	// vertex.setVertex(1,2)
	// Println(vertex)

	// s := make([]int, 5, 5)
	// fmt.Printf("%p", s)
	// Println()
	// s = append(s, 1)
	// fmt.Printf("%p, cap = %d", s, cap(s))
	// Println(s)

	// pic.Show(Pic)
	// for k, v := range m {
	// 	Println(k, v)
	// }
	wc.Test(WordCount)
	pos, neg := adder(), adder()
	for i := 1; i < 10; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}
	Println("==============")
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}

	v := Vertex{3, 4}
	Println(v.Abs())

	stu := Student{name: "kanon", age: 26, score: 100}
	Println(stu)

	addrs := map[string]IPAddr{
		"loopbak":   {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for n, a := range addrs {
		fmt.Printf("%v: %v\n", n, a)
	}

	// if err := run(); err != nil {
	// 	Println(err)
	// }
	// r := strings.NewReader("hello kanonlee!")
	// b := make([]byte, 8)
	// for {
	// 	n, err := r.Read(b)
	// 	fmt.Printf("n = %v, err = %v, b = %v\n", n, err, b)
	// 	fmt.Printf("b[:n] = %q\n", b[:n])
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }

	reader.Validate(MyReader{})
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
	Println()
	// var h hello
	// err := http.ListenAndServe("localhost:4000", &h)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// http.Handle("/string", String("i am fucking you bitch!!"))
	// http.Handle("/struct", &Struct{"fuck you", ":", "kanon"})
	// log.Fatal(http.ListenAndServe("localhost:4000", nil))
	// Println()
	// go say("world")
	// say("hello")

	// a := []int{7, 2, 8, -9, 4, 0}
	// c := make(chan int)
	// go sum(a[:len(a)/2], c)
	// go sum(a[len(a)/2:], c)
	// x, y := <-c, <-c
	// x := <-c
	// Println(x)
	// fmt.Println(x, y, x+y)
	// c := make(chan int, 10)
	// go gofibonacci(cap(c), c)
	// for i := range c {
	// 	Println(i)
	// }
	// tick := time.Tick(100 * time.Millisecond)
	// boom := time.After(500 * time.Millisecond)
	// for {
	// 	select {
	// 	case <-tick:
	// 		Println("Tick.")
	// 	case <-boom:
	// 		Println("BOOM!")
	// 	default:
	// 		Println("     .")
	// 		time.Sleep(50 * time.Millisecond)
	// 	}
	// }
	// for _, v := range rand.Perm(10) {
	// 	Println(v)
	// }
	// Println(Same(tree.New(11), tree.New(12)))
	// println("======quicksort========")
	// values := []int{99, 9, 1, 2, 5, 12, 12, 0, 63, 51, 73, 89, 10, 99}
	// values2 := []int{99, 9, 1, 2, 5, 12, 12, 0, 63, 51, 73, 89, 10, 99}
	// QuickSort(values)
	// Println(values)
	// Quick2Sort(values2)
	// Println(values2)
	// ff := float32(16.7)
	// sl := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	// sl2 := sl[:0]
	// sl2[0] = 11
	// Println(sl2)	Println(base.Int16ToBytes(256))
	// Println(2 | 4<<3)
	// Kbcs(1, "22")

	// Kbcs(1, 2, 3, 4)
	// // si := []int{1, 2, 3}
	// // ptype := reflect.5TypeOf(si)
	// ss := []Student{}
	// // num := 425
	// ptype2 := reflect.TypeOf(ss)
	// println(ptype2.String(), ptype2.Kind().String())

	// ts := make([]int, 5, 10)
	// println(len(ts), ts)

	// ts = append(ts, make([]int, 3)...)
	// println(len(ts), ts)

	db, err := sql.Open("mysql", "root:ljnleon042!@#@tcp(localhost:3306)/kanet")
	if err != nil {
		panic(err.Error())
	}
	println("mysql connect success")
	defer db.Close()

	rows, err := db.Query("SELECT * From `title`")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bbsID string
		var title string
		var category string
		var lastTime string
		var author string
		var lastReplyUser string
		var replyCount string

		if err := rows.Scan(&bbsID, &title, &category, &lastTime, &author, &lastReplyUser, &replyCount); err != nil {
			log.Fatal(err)
		}
		fmt.Println(rows.Columns())
		fmt.Println(bbsID, title, category, lastTime, author, lastReplyUser, replyCount)
	}

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// res, err := tx.Exec("INSERT INTO `title` VALUES(?,?,?,?,?,?,?)", 100005, "调理农务系", 3, 1552323, 100003, 100003, 2)
	// if err != nil {
	// 	log.Fatalln(err)
	// } else {
	// 	n, err := res.RowsAffected()
	// 	if err != nil {
	// 		log.Fatalln()
	// 	}
	// 	println("row affected :", n)
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	err = tx.Rollback()
	// 	println("row back")
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

}
