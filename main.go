package main

import (
	"fmt"
	"log"
	"myapp/ctr"
	"myapp/gra"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

const BLINK_TIME = 300

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
			if ctr.CursorX < 8 { // Impido que el cursor se salga del tablero.
				deleteCrusor()                                                    // Borro el cursor en su posición actual por moverse a una nueva posición.
				gra.CursorInv = string(gra.Boardgame[ctr.CursorY][ctr.CursorX+4]) // Almaceno el character que ya esta en el lugar adonde se va a colocar el cursor. Lo convertirlo a string porque es un byte.
				ctr.CursorX = ctr.CursorX + 4                                     // Actualizo la posición del cursor segun la tecla presionada.
			}
		case keyboard.KeyArrowLeft:
			if ctr.CursorX > 0 {
				deleteCrusor()
				gra.CursorInv = string(gra.Boardgame[ctr.CursorY][ctr.CursorX-4])
				ctr.CursorX = ctr.CursorX - 4
			}
		case keyboard.KeyArrowDown:
			if ctr.CursorY < 4 {
				deleteCrusor()
				gra.CursorInv = string(gra.Boardgame[ctr.CursorY+2][ctr.CursorX])
				ctr.CursorY = ctr.CursorY + 2
			}
		case keyboard.KeyArrowUp:
			if ctr.CursorY > 0 {
				deleteCrusor()
				gra.CursorInv = string(gra.Boardgame[ctr.CursorY-2][ctr.CursorX])
				ctr.CursorY = ctr.CursorY - 2
			}
		case keyboard.KeySpace:
			if gra.CursorInv == " " { // Si esta vacio el lugar que ocupa el cursor, poner la pieza del jugador.
				winControl()
				gra.CursorInv = gra.PLAYER
			}
		case keyboard.KeyEnter:
			if gra.CursorInv == " " {
				gra.CursorInv = gra.COMPUTER
			}
		default:
			fmt.Print(string(char))
		}
	}
}

func winControl() {
	// Control de filas ─
	// Si en la fila en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if strings.Count(gra.Boardgame[ctr.CursorY], gra.PLAYER) == 2 {
		winnerIs = "Player Wins!!"
		return
	}

	// Control de columnas |
	//  Si en la columna en la que voy a poner una pieza ya hay dos piezas puestas, se gana la partida.
	count := 0
	for y := 0; y <= len(gra.Boardgame); y = y + 2 {
		if string(gra.Boardgame[y][ctr.CursorX]) == gra.PLAYER {
			count++
		}
	}
	if count == 2 {
		winnerIs = "Player Wins!!"
		return
	}

	// Control de diagonal \
	// Si en la diagonal en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if (ctr.CursorY == 0 && ctr.CursorX == 0) && string(gra.Boardgame[2][4]) == gra.PLAYER && string(gra.Boardgame[4][8]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (ctr.CursorY == 2 && ctr.CursorX == 4) && string(gra.Boardgame[0][0]) == gra.PLAYER && string(gra.Boardgame[4][8]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (ctr.CursorY == 4 && ctr.CursorX == 8) && string(gra.Boardgame[0][0]) == gra.PLAYER && string(gra.Boardgame[2][4]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	}

	//Control de diagonal /
	// Si en la diagonal en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if (ctr.CursorY == 0 && ctr.CursorX == 8) && string(gra.Boardgame[2][4]) == gra.PLAYER && string(gra.Boardgame[4][0]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (ctr.CursorY == 2 && ctr.CursorX == 4) && string(gra.Boardgame[0][8]) == gra.PLAYER && string(gra.Boardgame[4][0]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	} else if (ctr.CursorY == 4 && ctr.CursorX == 0) && string(gra.Boardgame[0][8]) == gra.PLAYER && string(gra.Boardgame[2][4]) == gra.PLAYER {
		winnerIs = "Player Wins!!"
		return
	}

}

func deleteCrusor() { // Borro el cursor en su posición actual por moverse a una nueva posición.
	gra.Boardgame[ctr.CursorY] = fmt.Sprintf("%v%v%v", gra.Boardgame[ctr.CursorY][:ctr.CursorX], gra.CursorInv, gra.Boardgame[ctr.CursorY][ctr.CursorX+1:])
	clearScreen()
	renderBoardgame()
}

func render() {
	for {
		renderBoardgame()
		time.Sleep(time.Millisecond * BLINK_TIME)
		gra.Boardgame[ctr.CursorY] = fmt.Sprintf("%v%v%v", gra.Boardgame[ctr.CursorY][:ctr.CursorX], gra.Cursor, gra.Boardgame[ctr.CursorY][ctr.CursorX+1:])
		clearScreen()

		renderBoardgame()
		time.Sleep(time.Millisecond * BLINK_TIME)
		gra.Boardgame[ctr.CursorY] = fmt.Sprintf("%v%v%v", gra.Boardgame[ctr.CursorY][:ctr.CursorX], gra.CursorInv, gra.Boardgame[ctr.CursorY][ctr.CursorX+1:])
		clearScreen()
	}
}

func renderBoardgame() {
	for i := 0; i < len(gra.Boardgame); i++ {
		fmt.Print(gra.MARGIN)
		fmt.Println(gra.Boardgame[i])
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
