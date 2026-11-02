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
	mu    sync.Mutex
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
	// TODO
}

func NovaCache() *Cache {
	// TODO
	return nil
}

func (c *Cache) Pôr(chave string, valor int) {
	// TODO
}

func (c *Cache) Obter(chave string) (int, bool) {
	// TODO
	return 0, false
}

// =====================================================================
// Exercício 3 — select com timeout
// EsperarOuTimeout espera um valor de ch durante no máximo d.
// Devolve (valor, nil) se chegar a tempo, ou (0, ErrTimeout) se não.
// Dica: select { case v := <-ch: ... case <-time.After(d): ... }
// =====================================================================
var ErrTimeout = errors.New("timeout")

func EsperarOuTimeout(ch <-chan int, d time.Duration) (int, error) {
	// TODO
	return 0, nil
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
	// TODO
	return ""
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
	// TODO
	return nil, nil
}

func main() {}
