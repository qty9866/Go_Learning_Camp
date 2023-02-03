package main

import "fmt"

type Header map[string][]string

// USB 定义一个USB接口类型
type USB interface {
	read()
	write()
}

// 定义phone结构体
type phone struct {
	brand string
}

// 定义pad结构体
type pad struct {
	brand string
}

// 定义一个computer结构体
type computer struct{}

func (c *computer) working(u USB) {
	u.read()
	u.write()
	if p, ok := u.(*phone); ok {
		p.call()
	}
}
func (p *phone) read() {
	fmt.Println(p.brand, "start to read")
}
func (p *phone) write() {
	fmt.Println(p.brand, "start to write")
}

func (p *phone) call() {
	fmt.Println(p.brand, "start calling...")
}

func (p *pad) read() {
	fmt.Println(p.brand, "start to read")
}
func (p *pad) write() {
	fmt.Println(p.brand, "start to read")
}

func main() {
	var Arr [3]USB
	Arr[0] = &pad{"Huawei"}
	Arr[1] = &pad{"Oppo"}
	Arr[2] = &phone{"Apple"}
	var c computer
	for _, i := range Arr {
		c.working(i)
	}

}
