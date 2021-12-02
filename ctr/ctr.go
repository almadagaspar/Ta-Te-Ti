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

var BoardGame = []string{
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
	"──┼───┼──",
	"  |   |  ",
}

// var places = map[int]string{
// 	1: "00",
// 	2: "04",
// 	3: "08",
// 	4: "20",
// 	5: "24",
// 	6: "28",
// 	7: "40",
// 	8: "42",
// 	9: "48",
// }

const MARGIN = "            "
const PLAYER = "O"
const COMPUTER = "X"

var Cursor = "*"
var CursorInv = string(BoardGame[CursorY][CursorX]) // Guardo el character que se debe mostrar cuando el cursor se hace invisible al parpadear.

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
	if strings.Count(BoardGame[CursorY], PLAYER) == 2 {
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
	if count == 2 {
		showWinner(PLAYER)
		return
	}

	// Si en la DIAGONAL (\) en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if ((CursorY == 0 && CursorX == 0) && string(BoardGame[2][4]) == PLAYER && string(BoardGame[4][8]) == PLAYER) ||
		((CursorY == 2 && CursorX == 4) && string(BoardGame[0][0]) == PLAYER && string(BoardGame[4][8]) == PLAYER) ||
		((CursorY == 4 && CursorX == 8) && string(BoardGame[0][0]) == PLAYER && string(BoardGame[2][4]) == PLAYER) {
		showWinner(PLAYER)
		return
	}

	// Si en la DIAGONAL (/) en la que voy a poner una pieza ya hay dos piezas puesas, se gana la partida.
	if ((CursorY == 0 && CursorX == 8) && string(BoardGame[2][4]) == PLAYER && string(BoardGame[4][0]) == PLAYER) ||
		((CursorY == 2 && CursorX == 4) && string(BoardGame[0][8]) == PLAYER && string(BoardGame[4][0]) == PLAYER) ||
		((CursorY == 4 && CursorX == 0) && string(BoardGame[0][8]) == PLAYER && string(BoardGame[2][4]) == PLAYER) {
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

func DeleteCrusor() { // Borro el cursor en su posición actual por estar por moverse a una nueva posición.
	BoardGame[CursorY] = fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], CursorInv, BoardGame[CursorY][CursorX+1:])
}

func Render() {
	rand.Seed(time.Now().UnixNano()) // Preparo la variable rand para que pueda generar numeros randoms.

	for {

		if Message == TURN_PLAYER { // El cursor solo debe parpadear si es el turno del jugador
			RenderBoardgame()
			time.Sleep(time.Millisecond * BLINK_TIME)                                                                        // Genero una demora para que se vea el cursor
			BoardGame[CursorY] = fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], Cursor, BoardGame[CursorY][CursorX+1:]) // Actualizo la posicion del cursor segun su ultima posición definida por el jugador.
			ClearScreen()
		}

		RenderBoardgame()
		time.Sleep(time.Millisecond * BLINK_TIME)                                                                           // Genero una demora para que se vea lo que esta debajo del cursor
		BoardGame[CursorY] = fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], CursorInv, BoardGame[CursorY][CursorX+1:]) // Actualizo la posicion del cursor segun su ultima posición definida por el jugador.
		ClearScreen()

	}
}

func RenderBoardgame() { // Renderizo el tablero de juego junto con información del juego.
	for i := 0; i < len(BoardGame); i++ {
		fmt.Print(MARGIN)
		fmt.Println(BoardGame[i])
	}
	fmt.Println(Message, lastCompY, lastCompX)
}

func ClearScreen() {
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

		// ENTRA A ESTE IF SIN QUE DEBA HACERLO
		if string(BoardGame[y][x]) != string(PLAYER) &&
			string(BoardGame[y][x]) != string(COMPUTER) &&
			string(BoardGame[y][x]) != string(Cursor) { // Poner != " " lleva que la computadora pueda elegir el ultimo lugar que eligio el jugador.
			// lastCompX = strconv.Itoa(x)
			// lastCompY = strconv.Itoa(y)
			BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
			ClearScreen()
			RenderBoardgame()
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
	// 1º columna OK
	// 2º columna OK 3
	// 3º columna OK 3

	// 1º fila    OK 2
	// 2º fila    OK 3
	// 3º fila    OK 2

	// Diagonal /  OK
	// Diagonal \  OK
}
