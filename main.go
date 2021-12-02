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
	go ctr.Render()                 // Ejecuto el proceso de renderizado del jugo como una co-routina

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
					ctr.DeleteCrusor()                                                // Borro el cursor en su posición actual por moverse a una nueva posición.
					ctr.CursorInv = string(ctr.BoardGame[ctr.CursorY][ctr.CursorX+4]) // Almaceno el character que ya esta en el lugar adonde se va a colocar el cursor. Lo convertirlo a string porque es un byte.
					ctr.CursorX = ctr.CursorX + 4                                     // Actualizo la posición del cursor segun la tecla presionada.
				}
			case keyboard.KeyArrowLeft:
				if ctr.CursorX > 0 {
					ctr.DeleteCrusor()
					ctr.CursorInv = string(ctr.BoardGame[ctr.CursorY][ctr.CursorX-4])
					ctr.CursorX = ctr.CursorX - 4
				}
			case keyboard.KeyArrowDown:
				if ctr.CursorY < 4 {
					ctr.DeleteCrusor()
					ctr.CursorInv = string(ctr.BoardGame[ctr.CursorY+2][ctr.CursorX])
					ctr.CursorY = ctr.CursorY + 2
				}
			case keyboard.KeyArrowUp:
				if ctr.CursorY > 0 {
					ctr.DeleteCrusor()
					ctr.CursorInv = string(ctr.BoardGame[ctr.CursorY-2][ctr.CursorX])
					ctr.CursorY = ctr.CursorY - 2
				}
			case keyboard.KeySpace: // Si el jugador quiere poner su pieza y si esta vacio el lugar que ocupa el cursor, la pieza del jugador es colocada.
				if ctr.CursorInv == " " {
					ctr.PlayerWinControl()
					ctr.CursorInv = ctr.PLAYER
				}
			default:
				fmt.Print(string(char))
			}
		}
	}
}
