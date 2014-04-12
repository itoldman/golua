package main

import (
	"../lua"
	"fmt"
	"testing"
	//"github.com/dongaihua/golua/lua"
	//"runtime"
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

func main() {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	err := L.DoFile("test.lua")
	if err != nil {
		p(err)
		return
	}

	// L.Exec("print", "Hello World!")
	// result, err := L.Exec("calculate", 1.1, 2)
	// result, err := L.Exec("calculate",  1.1, 2.2)
	// result, err := L.Exec("calucuate", 1, 2)

	//readConfig(L)
	exec(L)
}

func BenchmarkExec(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	err := L.DoFile("test.lua")
	if err != nil {
		p(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Exec("calculate", 1, 2, 3)
	}
}
