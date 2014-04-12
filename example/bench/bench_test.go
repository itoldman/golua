package test

import (
	"../../lua"
	"fmt"
	"testing"
)

func p(p interface{}) {
	fmt.Printf("%v\n", p)

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

func BenchmarkExecLua(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	err := L.DoFile("../test.lua")
	if err != nil {
		p(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.Exec("calculate", 1, 2, 3)
	}
}

func calculate(a int, b int) int {
	return a + b
}

func BenchmarkExecGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculate(2, 3)
	}
}
