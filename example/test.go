package main

import (
	"../lua"
	"fmt"
	"runtime"
	//"github.com/dongaihua/golua/lua"
	//"runtime"
	"time"
)

// func adder(L *lua.State) int {
// 	a := L.ToInteger(1)
// 	b := L.ToInteger(2)
// 	L.PushInteger(a + b)
// 	return 1 // number of return values
// }

func p(p interface{}) {
	fmt.Printf("%v\n", p)

}

var x = 0
var canExit = false

func inc_x() { //test
	for {
		x += 1
	}
}

func print() {
	for {
		x = x + 1
		fmt.Println(x)
	}
}

func print1(p interface{}, times int) {
	for i := 0; i < times; i++ {
		fmt.Printf("%v, %v\n", i, p)
		//runtime.Gosched()
	}
}

func print2(p interface{}, times int) {
	for i := 0; i < times; i++ {
		fmt.Printf("%v, %v\n", i, p)
		//runtime.Gosched()
	}
	canExit = true
}
func readConfig(L *lua.State) {
	// L.GetGlobal("width")
	// width := L.ToNumber(-1)
	// p(width)
	// L.DumpStack()
	// L.Pop(1)
	// L.DumpStack()

	// L.GetGlobal("height")
	// height := L.ToNumber(-1)
	// p(height)
	// L.Pop(1)
	// L.DumpStack()

	// L.GetGlobal("background")
	// L.DumpStack()
	// L.PushString("r")
	// L.GetTable(-2)
	// red := L.ToNumber(-1)
	// p(red)
	// L.Pop(2)
	// L.DumpStack()

	L.DumpVar("array")

	L.GetGlobal("array")
	p(L.IsTable(-1))
	L.PushValue(-1)
	L.DumpStack()
	// L.DumpStack()
	// L.PushInteger(1)
	// L.GetTable(-2)
	// value := L.ToString(-1)
	// p(value)
	// L.Pop(2)
	// L.DumpStack()

}

func readConfig2(L *lua.State) {
	//table := L.GetTableX("background")
	//p(table)

}

func exec(L *lua.State) {
	L.Exec("print", 0, "hello world")
	i, err := L.Exec("calculate", 1, 2, 3)
	if err != nil {
		p(err)
		return
	}
	p(i)

}

func basicTest(L *lua.State) {
	// push "print" function on the stack
	L.GetField(lua.LUA_GLOBALSINDEX, "print")
	// push the string "Hello World!" on the stack
	L.PushString("Hello World!")
	L.DumpStack()

	// call print with one argument, expecting no results
	L.Call(1, 0)
	L.DumpStack()

	L.GetField(lua.LUA_GLOBALSINDEX, "calculate")
	L.PushNumber(1.1)
	L.PushInteger(2)
	L.DumpStack()
	L.Call(2, 1)
	L.DumpStack()
	result := L.ToNumber(-1)
	L.DumpStack()
	p(result)

	L.SetTop(0)
	L.GetField(lua.LUA_GLOBALSINDEX, "calculate")
	L.PushNumber(1.1)
	L.PushNumber(2.2)
	L.Call(2, 1)
	result = L.ToNumber(-1)
	p(result)

	L.SetTop(0)
	L.GetField(lua.LUA_GLOBALSINDEX, "calculate")
	L.PushNumber(1.1)
	L.PushNumber(3)
	L.Call(2, 1)
	result = L.ToNumber(-1)
	p(result)

}

var count = 1

func execTest(number int) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	err := L.DoFile("test.lua")
	if err != nil {
		p(err)
		return
	}

	L.Exec("print", 0, number)
	i, err := L.Exec("calculate", 1, 2, 3)
	if err != nil {
		p(err)
		return
	}
	p(i)
	count--

}

func main() {
	runtime.GOMAXPROCS(100)
	// execTest(1)
	// execTest(2)
	// execTest(3)
	// execTest(4)
	// execTest(5)
	// execTest(6)

	go execTest(1)
	// go execTest(2)
	// go execTest(3)
	// go execTest(4)
	// go execTest(5)
	// go execTest(6)
	// go execTest(7)
	// go execTest(8)
	// go execTest(9)
	// go execTest(10)
	// go execTest(11)
	// go execTest(12)
	// go execTest(13)
	// go execTest(14)
	// go execTest(15)
	for count > 0 {
		fmt.Printf("count:%v\n", count)
		time.Sleep(3000 * time.Millisecond)

	}

}
