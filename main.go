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

	ctr.ClearScreen()               // Limpio la terminal
	ctr.HideTerminalCursor("civis") // Desactivo el cursor de la terminal
	go ctr.BlinkCursor()            // Ejecuto como una co-routina el parpadeo del cursor del juego

	for {
		// Registramos que tecla presiono el usuario.
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}

		// Detectamos si el usuario presiono la tecla para salir del programa:
		if key == keyboard.KeyEsc {
			break
		}

		if ctr.Message == ctr.TURN_PLAYER {
			switch key { // Realizo una acción segun la tecla presionada por el jugador.
			case keyboard.KeyArrowRight:
				if ctr.CursorX < 8 { // Impido que el cursor se salga del tablero.
					ctr.CursorX = ctr.CursorX + 4 // Modifico la posición que tendrá el cursor segun la tecla presionada.
					ctr.ShowCursor(ctr.CursorY)   // Muestro el cursor en la nueva posición.
				}
			case keyboard.KeyArrowLeft:
				if ctr.CursorX > 0 {
					ctr.CursorX = ctr.CursorX - 4
					ctr.ShowCursor(ctr.CursorY)
				}
			case keyboard.KeyArrowDown:
				if ctr.CursorY < 4 {
					ctr.CursorY = ctr.CursorY + 2
					ctr.ShowCursor(ctr.CursorY)
				}
			case keyboard.KeyArrowUp:
				if ctr.CursorY > 0 {
					ctr.CursorY = ctr.CursorY - 2
					ctr.ShowCursor(ctr.CursorY)
				}
			case keyboard.KeySpace: // Si el jugador quiere poner su pieza y si esta vacio el lugar que ocupa el cursor, la pieza del jugador es colocada.
				if string(ctr.BoardGame[ctr.CursorY][ctr.CursorX]) == " " {
					ctr.BoardGame[ctr.CursorY] = fmt.Sprintf("%s%s%s", ctr.BoardGame[ctr.CursorY][:ctr.CursorX], ctr.PLAYER, ctr.BoardGame[ctr.CursorY][ctr.CursorX+1:])
					ctr.PlayerWinControl()
				}
			default:
				fmt.Print(string(char))
			}
		}
	}
}
