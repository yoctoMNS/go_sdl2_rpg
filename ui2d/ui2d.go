package ui2d

import (
	"bufio"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/yoctoMNS/rpg/game"
)

const (
	WinWidth      int32 = 1280
	WinHeight     int32 = 720
	TextureWidth  int32 = 32
	TextureHeight int32 = 32
)

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile][]sdl.Rect
var keyboardState []uint8
var prevKeyboardState []uint8
var centerX int
var centerY int

func init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WinWidth, WinHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	textureAtlas = imgFileToTexture("ui2d/assets/tiles.png")
	loadTextureIndex()

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	copy(prevKeyboardState, keyboardState)

	centerX = -1
	centerY = -1
}

type UI2d struct {
}

func (ui *UI2d) Draw(level *game.Level) {
	if centerX == -1 && centerY == -1 {
		centerX = level.Player.X
		centerY = level.Player.Y
	}

	limit := 5
	if level.Player.X > centerX+limit {
		centerX++
	} else if level.Player.X < centerX-limit {
		centerX--
	} else if level.Player.Y > centerY+limit {
		centerY++
	} else if level.Player.Y < centerY-limit {
		centerY--
	}

	offsetX := WinWidth/2 - int32(centerX)*TextureWidth
	offsetY := WinHeight/2 - int32(centerY)*TextureHeight

	renderer.Clear()
	rand.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile != game.Blank {
				srcRects := textureIndex[tile]
				srcRect := srcRects[rand.Intn(len(srcRects))]
				dstRect := sdl.Rect{
					X: int32(x)*TextureWidth + offsetX,
					Y: int32(y)*TextureHeight + offsetY,
					W: TextureWidth,
					H: TextureHeight,
				}
				renderer.Copy(textureAtlas, &srcRect, &dstRect)
			}
		}
	}

	srcRect := sdl.Rect{
		X: 21 * TextureWidth,
		Y: 59 * TextureWidth,
		W: TextureWidth,
		H: TextureHeight,
	}
	dstRect := sdl.Rect{
		X: int32(level.Player.X)*TextureWidth + offsetX,
		Y: int32(level.Player.Y)*TextureHeight + offsetY,
		W: TextureWidth,
		H: TextureHeight,
	}
	renderer.Copy(textureAtlas, &srcRect, &dstRect)

	renderer.Present()
}

func imgFileToTexture(fileName string) *sdl.Texture {
	infile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	texture.Update(nil, pixels, w*4)

	err = texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return texture
}

func loadTextureIndex() {
	textureIndex = make(map[game.Tile][]sdl.Rect)

	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := game.Tile(line[0])
		xy := line[1:]
		splitXYC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		if err != nil {
			panic(err)
		}
		var rects []sdl.Rect
		for i := 0; int64(i) < variationCount; i++ {
			rect := sdl.Rect{
				X: int32(x) * TextureWidth,
				Y: int32(y) * TextureHeight,
				W: TextureWidth,
				H: TextureHeight,
			}
			rects = append(rects, rect)
			x++
			if x > 62 {
				x = 0
				y++
			}
		}

		textureIndex[tileRune] = rects
	}
}

func (ui *UI2d) GetInput() *game.Input {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{
					Typ: game.Quit,
				}
			}
		}

		var input game.Input
		if keyboardState[sdl.SCANCODE_UP] == 0 && prevKeyboardState[sdl.SCANCODE_UP] != 0 {
			input.Typ = game.Up
		}
		if keyboardState[sdl.SCANCODE_DOWN] == 0 && prevKeyboardState[sdl.SCANCODE_DOWN] != 0 {
			input.Typ = game.Down
		}
		if keyboardState[sdl.SCANCODE_LEFT] == 0 && prevKeyboardState[sdl.SCANCODE_LEFT] != 0 {
			input.Typ = game.Left
		}
		if keyboardState[sdl.SCANCODE_RIGHT] == 0 && prevKeyboardState[sdl.SCANCODE_RIGHT] != 0 {
			input.Typ = game.Right
		}
		if keyboardState[sdl.SCANCODE_ESCAPE] == 0 && prevKeyboardState[sdl.SCANCODE_ESCAPE] != 0 {
			input.Typ = game.Quit
		}

		copy(prevKeyboardState, keyboardState)

		if input.Typ != game.None {
			return &input
		}
	}
}
