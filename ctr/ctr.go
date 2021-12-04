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

//
//    0 | 4 | 8   ─► x index
//    ──┼───┼──
//    2 |   |
//    ──┼───┼──
//    4 |   |
//
//    |
//    ▼
//
//    y index
//

var BoardGame = []string{"", "", "", "", ""}     // Tablero de juego para el desarrollo de cada partida.
var BoardGameTest = []string{"", "", "", "", ""} // Tablero de juego para testear si la computadora o el jugador puede ganar al poner su proxima ficha.

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
const EASY = "easy"
const HARD = "hard"

var CursorY = 2
var CursorX = 4
var Status = TURN_PLAYER
var remainingTurns = 8 // Maxima cantidad de veces que se puede cambiar de turno antes de un empate.
var playerWins = 0
var computerWins = 0

var difficuty = HARD

// var difficuty = EASY

var next_y_1 int = -1 // Posicion de y en una conveniente futura jugada
var next_x_1 int = -1 // Posicion de x en una conveniente futura jugada

const BLINK_TIME = 250

//
//

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

//
//

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

//
//

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

//
//

func HideTerminalCursor(option string) { // Oculto el cursor de la terminal
	cmd := exec.Command("tput", option)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//
//

func ShowWinner(winner string) { // Muestro un mensaje indicando quen gano la partida.
	if winner == PLAYER {
		Status = "PLAYER WINS!!"
		playerWins++
	} else {
		Status = "COMPUTER WINS!!"
		computerWins++
	}
}

//
//

func computerElection() {
	time.Sleep(time.Millisecond * COMPUTER_THINK_TIME) // Genero una demora para simular el pensar de la computadora

	if difficuty == EASY {
		var y int
		var x int
		for {
			// Genero números random para la posible próxima posición para la ficha de la computadora.
			var randNumY = rand.Intn(3)
			var randNumX = rand.Intn(3)

			// Adapto esos números a las posiciones de mi tablero de juego.
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

			if string(BoardGame[y][x]) == EMPTY {
				BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
				if !WinControl(COMPUTER, y, x, BoardGame) { // Si con la ultima pieza que puso la computadora NO ganó la partida...
					ChangeTurn()
				} else {
					ShowWinner(COMPUTER)
				}
				break
			}

		}
	} else { // Si la dificultad es HARD
		switch remainingTurns {
		case 8: // Primer turno de la computadora, pone su ficha en el centro.
			BoardGame[2] = fmt.Sprintf("%s%s%s", BoardGame[2][:4], COMPUTER, BoardGame[2][5:])
			ChangeTurn()
			return
		case 6: // Segundo turno de la computadora
			if string(BoardGame[0][4]) == PLAYER { // Si el jugador puso su primer ficha en la primer fila y segunda columna...
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda
				next_y_1 = 4
				next_x_1 = 0
			} else if string(BoardGame[2][0]) == PLAYER { // Si el jugador puso su primer ficha en la segunda fila y la primer columna...
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda
				next_y_1 = 0
				next_x_1 = 8
			} else if string(BoardGame[2][8]) == PLAYER { // Si el jugador puso su primer ficha en la segunda fila y la tercer columna...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha
				next_y_1 = 4
				next_x_1 = 0
			} else if string(BoardGame[4][4]) == PLAYER { // Si el jugador puso su primer ficha en la tercer fila y la segunda columna...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha
				next_y_1 = 0
				next_x_1 = 8
			} else if string(BoardGame[0][0]) == PLAYER { // Si el jugador puso su primer ficha en la esquina superior izquierda...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha
				// next_y_1 = 0
				// next_x_1 = 8
			} else if string(BoardGame[0][8]) == PLAYER { // Si el jugador puso su primer ficha en la esquina superior derecha...
				BoardGame[4] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[4][1:]) // ...la cumputadora pone su segunda ficha en la esquina inferior izquierda
				// next_y_1 = 0
				// next_x_1 = 8
			} else if string(BoardGame[4][0]) == PLAYER { // Si el jugador puso su primer ficha en la esquina inferior izquierda...
				BoardGame[0] = fmt.Sprintf("%s%s", BoardGame[0][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina superior derecha
				// next_y_1 = 0
				// next_x_1 = 8
			} else if string(BoardGame[4][8]) == PLAYER { // Si el jugador puso su primer ficha en la esquina inferior derecha...
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda
				// next_y_1 = 0
				// next_x_1 = 8
			}
			ChangeTurn()
		case 4: // Tercer turno de la computadora. Ya puede intentar ganar.
			if next_y_1 != -1 { // Si el jugador NO puso su primer ficha en una esquina...
				if !searchForBestPlay(COMPUTER) { // Si la computadora NO encontro una jugada que le de la victoria....
					BoardGame[next_y_1] = fmt.Sprintf("%s%s%s", BoardGame[next_y_1][:next_x_1], COMPUTER, BoardGame[next_y_1][next_x_1+1:]) // Ubica su proxima ficha en el mejor lugar predefinido, para intentar ganar en el proximo turno.
					ChangeTurn()
				} else {
					ShowWinner(COMPUTER)
				}
			} else { // Si el jugador puso su primer ficha en una esquina, la computadora busca si el jugador puede ganar en su proximo turno, y lo impide si es asi.
				searchForBestPlay(PLAYER)
				ChangeTurn()
			}
		case 2:
			if !searchForBestPlay(COMPUTER) {
				// LAS VARIABLES next DEBEN LLEGAR A ESTA INSTANCIA CON LOS VALORES CORRESPONDIENTES !!!!!!!!!!!!!!
				// BoardGame[next_y_1] = fmt.Sprintf("%s%s%s", BoardGame[next_y_1][:next_x_1], COMPUTER, BoardGame[next_y_1][next_x_1+1:])
				ChangeTurn()
			} else {
				ShowWinner(COMPUTER)
			}
		default:
		}
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}

//
//

func searchForBestPlay(CorP string) bool { // CorP = COMPUTER o PLAYER
	foundVictory := false
	for y := 0; y <= 4 && !foundVictory; y = y + 2 {
		for x := 0; x <= 8 && !foundVictory; x = x + 4 {
			copy(BoardGameTest, BoardGame)
			if string(BoardGameTest[y][x]) == EMPTY {
				BoardGameTest[y] = fmt.Sprintf("%s%s%s", BoardGameTest[y][:x], CorP, BoardGameTest[y][x+1:])
				if CorP == PLAYER {
					if WinControl(PLAYER, y, x, BoardGameTest) { // Si con la ultima pieza que se SIMULO poner la computadora ganaría la partida...
						foundVictory = true
						BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
					}
				}
				if CorP == COMPUTER {
					if WinControl(COMPUTER, y, x, BoardGameTest) { // Si con la ultima pieza que se SIMULO poner la computadora ganaría la partida...
						foundVictory = true
						BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
					}
				}
			}
		}
	}
	return foundVictory
}

// PlayerWinControl y computerWinControl PODRIAN FUCIONARSE EN UNA UNICA FUNCIÓN AGREGANDO UN PARAMETRO Who... !!!!!!!!!!!!!!!!
//

func WinControl(CorP string, y int, x int, boardGame []string) bool { // CorP = COMPUTER o PLAYER
	if strings.Count(boardGame[y], CorP) == 3 {
		return true
	}

	//  Si en la COLUMNA en la que voy a poner una pieza ya hay dos piezas puestas, se gana la partida.
	count := 0
	for y := 0; y <= len(boardGame); y = y + 2 {
		if string(boardGame[y][x]) == CorP {
			count++
		}
	}
	if count == 3 {
		return true
	}

	// Si en la DIAGONALES \ o / hay tres piezas del jugador puesas, se gana la partida.
	if (string(boardGame[0][0]) == CorP && string(boardGame[2][4]) == CorP && string(boardGame[4][8]) == CorP) || // Diagonal \
		(string(boardGame[0][8]) == CorP && string(boardGame[2][4]) == CorP && string(boardGame[4][0]) == CorP) { // Diagonal /
		return true
	}
	return false
}

//
//

func ChangeTurn() { // Cambio el mensaje de a quien le toca jugar por no haberse jenerado un ganador en el anterior turno
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

//
//

func ResetBoardGame() {
	copy(BoardGame, BoardGameInitial)

	CursorX = 4
	if difficuty == EASY {
		Status = TURN_PLAYER
		CursorY = 2
	} else { // Si la dificultad es HARD
		CursorY = 0
		next_y_1 = -1
		next_x_1 = -1
		Status = TURN_COMPUTER
		go computerElection() // Ejecuto esta funcion como una co-routine solo para que no aparezcan en la terminal las posibles teclas presionadas durante el turno de la computadora.
	}
	remainingTurns = 8
}
