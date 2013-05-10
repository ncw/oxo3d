package main

import (
	//"fmt"
	"./oxo3d"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	o := oxo3d.NewOxo3d(true)
	//p := oxo3d.NewOxo3dHeuristic(o, 2)
	//p := oxo3d.NewOxo3dMinimax(o, 4)
	p := oxo3d.NewOxo3dAlphaBeta(o, 5)
	for !o.GameOver() {
		if o.IsMyGo() {
			Go := p.CalculateMyGo()
			o.Play(Go, true)
		} else {
			o.Print()
			o.ReadYourGo()
		}
	}
	o.Print()

}
