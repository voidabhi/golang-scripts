package main

import "fmt"

type Car struct {
	speed    int
	distance int
}

func NewCar() *Car {
	return &Car{0, 0}
}

func (self *Car) SetSpeed(speed int) {
	self.speed = speed
}

func (self *Car) Run(minutes int) {
	self.distance = minutes
}

func (self *Car) GetDistance() int {
	return self.distance
}

type Driver struct {
	car *Car
}

func NewDriver(car *Car) *Driver {
	return &Driver{car}
}

func (self *Driver) PushPedal(speed int) {
	self.car.SetSpeed(speed)
}

func (self *Driver) Drive(minutes int) {
	self.car.Run(minutes)
}

type DrivingSimulator struct {
}

func (self *DrivingSimulator) Simulate() {
	c := NewCar()
	d := NewDriver(c)
	d.PushPedal(700)
	d.Drive(30)
	d.PushPedal(750)
	d.Drive(20)

	fmt.Println("The travel distance is ", c.GetDistance(), " m.")
}

func main() {
	d := &DrivingSimulator{}
	d.Simulate()
}
