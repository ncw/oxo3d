// Oxo3d tournament to compare the different players against each other
package main

import (
	"../oxo3d"
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Player struct {
	NewPlayer func(*oxo3d.Oxo3d, int) oxo3d.Player
	Level     int
	Name      string
	Score     int
	CpuTime   time.Duration
}

type Players []Player

// Sort interface
func (ps Players) Len() int {
	return len(ps)
}
func (ps Players) Less(i, j int) bool {
	return ps[i].Score < ps[j].Score
}
func (ps Players) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

var (
	rounds = flag.Int("rounds", 10, "Number of rounds")
	// FIXME do a random player
	players = Players{
		{NewPlayer: oxo3d.NewOxo3dRandom, Level: 0, Name: "Random"},
		{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 0, Name: "Heuristic0"},
		{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 1, Name: "Heuristic1"},
		{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 2, Name: "Heuristic2"},
		//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 0, Name: "Minimax0"},
		//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 1, Name: "Minimax1"},
		//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 2, Name: "Minimax2"},
		//{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 0, Name: "AlphaBeta0"},
		{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 1, Name: "AlphaBeta1"},
		{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 2, Name: "AlphaBeta2"},
		{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 3, Name: "AlphaBeta3"},
		//{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 4, Name: "AlphaBeta4"},
	}
)

// FIXME name those players - put them in a registry?

func tourney(rounds int) {
	for playerB := 0; playerB < len(players)-1; playerB++ {
		for playerA := playerB + 1; playerA < len(players); playerA++ {
			a := &players[playerA]
			b := &players[playerB]
			first := true
			wins := 0
			losses := 0
			draws := 0
			var a_time, b_time time.Duration
			a_gos := 0
			b_gos := 0
			Go := 0
			for round := 0; round < rounds; round++ {
				a_game := oxo3d.NewOxo3d(first)
				b_game := oxo3d.NewOxo3d(!first)
				a_player := a.NewPlayer(a_game, a.Level)
				b_player := b.NewPlayer(b_game, b.Level)
				// Make first two moves randomly
				Go = rand.Intn(64)
				a_game.Play(Go, a_game.IsMyGo())
				b_game.Play(Go, b_game.IsMyGo())
				for {
					Go = rand.Intn(64)
					if a_game.ValidMove(Go) {
						break
					}
				}
				a_game.Play(Go, a_game.IsMyGo())
				b_game.Play(Go, b_game.IsMyGo())
				for !a_game.GameOver() {
					start := time.Now()
					if a_game.IsMyGo() {
						Go = a_player.CalculateMyGo()
						a_time += time.Since(start)
						a_gos += 1
						a_game.MyGo(Go)
						b_game.YourGo(Go)
					} else {
						Go = b_player.CalculateMyGo()
						b_time += time.Since(start)
						b_gos += 1
						b_game.MyGo(Go)
						a_game.YourGo(Go)
					}
				}
				switch a_game.WhoWon() {
				case oxo3d.ME:
					if b_game.WhoWon() != oxo3d.YOU {
						panic("Bad ME")
					}
					wins++
				case oxo3d.YOU:
					if b_game.WhoWon() != oxo3d.ME {
						panic("Bad YOU")
					}
					losses++
				case oxo3d.DRAW:
					if b_game.WhoWon() != oxo3d.DRAW {
						panic("Bad DRAW")
					}
					draws++
				}
				first = !first
			}
			fmt.Printf("level %10s vs %10s, Wins %5d, Draws %5d, Losses %5d: %10d ns vs %10d ns (per go)\n", a.Name, b.Name, wins, draws, losses, a_time.Nanoseconds()/int64(a_gos), b_time.Nanoseconds()/int64(b_gos))
			a.Score += wins - losses
			b.Score += losses - wins
			a.CpuTime += a_time
			b.CpuTime += b_time
		}
	}

	fmt.Printf("\n\nScores\n")
	sort.Sort(players)
	for i := range players {
		player := players[i]
		fmt.Printf("Player %10s Score %5d CpuTime %20s\n", player.Name, player.Score, player.CpuTime)

	}
}

func main() {
	// Flag for seed?
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	tourney(*rounds)
}
