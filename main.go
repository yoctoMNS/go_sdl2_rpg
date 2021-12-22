package main

import (
	"github.com/yoctoMNS/rpg/game"
	"github.com/yoctoMNS/rpg/ui2d"
)

func main() {
	game.Run(&ui2d.UI2d{})
}
