package main

import (
	"math/rand"
	"sync"
	"time"
)

func iniciarInimigoTaticoZonaFixa(jogo *Jogo, mutex *sync.Mutex) {
	go func() {
		var zonaX, zonaY int
		encontrado := false

		// Detecta a posição original do símbolo ☣ e remove do mapa
		mutex.Lock()
		for y := range jogo.Mapa {
			for x := range jogo.Mapa[y] {
				if jogo.Mapa[y][x] == InimigoTatico {
					zonaX, zonaY = x, y
					jogo.Mapa[y][x] = Vazio
					encontrado = true
					break
				}
			}
			if encontrado {
				break
			}
		}
		mutex.Unlock()
		if !encontrado {
			return
		}

		// zona de patrulha
		minX := max(0, zonaX-5)
		maxX := min(len(jogo.Mapa[0])-1, zonaX+5)
		minY := max(0, zonaY-5)
		maxY := min(len(jogo.Mapa)-1, zonaY+5)

		// Encontra posição inicial vazia na zona
		var px, py int
		inicializado := false
		mutex.Lock()
		for y := minY; y <= maxY && !inicializado; y++ {
			for x := minX; x <= maxX; x++ {
				if jogo.Mapa[y][x] == Vazio {
					px, py = x, y
					jogo.Mapa[py][px] = InimigoTatico
					interfaceDesenharJogo(jogo)
					inicializado = true
					break
				}
			}
		}
		mutex.Unlock()
		if !inicializado {
			return
		}

		// Canais de comunicação
		alertaZona := make(chan bool)
		moedaNaZona := make(chan [2]int)

		// Monitor de presença do jogador
		go func() {
			for {
				time.Sleep(200 * time.Millisecond)
				mutex.Lock()
				naZona := jogo.PosX >= minX && jogo.PosX <= maxX &&
					jogo.PosY >= minY && jogo.PosY <= maxY
				mutex.Unlock()
				alertaZona <- naZona
			}
		}()

		// Monitor de moedas na zona
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				mutex.Lock()
				encontrada := false
				for y := minY; y <= maxY && !encontrada; y++ {
					for x := minX; x <= maxX; x++ {
						s := jogo.Mapa[y][x].simbolo
						if s == MoedaAmarela.simbolo || s == MoedaLaranja.simbolo || s == MoedaVermelhaVerde.simbolo || s == MoedaVermelhaNegativa.simbolo {
							moedaNaZona <- [2]int{x, y}
							encontrada = true
							break
						}
					}
				}
				mutex.Unlock()
			}
		}()

		// Estado do inimigo
		ultimo := Vazio
		perseguirJogador := false
		objetivoMoeda := [2]int{-1, -1}

		for {
			select {
			case jogadorNaZona := <-alertaZona:
				perseguirJogador = jogadorNaZona

			case posMoeda := <-moedaNaZona:
				objetivoMoeda = posMoeda

			case <-time.After(3 * time.Second):
				// Nenhum evento novo, mantém estado atual
			}

			mutex.Lock()
			if px == jogo.PosX && py == jogo.PosY {
				jogo.Pontos = 0
				jogo.StatusMsg = "Inimigo tático capturou você!"
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				time.Sleep(300 * time.Millisecond)
				continue
			}

			var dx, dy int

			if perseguirJogador {
				if jogo.PosX > px {
					dx = 1
				} else if jogo.PosX < px {
					dx = -1
				}
				if jogo.PosY > py {
					dy = 1
				} else if jogo.PosY < py {
					dy = -1
				}
			} else if objetivoMoeda != [2]int{-1, -1} {
				if objetivoMoeda[0] > px {
					dx = 1
				} else if objetivoMoeda[0] < px {
					dx = -1
				}
				if objetivoMoeda[1] > py {
					dy = 1
				} else if objetivoMoeda[1] < py {
					dy = -1
				}
			} else {
				direcoes := []struct{ dx, dy int }{
					{1, 0}, {-1, 0}, {0, 1}, {0, -1},
				}
				r := direcoes[rand.Intn(len(direcoes))]
				dx, dy = r.dx, r.dy
			}

			tx, ty := px+dx, py+dy
			dentroDaZona := tx >= minX && tx <= maxX && ty >= minY && ty <= maxY

			if jogoPodeMoverPara(jogo, tx, ty) &&
				jogo.Mapa[ty][tx] != Inimigo && jogo.Mapa[ty][tx] != InimigoTatico &&
				dentroDaZona {
				jogo.Mapa[py][px] = ultimo
				ultimo = jogo.Mapa[ty][tx]
				jogo.Mapa[ty][tx] = InimigoTatico
				px, py = tx, ty
				interfaceDesenharJogo(jogo)

				// Zera objetivo se alcançou
				if [2]int{tx, ty} == objetivoMoeda {
					objetivoMoeda = [2]int{-1, -1}
				}
			}

			mutex.Unlock()
			time.Sleep(300 * time.Millisecond)
		}
	}()
}
