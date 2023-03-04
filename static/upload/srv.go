package demo

import "fmt"

type Srv struct{}

func (s *Srv) Ser1() string {
	if !s.ser1() {
		return "Ser1 s.ser1"
	}
	a := 0
	if s.ser1() {
		a++
	}
	fmt.Println(a)
	if s.ser2() {
		return "Ser1 s.ser2()"
	} else if s.ser3() {
		return "Ser1 s.ser3()"
	} else {
		return "Ser1 "
	}
}
func (s *Srv) Ser2() string {
	if !s.ser1() {
		if s.ser3() {
			return " Ser2 1"
		}
		return " Ser2 2"
	}
	if s.ser2() && s.ser3() {
		return " Ser2 3"
	}

	return " Ser2 3"
}

func (s *Srv) ser1() bool {
	return false
}

func (s *Srv) ser2() bool {
	return true
}
func (s *Srv) ser3() bool {
	return true
}
