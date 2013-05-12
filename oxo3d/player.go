// Player interface for Oxo3d
package oxo3d

type Player interface {
	CalculateMyGo() int
	Play(int, bool)
}

type NewPlayer func(*Oxo3d) Player

// Registry of players
var Players = map[string]NewPlayer{}

// Check the interfaces
var _ Player = (*Oxo3dHeuristic)(nil)
var _ Player = (*Oxo3dMinimax)(nil)
var _ Player = (*Oxo3dAlphaBeta)(nil)
var _ Player = (*Oxo3dRandom)(nil)
