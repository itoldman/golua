// This package provides access to the excellent lua language interpreter from go code.
//
// Access to most of the functions in lua.h and lauxlib.h is provided as well as additional convenience functions to publish Go objects and functions to lua code.
//
// The documentation of this package is no substitute for the official lua documentation and in many instances methods are described only with the name of their C equivalent
package lua

/*
#cgo CFLAGS: -Ilua
#cgo llua LDFLAGS: -llua
#cgo luaa LDFLAGS: -llua -lm -ldl
#cgo linux,!llua,!luaa LDFLAGS: -llua5.1
#cgo darwin,!luaa LDFLAGS: -llua
#cgo freebsd,!luaa LDFLAGS: -llua

#include <lua.h>
#include <stdlib.h>

#include "golua.h"

*/
import "C"
import (
	"errors"
	"reflect"
)

func (L *State) DumpElement(idx int, depth int) {
	C.dump_element(L.s, C.int(idx), C.int(depth))
}

func (L *State) DumpElementRecurse(idx int, depth int) {
	C.dump_element_recurse(L.s, C.int(idx), C.int(depth))
}

func (L *State) DumpStack() {
	C.dump_stack(L.s)
}

func (L *State) DumpVar(var1 string) {
	C.dump_var(L.s, C.CString(var1))
}

func (L *State) DumpVarRecurse(var1 string) {
	C.dump_var_recurse(L.s, C.CString(var1))
}

func (L *State) GetInteger(key string) int {
	L.GetGlobal(key)
	value := L.ToInteger(-1)
	L.Pop(1)
	return value
}

func (L *State) GetNumber(key string) float64 {
	L.GetGlobal(key)
	value := L.ToNumber(-1)
	L.Pop(1)
	return value
}

func (L *State) GetString(key string) string {
	L.GetGlobal(key)
	value := L.ToString(-1)
	L.Pop(1)
	return value
}

func (L *State) GetBoolean(key string) bool {
	L.GetGlobal(key)
	value := L.ToBoolean(-1)
	L.Pop(1)
	return value
}

func (L *State) Exec(f string, nresults int, args ...interface{}) ([]interface{}, error) {
	n := L.GetTop()
	L.GetGlobal(f)

	nargs := len(args)
	if !L.CheckStack(nargs) {
		return nil, errors.New("no enough memory")
	}

	for _, v := range args {
		L.PushVal(v)
	}

	err := L.Call(nargs, nresults)
	if err != nil {
		return nil, err
	}

	retCnt := L.GetTop() - n
	ret := L.GetRetValue(retCnt)
	L.SetTop(n)
	return ret, nil
}

func (L *State) GetRetValue(nargs int) []interface{} {
	ret := make([]interface{}, nargs)

	for i, index := 1, 0; i <= nargs; i++ {
		ret[index] = L.Val(i)
		index++
	}
	return ret
}

func (L *State) ParseTable(index int, retMap map[interface{}]interface{}) {
	L.PushValue(index)
	L.PushNil()

	keyIndex := -2
	valueIndex := -1

	for L.Next(keyIndex) != 0 {
		var key interface{}

		kt := L.Type(keyIndex)
		switch kt {
		case C.LUA_TNUMBER:
			key = L.ToNumber(keyIndex)
		case C.LUA_TSTRING:
			key = L.ToString(keyIndex)
		default:
			return
		}

		retMap[key] = L.Val(valueIndex)
		L.SetTop(-2)
	}
	L.SetTop(-2)
}

func (L *State) Val(index int) interface{} {
	t := L.Type(index)
	switch t {
	case C.LUA_TNIL:
		return nil
	case C.LUA_TBOOLEAN:
		return L.ToBoolean(index)
	case C.LUA_TLIGHTUSERDATA:
		return nil
	case C.LUA_TNUMBER:
		return L.ToNumber(index)
	case C.LUA_TSTRING:
		return L.ToString(index)
	case C.LUA_TTABLE:
		retMap := make(map[interface{}]interface{})
		L.ParseTable(index, retMap)
		return retMap
	case C.LUA_TFUNCTION:
		return nil
	case C.LUA_TUSERDATA:
		return nil
	case C.LUA_TTHREAD:
		return nil
	}
	return nil
}

func (L *State) PushTable(v reflect.Value) {
	L.CreateTable(0, 0)
	for i := 0; i < v.Len(); i++ {
		L.PushVal(i + 1)
		L.PushVal(v.Index(i).Interface())
		L.SetTable(-3)
	}
}

func (L *State) PushMap(v reflect.Value) {
	L.CreateTable(0, 0)
	keys := v.MapKeys()
	for i := 0; i < len(keys); i++ {
		L.PushVal(keys[i].Interface())
		L.PushVal(v.MapIndex(keys[i]).Interface())
		L.SetTable(-3)
	}
}

func (L *State) PushVal(v interface{}) {
	switch v.(type) {
	case int:
		L.PushInteger(int64(v.(int)))
	case int64:
		L.PushInteger(v.(int64))
	case float64:
		L.PushNumber(v.(float64))
	case string:
		L.PushString(v.(string))
	case bool:
		L.PushBoolean(v.(bool))
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice:
			L.PushTable(rv)
		case reflect.Map:
			L.PushMap(rv)
		default:
			L.PushNil()
		}
	}
}
