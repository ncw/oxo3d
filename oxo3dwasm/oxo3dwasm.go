package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"

	"github.com/ncw/oxo3d/oxo3d"
)

// Globals
var (
	board    = oxo3d.NewOxo3d(true)
	player   = oxo3d.NewOxo3dHeuristic(board, 0)
	symbols  = [3]string{"O", "", "X"}
	square   [64]js.Value
	viewMap  [64]int // mapping for the view
	document js.Value
	messageP js.Value
	helpSpan js.Value
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Exit with the message
func fatalf(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	js.Global().Call("alert", text)
	panic(text)
}

// getElementById gets the element with the ID passed in or panics
func getElementById(ID string) js.Value {
	obj := document.Call("getElementById", ID)
	if obj == js.Undefined() {
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
		p := viewMap[i]
		square[i].Set("innerHTML", symbols[board.Board(p)+1])
		if board.Marked(p) {
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
func clickSquare(arg []js.Value) {
	event := arg[0]
	// event.Call("preventDefault")
	this := event.Get("target")
	id := this.Get("id").String()
	if len(id) != 4 || id[0] != 's' || id[1] != 'q' {
		fmt.Printf("Invalid square id %q\n", id)
		return
	}
	Go, err := strconv.Atoi(id[2:])
	if err != nil {
		fmt.Printf("Invalid square id %q: %v\n", id, err)
		return
	}
	fmt.Printf("click on %q = %d\n", id, Go)
	Go = viewMap[Go]
	if board.GameOver() {
		return
	}
	message("")
	if board.IsMyGo() {
		message("Not your go yet!")
		return
	}
	fmt.Println(fmt.Sprintf("click on %d", Go))
	if !board.ValidMove(Go) {
		message("You can't go there!")
		return
	}

	board.YourGo(Go)
	draw()
	checkState()
}

// Set viewMap to describe the current orientation
func setOrientation(orientation int) {
	stride := 1
	switch orientation {
	case 0:
		stride = 1
	case 1:
		stride = 4
	case 2:
		stride = 16
	default:
		message("unknown orientation")
		return
	}
	// swizzle the viewMap
	p := 0
	for i := range square {
		viewMap[i] = p
		p += stride
		if p >= len(square) {
			p -= len(square)
			p += 1
		}
	}
}

// Called when an orientation is clicked
func clickOrientation(arg []js.Value) {
	event := arg[0]
	// event.Call("preventDefault")
	this := event.Get("target")
	id := this.Get("id").String()
	orientation := 0
	switch id {
	case "orientation1":
		orientation = 1
	case "orientation2":
		orientation = 2
	}
	fmt.Printf("orientation = %q\n", id)
	setOrientation(orientation)
	draw()
}

const helpMessage = `
<h2>Oxo 3D</h1>

<p><em>Oxo 3D</em> is a 3 dimensional (4x4x4) noughts-and-crosses /
tic-tac-toe game for the web browser.  This is a game with
considerably more strategy than the traditional 3x3 version.</p>

<h3>Quick start</h3>

<p>You're O and you have to get 4 in a row.  The board is a cube
viewed in slices.  Imagine the 4 slices piled on top of each other.</p>

<p>Touch the left side of the screen to play.  Rotate the cube on the
right side by dragging.</p>

<p>Watch out for tricky diagonal lines! (Use the X, Y Z buttons to
view the cube in different orientations.)</p>

<p>Good luck!</p>

PS For source code and more info see the <a href="https://github.com/ncw/oxo3d" target="_blank">oxo3d project on github</a>.
`

// Called when help is clicked
func clickHelp(arg []js.Value) {
	event := arg[0]
	this := event.Get("target")
	status := this.Get("innerHTML").String()
	var help string
	if status == "Show" {
		status = "Hide"
		help = helpMessage
	} else {
		status = "Show"
		help = ""
	}
	helpSpan.Set("innerHTML", help)
	this.Set("innerHTML", status)
}

func newGame(_ []js.Value) {
	form := getElementById("newGameForm")
	level := form.Get("level").Get("value").Int()
	first := form.Get("first").Get("value").String()
	fmt.Println("new game", level, first)
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
}

// Set up the page ready to play - called when the DOM is loaded
func initialise(_ []js.Value) {
	fmt.Printf("initialise\n")
	// find the squares
	for i := range square {
		square[i] = getElementById(fmt.Sprintf("sq%02d", i))
	}

	// attach click handlers
	for i := range square {
		square[i].Call("addEventListener", "click", js.NewCallback(clickSquare))
	}

	// Find the message div
	messageP = getElementById("message")
	message("")

	// attach handler for new game
	getElementById("newGame").Call("addEventListener", "click", js.NewCallback(newGame))

	// attach handler for orientation changes
	for orientation := 0; orientation < 3; orientation++ {
		getElementById(fmt.Sprintf("orientation%d", orientation)).Call("addEventListener", "click", js.NewCallback(clickOrientation))
	}

	// attach handler for help
	getElementById("help").Call("addEventListener", "click", js.NewCallback(clickHelp))
	helpSpan = getElementById("showHelp")

	// draw the board
	setOrientation(0)
	draw()
	checkState()
}

// main entry point
func main() {
	fmt.Printf("main\n")
	// find the document
	document = js.Global().Get("document")
	if document == js.Undefined() {
		fatalf("couldn't find document")
	}

	// FIXME make it run immediately

	// Run initialise when the dom is loaded
	//document.Call("addEventListener", "DOMContentLoaded", js.NewCallback(initialise))
	initialise(nil)

	// Wait forever - everything is done on callbacks now
	fmt.Printf("wait\n")
	select {}
}
