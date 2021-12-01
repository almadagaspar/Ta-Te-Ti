package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

var boardgame = []string{
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
}

var cursorY = 2
var cursorX = 4

var cursor = "*"
var cursorEmpty = string(boardgame[cursorY][cursorX]) // Guardo el character que estaba en la ubicación central antes de aparecer el cursor ahí.

const LEFT_MARGIN = "            "
const BLINK_TIME = 300
const PLAYER = "O"
const COMPUTER = "X"

var winnerIs = ""

func main() {

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		hideTerminalCursor("cvvis")
		_ = keyboard.Close()
	}()

	clearScreen()
	hideTerminalCursor("civis")
	go render()

	for {
		// Registramos que tecla presiono el usuario.
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}

		// Detectamos si el usuario presiono la tecla para salir del proigrama:
		if key == keyboard.KeyEsc {
			break
		}

		switch key {
		case keyboard.KeyArrowRight:
			if cursorX < 8 { // Impido que el cursor se salga del tablero.
				deleteCrusor()                                      // Borro el cursor en su posición actual por moverse a una nueva posición.
				cursorEmpty = string(boardgame[cursorY][cursorX+4]) // Almaceno el character que ya esta en el lugar adonde se va a colocar el cursor. Lo convertirlo a string porque es un byte.
				cursorX = cursorX + 4                               // Actualizo la posición del cursor segun la tecla presionada.
			}
		case keyboard.KeyArrowLeft:
			if cursorX > 0 {
				deleteCrusor()
				cursorEmpty = string(boardgame[cursorY][cursorX-4])
				cursorX = cursorX - 4
			}
		case keyboard.KeyArrowDown:
			if cursorY < 4 {
				deleteCrusor()
				cursorEmpty = string(boardgame[cursorY+2][cursorX])
				cursorY = cursorY + 2
			}
		case keyboard.KeyArrowUp:
			if cursorY > 0 {
				deleteCrusor()
				cursorEmpty = string(boardgame[cursorY-2][cursorX])
				cursorY = cursorY - 2
			}
		case keyboard.KeySpace:
			if cursorEmpty == " " { // Si esta vacio el lugar que ocupa el cursor, poner la pieza del jugador.
				winControl()
				cursorEmpty = PLAYER
			}
		case keyboard.KeyEnter:
			if cursorEmpty == " " {
				cursorEmpty = "X"
			}
		default:
			fmt.Print(string(char))
		}
	}
}

func winControl() {
	// Control de filas ─
	// Si en la fila en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if strings.Count(boardgame[cursorY], PLAYER) == 2 {
		winnerIs = "Player Wins!!"
		return
	}

	// Control de columnas |
	//  Si en la columna en la que voy a poner una pieza ya hay dos piezas puestas, se gana la partida.
	count := 0
	for y := 0; y <= len(boardgame); y = y + 2 {
		if string(boardgame[y][cursorX]) == PLAYER {
			count++
		}
	}
	if count == 2 {
		winnerIs = "Player Wins!!"
		return
	}

	// Control de diagonal \
	// Si en la diagonal en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if (cursorY == 0 && cursorX == 0) && string(boardgame[2][4]) == PLAYER && string(boardgame[4][8]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (cursorY == 2 && cursorX == 4) && string(boardgame[0][0]) == PLAYER && string(boardgame[4][8]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (cursorY == 4 && cursorX == 8) && string(boardgame[0][0]) == PLAYER && string(boardgame[2][4]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	}

	//Control de diagonal /
	// Si en la diagonal en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if (cursorY == 0 && cursorX == 8) && string(boardgame[2][4]) == PLAYER && string(boardgame[4][0]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (cursorY == 2 && cursorX == 4) && string(boardgame[0][8]) == PLAYER && string(boardgame[4][0]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (cursorY == 4 && cursorX == 0) && string(boardgame[0][8]) == PLAYER && string(boardgame[2][4]) == PLAYER {
		winnerIs = "Player Wins!!"
		return
	}

}

func deleteCrusor() { // Borro el cursor en su posición actual por moverse a una nueva posición.
	boardgame[cursorY] = fmt.Sprintf("%v%v%v", boardgame[cursorY][:cursorX], cursorEmpty, boardgame[cursorY][cursorX+1:])
	clearScreen()
	renderBoardgame()
}

func render() {
	for {
		renderBoardgame()
		time.Sleep(time.Millisecond * BLINK_TIME)
		boardgame[cursorY] = fmt.Sprintf("%v%v%v", boardgame[cursorY][:cursorX], cursor, boardgame[cursorY][cursorX+1:])
		clearScreen()

		renderBoardgame()
		time.Sleep(time.Millisecond * BLINK_TIME)
		boardgame[cursorY] = fmt.Sprintf("%v%v%v", boardgame[cursorY][:cursorX], cursorEmpty, boardgame[cursorY][cursorX+1:])
		clearScreen()
	}
}

func renderBoardgame() {
	for i := 0; i < len(boardgame); i++ {
		fmt.Print(LEFT_MARGIN)
		fmt.Println(boardgame[i])
	}
	fmt.Println(winnerIs)
}

func clearScreen() {
	if strings.Contains(runtime.GOOS, "windows") { // windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") // linux or mac
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func hideTerminalCursor(option string) {
	cmd := exec.Command("tput", option)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
