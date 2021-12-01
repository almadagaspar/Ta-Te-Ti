package gra

import "myapp/ctr"

// Este package contiene los graficos del juego.

var Boardgame = []string{
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
}

var Cursor = "*"
var CursorInv = string(Boardgame[ctr.CursorY][ctr.CursorX]) // Guardo el character que se debe mostrar cuando el cursor se hace invisible al parpadear.

const MARGIN = "            "
const PLAYER = "O"
const COMPUTER = "X"
