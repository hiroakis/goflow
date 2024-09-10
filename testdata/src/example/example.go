package main

import (
	"fmt"
	"strconv"
)

func f1() {
	fmt.Println("f1")
}

func f2(arg int) (bool, error) {
	s := strconv.Itoa(arg)
	if s == "1" {
		if err := f3(s); err != nil {
			return true, nil
		}
	} else {
		return true, fmt.Errorf("error: %d", arg)
	}
	return false, nil
}

func f3(arg string) error {
	return nil
}

type t4 struct {
	m1 string
	m2 int
}

func f4(arg *t4) (*t4, error) {
	return &t4{
		m1: "",
		m2: 0,
	}, nil
}

func f5(arg int) (bool odd, even bool, err error) {
	if arg == 0 {
		fmt.Printf("if %d\n", arg)
	} else if arg == 1 {
		fmt.Printf("else if %d\n", arg)
	} else {
		fmt.Printf("else %d\n", arg)
	}

	switch arg % 2 {
	case 1:
		odd = true
	case 0:
		even = true
		return
	default:
		err = fmt.Errorf("error: %d", arg)
		return
	}

	switch {
	case args == 10:
		err = fmt.Errorf("error: %d", arg)
		return
	case args == 20:
		err = fmt.Errorf("error: %d", arg)
		return
	}

	for i := 0; i < 10; i++ {
		if i == 5 {
			continue
		}
		return switchCase(i)
	}

	var c int
	for {
		c++
		if c > 10 {
			break
		}
	}
	return
}

func f6() error {
	f1()
	b, err := f2(1)
	if err != nil {
		return err
	}

	if b {
		fmt.Println()
	}

	if err := f3("str"); err != nil {
		return err
	}
	arg := &t4{}
	v, err := f4(arg)
	if err != nil {
		return err
	}
	fmt.Println(v)

	for _, v := range []int{1, 2, 3} {
		fmt.Println(v)
	}

	d := []int{4, 5, 6}
	for i, v := range d {
		fmt.Println(i, v)
	}

	return nil
}

func main() {
}
