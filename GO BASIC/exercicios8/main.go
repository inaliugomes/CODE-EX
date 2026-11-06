package main

import (
	"errors"
	"sync"
	"time"
)

// Evita erros de "import não usado" enquanto os esqueletos não estão feitos.
// Apaga esta linha quando já não precisares.
var _ = []any{errors.New, time.Second, sync.Mutex{}}

// =====================================================================
// Exercício 1 — Contador com Mutex
// Contador é seguro para usar em várias goroutines ao mesmo tempo.
// Incrementar soma 1; Valor devolve o total atual.
// Dica: junta um sync.Mutex à struct; Lock/defer Unlock em cada método.
// Os métodos que mexem no estado partilhado precisam TODOS do lock
// (também a leitura). Receptor por ponteiro: copiar um Mutex é erro.
// =====================================================================
type Contador struct {
	mu    sync.RWMutex
	valor int
}

func (c *Contador) Incrementar() {

	c.mu.Lock()

	defer c.mu.Unlock()

	c.valor++

}

func (c *Contador) Valor() int {

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.valor
}

// =====================================================================
// Exercício 2 — Cache segura (map + Mutex)
// Cache guarda pares chave -> valor. Maps NÃO são seguros para acesso
// concorrente: escrever e ler ao mesmo tempo dá corrida (e até crash).
// Pôr(k, v) grava; Obter(k) devolve (valor, ok) como um map normal.
// Dica: NovaCache() inicializa o map (map nil não aceita escrita).
// Se quiseres, experimenta sync.RWMutex (RLock na leitura).
// =====================================================================
type Cache struct {
	mu    sync.Mutex
	dados map[string]int
}

//aqui estamos a falar de um contrutor 
func NovaCache() *Cache {
	return &Cache{
		dados: make(map[string]int),
	}
}

func (c *Cache) Pôr(chave string, valor int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dados[chave] = valor
}

func (c *Cache) Obter(chave string) (int, bool) {

	c.mu.Lock()

	defer c.mu.Unlock()
	valor, found := c.dados[chave]

	return valor, found
}

// =====================================================================
// Exercício 3 — select com timeout
// EsperarOuTimeout espera um valor de ch durante no máximo d.
// Devolve (valor, nil) se chegar a tempo, ou (0, ErrTimeout) se não.
// Dica: select { case v := <-ch: ... case <-time.After(d): ... }
// =====================================================================
var ErrTimeout = errors.New("timeout")

func EsperarOuTimeout(ch <-chan int, d time.Duration) (int, error) {

	select {

	//Na verdade o que estamos a dizer que é chain vai passar um valor e temos de receber este valor .
	case v := <-ch:
		return v, nil
	case <-time.After(d):
		return 0, ErrTimeout

	}

}

// =====================================================================
// Exercício 4 — fan-in com select
// PrimeiroAResponder lê de dois channels e devolve o primeiro valor que
// chegar (de a ou de b). Ambos podem demorar; o outro valor é ignorado.
// Cuidado: se o outro produtor enviar num channel SEM buffer e ninguém
// ler, essa goroutine fica presa para sempre (leak). Os testes usam
// channels com buffer 1, mas pensa porquê.
// =====================================================================
func PrimeiroAResponder(a, b <-chan string) string {

	select {

	case str := <-a:
		return str

	case str2 := <-b:
		return str2
	}

}

// =====================================================================
// Exercício 5 — descarregar "vários URLs" em paralelo
// Sem rede: recebe uma função buscar(url) (string, error) que simula o
// download. BuscarTodos chama-a numa goroutine por URL e devolve um map
// url -> conteúdo com os que correram bem, mais a lista de URLs que
// falharam (ordem qualquer).
// Dica: usa WaitGroup + Mutex para proteger o map e a slice de falhas
// (ou um channel de resultados; escolhe uma e justifica).
// Rodar com -race: se esquecer o Mutex, o detetor apanha.
// =====================================================================
func BuscarTodos(urls []string, buscar func(string) (string, error)) (map[string]string, []string) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	final := make(map[string]string)
	var erro []string

	for _, url := range urls {

		wg.Add(1)

		go func() {

			defer wg.Done()

			resulta, err := buscar(url)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				final[url] = resulta

			} else {

				erro = append(erro, url)

			}
		}()

	}
	//Vai permitir que main no acabe e deixe funcionalidad a meio , 
	//ou seja a main vai precisar de ficar a espera que esta wait finalize y todos estejam marcado
	
	wg.Wait()

	// TODO
	return final, erro
}

func main() {}
