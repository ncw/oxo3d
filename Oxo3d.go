package main

// FIXME make it so can serialize and unserialize the game for save/load

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// We keep count of how many X and O there are on each line by encoding them like this
	// numberOfXs * 5 + numberOfOs
	// Since Xs and Os are both 0..4 the max value here is 24
	XS                = 5
	OS                = 1
	maxEncodedOxCount = XS * XS

	// players
	NOONE = 0
	ME    = 1
	YOU   = -1
	DRAW  = 2

	// names for numbers of lines
	NONE = 0
	O1   = 1
	O2   = 2
	O3   = 3
	O4   = 4
	X1   = 5
	X2   = 6
	X3   = 7
	X4   = 8
	O11  = 9
	O12  = 10
	O22  = 11
	X11  = 12
	X12  = 13
	X22  = 14
	JUNK = 15

	// A score larger than any possible score
	INF = 0x7FFFFFFF
)

// Global variables
var (
	// movet holds the positions of the 76 possible lines
	movet [76][4]int

	// for each space on the board, moveb records the number of all the
	// lines which go through that square (indexed into movet)
	moveb [64][]int

	// Globals for heuristic player
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

	// Globals for MinMax player
	cost [maxEncodedOxCount]int
)

// Initialise the global variables
func init() {
	// movet holds the positions of the 76 possible lines
	for a := 0; a < 16; a++ {
		for b := 0; b < 4; b++ {
			movet[a][b] = a*4 + b
			movet[a+16][b] = a + 4*b + 3*(a&0xC)
			movet[a+32][b] = a + 16*b
		}
	}
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			movet[a+48][b] = 16*a + 5*b
			movet[a+52][b] = 16*a + 3*b + 3
			movet[a+56][b] = a + 20*b
			movet[a+60][b] = a + 12*b + 12
			movet[a+64][b] = 4*a + 17*b
			movet[a+68][b] = 4*a + 15*b + 3
		}
	}
	for a := 0; a < 4; a++ {
		movet[72][a] = 21 * a
		movet[73][a] = 19*a + 3
		movet[74][a] = 13*a + 12
		movet[75][a] = 11*a + 15
	}

	// for each space on the board, moveb records the index of all the
	// lines which go through that square (indexed into movet)
	for line := range movet {
		for b := 0; b < 4; b++ {
			pos := movet[line][b]
			moveb[pos] = append(moveb[pos], line)
		}
	}

	// One time only initialisation for heuristic player
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

	// One time only initialisation for heuristic player
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

	// Setup for MinMax player
	// Evaluation of lines.  We are X and these count the score of the encoded O & X count
	// Player    MinMax1 Score   -26
	// Player Heuristic1 Score    -3
	// Player Heuristic2 Score    -2
	// Player    MinMax2 Score     3
	// Player    MinMax3 Score    28
	cost[OS*4] = -0x10000000
	cost[OS*3] = -0x100000
	cost[OS*2] = -0x1000
	cost[OS*1] = -0x1
	cost[XS*1] = 0x10
	cost[XS*2] = 0x1000
	cost[XS*3] = 0x100000
	cost[XS*4] = 0x10000000

}

type Oxo3d struct {
	plays   []int
	board   [64]int
	lines   [76]int
	lastGo  int
	winline int
	moves   int
	first   bool // who went first
	isMyGo  bool
	winner  int
}

// Setup all the things which need to be done for every game.
//  Note that the state of first is flipped for the next game

func (o *Oxo3d) newGame(first bool) {
	o.plays = make([]int, 0)
	for i := range o.board {
		o.board[i] = 0
	}
	for i := range o.lines {
		o.lines[i] = 0
	}
	o.lastGo = -1
	o.winline = -1
	o.moves = 0
	o.isMyGo = first
	o.first = !first
	o.winner = NOONE
}

// Create a new Oxo3d instance
func NewOxo3d(first bool) *Oxo3d {
	o := &Oxo3d{}
	o.newGame(first)
	return o
}

// Construct a new instance with the game in the same state as the instance
//passed in
/*
FIXME
func (o *Oxo3d) CopyOxo3d(that *Oxo3d, level int) {
    	this(!that.first, level)
        // Replay the game
        for (int go : that.plays) {
        	play(go, isMyGo)
        }
        if (isMyGo != that.isMyGo) {
            panic("Game replay went wrong: isMyGo")
        }
        if (moves != that.moves) {
            panic("Game replay went wrong: moves")
        }
        if (!Arrays.equals(o.board, that.board)) {
            panic("Game replay went wrong: board")
        }
        if (!plays.equals(that.plays)) {
            panic("Game replay went wrong: plays")
        }
        if (winline != that.winline) {
            panic("Game replay went wrong: winline")
        }
        if (!Arrays.equals(getBoard(), that.getBoard())) {
            panic("Game replay went wrong: getBoard")
        }
        if (!Arrays.equals(marks(), that.marks())) {
            panic("Game replay went wrong: marks")
        }
}
*/

// Returns NOONE, ME, YOU or DRAW for who won
func (o *Oxo3d) whoWon() int {
	if o.winner != NOONE {
		return o.winner
	}
	if o.moves == 64 {
		return DRAW
	}
	return NOONE
}

// Print the board
func (o *Oxo3d) prn() {
	win := o.whoWon()
	if win == ME {
		fmt.Println("I win")
	} else if win == YOU {
		fmt.Println("You win")
	} else if win == DRAW {
		fmt.Println("A draw")
	}
	mark := o.marks()
	symbols := [3]string{"O", ".", "X"}
	for a := 0; a < 64; a += 4 {
		if (a & 0xC) == 0 {
			fmt.Println("+-------------+-------------+")
		}
		fmt.Print("|")
		for b := a; b < a+4; b++ {
			fmt.Print(" ")
			fmt.Print(symbols[o.board[b]+1])
			if mark[b] {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print(" |")
		for b := a; b < a+4; b++ {
			fmt.Printf(" %2d", b)
		}
		fmt.Println(" |")
	}
}

// return an array telling whether each board position should be marked
func (o *Oxo3d) marks() [64]bool {
	var mark [64]bool
	if o.lastGo >= 0 {
		mark[o.lastGo] = true
	}
	if o.winline >= 0 {
		for _, a := range movet[o.winline] {
			mark[a] = true
		}
	}
	return mark
}

// Test to see if someone won or it was a draw
func (o *Oxo3d) gameOver() bool {
	return o.whoWon() != NOONE
}

// Override in base classes if desired to be called after play or unplay has executed
// func (o *Oxo3d) updateEvaluation(Go int, who int, encodedWho int, unplay bool) {
// }

// Keep internal state up to date when a piece is played
func (o *Oxo3d) play(Go int, myGo bool) {
	if Go < 0 || Go >= 64 {
		panic("Illegal move: out of range")
	}
	if o.board[Go] != NOONE {
		panic("Illegal move: board position occupied")
	}
	who := ME
	encodedWho := XS
	if !myGo {
		who = YOU
		encodedWho = OS
	}
	fourInARow := 4 * encodedWho
	o.lastGo = Go
	o.plays = append(o.plays, Go)
	o.moves += 1
	o.board[Go] = who
	o.winner = NOONE
	o.winline = -1
	for _, a := range moveb[Go] {
		ox := o.lines[a] + encodedWho
		o.lines[a] = ox
		if ox == fourInARow {
			o.winner = who
			o.winline = a
		}
	}
	if myGo != o.isMyGo {
		panic("Not my go!")
	}
	o.isMyGo = !o.isMyGo
	//isMyGo = false; // FIXME
	//o.updateEvaluation(Go, who, encodedWho, false)
}

// Unplay the last move keeping internal state up to date
func (o *Oxo3d) unplay() {
	// FIXME some assertions?
	if len(o.plays) <= 0 {
		panic("No moves to unplay")
	}
	Go := o.plays[len(o.plays)-1]
	o.plays = o.plays[:len(o.plays)-1]
	who := o.board[Go]
	if who == NOONE {
		panic("Unplaying no one")
	}
	encodedWho := XS
	if who != ME {
		encodedWho = OS
	}
	if len(o.plays) <= 0 {
		o.lastGo = -1
	} else {
		o.lastGo = o.plays[len(o.plays)-1]
	}
	o.moves -= 1
	o.board[Go] = NOONE
	o.winner = NOONE
	o.winline = -1
	for _, a := range moveb[Go] {
		o.lines[a] -= encodedWho
	}
	o.isMyGo = !o.isMyGo
	// o.updateEvaluation(Go, who, encodedWho, true)
}

// Calculate and do the computer move
// func (o *Oxo3d) myGo() int {
//     	// Make copies of marks and board
//     	o.thinking = true
//         Go := o.calculateMyGo()
//         o.play(Go, true)
//     	o.thinking = false
//         return Go
// }

// Computer move
//func (o *Oxo3d) calculateMyGo() int {
//	return 0
//}

// Do the human move
func (o *Oxo3d) readYourGo() int {
	var Go int
	for {
		fmt.Print("Go (0..63) ")
		n, err := fmt.Scanf("%d", &Go)
		if n != 1 || err != nil {
			continue
		}
		if Go >= 0 && Go <= 63 && o.board[Go] == 0 {
			break
		}
	}
	fmt.Printf("Going at %d\n", Go)
	o.yourGo(Go)
	return Go
}

// Do the opponents move
func (o *Oxo3d) yourGo(Go int) {
	o.play(Go, false)
}

// Oxo3d Heuristic player
type Oxo3dHeuristic struct {
	sfork [64]int
	cost  [16]int // initialised from defaultCost
	o     *Oxo3d
	level int
}

// Initialise the player
func NewOxo3dHeuristic(o *Oxo3d, level int) *Oxo3dHeuristic {
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
func (p *Oxo3dHeuristic) calculateMyGo() int {
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

type Oxo3dMinMax struct {
	o          *Oxo3d
	bestMoves  []int
	level      int
	evaluation int
}

// Initialise the player
func NewOxo3dMinMax(o *Oxo3d, level int) *Oxo3dMinMax {
	p := &Oxo3dMinMax{
		level:     level,
		o:         o,
		bestMoves: []int{},
	}
	return p
}

// Finds the best move for the current board position.  Returns
// the evaluation of that move.  Sets the bestMoves array to be a
// list of possible moves with that evaluation
func (p *Oxo3dMinMax) findBestMove(myGo bool, level int) int {
	newMoves := []int{}
	bestScore := -INF
	if !myGo {
		bestScore = INF
	}
	for Go := range p.o.board {
		if p.o.board[Go] == 0 {
			old_eval := p.evaluation
			p.o.play(Go, myGo)
			p.updateEvaluation(Go, myGo, false)
			score := 0
			if level <= 1 || p.o.whoWon() != NOONE {
				score = p.evaluation
				// Randomly overlook squares on level 0
				if level == 0 && rand.Intn(2) == 0 {
					score = 0
				}
			} else {
				score = p.findBestMove(!myGo, level-1)
			}
			p.o.unplay()
			p.updateEvaluation(Go, myGo, true)
			if p.evaluation != old_eval {
				panic("Evaluation is wrong")
			}
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
func (p *Oxo3dMinMax) updateEvaluation(Go int, myGo bool, unplay bool) {
	encodedWho := XS
	if !myGo {
		encodedWho = OS
	}
	if unplay {
		encodedWho = -encodedWho
	}
	for _, a := range moveb[Go] {
		ox := p.o.lines[a]
		p.evaluation += cost[ox] - cost[ox-encodedWho]
	}
}

// Do the computer move
func (p *Oxo3dMinMax) calculateMyGo() int {
	p.findBestMove(true, p.level)
	return p.bestMoves[rand.Intn(len(p.bestMoves))]
}

type Oxo3dAlphaBeta struct {
	Oxo3dMinMax
}

// Initialise the player
func NewOxo3dAlphaBeta(o *Oxo3d, level int) *Oxo3dAlphaBeta {
	p := &Oxo3dAlphaBeta{
		Oxo3dMinMax: *NewOxo3dMinMax(o, level),
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
			old_eval := p.evaluation
			p.o.play(Go, myGo)
			p.updateEvaluation(Go, myGo, false)
			score := 0
			if level <= 1 || p.o.whoWon() != NOONE {
				score = p.evaluation
				// Randomly overlook squares on level 0
				if level == 0 && rand.Intn(2) == 0 {
					score = 0
				}
			} else {
				score = p.findBestMove(!myGo, level-1, alpha, beta)
			}
			p.o.unplay()
			p.updateEvaluation(Go, myGo, true)
			if p.evaluation != old_eval {
				panic("Evaluation is wrong")
			}
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
func (p *Oxo3dAlphaBeta) calculateMyGo() int {
	p.findBestMove(true, p.level, -INF, INF)
	return p.bestMoves[rand.Intn(len(p.bestMoves))]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// FIXME updateEvaluation not being called properly
	o := NewOxo3d(true)
	//p := NewOxo3dHeuristic(o, 2)
	//p := NewOxo3dMinMax(o, 4)
	p := NewOxo3dAlphaBeta(o, 5)
	for !o.gameOver() {
		if o.isMyGo {
			Go := p.calculateMyGo()
			o.play(Go, true)
		} else {
			o.prn()
			o.readYourGo()
		}
	}
	o.prn()

}
