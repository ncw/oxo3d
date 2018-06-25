# Oxo3d

This is a suite of Go programs for playing and testing 3D noughts and crosses / tic tac toe.

Try [Oxo3d online](https://www.craig-wood.com/nick/oxo3d/).

Or try it in your terminal in old school ASCII:

    go get -u github.com/ncw/oxo3d
    oxo3d

```
+-------------+-------------+
| .  .  .  .  |  0  1  2  3 |
| .  .  .  .  |  4  5  6  7 |
| .  .  .  .  |  8  9 10 11 |
| .  .  .  .  | 12 13 14 15 |
+-------------+-------------+
| .  .  .  .  | 16 17 18 19 |
| .  .  .  .  | 20 21 22 23 |
| .  .  .  .  | 24 25 26 27 |
| .  .  .  X* | 28 29 30 31 |
+-------------+-------------+
| .  .  .  .  | 32 33 34 35 |
| .  .  .  .  | 36 37 38 39 |
| .  .  .  .  | 40 41 42 43 |
| .  .  .  .  | 44 45 46 47 |
+-------------+-------------+
| .  .  .  .  | 48 49 50 51 |
| .  .  .  .  | 52 53 54 55 |
| .  .  .  .  | 56 57 58 59 |
| .  .  .  .  | 60 61 62 63 |
Go (0..63) 
```

Or try [Oxo3d on your Android
device](https://play.google.com/store/apps/details?id=com.craig_wood.Oxo3d).
The Go code here is a port of the Java code used in this game.

## Playing the game

Oxo 3D is a 3 dimensional (4x4x4) noughts-and-crosses / tic-tac-toe
game.  This is a game with considerably more strategy than the
traditional 3x3 version.

You're O and you have to get 4 in a row.  The board is a cube
viewed in slices.  Imagine the 4 slices piled on top of each
other.

Watch out for tricky diagonal lines!

Good luck!

## How to play 3D Noughts-and-Crosses

Each group of 4x4 positions represents one plane on the cube.

On each layer there are 10 possible winning lines, 4 horizontal, 4
vertical and two diagonal.

However lines may run from layer to layer also.

There are 76 lines in total possible, these are 4x4 = north-south
lines, 16 east-west lines, 16 vertical lines, 2x4x3 = 24 diagonal
lines from one edge to another, and 4 diagonal ines from corner to
corner.
    
## Levels
    
There are 6 carefully graduated levels with interesting names.

The levels get increasingly hard.  *Easy* and *Simple*
use the same simple heuristic player, but *Easy* makes random
mistakes.  *Sneaky* uses a two level lookahead player.  It is
called *Sneaky* for good reason as you'll see when you play it.
*Relentless* plays very slightly better than *Sneaky*
but with a completely different style.  *Hard* is exactly that
and *Very Hard* plays so well that I haven't beaten it yet.

| Level name | Internal name | Description | Rank |
| ---------- | --------------| ----------- | ---- |
| Easy       | Heuristic0    | 1 level lookahead heuristic player with random errors | 1 |
| Simple     | Heuristic1    | 1 level lookahead heuristic player | 2 |
| Sneaky     | Heuristic2    | 2 level lookahead heuristic player | 3 |
| Relentless | AlphaBeta2    | 2 level lookahead minimax player   | 4 |
| Hard       | AlphaBeta3    | 3 level lookahead minimax player   | 5 |
| Very Hard  | AlphaBeta4    | 4 level lookahead minimax player   | 6 |
    
If you play 2,500 rounds of each one against the other this is what
the scores look like.  A Random player is included too as a
control. There are some interesting reversals of form in the table

| Player 1   | Player 2   | Wins  | Draws | Losses | Player 1 time per go | Player 2 time per go |
| :--------- | :--------- | ----: | ----: | -----: | ------------: | ------------: |
| Heuristic0 |     Random |  2500 |     0 |      0 |     182041 ns |       3506 ns |
| Heuristic1 |     Random |  2500 |     0 |      0 |      17193 ns |       1502 ns |
| Heuristic2 |     Random |  2500 |     0 |      0 |      24100 ns |        789 ns |
| AlphaBeta1 |     Random |  2500 |     0 |      0 |      18244 ns |       1320 ns |
| AlphaBeta2 |     Random |  2500 |     0 |      0 |     220337 ns |        307 ns |
| AlphaBeta3 |     Random |  2500 |     0 |      0 |    6296231 ns |        318 ns |
| AlphaBeta4 |     Random |  2500 |     0 |      0 |   88652221 ns |        413 ns |
| Heuristic1 | Heuristic0 |  2470 |     0 |     30 |      10722 ns |     145330 ns |
| Heuristic2 | Heuristic0 |  2274 |    22 |    204 |      17717 ns |     120008 ns |
| AlphaBeta1 | Heuristic0 |  2185 |     0 |    315 |      14914 ns |     146380 ns |
| AlphaBeta2 | Heuristic0 |  2465 |     1 |     34 |     248772 ns |      51983 ns |
| AlphaBeta3 | Heuristic0 |  2470 |     0 |     30 |    6134420 ns |      19663 ns |
| AlphaBeta4 | Heuristic0 |  2489 |     1 |     10 |  133985769 ns |      15304 ns |
| Heuristic2 | Heuristic1 |  1002 |   642 |    856 |      16801 ns |       8780 ns |
| AlphaBeta1 | Heuristic1 |   242 |     2 |   2256 |      15445 ns |      10625 ns |
| AlphaBeta2 | Heuristic1 |  1187 |   486 |    827 |     184999 ns |       8261 ns |
| AlphaBeta3 | Heuristic1 |  1470 |   124 |    906 |    4577929 ns |       8403 ns |
| AlphaBeta4 | Heuristic1 |  1775 |   216 |    509 |   93224450 ns |       8681 ns |
| AlphaBeta1 | Heuristic2 |   864 |    45 |   1591 |      13980 ns |      20328 ns |
| AlphaBeta2 | Heuristic2 |  1326 |   621 |    553 |     156175 ns |      16215 ns |
| AlphaBeta3 | Heuristic2 |  1848 |    81 |    571 |    4485952 ns |      19077 ns |
| AlphaBeta4 | Heuristic2 |  1599 |   660 |    241 |   49535585 ns |      16807 ns |
| AlphaBeta2 | AlphaBeta1 |  2220 |     1 |    279 |     234968 ns |      13948 ns |
| AlphaBeta3 | AlphaBeta1 |  2095 |     9 |    396 |    6054260 ns |      13851 ns |
| AlphaBeta4 | AlphaBeta1 |  2388 |     0 |    112 |  151574104 ns |      14706 ns |
| AlphaBeta3 | AlphaBeta2 |  1266 |    59 |   1175 |    4487537 ns |     237976 ns |
| AlphaBeta4 | AlphaBeta2 |  1319 |   193 |    988 |   86538385 ns |     223163 ns |
| AlphaBeta4 | AlphaBeta3 |  1740 |    85 |    675 |   82880581 ns |    4723343 ns |
    
...and here is the league table summary of the above
    
| Player     | Score  | CPU Time      |
| :--------- | -----: | ------------: |
|     Random | -17500 |          72ms |
| Heuristic0 | -11230 |        11.56s |
| AlphaBeta1 |  -4287 |         1.98s |
| Heuristic2 |   2035 |         4.54s |
| Heuristic1 |   4618 |         2.04s |
| AlphaBeta3 |   7506 |     15m51.07s |
| AlphaBeta2 |   7583 |        45.25s |
| AlphaBeta4 |  11275 |   5h19m56.46s |

Interestingly the Heuristic1 player beats the Heuristic2 player
overall, however you'll see in the detailed breakdown that the
Heuristic2 player beats the Heuristic1 player.

## History

  * Converted to Web with Gopherjs (2018)
  * Converted to Go (2013)
  * Converted to Java on Android (2010)
  * Converted to Python (2004)
  * Converted to Perl (2000)
  * Converted to C (1996)
  * Converted to Psion 3a OPL (1994)
  * Converted to BBC Basic (1985)
  * Converted to QL BASIC(1985)
  * Writen in ZX81 BASIC (1983)
  * 35 Years in the Making!

I used to play this game with my father when I was a boy.  We always
used to play on squared paper while sitting round the kitchen table.
When I was old enough (and the home computer had been invented) this
inspired me to write a computer player for the game.

The very first version of this game was written in 1983 on a Sinclair
ZX81 in BASIC.  It used to take 60 seconds to think of a move, and it
used to beat the author nearly all of the time.  It used the
Heuristic1 player.  However computing has moved on, and the current
version plays at a much higher strength almost instantly.  The author
has improved too!
    
## Contacting the Author

If you have a problem with *Oxo 3D*, it goes wrong in some fashion
then file an issue.  If you just want to say hello then you'll find my
email address below!

Nick Craig-Wood
nick@craig-wood.com