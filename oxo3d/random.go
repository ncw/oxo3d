// Random player for Oxo3d
package oxo3d

import (
	"math/rand"
)

func init () {
	// Register this player
	Players["Random"] = func(o *Oxo3d) Player { return NewOxo3dRandom(o) }
}

type Oxo3dRandom struct {
	o *Oxo3d
}

// Initialise the player
func NewOxo3dRandom(o *Oxo3d) Player {
	p := &Oxo3dRandom{
		o: o,
	}
	return p
}

// Do the computer move
func (p *Oxo3dRandom) CalculateMyGo() int {
	for {
		Go := rand.Intn(64)
		if p.o.ValidMove(Go) {
			return Go
		}
	}
	return -1
}

// Do a move
func (p *Oxo3dRandom) Play(Go int, myGo bool) {
	p.o.Play(Go, myGo)
}
