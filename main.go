package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"time"
)

var listeners []*lua.LFunction

func callListeners(state *lua.LState) {
	event := state.NewTable()
	event.RawSetString("trigger", lua.LString("manual"))
	event.RawSetString("timestamp", lua.LString(time.Now().UTC().String()))

	for _, listener := range listeners {
		err := state.CallByParam(lua.P{
			Fn:      listener,
			NRet:    0,
			Protect: true,
		}, event)

		if err != nil {
			panic(err)
		}
	}
}

func main() {
	state := lua.NewState()
	defer state.Close()

	state.SetGlobal("debug", state.NewFunction(func(l *lua.LState) int {
		format := l.ToString(1)
		var args []any

		for i := 2; i <= l.GetTop(); i++ {
			arg := l.CheckAny(i)
			args = append(args, arg)
		}

		fmt.Printf(fmt.Sprintf("%s\n", format), args...)
		return 0
	}))

	eventsModule := state.NewTable()
	eventsModule.RawSetString("add_listener", state.NewFunction(func(l *lua.LState) int {
		handler := l.ToFunction(1)
		if handler == nil {
			return 1
		}

		listeners = append(listeners, handler)

		return 0
	}))
	state.SetGlobal("events", eventsModule)

	if err := state.DoFile("main.lua"); err != nil {
		panic(err)
	}

	callListeners(state)
}
