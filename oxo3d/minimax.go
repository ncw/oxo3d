// Minimax player for Oxo3d
package oxo3d

import (
	"math/rand"
)

// Global variables
var (
	cost     [maxEncodedOxCount]int
	yourCost [maxEncodedOxCount]int
)

// Initialise the global variables
func init() {
	// Register this player
	Players["Minimax0"] = func(o *Oxo3d) Player { return NewOxo3dMinimax(o, 0) }
	Players["Minimax1"] = func(o *Oxo3d) Player { return NewOxo3dMinimax(o, 1) }
	Players["Minimax2"] = func(o *Oxo3d) Player { return NewOxo3dMinimax(o, 2) }
	Players["Minimax3"] = func(o *Oxo3d) Player { return NewOxo3dMinimax(o, 3) }
	Players["Minimax4"] = func(o *Oxo3d) Player { return NewOxo3dMinimax(o, 4) }

	// Evaluation of lines.  We are X and these count the score of the encoded O & X count

	// My X Go
	cost[OS*4] = -0x10000000
	cost[OS*3] = -0x100000
	cost[OS*2] = -0x1000
	cost[OS*1] = -0x1
	cost[XS*1] = 0x10 // prefer to make new lines over blocking old ones
	cost[XS*2] = 0x1000
	cost[XS*3] = 0x100000
	cost[XS*4] = 0x10000000

	// Your O Go
	yourCost[XS*4] = -cost[OS*4]
	yourCost[XS*3] = -cost[OS*3]
	yourCost[XS*2] = -cost[OS*2]
	yourCost[XS*1] = -cost[OS*1]
	yourCost[OS*1] = -cost[XS*1]
	yourCost[OS*2] = -cost[XS*2]
	yourCost[OS*3] = -cost[XS*3]
	yourCost[OS*4] = -cost[XS*4]
}

type Oxo3dMinimax struct {
	o          *Oxo3d
	bestMoves  []int
	level      int
	evaluation int
}

// Initialise the player
func NewOxo3dMinimax(o *Oxo3d, level int) Player {
	p := &Oxo3dMinimax{
		level:     level,
		o:         o,
		bestMoves: []int{},
	}
	return p
}

// Finds the best move for the current board position.  Returns
// the evaluation of that move.  Sets the bestMoves array to be a
// list of possible moves with that evaluation
func (p *Oxo3dMinimax) findBestMove(myGo bool, level int) int {
	newMoves := []int{}
	bestScore := -INF
	if !myGo {
		bestScore = INF
	}
	for Go := range p.o.board {
		if p.o.board[Go] == 0 {
			old_evaluation := p.evaluation
			p.o.Play(Go, myGo)
			p.updateEvaluation(Go, myGo)
			score := 0
			if level <= 1 || p.o.WhoWon() != NOONE {
				score = p.evaluation
				// Randomly overlook squares on level 0
				if level == 0 && rand.Intn(2) == 0 {
					score = 0
				}
			} else {
				score = p.findBestMove(!myGo, level-1)
			}
			p.o.UnPlay()
			p.evaluation = old_evaluation
			comparison := true
			if myGo {
				comparison = (score > bestScore)
			} else {
				comparison = (score < bestScore)
			}
			if comparison {
				bestScore = score
				newMoves = []int{Go}
			} else if score == bestScore {
				newMoves = append(newMoves, Go)
			}
			// if score >= 0 {
			// 	fmt.Printf("..Level %d, score  0x%08X at %d\n", level, score, Go)
			// } else {
			// 	fmt.Printf("..Level %d, score -0x%08X at %d\n", level, -score, Go)
			// }

		}
	}
	if len(newMoves) == 0 {
		panic("Should have found some moves")
	}
	p.bestMoves = newMoves
	return bestScore
}

// update evaluation after a go by who
func (p *Oxo3dMinimax) updateEvaluation(Go int, myGo bool) {
	encodedWho := XS
	pcost := &cost
	if !myGo {
		encodedWho = OS
		pcost = &yourCost
	}
	for _, a := range moveb[Go] {
		ox := p.o.lines[a]
		p.evaluation += pcost[ox] - pcost[ox-encodedWho]
	}
	/*
		// Check the evaluation shortcut is working
		evaluation := 0
		for _, ox := range p.o.lines {
			evaluation += pcost[ox]
		}
		if evaluation != p.evaluation {
			fmt.Printf("evaluation %d should be %d\n", p.evaluation, evaluation)
		}
	*/
}

// Do the computer move
func (p *Oxo3dMinimax) CalculateMyGo() int {
	p.findBestMove(true, p.level)
	return p.bestMoves[rand.Intn(len(p.bestMoves))]
}

// Do a move
func (p *Oxo3dMinimax) Play(Go int, myGo bool) {
	p.o.Play(Go, myGo)
	p.updateEvaluation(Go, myGo)
}
