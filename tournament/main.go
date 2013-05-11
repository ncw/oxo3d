// Oxo3d tournament to compare the different players against each other
package main

// FIXME name those players - put them in a registry? Then can mention
// them on the command line

// FIXME first isn't changing is it?
// Implement reset the board newGame so it does!

import (
	"../oxo3d"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"sync"
	"time"
)

// Flags
var (
	rounds     = flag.Int("rounds", 10, "Number of rounds")
	cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
	threads    = flag.Int("threads", 0, "Number of threads to use (default 1 per CPU)")
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

var players = Players{
	//{NewPlayer: oxo3d.NewOxo3dRandom, Level: 0, Name: "Random"},
	//{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 0, Name: "Heuristic0"},
	//{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 1, Name: "Heuristic1"},
	//{NewPlayer: oxo3d.NewOxo3dHeuristic, Level: 2, Name: "Heuristic2"},
	//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 0, Name: "Minimax0"},
	//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 1, Name: "Minimax1"},
	//{NewPlayer: oxo3d.NewOxo3dMinimax, Level: 2, Name: "Minimax2"},
	//{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 0, Name: "AlphaBeta0"},
	{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 1, Name: "AlphaBeta1"},
	{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 2, Name: "AlphaBeta2"},
	{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 3, Name: "AlphaBeta3"},
	//{NewPlayer: oxo3d.NewOxo3dAlphaBeta, Level: 4, Name: "AlphaBeta4"},
}

type stats struct {
	sync.Mutex
	queue               chan bool
	wg                  sync.WaitGroup
	wins, losses, draws int
	a_time, b_time      time.Duration
	a_gos, b_gos        int
}

// Do a round between player a and b recording the answer in stats
func doRound(a, b *Player, first bool, s *stats) {
	defer s.wg.Done()
	a_game := oxo3d.NewOxo3d(first)
	b_game := oxo3d.NewOxo3d(!first)
	a_player := a.NewPlayer(a_game, a.Level)
	b_player := b.NewPlayer(b_game, b.Level)
	// Make first two moves randomly
	Go := rand.Intn(64)
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
			s.Lock()
			s.a_time += time.Since(start)
			s.a_gos += 1
			s.Unlock()
			a_game.MyGo(Go)
			b_game.YourGo(Go)
		} else {
			Go = b_player.CalculateMyGo()
			s.Lock()
			s.b_time += time.Since(start)
			s.b_gos += 1
			s.Unlock()
			b_game.MyGo(Go)
			a_game.YourGo(Go)
		}
	}
	s.Lock()
	switch a_game.WhoWon() {
	case oxo3d.ME:
		if b_game.WhoWon() != oxo3d.YOU {
			panic("Bad ME")
		}
		s.wins++
	case oxo3d.YOU:
		if b_game.WhoWon() != oxo3d.ME {
			panic("Bad YOU")
		}
		s.losses++
	case oxo3d.DRAW:
		if b_game.WhoWon() != oxo3d.DRAW {
			panic("Bad DRAW")
		}
		s.draws++
	}
	s.Unlock()
	<-s.queue
}

func tourney(rounds int) {
	for playerB := 0; playerB < len(players)-1; playerB++ {
		for playerA := playerB + 1; playerA < len(players); playerA++ {
			a := &players[playerA]
			b := &players[playerB]
			first := true
			s := stats{
				queue: make(chan bool, *threads),
			}
			for round := 0; round < rounds; round++ {
				s.wg.Add(1)
				s.queue <- true
				go doRound(a, b, first, &s)
				first = !first
			}
			s.wg.Wait()
			fmt.Printf("level %10s vs %10s, Wins %5d, Draws %5d, Losses %5d: %10d ns vs %10d ns (per go)\n", a.Name, b.Name, s.wins, s.draws, s.losses, s.a_time.Nanoseconds()/int64(s.a_gos), s.b_time.Nanoseconds()/int64(s.b_gos))
			a.Score += s.wins - s.losses
			b.Score += s.losses - s.wins
			a.CpuTime += s.a_time
			b.CpuTime += s.b_time
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
	if *threads == 0 {
		*threads = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*threads)

	// Setup profiling if desired
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	tourney(*rounds)
}
