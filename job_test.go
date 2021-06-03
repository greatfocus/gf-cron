package gfcron

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func myFunc() {
	fmt.Println("Helo, world")
}

func myFunc2(s string, n int) {
	fmt.Println("We have params here, string", s, "and number", n)
}

type MyTypeInterface struct {
	ID   int
	Name string
}

func (m MyTypeInterface) Bar() string {
	return "OK"
}

type MyTypeNoInterface struct {
	ID   int
	Name string
}

func myFuncStruct(m MyTypeInterface) {
	fmt.Println("Custom type as param")
}

func myFuncInterface(i Foo) {
	i.Bar()
}

type Foo interface {
	Bar() string
}

func TestJobError(t *testing.T) {

	c := New()

	if err := c.AddJob("* * * * *", myFunc, 10); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := c.AddJob("* * * * *", nil); err == nil {
		t.Error("This AddJob should return Error, fn is nil")
	}

	var x int
	if err := c.AddJob("* * * * *", x); err == nil {
		t.Error("This AddJob should return Error, fn is not func kind")
	}

	if err := c.AddJob("* * * * *", myFunc2, "s", 10, 12); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := c.AddJob("* * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, args are not the correct type")
	}

	if err := c.AddJob("* * * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, syntax error")
	}

	// custom types and interfaces as function params
	var m MyTypeInterface
	if err := c.AddJob("* * * * *", myFuncStruct, m); err != nil {
		t.Error(err)
	}

	if err := c.AddJob("* * * * *", myFuncInterface, m); err != nil {
		t.Error(err)
	}

	var mwo MyTypeNoInterface
	if err := c.AddJob("* * * * *", myFuncInterface, mwo); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	c.Shutdown()
}

var testN int
var testS string

func TestCron(t *testing.T) {
	testN = 0
	testS = ""

	c := Fake(2) // fake cron wiht 2sec timer to speed up test

	var wg sync.WaitGroup
	wg.Add(2)

	if err := c.AddJob("* * * * *", func() { testN++; wg.Done() }); err != nil {
		t.Fatal(err)
	}

	if err := c.AddJob("* * * * *", func(s string) { testS = s; wg.Done() }, "param"); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}

	if testN != 1 {
		t.Error("func 1 not executed as scheduled")
	}

	if testS != "param" {
		t.Error("func 2 not executed as scheduled")
	}
	c.Shutdown()
}

func TestRunAll(t *testing.T) {
	testN = 0
	testS = ""

	c := New()

	if err := c.AddJob("* * * * *", func() { testN++ }); err != nil {
		t.Fatal(err)
	}

	if err := c.AddJob("* * * * *", func(s string) { testS = s }, "param"); err != nil {
		t.Fatal(err)
	}

	c.RunAll()
	time.Sleep(time.Second)

	if testN != 1 {
		t.Error("func not executed on RunAll()")
	}

	if testS != "param" {
		t.Error("func not executed on RunAll() or arg not passed")
	}

	c.Clear()
	c.RunAll()

	if testN != 1 {
		t.Error("Jobs not cleared")
	}

	if testS != "param" {
		t.Error("Jobs not cleared")
	}

	c.Shutdown()
}
