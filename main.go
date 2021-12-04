package main

import (
	"fmt"
	"log"
	"myapp/ctr"

	"github.com/eiannone/keyboard"
)

func main() {
	err := keyboard.Open() // Activo la funcionalidad que lee las teclas presionadas.
	if err != nil {
		log.Fatal(err)
	}

	defer func() { // Me aseguro que se desactive la funcionalidad que lee las teclas presionadas y que se reactive el cursor de la terminal al finalizar esta aplicación.
		ctr.HideTerminalCursor("cvvis")
		_ = keyboard.Close()
	}()

	ctr.ClearTerminal()             // Limpio la terminal
	ctr.HideTerminalCursor("civis") // Desactivo el cursor de la terminal
	go ctr.BlinkCursor()            // Ejecuto como una co-routina el parpadeo del cursor del juego

	for {
		// Registramos que tecla presiono el usuario.
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc { // Detectamos si el usuario presiono la tecla para salir del programa:
			break
		}

		if ctr.Status != ctr.TURN_PLAYER && ctr.Status != ctr.TURN_COMPUTER {
			if string(char) == "n" || string(char) == "N" {
				break
			} else if string(char) == "y" || string(char) == "Y" {
				ctr.ResetBoardGame()
			}
		}

		if ctr.Status == ctr.TURN_PLAYER { // Solo quiero que se mueva el cursor si es el turno del jugador.
			switch key { // Realizo una acción segun la tecla presionada por el jugador.
			case keyboard.KeyArrowRight:
				if ctr.CursorX < 8 { // Impido que el cursor se salga del tablero.
					ctr.CursorX = ctr.CursorX + 4 // Modifico la posición que tendrá el cursor segun la tecla presionada.
					ctr.ShowCursor(true)          // Muestro inmediatamente el cursor en la nueva posición.
				}
			case keyboard.KeyArrowLeft:
				if ctr.CursorX > 0 {
					ctr.CursorX = ctr.CursorX - 4
					ctr.ShowCursor(true)
				}
			case keyboard.KeyArrowDown:
				if ctr.CursorY < 4 {
					ctr.CursorY = ctr.CursorY + 2
					ctr.ShowCursor(true)
				}
			case keyboard.KeyArrowUp:
				if ctr.CursorY > 0 {
					ctr.CursorY = ctr.CursorY - 2
					ctr.ShowCursor(true)
				}
			case keyboard.KeySpace: // Si el jugador quiere poner su pieza y si esta vacio el lugar que ocupa el cursor, la pieza del jugador es colocada.
				if string(ctr.BoardGame[ctr.CursorY][ctr.CursorX]) == ctr.EMPTY {
					ctr.BoardGame[ctr.CursorY] = fmt.Sprintf("%s%s%s", ctr.BoardGame[ctr.CursorY][:ctr.CursorX], ctr.PLAYER, ctr.BoardGame[ctr.CursorY][ctr.CursorX+1:])
					if !ctr.WinControl(ctr.PLAYER, ctr.CursorY, ctr.CursorX, ctr.BoardGame) { // Si con la ultima ficha que puso el jugador NO gano...
						ctr.ChangeTurn()
					} else {
						ctr.ShowWinner(ctr.PLAYER)
					}
				}
			default:

			}
		}

	}
}
