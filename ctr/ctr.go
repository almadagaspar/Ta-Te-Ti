package ctr // Controladores

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
//    0 | 4 | 8   ─>  x index
//    ──┼───┼──
//    2 |   |
//    ──┼───┼──
//    4 |   |
//
//    |
//    v
//
//    y index
//
//

var BoardGame = []string{"", "", "", "", ""}     // Tablero de juego para el desarrollo de cada partida.
var BoardGameTest = []string{"", "", "", "", ""} // Tablero de juego para testear si la computadora o el jugador pueden ganar al poner su próxima ficha.

const PLAYER = "O"
const COMPUTER = "X"
const EMPTY = " "
const CURSOR = "*"
const PLAY_AGAIN = "Play Again? (y/n)"
const COMPUTER_THINK_TIME = 900
const EASY = "EASY"
const HARD = "HARD"
const BLINK_TIME = 250

// Estados en los que puede estar el juego.
const TURN_PLAYER = "Turn: Player"
const TURN_COMPUTER = "Turn: Computer"
const DRAW = "IT'S A DRAW!"
const MENU = "Menu"

var Status = MENU
var CursorY = 2
var CursorX = 4
var remainingTurns = 8 // Turnos restantes antes de un empate.
var PlayerWins = 0
var ComputerWins = 0
var Draws = 0

var Difficuty = ""

var next_y_1 int = -1 // La posición de y en una conveniente futura jugada para la computadora.
var next_x_1 int = -1 // La posición de x en una conveniente futura jugada para la computadora.

func RunGame() {
	ResetBoardGame()
	rand.Seed(time.Now().UnixNano()) // Preparo la variable rand para generar posteriormente un número random.
	for {
		if Status == TURN_PLAYER {
			RenderGame(true)
			time.Sleep(time.Millisecond * BLINK_TIME)
		}
		RenderGame(false)
		time.Sleep(time.Millisecond * BLINK_TIME)
		if Status == MENU {
			break
		}
	}
}

//
//

func RenderGame(show bool) { // Renderizo el tablero junto con información del juego.
	ClearTerminal()
	fmt.Println("============================================================================================")
	fmt.Println()
	fmt.Println("Controls:      Use the ARROWS KEYS to move the cursor, and the SPACE BAR to select a place.")
	fmt.Println()
	fmt.Println("Difficulty:   ", Difficuty)
	fmt.Println("Player Wins:  ", PlayerWins)
	fmt.Println("Computer Wins:", ComputerWins)
	fmt.Println("Draws:  ", Draws)
	for y := 0; y < len(BoardGame); y++ {
		// Si es el turno del jugador, y se esta por renderizar la fila donde esta el cursor, y el cursor debe mostarse, renderizar esa fila con el cursor.
		if Status == TURN_PLAYER && y == CursorY && show {
			rowWithCursor := fmt.Sprintf("%s%s%s", BoardGame[CursorY][:CursorX], CURSOR, BoardGame[CursorY][CursorX+1:])
			fmt.Println("\t\t\t\t", rowWithCursor)
		} else {
			fmt.Println("\t\t\t\t", BoardGame[y]) // Renderizar sin el cursor.
		}
	}
	fmt.Println()
	fmt.Println(Status)
	if Status != TURN_PLAYER && Status != TURN_COMPUTER { // Si alguien ganó, o si hubo empate...
		fmt.Println(PLAY_AGAIN)
	} else {
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("============================================================================================")
}

//
//

func ClearTerminal() {
	if strings.Contains(runtime.GOOS, "windows") { // Limpiar la terminal en Windows.
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") // Limpiar la terminal en Linux o Mac.
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//
//

func HideTerminalCursor(option string) { // Oculto el cursor de la terminal.
	cmd := exec.Command("tput", option)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//
//

func ShowWinner(winner string) { // Muestro un mensaje indicando quen ganó la partida.
	if winner == PLAYER {
		Status = "PLAYER WINS!!"
		PlayerWins++
	} else if winner == COMPUTER {
		Status = "COMPUTER WINS!!"
		ComputerWins++
	} else {
		Status = DRAW
		Draws++
	}
}

//
//

func computerElection() {
	time.Sleep(time.Millisecond * COMPUTER_THINK_TIME) // Genero una demora para simular el pensar de la computadora

	if Difficuty == EASY {
		for {
			// Genero números random para la posición de la próxima ficha de la computadora, y adapto esos números a los indices de mi tablero de juego.
			y := rand.Intn(3) * 2
			x := rand.Intn(3) * 4

			if string(BoardGame[y][x]) == EMPTY {
				BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:])
				if !WinControl(COMPUTER, y, x, BoardGame) { // Si con la última pieza que puso la computadora NO ganó la partida...
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
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda.
				next_y_1 = 4
				next_x_1 = 0
			} else if string(BoardGame[2][0]) == PLAYER { // Si el jugador puso su primer ficha en la segunda fila y la primer columna...
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda.
				next_y_1 = 0
				next_x_1 = 8
			} else if string(BoardGame[2][8]) == PLAYER { // Si el jugador puso su primer ficha en la segunda fila y la tercer columna...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha.
				next_y_1 = 4
				next_x_1 = 0
			} else if string(BoardGame[4][4]) == PLAYER { // Si el jugador puso su primer ficha en la tercer fila y la segunda columna...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha.
				next_y_1 = 0
				next_x_1 = 8
			} else if string(BoardGame[0][0]) == PLAYER { // Si el jugador puso su primer ficha en la esquina superior izquierda...
				BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina inferior derecha.
			} else if string(BoardGame[0][8]) == PLAYER { // Si el jugador puso su primer ficha en la esquina superior derecha...
				BoardGame[4] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[4][1:]) // ...la cumputadora pone su segunda ficha en la esquina inferior izquierda.
			} else if string(BoardGame[4][0]) == PLAYER { // Si el jugador puso su primer ficha en la esquina inferior izquierda...
				BoardGame[0] = fmt.Sprintf("%s%s", BoardGame[0][:8], COMPUTER) // ...la cumputadora pone su segunda ficha en la esquina superior derecha.
			} else if string(BoardGame[4][8]) == PLAYER { // Si el jugador puso su primer ficha en la esquina inferior derecha...
				BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:]) // ...la cumputadora pone su segunda ficha en la esquina superior izquierda.
			}
			ChangeTurn()
		case 4: // Tercer turno de la computadora. Ya puede intentar ganar.
			if next_y_1 != -1 { // Si el jugador NO puso su primer ficha en una esquina...
				if !searchForBestPlay(COMPUTER) { // Si la computadora NO encontro una jugada que le de la victoria....
					BoardGame[next_y_1] = fmt.Sprintf("%s%s%s", BoardGame[next_y_1][:next_x_1], COMPUTER, BoardGame[next_y_1][next_x_1+1:]) // Ubica su proxima ficha en el mejor lugar predefinido, para intentar ganar en el próximo turno.
					ChangeTurn()
				} else {
					ShowWinner(COMPUTER)
				}
			} else { // Si el jugador puso su primer ficha en una esquina, la computadora busca si el jugador puede ganar en su próximo turno, y lo impide si es asi.
				if !searchForBestPlay(PLAYER) { // Si el jugador no tuvo la oportunidad de ganar en su anterior turno...
					// Segun donde puso el jugador sus dos primeras fichas, la computadora ubicará su próxima ficha en el mejor lugar para intentar ganar en el próximo turno.
					if (string(BoardGame[0][0]) == PLAYER && string(BoardGame[2][8]) == PLAYER) ||
						(string(BoardGame[0][4]) == PLAYER && string(BoardGame[4][8]) == PLAYER) {
						BoardGame[4] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[4][1:])
					} else if (string(BoardGame[0][0]) == PLAYER && string(BoardGame[4][4]) == PLAYER) ||
						(string(BoardGame[2][0]) == PLAYER && string(BoardGame[4][8]) == PLAYER) {
						BoardGame[0] = fmt.Sprintf("%s%s", BoardGame[0][:8], COMPUTER)
					} else if (string(BoardGame[0][8]) == PLAYER && string(BoardGame[2][0]) == PLAYER) ||
						(string(BoardGame[0][4]) == PLAYER && string(BoardGame[4][0]) == PLAYER) {
						BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER)
					} else if (string(BoardGame[0][8]) == PLAYER && string(BoardGame[4][4]) == PLAYER) ||
						(string(BoardGame[2][8]) == PLAYER && string(BoardGame[4][0]) == PLAYER) {
						BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:])
					}
				}
				ChangeTurn()
			}
		case 2: // Cuarto turno de la computadora.
			if !searchForBestPlay(COMPUTER) {
				if string(BoardGame[0][4]) == EMPTY {
					BoardGame[0] = fmt.Sprintf("%s%s%s", BoardGame[0][:4], COMPUTER, BoardGame[0][4+1:])
				} else {
					BoardGame[2] = fmt.Sprintf("%s%s%s", BoardGame[2][:0], COMPUTER, BoardGame[2][0+1:])
				}
				ChangeTurn()
			} else {
				ShowWinner(COMPUTER)
			}
		case 0: // Quinto y último turno de la computadora.
			if !searchForBestPlay(COMPUTER) {
				if string(BoardGame[0][0]) == EMPTY { // La computadora busca la última esquina vacia para poner su ficha.
					BoardGame[0] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[0][1:])
				} else if string(BoardGame[0][8]) == EMPTY {
					BoardGame[0] = fmt.Sprintf("%s%s", BoardGame[0][:8], COMPUTER)
				} else if string(BoardGame[4][0]) == EMPTY {
					BoardGame[4] = fmt.Sprintf("%s%s", COMPUTER, BoardGame[4][1:])
				} else if string(BoardGame[4][8]) == EMPTY {
					BoardGame[4] = fmt.Sprintf("%s%s", BoardGame[4][:8], COMPUTER)
				}
				ChangeTurn()
			} else {
				ShowWinner(COMPUTER)
			}
		default:
		}
	}
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
					if WinControl(PLAYER, y, x, BoardGameTest) { // Si con la última pieza que se SIMULÓ poner, el jugador ganaría la partida...
						foundVictory = true
						BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:]) // ...la computadora pone su ficha en ese lugar para impedir perder.
					}
				}
				if CorP == COMPUTER {
					if WinControl(COMPUTER, y, x, BoardGameTest) { // Si con la última pieza que se SIMULÓ poner, la computadora ganaría la partida...
						foundVictory = true
						BoardGame[y] = fmt.Sprintf("%s%s%s", BoardGame[y][:x], COMPUTER, BoardGame[y][x+1:]) // ...la computadora pone su ficha en ese lugar.
					}
				}
			}
		}
	}
	return foundVictory
}

//
//

func WinControl(CorP string, y int, x int, boardGame []string) bool { // CorP = COMPUTER o PLAYER. Este función controla si la computadora o el jugador ganaron en su última jugada, o en una jugada de prueba.
	// Se busca en cada fila.
	if strings.Count(boardGame[y], CorP) == 3 {
		return true
	}

	// Se busca en cada columna.
	count := 0
	for y := 0; y <= len(boardGame); y = y + 2 {
		if string(boardGame[y][x]) == CorP {
			count++
		}
	}
	if count == 3 {
		return true
	}

	// Se busca en cada diagonal.
	if (string(boardGame[0][0]) == CorP && string(boardGame[2][4]) == CorP && string(boardGame[4][8]) == CorP) || // Diagonal \
		(string(boardGame[0][8]) == CorP && string(boardGame[2][4]) == CorP && string(boardGame[4][0]) == CorP) { // Diagonal /
		return true
	}
	return false
}

//
//

func ChangeTurn() { // Cambio el mensaje de a quién le toca jugar por no haberse generado un ganador en el anterior turno.
	if remainingTurns == 0 {
		ShowWinner(DRAW)
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
	if Difficuty == EASY {
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

//
//

func MenuScreen() {
	ClearTerminal()
	fmt.Println("============================================================================================")
	fmt.Println()
	fmt.Println("\t\t\t\t--------------")
	fmt.Println("\t\t\t\t-  Ta-Te-Ti  -")
	fmt.Println("\t\t\t\t--------------")
	fmt.Println()
	fmt.Println()
	fmt.Println("\t\t\t    Select a Difficulty:")
	fmt.Println("\t\t\t     1- EASY")
	fmt.Println("\t\t\t     2- HARD")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("  Press ESC to quit.")
	fmt.Println("  Please, adjust the size of your terminal to ensure that you see the entire playing area.")
	fmt.Println("  Developed by Gaspar Almada - 2021")
	fmt.Println()
	fmt.Println("============================================================================================")
}
