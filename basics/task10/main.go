package main

import (
	"fmt"
	"math"
)

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Shape interface {
	Area() float64
}

func main() {
	//var circleOne []Circle
	var one []Shape

	one = append(one, Circle{Radius: 1})
	one = append(one, Rectangle{Width: 2, Height: 3})

	for _, y := range one {
		fmt.Println(y.Area())
	}

	//circleOne = append(circleOne, Circle{Radius: 1})
	//for _, x := range circleOne {
	//	fmt.Println(x.Area())
	//	}
	//var rectangleOne []Rectangle
	//rectangleOne = append(rectangleOne, Rectangle{Width: 2, Height: 3})

	//for _, y := range rectangleOne {
	//	fmt.Println(y.Area())
	//}

}

func (a Circle) Area() float64 {
	n := math.Pi * a.Radius * a.Radius
	return n

}

func (a Rectangle) Area() float64 {
	n := a.Width * a.Height
	return n

}
