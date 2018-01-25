package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/ncw/oxo3d/oxo3d"
)

// Globals
var (
	board    = oxo3d.NewOxo3d(true)
	player   = oxo3d.NewOxo3dHeuristic(board, 0)
	symbols  = [3]string{"O", "", "X"}
	square   [64]*js.Object
	document *js.Object
	messageP *js.Object
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Exit with the message
func fatalf(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	js.Global.Call("alert", text)
	panic(text)
}

// getElementById gets the element with the ID passed in or panics
func getElementById(ID string) *js.Object {
	obj := document.Call("getElementById", ID)
	if obj == js.Undefined {
		fatalf("couldn't find ID %q", ID)
	}
	return obj
}

func message(text string) {
	messageP.Set("innerHTML", text)
}

// Fills in the board with the state of the game
func draw() {
	for i := range square {
		square[i].Set("innerHTML", symbols[board.Board(i)+1])
		if board.Marked(i) {
			square[i].Call("setAttribute", "class", "marked")
		} else {
			square[i].Call("removeAttribute", "class")
		}
	}
	switch board.WhoWon() {
	case oxo3d.ME:
		message("I win")
	case oxo3d.YOU:
		message("You win")
	case oxo3d.DRAW:
		message("A draw")
	}
}

// Checks the state of the game after
func checkState() {
	// do computer move if necessary
	if !board.GameOver() && board.IsMyGo() {
		message("thinking...")
		Go := player.CalculateMyGo()
		board.MyGo(Go)
		message("")
		draw()
	}
}

// Called when a square is clicked
func clickSquare(Go int) bool {
	if board.GameOver() {
		return false
	}
	message("")
	if board.IsMyGo() {
		message("Not your go yet!")
		return false
	}
	println(fmt.Sprintf("click on %d", Go))
	if !board.ValidMove(Go) {
		message("You can't go there!")
		return false
	}

	board.YourGo(Go)
	draw()
	checkState()
	return false
}

func newGame() bool {
	form := getElementById("newGameForm")
	level := form.Get("level").Get("value").Int()
	first := form.Get("first").Get("value").String()
	println("new game", level, first)
	board = oxo3d.NewOxo3d(first == "no")
	switch level {
	case 0:
		player = oxo3d.NewOxo3dHeuristic(board, 0)
	case 1:
		player = oxo3d.NewOxo3dHeuristic(board, 1)
	case 2:
		player = oxo3d.NewOxo3dHeuristic(board, 2)
	case 3:
		player = oxo3d.NewOxo3dAlphaBeta(board, 2)
	case 4:
		player = oxo3d.NewOxo3dAlphaBeta(board, 3)
	case 5:
		player = oxo3d.NewOxo3dAlphaBeta(board, 4)
	default:
		fatalf("unknown level %d", level)
	}
	// draw the board
	draw()
	checkState()
	return false
}

// Set up the page ready to play - called when the DOM is loaded
func initialise() int {
	// find the squares
	for i := range square {
		square[i] = getElementById(fmt.Sprintf("sq%02d", i))
	}

	// attach click handlers
	for i := range square {
		square[i].Call("setAttribute", "onclick", fmt.Sprintf("oxo3dweb.clickSquare(%d)", i))
	}

	// Find the message div
	messageP = getElementById("message")
	message("")

	// attach handler for new game
	getElementById("newGame").Call("setAttribute", "onclick", "oxo3dweb.newGame(); return false;")

	// draw the board
	draw()
	checkState()
	return 0
}

// main entry point
func main() {
	if js.Global == nil {
		fmt.Println("Not running from browser - exiting")
		return
	}

	// expose functions to javascript
	js.Global.Set("oxo3dweb", map[string]interface{}{
		"clickSquare": clickSquare,
		"newGame":     newGame,
	})

	// find the document
	document = js.Global.Get("document")
	if document == js.Undefined {
		fatalf("couldn't find document")
	}

	// Run initialise when the dom is loaded
	document.Call("addEventListener", "DOMContentLoaded", initialise, false)
}
