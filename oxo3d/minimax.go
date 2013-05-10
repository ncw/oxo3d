// Minimax player for Oxo3d
package oxo3d

import (
	"math/rand"
)

// Global variables
var (
	cost [maxEncodedOxCount]int
)

// Initialise the global variables
func init() {
	// Evaluation of lines.  We are X and these count the score of the encoded O & X count
	// Player    Minimax1 Score   -26
	// Player Heuristic1 Score    -3
	// Player Heuristic2 Score    -2
	// Player    Minimax2 Score     3
	// Player    Minimax3 Score    28
	cost[OS*4] = -0x10000000
	cost[OS*3] = -0x100000
	cost[OS*2] = -0x1000
	cost[OS*1] = -0x1
	cost[XS*1] = 0x10
	cost[XS*2] = 0x1000
	cost[XS*3] = 0x100000
	cost[XS*4] = 0x10000000
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
	if !myGo {
		encodedWho = OS
	}
	for _, a := range moveb[Go] {
		ox := p.o.lines[a]
		p.evaluation += cost[ox] - cost[ox-encodedWho]
	}
}

// Do the computer move
func (p *Oxo3dMinimax) CalculateMyGo() int {
	p.findBestMove(true, p.level)
	return p.bestMoves[rand.Intn(len(p.bestMoves))]
}
