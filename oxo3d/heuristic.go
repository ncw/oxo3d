// Heuristic player for Oxo3d
package oxo3d

import (
	"math/rand"
)

// Global variables
var (
	xlinesx     [maxEncodedOxCount]int
	defaultCost [16]int
	nonel       = []int{
		0x1,   // level 0
		0x1,   // level 1
		0x300, // level 2 0x300
	}
	x2l = []int{
		0x10000, // level 0
		0x10000, // level 1
		0x40,    // level 2 // 0x40
	}
)

// Initialise the global variables
func init() {
	// This keeps score when counting lines
	for i := range xlinesx {
		xlinesx[i] = JUNK
	}
	xlinesx[0] = NONE
	xlinesx[OS*1] = O1
	xlinesx[OS*2] = O2
	xlinesx[OS*3] = O3
	xlinesx[OS*4] = O4
	xlinesx[XS*1] = X1
	xlinesx[XS*2] = X2
	xlinesx[XS*3] = X3
	xlinesx[XS*4] = X4

	// Scores: X has just gone
	defaultCost[NONE] = nonel[2]
	defaultCost[O1] = 0x10
	defaultCost[X1] = 0x100
	defaultCost[O2] = 0x1000
	defaultCost[X2] = x2l[2]
	defaultCost[O11] = 0x40
	defaultCost[X11] = 0x1000
	defaultCost[O12] = 0x10000
	defaultCost[X12] = 0x10000
	defaultCost[O22] = 0x100000
	defaultCost[X22] = 0x400000
	defaultCost[O3] = 0x1000000
	defaultCost[X3] = 0x10000000
	defaultCost[O4] = 0
	defaultCost[X4] = 0
}

// Oxo3d Heuristic player
type Oxo3dHeuristic struct {
	sfork [64]int
	cost  [16]int // initialised from defaultCost
	o     *Oxo3d
	level int
}

// Initialise the player
func NewOxo3dHeuristic(o *Oxo3d, level int) Player {
	p := &Oxo3dHeuristic{
		level: level,
		o:     o,
		cost:  defaultCost,
	}
	// set level by adjusting cost
	p.cost[X2] = x2l[p.level]
	p.cost[NONE] = nonel[p.level]
	return p
}

// Evaluate the cross lines potential of a given blank square
func (p *Oxo3dHeuristic) evalmovex(Go int) int {
	var xlines [16]int
	for _, a := range moveb[Go] {
		b := xlinesx[p.o.lines[a]]
		xlines[b] += 1
	}
	return p.cost[X11]*max(xlines[X1]-1, 0) +
		p.cost[X12]*min(xlines[X2], xlines[X1])
}

// Fill up the fork array
func (p *Oxo3dHeuristic) calcfork() {
	for i := range p.sfork {
		p.sfork[i] = 0
	}
	for j := range p.o.board {
		if p.o.board[j] == 0 {
			p.sfork[j] = p.evalmovex(j)
		}
	}
}

// Investigate what would be the consequences of playing at Go with piece type
func (p *Oxo3dHeuristic) evalmove(Go int, who int) int {
	var xlines [16]int
	s := 0
	encodedWho := XS
	if who != ME {
		encodedWho = OS
	}
	if p.level == 0 && rand.Intn(2) == 0 {
		return s
	}
	for _, a := range moveb[Go] {
		ox := p.o.lines[a]
		b := xlinesx[ox]
		s += p.cost[b]
		xlines[b] += 1
		if ox == encodedWho && p.level >= 2 {
			for _, bb := range movet[a] {
				if Go != bb {
					s += p.sfork[bb]
				}
			}
		}
	}
	if p.level >= 2 {
		s -= p.sfork[Go]
		s += p.cost[X22] * max(xlines[X2]-1, 0)
		s += p.cost[O22] * max(xlines[O2]-1, 0)
		s += p.cost[O11] * max(xlines[O1]-1, 0)
		s += p.cost[O12] * min(xlines[O1], xlines[O2])
	}
	return s
}

// Calculate the computer move
func (p *Oxo3dHeuristic) CalculateMyGo() int {
	if p.level >= 2 {
		p.calcfork()
	}
	myPlay := []int{}
	maxs := -0x7FFFFFFF
	for a := range p.o.board {
		if p.o.board[a] == 0 {
			s := p.evalmove(a, 1)
			if s > maxs {
				maxs = s
				myPlay = []int{a}
			} else if s == maxs {
				myPlay = append(myPlay, a)
			}
		}
	}
	if len(myPlay) == 0 {
		panic("Should have found some moves")
	}
	return myPlay[rand.Intn(len(myPlay))]
}
