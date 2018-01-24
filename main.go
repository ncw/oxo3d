package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ncw/oxo3d/oxo3d"
)

// Do the human move
func readYourGo(o *oxo3d.Oxo3d) int {
	var Go int
	for {
		fmt.Print("Go (0..63) ")
		n, err := fmt.Scanf("%d", &Go)
		if n != 1 || err != nil {
			continue
		}
		if o.ValidMove(Go) {
			break
		}
	}
	fmt.Printf("Going at %d\n", Go)
	o.YourGo(Go)
	return Go
}

func printBoard(o *oxo3d.Oxo3d) {
	fmt.Print(o.String())
}

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
			printBoard(o)
			readYourGo(o)
		}
	}
	printBoard(o)
}
