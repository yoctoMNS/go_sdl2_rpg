package ui2d

import (
	"bufio"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/yoctoMNS/rpg/game"
)

const (
	WinWidth  int32 = 1280
	WinHeight int32 = 720
)

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile]sdl.Rect

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
}

type UI2d struct {
}

func (ui *UI2d) Draw(level *game.Level) {
	renderer.Copy(textureAtlas, nil, nil)
	renderer.Present()
	sdl.Delay(5000)
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
		splitXY := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXY[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXY[1]), 10, 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(tileRune, x, y)
	}
}
