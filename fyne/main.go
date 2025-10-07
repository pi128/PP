package main

import (
	"image/color"
	"math/rand"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// ---- Maze types ----

type Cell struct {
	Visited bool
	// walls: N, E, S, W
	Wall [4]bool
}

type Maze struct {
	Cols, Rows int
	Grid       [][]Cell
}

const (
	N = 0
	E = 1
	S = 2
	W = 3
)

func newMaze(cols, rows int) *Maze {
	m := &Maze{
		Cols: cols,
		Rows: rows,
		Grid: make([][]Cell, rows),
	}
	for y := 0; y < rows; y++ {
		m.Grid[y] = make([]Cell, cols)
		for x := 0; x < cols; x++ {
			m.Grid[y][x] = Cell{Wall: [4]bool{true, true, true, true}}
		}
	}
	rand.Seed(time.Now().UnixNano())
	m.carvePassagesFrom(0, 0)
	return m
}

func (m *Maze) carvePassagesFrom(cx, cy int) {
	type pt struct{ x, y int }
	stack := []pt{{cx, cy}}
	m.Grid[cy][cx].Visited = true

	dirs := []int{N, E, S, W}

	for len(stack) > 0 {
		cur := stack[len(stack)-1]

		// shuffle directions
		rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

		moved := false
		for _, d := range dirs {
			nx, ny := neighbor(cur.x, cur.y, d)
			if nx >= 0 && nx < m.Cols && ny >= 0 && ny < m.Rows && !m.Grid[ny][nx].Visited {
				// knock down walls both sides
				m.Grid[cur.y][cur.x].Wall[d] = false
				m.Grid[ny][nx].Wall[opp(d)] = false

				m.Grid[ny][nx].Visited = true
				stack = append(stack, pt{nx, ny})
				moved = true
				break
			}
		}
		if !moved {
			stack = stack[:len(stack)-1]
		}
	}
}

func neighbor(x, y, dir int) (int, int) {
	switch dir {
	case N:
		return x, y - 1
	case E:
		return x + 1, y
	case S:
		return x, y + 1
	case W:
		return x - 1, y
	}
	return x, y
}

func opp(dir int) int {
	return (dir + 2) % 4
}

// ---- UI / Game ----

type Game struct {
	maze          *Maze
	playerX       int
	playerY       int
	goalX         int
	goalY         int
	cellSize      float32
	wallThickness float32
	padding       float32

	// Drawn objects
	wallLines []fyne.CanvasObject
	player    *canvas.Circle
	goal      *canvas.Rectangle
	layer     *fyne.Container
}

func newGame(cols, rows int, cellSize float32) *Game {
	g := &Game{
		maze:          newMaze(cols, rows),
		playerX:       0,
		playerY:       0,
		goalX:         cols - 1,
		goalY:         rows - 1,
		cellSize:      cellSize,
		wallThickness: 2,
		padding:       8,
	}
	return g
}

func (g *Game) buildUI(win fyne.Window) fyne.CanvasObject {
	// Layers: background + walls + goal + player
	bg := canvas.NewRectangle(color.RGBA{240, 240, 240, 255})

	// container without layout so we can position precisely
	g.layer = container.NewWithoutLayout(bg)

	// Draw cell backgrounds (optional subtle grid)
	gridColor := color.RGBA{250, 250, 250, 255}
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			r := canvas.NewRectangle(gridColor)
			r.Resize(fyne.NewSize(g.cellSize-1, g.cellSize-1))
			r.Move(g.cellTopLeft(x, y))
			g.layer.Add(r)
		}
	}

	// Draw walls
	g.wallLines = g.drawWalls()
	for _, w := range g.wallLines {
		g.layer.Add(w)
	}

	// Goal (finish) cell
	g.goal = canvas.NewRectangle(color.RGBA{200, 255, 200, 255})
	g.goal.Resize(fyne.NewSize(g.cellSize-6, g.cellSize-6))
	g.goal.Move(g.cellTopLeft(g.goalX, g.goalY).AddXY(3, 3))
	g.layer.Add(g.goal)

	// Player
	g.player = canvas.NewCircle(color.RGBA{50, 120, 255, 255})
	size := g.cellSize * 0.6
	g.player.Resize(fyne.NewSize(size, size))
	g.updatePlayerPos()
	g.layer.Add(g.player)

	// Background should span full board
	boardW := g.padding*2 + float32(g.maze.Cols)*g.cellSize
	boardH := g.padding*2 + float32(g.maze.Rows)*g.cellSize
	bg.Resize(fyne.NewSize(boardW, boardH))

	// Key handling: arrows + WASD
	win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyUp:
			g.tryMove(0, -1, win)
		case fyne.KeyDown:
			g.tryMove(0, 1, win)
		case fyne.KeyLeft:
			g.tryMove(-1, 0, win)
		case fyne.KeyRight:
			g.tryMove(1, 0, win)
		}
	})
	win.Canvas().SetOnTypedRune(func(r rune) {
		switch unicode.ToLower(r) {
		case 'w':
			g.tryMove(0, -1, win)
		case 's':
			g.tryMove(0, 1, win)
		case 'a':
			g.tryMove(-1, 0, win)
		case 'd':
			g.tryMove(1, 0, win)
		}
	})

	return g.layer
}

func (g *Game) cellTopLeft(x, y int) fyne.Position {
	return fyne.NewPos(g.padding+float32(x)*g.cellSize, g.padding+float32(y)*g.cellSize)
}

func (g *Game) drawWalls() []fyne.CanvasObject {
	var lines []fyne.CanvasObject
	w := g.wallThickness

	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			cell := g.maze.Grid[y][x]
			topLeft := g.cellTopLeft(x, y)
			x0 := topLeft.X
			y0 := topLeft.Y
			x1 := x0 + g.cellSize
			y1 := y0 + g.cellSize

			col := color.RGBA{0, 0, 0, 255}

			if cell.Wall[N] {
				r := canvas.NewRectangle(col)
				r.Resize(fyne.NewSize(g.cellSize, w))
				r.Move(fyne.NewPos(x0, y0))
				lines = append(lines, r)
			}
			if cell.Wall[E] {
				r := canvas.NewRectangle(col)
				r.Resize(fyne.NewSize(w, g.cellSize))
				r.Move(fyne.NewPos(x1-w, y0))
				lines = append(lines, r)
			}
			if cell.Wall[S] {
				r := canvas.NewRectangle(col)
				r.Resize(fyne.NewSize(g.cellSize, w))
				r.Move(fyne.NewPos(x0, y1-w))
				lines = append(lines, r)
			}
			if cell.Wall[W] {
				r := canvas.NewRectangle(col)
				r.Resize(fyne.NewSize(w, g.cellSize))
				r.Move(fyne.NewPos(x0, y0))
				lines = append(lines, r)
			}
		}
	}
	return lines
}

func (g *Game) updatePlayerPos() {
	// center the circle in the cell
	center := g.cellTopLeft(g.playerX, g.playerY)
	center = center.AddXY(g.cellSize/2, g.cellSize/2)
	// move circle so its center aligns
	g.player.Move(center.Subtract(fyne.NewPos(g.player.Size().Width/2, g.player.Size().Height/2)))
	canvas.Refresh(g.player)
}

func (g *Game) tryMove(dx, dy int, win fyne.Window) {
	nx := g.playerX + dx
	ny := g.playerY + dy
	if nx < 0 || ny < 0 || nx >= g.maze.Cols || ny >= g.maze.Rows {
		return
	}
	// check wall between current and target
	var dir int
	switch {
	case dy == -1 && dx == 0:
		dir = N
	case dy == 1 && dx == 0:
		dir = S
	case dx == -1 && dy == 0:
		dir = W
	case dx == 1 && dy == 0:
		dir = E
	default:
		return
	}
	if g.maze.Grid[g.playerY][g.playerX].Wall[dir] {
		// blocked by wall
		return
	}

	g.playerX = nx
	g.playerY = ny
	g.updatePlayerPos()

	// Check win
	if g.playerX == g.goalX && g.playerY == g.goalY {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "You escaped!",
			Content: "Generating a new maze…",
		})
		g.regenerate(win)
	}
}

func (g *Game) regenerate(win fyne.Window) {
	// new maze, reset pos
	g.maze = newMaze(g.maze.Cols, g.maze.Rows)
	g.playerX, g.playerY = 0, 0
	g.goalX, g.goalY = g.maze.Cols-1, g.maze.Rows-1

	// remove old walls
	for _, w := range g.wallLines {
		g.layer.Remove(w)
	}
	g.wallLines = g.drawWalls()
	for _, w := range g.wallLines {
		g.layer.Add(w)
	}

	// move goal + player
	g.goal.Move(g.cellTopLeft(g.goalX, g.goalY).AddXY(3, 3))
	g.updatePlayerPos()
	canvas.Refresh(g.layer)
}

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Maze — WASD / Arrows")
	// tweak these to change maze size
	cols, rows := 25, 18
	cell := float32(28)

	game := newGame(cols, rows, cell)
	content := game.buildUI(w)

	boardW := int(game.padding*2 + float32(cols)*game.cellSize)
	boardH := int(game.padding*2 + float32(rows)*game.cellSize)
	w.SetContent(content)
	w.Resize(fyne.NewSize(float32(boardW), float32(boardH)))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
