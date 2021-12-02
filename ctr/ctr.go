package ctr

// Este package contiene los controladores del juego.

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var BoardGame = []string{
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
}

const MARGIN = "            "
const PLAYER = "O"
const COMPUTER = "X"

const CURSOR = "*"

var CursorY = 2
var CursorX = 4

const TURN_PLAYER = "Turn: Player"
const TURN_COMPUTER = "Turn: Computer"

var Message = TURN_PLAYER
var lastCompY = ""
var lastCompX = ""

const BLINK_TIME = 300

func PlayerWinControl() { // Controlo si el Jugador o la Computadora (PLAYER = PLAYER o COMPUTER) ganó tras su ultima jugada.

	// Si en la FILA en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if strings.Count(BoardGame[CursorY], PLAYER) == 3 {
		showWinner(PLAYER)
		return
	}

	//  Si en la COLUMNA en la que voy a poner una pieza ya hay dos piezas puestas, se gana la partida.
	count := 0
	for y := 0; y <= len(BoardGame); y = y + 2 {
		if string(BoardGame[y][CursorX]) == PLAYER {
			count++
		}
	}
	if count == 3 {
		showWinner(PLAYER)
		return
	}

	// Si en la DIAGONALES \ o / hay tres piezas del jugador puesas, se gana la partida.
	if (string(BoardGame[0][0]) == PLAYER && string(BoardGame[2][4]) == PLAYER && string(BoardGame[4][8]) == PLAYER) ||
		(string(BoardGame[0][8]) == PLAYER && string(BoardGame[2][4]) == PLAYER && string(BoardGame[4][0]) == PLAYER) {
		showWinner(PLAYER)
		return
	}
	changeTurn()
}

func changeTurn() { // Cambio el mensaje de a quien le toca jugar por no haberse jenerado un ganador en el anterior turno
	if Message == TURN_PLAYER {
		Message = TURN_COMPUTER
		computerElection()
	} else {
		Message = TURN_PLAYER
	}
}

func BlinkCursor() {
	for {
		if Message == TURN_PLAYER {
			ShowCursor(true)
			time.Sleep(time.Millisecond * BLINK_TIME)
		}
		ShowCursor(false)
		time.Sleep(time.Millisecond * BLINK_TIME)
	}
}

func ShowCursor(show bool) { // Renderizo el tablero de juego junto con información del juego.
	ClearScreen()
	for y := 0; y < len(BoardGame); y++ {
		fmt.Print(MARGIN)
		// Si es el turno del jugador, y se esta por renderizar la fila donde esta el cursor, y el cursor debe mostarse, renderizo el tablero con el cursor, si no, sin él.
		if Message == TURN_PLAYER && y == CursorY && show {
			rowWithCursor := fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], CURSOR, BoardGame[CursorY][CursorX+1:])
			fmt.Println(rowWithCursor)
		} else {
			fmt.Println(BoardGame[y])
		}
	}
	fmt.Println(Message, lastCompY, lastCompX)
}

func ClearScreen() {
	if strings.Contains(runtime.GOOS, "windows") { // windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") // linux o mac
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func HideTerminalCursor(option string) { // Oculto el cursor de la terminal
	cmd := exec.Command("tput", option)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showWinner(winner string) { // Muestro un mensaje indicando quen gano la partida.
	if winner == PLAYER {
		Message = "Player Wins!!"
	} else {
		Message = "Computer Wins!!"
	}
}

func computerElection() {
	time.Sleep(time.Millisecond * BLINK_TIME * 5) // Genero una demora para simular el pensar de la computadora

	var y int
	var x int

	for {

		var randNumY = rand.Intn(3)
		var randNumX = rand.Intn(3)

		if randNumY == 0 {
			y = 0
		} else if randNumY == 1 {
			y = 2
		} else if randNumY == 2 {
			y = 4
		}

		if randNumX == 0 {
			x = 0
		} else if randNumX == 1 {
			x = 4
		} else if randNumX == 2 {
			x = 8
		}

		// fmt.Println(y, x)
		if string(BoardGame[y][x]) == " " {
			lastCompX = strconv.Itoa(x)
			lastCompY = strconv.Itoa(y)
			BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
			ClearScreen()
			ShowCursor(true)
			computerWinControl(y, x)
			break
		}

	}
}

func computerWinControl(y, x int) {
	if strings.Count(BoardGame[y], COMPUTER) == 3 {
		showWinner(COMPUTER)
		return
	}

	count := 0
	for i := 0; i <= len(BoardGame); i = i + 2 {
		if string(BoardGame[i][x]) == COMPUTER {
			count++
		}
	}
	if count == 3 {
		showWinner(COMPUTER)
		return
	}

	if (string(BoardGame[0][0]) == COMPUTER && string(BoardGame[2][4]) == COMPUTER && string(BoardGame[4][8]) == COMPUTER) || // Diagonal \
		(string(BoardGame[0][8]) == COMPUTER && string(BoardGame[2][4]) == COMPUTER && string(BoardGame[4][0]) == COMPUTER) { // Diagonal /
		showWinner(COMPUTER)
		return
	}

	changeTurn()

}
