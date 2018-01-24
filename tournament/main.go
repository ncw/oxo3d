// Oxo3d tournament to compare the different players against each other
package main

// FIXME name those players - put them in a registry? Then can mention
// them on the command line

// FIXME first isn't changing is it?
// Implement reset the board newGame so it does!

import (
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

	"github.com/ncw/oxo3d/oxo3d"
)

// Flags
var (
	rounds     = flag.Int("rounds", 10, "Number of rounds")
	cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
	threads    = flag.Int("threads", 0, "Number of threads to use (default 1 per CPU)")
	list       = flag.Bool("list", false, "List all the possible players")
)

type Player struct {
	Name      string
	NewPlayer oxo3d.NewPlayer
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

var defaultPlayerNames = []string{
	"Random",
	"Heuristic0",
	"Heuristic1",
	"Heuristic2",
	// "Minimax0",
	// "Minimax1",
	// "Minimax2",
	// "Minimax3",
	// "AlphaBeta0",
	"AlphaBeta1",
	"AlphaBeta2",
	"AlphaBeta3",
	"AlphaBeta4",
}

var players Players

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
	a_player := a.NewPlayer(a_game)
	b_player := b.NewPlayer(b_game)
	Go := -1
	// Make first two moves randomly
	for i := 0; i < 2; {
		Go = rand.Intn(64)
		if a_game.ValidMove(Go) {
			a_player.Play(Go, a_game.IsMyGo())
			b_player.Play(Go, b_game.IsMyGo())
			i++
		}
	}
	for !a_game.GameOver() {
		start := time.Now()
		asGo := a_game.IsMyGo()
		if asGo {
			Go = a_player.CalculateMyGo()
			s.Lock()
			s.a_time += time.Since(start)
			s.a_gos += 1
			s.Unlock()
		} else {
			Go = b_player.CalculateMyGo()
			s.Lock()
			s.b_time += time.Since(start)
			s.b_gos += 1
			s.Unlock()
		}
		a_player.Play(Go, asGo)
		b_player.Play(Go, !asGo)
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

// usage prints the syntax
func usage() {
	fmt.Fprintf(os.Stderr, `Oxo3d Tournament 

Syntax: %s [Options] [<Player>]+

Supply a list of players to engage in a tournament.  If you supply
none then you will get a default set of players.

Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Flag for seed?
	flag.Usage = usage
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

	// List players if desired
	if *list {
		fmt.Printf("Players are :-\n")
		names := []string{}
		for name := range oxo3d.Players {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("  %s\n", name)
		}
		return
	}

	// Work out who is playing
	playerNames := flag.Args()
	if len(playerNames) == 0 {
		playerNames = defaultPlayerNames
	}
	for _, Name := range playerNames {
		NewPlayer, ok := oxo3d.Players[Name]
		if !ok {
			log.Fatalf("Couldn't find player %v", Name)
		}
		player := Player{
			Name:      Name,
			NewPlayer: NewPlayer,
		}
		players = append(players, player)
	}
	tourney(*rounds)
}
