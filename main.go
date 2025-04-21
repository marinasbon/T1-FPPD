// main.go - Loop principal do jogo
package main

import (
	"os"
	"sync"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	var jogo Jogo

	// Inicializa o jogo
	jogo = jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	
	// IMPLEMENTAÇÕES INIMIGO E MOEDAS
	// Mutex global usado para moedas e inimigos
	var mutex sync.Mutex

	// Inicia gerador de moedas (amarela, laranja e vermelha)
	iniciarGeradorDeMoedas(&jogo, &mutex)
	
	// Inicia movimento para cada inimigo que no mapa
	for y := 0; y < len(jogo.Mapa); y++ {
		for x := 0; x < len(jogo.Mapa[y]); x++ {
			if jogo.Mapa[y][x] == Inimigo {
				inimigoMoverInteligente(&jogo, x, y, &mutex)
			} else if jogo.Mapa[y][x] == InimigoTatico {
				iniciarInimigoTaticoZonaFixa(&jogo, &mutex)
			}
		}
	}	
	// IMPLEMENTAÇÕES INIMIGO E MOEDAS


	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}