package main

import (
	"sync"
	"time"
)

func inimigoMoverInteligente(jogo *Jogo, inimigoX, inimigoY int, mutex *sync.Mutex) {
	go func() {
		nx, ny := inimigoX, inimigoY
		ativo := true
		inicio := time.Now()
		ultimo := Vazio

		for {
			duracao := time.Since(inicio)
			if ativo && duracao > 14*time.Second {
				ativo = false
				inicio = time.Now()
			} else if !ativo && duracao > 3*time.Second {
				ativo = true
				inicio = time.Now()
			}

			if ativo {
				mutex.Lock()

				// Verifica se encostou no personagem
				if nx == jogo.PosX && ny == jogo.PosY {
					jogo.Pontos = 0
					jogo.StatusMsg = "Você foi pego pelo inimigo! Pontos zerados."
					interfaceDesenharJogo(jogo)
					mutex.Unlock()
					time.Sleep(300 * time.Millisecond)
					continue
				}

				dx, dy := 0, 0
				if jogo.PosX > nx {
					dx = 1
				} else if jogo.PosX < nx {
					dx = -1
				}
				if jogo.PosY > ny {
					dy = 1
				} else if jogo.PosY < ny {
					dy = -1
				}

				tx, ty := nx+dx, ny+dy
				if jogoPodeMoverPara(jogo, tx, ty) && jogo.Mapa[ty][tx] != Inimigo {
					jogo.Mapa[ny][nx] = ultimo
					ultimo = jogo.Mapa[ty][tx]
					jogo.Mapa[ty][tx] = Inimigo
					nx, ny = tx, ty
					interfaceDesenharJogo(jogo)
				}

				mutex.Unlock()
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()
}


