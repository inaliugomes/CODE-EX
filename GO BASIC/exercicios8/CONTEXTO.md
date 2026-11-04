# Exercicios8 — Concorrência II (select, Mutex, -race)

Pré-requisito: exercicios7 concluído (goroutines, channels, close/range,
WaitGroup, worker pool).

## Ideias-chave

**O problema: estado partilhado.** Channels passam dados entre goroutines.
Mas quando várias goroutines mexem na MESMA variável/map, há corrida de dados
(data race): resultado imprevisível. `x++` não é atómico (ler, somar, gravar).

**sync.Mutex** = uma chave. `Lock()` espera pela chave; `Unlock()` devolve-a.
- Só uma goroutine de cada vez na zona protegida.
- Padrão: `mu.Lock(); defer mu.Unlock()`.
- Protege-se o DADO, em TODOS os sítios que lhe tocam (leitura incluída).
- Nunca copiar um Mutex (receptor por ponteiro, struct passada por ponteiro).
- `sync.RWMutex`: muitos leitores (`RLock`) OU um escritor (`Lock`).

**select** = `switch` para channels. Espera por vários e executa o primeiro
que estiver pronto (se vários, escolhe um ao acaso).
- `time.After(d)` devolve um channel que recebe um valor passado `d`:
  serve de timeout.
- `default` torna o select não-bloqueante.
- Um channel nil nunca fica pronto (útil para "desligar" um case).

**Channel ou Mutex?** Passar dados / coordenar fluxo -> channel. Proteger
estado simples partilhado (contador, cache) -> Mutex. Ambos são válidos.

**Detetor de corridas:** `go test -race ./...`. Só apanha corridas que
realmente acontecem durante o teste, mas quando apanha mostra as duas linhas
em conflito. Corre sempre com `-race`.

**Goroutine leak:** goroutine presa a enviar/receber para sempre. Não dá erro,
só consome memória. Channel com buffer ou fechar/cancelar evita.

## Exercícios (main.go)

- [x] 1. `Contador`: Mutex à volta de Incrementar/Valor
- [x] 2. `Cache`: map + Mutex (ou RWMutex); `NovaCache` inicializa o map
- [ ] 3. `EsperarOuTimeout`: select com `time.After`
- [ ] 4. `PrimeiroAResponder`: select sobre dois channels
- [ ] 5. `BuscarTodos`: vários "downloads" em paralelo (WaitGroup + Mutex)

Validar: `cd exercicios8 && go test -race ./...` (timeout de 2s por teste).
Os testes falham de início, é normal: são esqueletos com TODO.
Apaga a linha `var _ = ...` de main.go quando já usares os imports.

## Fontes

- Tour of Go: Select, Default Selection, sync.Mutex
- Go by Example: Select, Timeouts, Non-Blocking Channel Operations, Mutexes
- pkg.go.dev: `sync`, `time`

## Perguntas para pensares (sem pressa)

1. Porque é que `Valor()` também precisa de Lock, se só lê?
2. O que acontece num select se dois cases estiverem prontos?
3. No exercício 4, o que acontece à goroutine lenta depois de teres a resposta?
4. No exercício 5, quando preferes um channel de resultados em vez de Mutex?
