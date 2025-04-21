package main

import (
	"math/rand"
	"sync"
	"time"
)

func iniciarGeradorDeMoedas(jogo *Jogo, mutex *sync.Mutex) {
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(4)+2) * time.Second)
			x := rand.Intn(len(jogo.Mapa[0]))
			y := rand.Intn(len(jogo.Mapa))
			if jogo.Mapa[y][x] == Vazio {
				go iniciarMoedaAmarela(jogo, x, y, mutex)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(4)+2) * time.Second)
			x := rand.Intn(len(jogo.Mapa[0]))
			y := rand.Intn(len(jogo.Mapa))
			if jogo.Mapa[y][x] == Vazio {
				go iniciarMoedaLaranja(jogo, x, y, mutex)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(4)+2) * time.Second)
			x := rand.Intn(len(jogo.Mapa[0]))
			y := rand.Intn(len(jogo.Mapa))
			if jogo.Mapa[y][x] == Vazio {
				go iniciarMoedaVermelha(jogo, x, y, mutex)
			}
		}
	}()
}

func iniciarMoedaAmarela(jogo *Jogo, x, y int, mutex *sync.Mutex) {
	mutex.Lock()
	jogo.Mapa[y][x] = MoedaAmarela
	interfaceDesenharJogo(jogo)
	mutex.Unlock()

	go func(x, y int) {
		start := time.Now()
		for {
			time.Sleep(100 * time.Millisecond)
			mutex.Lock()
			if jogo.PosX == x && jogo.PosY == y {
				jogo.Pontos += 1
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				jogo.StatusMsg = "Coletou moeda amarela (+1)"
				mutex.Unlock()
				return
			}
			if time.Since(start) > 9*time.Second {
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				return
			}
			mutex.Unlock()
		}
	}(x, y)
}

func iniciarMoedaLaranja(jogo *Jogo, x, y int, mutex *sync.Mutex) {
	mutex.Lock()
	jogo.Mapa[y][x] = MoedaLaranja
	interfaceDesenharJogo(jogo)
	mutex.Unlock()

	go func(x, y int) {
		ativada := false
		start := time.Now()
		destruirEm := 15 * time.Second
		fugaAte := time.Time{}

		for {
			time.Sleep(200 * time.Millisecond)
			mutex.Lock()
			dist := abs(jogo.PosX-x) + abs(jogo.PosY-y)
			if dist <= 2 && !ativada {
				ativada = true
				fugaAte = time.Now().Add(5 * time.Second)
			}

			if ativada && time.Now().Before(fugaAte) {
				dx, dy := 0, 0
				if jogo.PosX < x {
					dx = 1
				} else if jogo.PosX > x {
					dx = -1
				}
				if jogo.PosY < y {
					dy = 1
				} else if jogo.PosY > y {
					dy = -1
				}
				nx, ny := x+dx, y+dy
				if jogoPodeMoverPara(jogo, nx, ny) {
					jogo.Mapa[y][x] = Vazio
					x, y = nx, ny
					jogo.Mapa[y][x] = MoedaLaranja
					interfaceDesenharJogo(jogo)					
				}
			}

			if jogo.PosX == x && jogo.PosY == y {
				jogo.Pontos += 5
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				jogo.StatusMsg = "Coletou moeda laranja (+5)"
				mutex.Unlock()
				return
			}

			if !ativada && time.Since(start) > destruirEm {
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				return
			}

			if ativada && time.Now().After(fugaAte) {
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				return
			}
			mutex.Unlock()
		}
	}(x, y)
}

func iniciarMoedaVermelha(jogo *Jogo, x, y int, mutex *sync.Mutex) {
	go func(x, y int) {
		start := time.Now()
		for {
			time.Sleep(100 * time.Millisecond)
			mutex.Lock()

			proximoInimigo := false
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					ix, iy := x+dx, y+dy
					if iy >= 0 && iy < len(jogo.Mapa) && ix >= 0 && ix < len(jogo.Mapa[iy]) && jogo.Mapa[iy][ix] == Inimigo {
						proximoInimigo = true
					}
				}
			}
			if proximoInimigo {
				jogo.Mapa[y][x] = MoedaVermelhaNegativa
			} else {
				jogo.Mapa[y][x] = MoedaVermelhaVerde
			}

			interfaceDesenharJogo(jogo)

			if jogo.PosX == x && jogo.PosY == y {
				if proximoInimigo {
					jogo.Pontos -= 2
					jogo.StatusMsg = "Coletou moeda vermelha (-2)"
				} else {
					jogo.Pontos += 2
					jogo.StatusMsg = "Coletou moeda verde (+2)"
				}
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				return
			}

			if time.Since(start) > 15*time.Second {
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				mutex.Unlock()
				return
			}

			mutex.Unlock()
		}
	}(x, y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
