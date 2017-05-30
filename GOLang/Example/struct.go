package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	male = iota
	female
)

type Person struct {
	Name   string
	Age    uint8
	Gender int
}

func (p *Person) Work(content string) {
	fmt.Printf("%s %s\n", p.Name, content)
}

type Student struct {
	Person
	School string
}

type Tool interface {
	Use()
}

type Car struct {
}

func (c *Car) Use() {
	fmt.Println("Drive Drive Run")
}

type Gun struct {
}

func (g *Gun) Use() {
	fmt.Println("Shoot Shoot Win")
}

type Gamer struct {
	tool Tool
}

func (g *Gamer) Get(t Tool) {
	g.tool = t
}

func (g *Gamer) UseTool() {
	if g.tool == nil {
		fmt.Println("No tool")
		return
	}
	g.tool.Use()
}

type FLVHeader struct {
	Signature  [3]byte
	Version    uint8
	Flags      uint8
	DataOffset uint32
}

func main() {
	alex := new(Person)
	alex.Name = "Alex"
	alex.Age = 33
	alex.Gender = male

	lilei := Person{"LiLei", 20, male}

	hanmeimei := &Person{"HanMeimei", 19, female}

	fmt.Printf("Use fmt.Printf:\n%v\n%v\n%v\n", alex, lilei, hanmeimei)

	alex.Work("Just do IT")

	xiaoming := Student{Person{"Xiaoming", 17, male}, "First School"}
	fmt.Println(xiaoming, xiaoming.Name, xiaoming.School)
	xiaoming.Work("Study")

	car := new(Car)
	gun := new(Gun)
	gamer := new(Gamer)

	gamer.UseTool()
	gamer.Get(car)
	gamer.UseTool()
	gamer.Get(gun)
	gamer.UseTool()

	f, _ := os.Open("test.flv")
	flvHeader := new(FLVHeader)
	binary.Read(f, binary.BigEndian, flvHeader)
	fmt.Println(flvHeader)
}
