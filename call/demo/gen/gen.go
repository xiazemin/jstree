package demo

import . "call/demo/model"

type Srv struct{}

func (s *Srv) Ser1() string {
	Xzm.Count[0]++
	if !s.ser1() {
		Xzm.Count[2]++
		return "Ser1 s.ser1"
	}

	Xzm.Count[1]++
	if s.ser2() {
		Xzm.Count[3]++
		return
		Xzm.Count[4]++
		if s.ser3() {
			Xzm.Count[5]++
			return "Ser1 s.ser3()"
		} else {
			Xzm.Count[6]++
			{
				return "Ser1 "
			}
		}
	}
}
func (s *Srv) Ser2() string {
	Xzm.Count[7]++
	if !s.ser1() {
		Xzm.Count[10]++
		if s.ser3() {
			Xzm.Count[12]++
			return " Ser2 1"
		}
		Xzm.Count[11]++
		return " Ser2 2"
	}
	Xzm.Count[8]++
	if s.ser2() && s.ser3() {

		Xzm.Count[8]++
		if s.ser2() && s.ser3() {
			Xzm.Count[13]++
			return " Ser2 3"
		}

		Xzm.Count[9]++
		return " Ser2 3"
	}
}

func (s *Srv) ser1() bool {
	Xzm.Count[14]++
	return false
}

func (s *Srv) ser2() bool {
	Xzm.Count[15]++
	return true
}

func (s *Srv) ser3() bool {
	Xzm.Count[16]++
	return true
}
