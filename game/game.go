package game

import (
	"bufio"
	"math"
	"os"
	"sort"
	"time"
)

// Part25 49:00
type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
	Search
)

type Input struct {
	Typ InputType
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	CloseDoor Tile = '|'
	OpenDoor  Tile = '/'
	Blank     Tile = 0
	Pending   Tile = -1
)

type Pos struct {
	X int
	Y int
}

type Entity struct {
	Pos
}

type Player struct {
	Entity
}

type priorityPos struct {
	Pos
	priority int
}

type priorityArray []priorityPos

func (p priorityArray) Len() int {
	return len(p)
}

func (p priorityArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p priorityArray) Less(i, j int) bool {
	return p[i].priority < p[j].priority
}

type Level struct {
	Map    [][]Tile
	Player Player
	Debug  map[Pos]bool
}

func LoadLevelFromFile(fileName string) *Level {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}

	level := &Level{}
	level.Map = make([][]Tile, len(levelLines))
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = CloseDoor
			case '/':
				t = OpenDoor
			case '.':
				t = DirtFloor
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			default:
				panic("Invalid character in map")
			}

			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchX]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}

func canWalk(level *Level, pos Pos) bool {
	t := level.Map[pos.Y][pos.X]
	switch t {
	case StoneWall, CloseDoor, Blank:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]
	if t == CloseDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func handleInput(ui GameUI, level *Level, input *Input) {
	p := level.Player
	switch input.Typ {
	case Up:
		if canWalk(level, Pos{p.X, p.Y - 1}) {
			level.Player.Y--
		} else {
			checkDoor(level, Pos{p.X, p.Y - 1})
		}
	case Down:
		if canWalk(level, Pos{p.X, p.Y + 1}) {
			level.Player.Y++
		} else {
			checkDoor(level, Pos{p.X, p.Y + 1})
		}
	case Left:
		if canWalk(level, Pos{p.X - 1, p.Y}) {
			level.Player.X--
		} else {
			checkDoor(level, Pos{p.X - 1, p.Y})
		}
	case Right:
		if canWalk(level, Pos{p.X + 1, p.Y}) {
			level.Player.X++
		} else {
			checkDoor(level, Pos{p.X + 1, p.Y})
		}
	case Search:
		astar(ui, level, p.Pos, Pos{3, 2})
		// bfs(ui, level, p.Pos)
	}
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{X: pos.X - 1, Y: pos.Y}
	right := Pos{X: pos.X + 1, Y: pos.Y}
	up := Pos{X: pos.X, Y: pos.Y - 1}
	down := Pos{X: pos.X, Y: pos.Y + 1}

	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func astar(ui GameUI, level *Level, start, goal Pos) []Pos {
	frontier := make(priorityArray, 0, 8)
	frontier = append(frontier, priorityPos{Pos: start, priority: 1})
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0
	level.Debug = make(map[Pos]bool)

	for len(frontier) > 0 {
		sort.Stable(frontier)
		current := frontier[0]

		if current.Pos == goal {
			path := make([]Pos, 0)
			p := current.Pos
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}

			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			for _, pos := range path {
				level.Debug[pos] = true
				ui.Draw(level)
				time.Sleep(100 * time.Millisecond)
			}

			return path
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current.Pos) {
			newCost := costSoFar[current.Pos] + 1
			_, exists := costSoFar[next]

			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = append(frontier, priorityPos{Pos: next, priority: priority})
				cameFrom[next] = current.Pos
			}
		}
	}

	return nil
}

func bfs(ui GameUI, level *Level, start Pos) {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
				ui.Draw(level)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")
	for {
		ui.Draw(level)
		if input := ui.GetInput(); input != nil {
			if input.Typ == Quit {
				return
			}

			handleInput(ui, level, input)
		}
	}
}
