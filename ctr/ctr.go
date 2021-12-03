package ctr

// Este package contiene los controladores del juego.

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var BoardGameInitial = []string{
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
}

var BoardGame = []string{"", "", "", "", ""}

const MARGIN = "                "
const PLAYER = "O"
const COMPUTER = "X"
const EMPTY = " "
const CURSOR = "*"
const TURN_PLAYER = "Turn: Player"
const TURN_COMPUTER = "Turn: Computer"
const DRAW = "It's a Draw!"
const PLAY_AGAIN = "Play Again? (y/n)"
const COMPUTER_THINK_TIME = 900

var CursorY = 2
var CursorX = 4
var Status = TURN_PLAYER
var remainingTurns = 8 // Maxima cantidad de veces que se puede cambiar de turno antes de un empate.
var playerWins = 0
var computerWins = 0

// var lastCompY = ""
// var lastCompX = ""

const BLINK_TIME = 250

func BlinkCursor() {
	ResetBoardGame()
	rand.Seed(time.Now().UnixNano()) // Preparo la variable rand para generar posteriormente un número random
	for {
		if Status == TURN_PLAYER {
			ShowCursor(true)
			time.Sleep(time.Millisecond * BLINK_TIME)
		}
		ShowCursor(false)
		time.Sleep(time.Millisecond * BLINK_TIME)
	}
}

func ShowCursor(show bool) { // Renderizo el tablero de juego junto con información del juego.
	ClearTerminal()
	fmt.Println("Player Wins:  ", playerWins)
	fmt.Println("Computer Wins:", computerWins)
	fmt.Println()
	for y := 0; y < len(BoardGame); y++ {
		fmt.Print(MARGIN)
		// Si es el turno del jugador, y se esta por renderizar la fila donde esta el cursor, y el cursor debe mostarse, renderizar esa fila con el cursor
		if Status == TURN_PLAYER && y == CursorY && show {
			rowWithCursor := fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], CURSOR, BoardGame[CursorY][CursorX+1:])
			fmt.Println(rowWithCursor)
		} else {
			fmt.Println(BoardGame[y]) // Renderizar sin el cursor.
		}
	}
	fmt.Println()

	fmt.Println(Status)
	if Status != TURN_PLAYER && Status != TURN_COMPUTER { // Si alguien ganó, o si hubo empate
		fmt.Println(PLAY_AGAIN)
	}

}

func ClearTerminal() {
	if strings.Contains(runtime.GOOS, "windows") { // Limpia la terminal en Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") // Limpiar la terminal en Linux o Mac
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
		Status = "PLAYER WINS!!"
		playerWins++
	} else {
		Status = "COMPUTER WINS!!"
		computerWins++
	}
}

func computerElection() {
	time.Sleep(time.Millisecond * COMPUTER_THINK_TIME) // Genero una demora para simular el pensar de la computadora

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
		if string(BoardGame[y][x]) == EMPTY {
			// lastCompX = strconv.Itoa(x)
			// lastCompY = strconv.Itoa(y)
			BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
			computerWinControl(y, x)
			ShowCursor(true)
			break
		}

	}
}

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

func changeTurn() { // Cambio el mensaje de a quien le toca jugar por no haberse jenerado un ganador en el anterior turno
	if remainingTurns == 0 {
		Status = DRAW
	} else if Status == TURN_PLAYER {
		Status = TURN_COMPUTER
		go computerElection() // Ejecuto esta funcion como una co-routine solo para que no aparezcan en la terminal las posibles teclas presionadas durante el turno de la computadora.
	} else {
		Status = TURN_PLAYER
	}
	remainingTurns--
}

func ResetBoardGame() {
	copy(BoardGame, BoardGameInitial)
	CursorY = 2
	CursorX = 4
	Status = TURN_PLAYER
	remainingTurns = 8
}
