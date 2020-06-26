package main

import (
	"encoding/json"
	"fmt"
	"math"
)

const X float64 = 750.0 // 850
const Y float64 = 270.0

var Ix int
var Iy int

type Row struct {
	*Line
	*Grape
	I int  `json:"row"`
	Length float64 `json:"length"`
	MaxVines int `json:"max-vines"`
}

func (l *Line) Length() float64 {
	xdiff := l.X2 - l.X1
	ydiff := l.Y2 - l.Y1

	return math.Sqrt(xdiff*xdiff + ydiff*ydiff)
}


type Line struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

func NewRow(i int) Row {
	ixoff := i-Ix
	iyoff := i-Iy
	l := &Line{
		X1: math.Max(0, 19.8*float64(ixoff)),
		Y1: math.Min(Y, 8.5*float64(i) + 34.2),
		X2: math.Min(X, 23.5*float64(i) + 94),
		Y2: math.Max(0, 8.5*float64(iyoff)),
	}
	length := l.Length()

	return Row{
		Line: l,
		I: i,
		Length: length,
		MaxVines: int(length)/4,
	}
}

type Grape struct {
	Varietal string `json:"grape"`
	Root string `json:"root"`
	Quantity int `json:"amount"`
}

func G(v, r string, q int) *Grape {
	return &Grape{
		Varietal: v,
		Root: r,
		Quantity: q,
	}
}

func (g *Grape) Plant(i int) *Grape {
	_q := math.Min(float64(i), float64(g.Quantity))
	q := int(_q)
	g.Quantity = g.Quantity - q
	return G(g.Varietal, g.Root, q)
}

var plots = []*Grape{
	G("Marquette", "SO4", 685),
	G("Marquette", "3309", 685),
	G("Pinot Noir", "SO4", 685),
	G("Pinot Noir", "3309", 685),
	G("Riesling", "SO4", 685),
	G("Riesling", "3309", 2055),
	G("Chardonnay", "SO4", 685),
	G("Chardonnay", "3309", 685),
}

func main() {
	_Ix := X / 23.5
	_Iy := Y / 8.5

	Ix = int(_Ix)
	Iy = int(_Iy)

	rows := make([]Row, 0, 200)
	plot := &Grape{}
	for i:=0; plot.Quantity > 0 || len(plots) > 0; i++ {
		for plot.Quantity == 0 && len(plots) > 0 {
			plot, plots = plots[0], plots[1:]
		}
		if plot.Quantity == 0 {
			break
		}

		r := NewRow(i)
		r.Grape = plot.Plant(r.MaxVines)
		rows = append(rows, r)
	}

	data, _ := json.MarshalIndent(rows, "", "    ")
	fmt.Printf("%s\n", string(data))
}