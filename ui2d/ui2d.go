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
	rand.Seed(1)

	for y, row := range level.Map {
		for x, tile := range row {
			if tile != game.Blank {
				srcRects := textureIndex[tile]
				srcRect := srcRects[rand.Intn(len(srcRects))]
				dstRect := sdl.Rect{
					X: int32(x) * TextureWidth,
					Y: int32(y) * TextureHeight,
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
		X: int32(level.Player.X) * TextureWidth,
		Y: int32(level.Player.Y) * TextureHeight,
		W: TextureWidth,
		H: TextureHeight,
	}
	renderer.Copy(textureAtlas, &srcRect, &dstRect)

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
