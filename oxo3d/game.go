package oxo3d

/*
Shouldn't really need state in the players

Should initialise from the board from flat each time

Which makes them functions probably with subfunctions

Note that the evaluation isn't being updated properly because of that - when the human plays it goes wrong

*/

// FIXME make it so can serialize and unserialize the game for save/load

import (
	"fmt"
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
	o.plays = make([]int, 0, 64)
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
        	Play(go, isMyGo)
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
func (o *Oxo3d) WhoWon() int {
	if o.winner != NOONE {
		return o.winner
	}
	if o.moves == 64 {
		return DRAW
	}
	return NOONE
}

// Returns true if it is my go
func (o *Oxo3d) IsMyGo() bool {
	return o.isMyGo
}

// Returns true if the move is to an empty square
// Returns false if occupied or invalid
func (o *Oxo3d) ValidMove(Go int) bool {
	if Go < 0 || Go >= 64 {
		return false
	}
	if o.board[Go] != NOONE {
		return false
	}
	return true
}

// Print the board
func (o *Oxo3d) Print() {
	win := o.WhoWon()
	if win == ME {
		fmt.Println("I win")
	} else if win == YOU {
		fmt.Println("You win")
	} else if win == DRAW {
		fmt.Println("A draw")
	}
	mark := o.Marks()
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
func (o *Oxo3d) Marks() [64]bool {
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
func (o *Oxo3d) GameOver() bool {
	return o.WhoWon() != NOONE
}

// Keep internal state up to date when a piece is played
func (o *Oxo3d) Play(Go int, myGo bool) {
	if !o.ValidMove(Go) {
		panic("Illegal move")
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
}

// Unplay the last move keeping internal state up to date
func (o *Oxo3d) UnPlay() {
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
}

// Computer move
//func (o *Oxo3d) CalculateMyGo() int {
//	return 0
//}

// Do the human move
func (o *Oxo3d) ReadYourGo() int {
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
	o.YourGo(Go)
	return Go
}

// Do the computer move
func (o *Oxo3d) MyGo(Go int) {
	o.Play(Go, true)
}

// Do the opponents move
func (o *Oxo3d) YourGo(Go int) {
	o.Play(Go, false)
}
