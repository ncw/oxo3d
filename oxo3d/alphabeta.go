// Minimax player with Alpha-beta pruning for Oxo3d
package oxo3d

import (
	"math/rand"
)

type Oxo3dAlphaBeta struct {
	Oxo3dMinimax
}

// Initialise the player
func NewOxo3dAlphaBeta(o *Oxo3d, level int) *Oxo3dAlphaBeta {
	p := &Oxo3dAlphaBeta{
		Oxo3dMinimax: *NewOxo3dMinimax(o, level),
	}
	return p
}

// Finds the best move for the current board position.  Returns
// the evaluation of that move.  Sets the bestMoves array to be a
// list of possible moves with that evaluation
// FIXME should order the search for Alpha/Beta
func (p *Oxo3dAlphaBeta) findBestMove(myGo bool, level int, alpha int, beta int) int {
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
				score = p.findBestMove(!myGo, level-1, alpha, beta)
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
				// Alpha/Beta pruning
				if myGo {
					alpha = bestScore
				} else {
					beta = bestScore
				}
				if alpha > beta {
					break
				}
			} else if score == bestScore {
				newMoves = append(newMoves, Go)
			}
		}
	}
	if len(newMoves) == 0 {
		panic("Should have found some moves")
	}
	p.bestMoves = newMoves
	return bestScore
}

// Do the computer move
func (p *Oxo3dAlphaBeta) CalculateMyGo() int {
	p.findBestMove(true, p.level, -INF, INF)
	return p.bestMoves[rand.Intn(len(p.bestMoves))]
}
