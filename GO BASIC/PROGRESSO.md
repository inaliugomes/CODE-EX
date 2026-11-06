# Estudo de Go — Objetivo e Estado

Última atualização: 2026-11-04 (exercicios8 concluído)

Próximo: projeto 1 (CLI de tarefas, JSON + testes). A exercicios7 e a exercicios8 (concorrência) estão concluídas.

## Plano da semana (26–30 out 2026)

Ritmo: 5 dias, 1 a 2h por dia.

| Dia | Tema | Tempo | O que fazer |
|---|---|---|---|
| Seg 26 | JSON | 1,5h | Ler a secção "JSON e ficheiros" abaixo. Go by Example: *JSON*. Exercícios 1 (`Guardar`) e 2 (`Carregar`). |
| Ter 27 | Ficheiros | 1,5–2h | Go by Example: *Reading Files*, *Writing Files*, *Line Filters*. Exercícios 3 (`Gravar`), 4 (`Acrescentar`), 5 (`LerLinhas`). `go test ./...` e rever casos-limite (inexistente, vazio). |
| Qua 28 | Concorrência I | 1,5h | Tour of Go: *Concurrency* (Goroutines, Channels, Buffered Channels). Go by Example: *Goroutines*, *Channels*, *WaitGroups*. Exemplos próprios. |
| Qui 29 | Concorrência II | 1,5–2h | Go by Example: *Select*, *Mutexes*. Tour: *Select*, *sync.Mutex*. Exercício: descarregar vários URLs em paralelo. Rodar `go test -race`. |
| Sex 30 | Projeto 1 (início) | 1–2h | Desenhar a CLI de tarefas: struct `Tarefa`, comandos (`add`, `list`, `done`), ficheiro JSON. Testes primeiro, depois `Carregar`/`Guardar`. Rever o que ficou por fazer. |

Fontes: A Tour of Go (go.dev/tour), Go by Example (gobyexample.com),
pkg.go.dev (`encoding/json`, `os`, `bufio`, `sync`), Effective Go (só a
secção de concorrência, se sobrar tempo).

Notas: se atrasar na terça, o exercício 5 passa para quarta; não comprimir
a concorrência. Exercícios com pistas primeiro, solução só se pedir.

## Objetivo

Aprender Go (Golang) de forma prática: ler um resumo, fazer exercícios com
testes (`go test`) e validar a resolução. Meta final: conseguir construir
um projeto pequeno sozinho (CLI de tarefas em JSON ou API HTTP simples).

Preferências: explicações em português, exemplos sempre em Go, sem dar a
solução logo (primeiro pistas e modelo mental).

**Regra: nunca fazer commit.** Só o utilizador faz commits (git commit/push).
O Claude só os faz se houver um pedido direto e explícito do utilizador.

## Estrutura do projeto

```
GO BASIC/   (antes chapter9)
├── main.go          # slices e arrays (experiências soltas)
├── resumo/          # resumo do básico de Go (main.go, main2.go)
├── exercicios/      # ronda 1: slices, ciclos  -> CONCLUÍDO
├── exercicios2/     # ronda 2: maps e erros    -> CONCLUÍDO
├── exercicios3/     # ronda 3: ponteiros e pacotes -> CONCLUÍDO
│   ├── main.go          # ponteiros (4 exercícios)
│   └── conta/           # pacote próprio (conta.go, erros.go)
├── exercicios4/     # ronda 4: defer -> CONCLUÍDO
├── exercicios5/     # ronda 5: switch e variádicas -> CONCLUÍDO
├── exercicios6/     # ronda 6: JSON + ficheiros -> CONCLUÍDO
│   └── main.go          # 5 exercícios + main_test.go
├── exercicios7/     # ronda 7: concorrência I -> CONCLUÍDO
│   └── main.go          # 5 exercícios + main_test.go (com timeout)
└── exercicios8/     # ronda 8: concorrência II -> CONCLUÍDO
    ├── CONTEXTO.md      # conceitos, fontes e perguntas
    └── main.go          # 5 exercícios (TODO) + main_test.go
```

Para testar: `cd <pasta> && go test ./...`

## Já dominado

- Constantes, structs, métodos, interfaces (implícitas), arrays, slices
- Maps: contagem com `m[k]++`, leitura com `v, ok := m[k]`
- Erros: `errors.New`, `fmt.Errorf` com `%w`, `errors.Is`, `errors.As`,
  tipo de erro próprio
- Testes em tabela

## Estado atual: exercicios4 (defer) — CONCLUÍDO

exercicios3 (ponteiros + pacote `conta`) concluído e validado.
- Lição: o teste não apanhava `valor < 0` em vez de `<= 0` em `Levantar`;
  testar sempre o caso-limite `0`. Usar `errors.Is` para verificar o erro.
- Preferir `t.Errorf` a `t.Fatalf` quando os passos são independentes.

**exercicios4 — CONCLUÍDO**, `main.go` com:
- [x] 1. `Ordem`: defer em pilha (inicio, meio, fim-B, fim-A)
- [x] 2. `Usar`: `Fechar` uma só vez com defer, mesmo com return antecipado
- [x] 3. `Dobrar`: resultado com nome + defer com closure altera o retorno
- [x] 4. `DividirSeguro`: `recover()` dentro de defer, devolve erro em vez de panic
- Validar com `cd exercicios4 && go test ./...`

### Conceito fixado (defer)

- `defer` corre quando a FUNÇÃO onde está termina (não "no fim do código").
- Vários defer: pilha LIFO (o último registado corre primeiro).
- `defer f(x)` avalia os argumentos NA HORA ("foto"); `defer func(){...}()`
  lê a variável só quando corre ("janela"). Os `()` finais são obrigatórios.
- Ciclo com `defer` e `i`: Go >= 1.22 (versão em `go.mod`) dá uma variável
  `i` por iteração -> `2,1,0`; com `go 1.21` ou menos, uma só `i` -> `3,3,3`
  (correção antiga: `i := i`).
- Uso típico: `defer f.Close()` logo após abrir; `defer mu.Unlock()`;
  `recover()` só funciona dentro de uma função adiada.

### Conceito fixado (ponteiros)

- `x` é a variável (a "casa" com o valor); o ponteiro é `p`, que guarda um endereço.
- `&x` = endereço de x; `*p` = valor para onde p aponta (ler ou escrever).
- `p = x` com `var p int` COPIA o valor (duas casas independentes).
  Com `p *int` só se aceita `p = &x` (as duas chegam à mesma casa).
- Em Go tudo é cópia: argumentos e receptores. O ponteiro é uma cópia do
  endereço; `p = ...` só muda a cópia, `*p = ...` altera o original.
- `&a` com `a *int` dá `**int` (endereço do ponteiro), não o valor.
- Para trocar: `temp := *a; *a = *b; *b = temp`.
- Quando usar ponteiros: alterar o original numa função/método, evitar
  copiar structs grandes, representar ausência com `nil`, partilhar o
  mesmo objeto. Slices/maps/channels já guardam referência internamente.
- Desreferenciar `nil` provoca panic: verificar `p == nil` antes.
- Maiúscula = exportado (público); minúscula = privado ao pacote.

### Exercicios6 — CONCLUÍDO (28 out)

- [x] 1. `Guardar`: `json.MarshalIndent` + `os.WriteFile` (0644), erro com `%w`
- [x] 2. `Carregar`: `os.ReadFile` + `Unmarshal`; inexistente -> lista vazia, nil
- [x] 3. `Gravar`: `os.Create` + `defer` com Close que devolve erro (resultado com nome)
- [x] 4. `Acrescentar`: `os.OpenFile` com `O_APPEND|O_CREATE|O_WRONLY`
- [x] 5. `LerLinhas`: `os.Open` + `bufio.Scanner`, `sc.Err()`, inexistente -> `os.ErrNotExist`
- Validar com `cd exercicios6 && go test -race ./...` (9 testes passam).
- Lições: `Unmarshal` precisa de `&var` (escreve na variável); `errors.Is`
  devolve bool -> usar DENTRO de `if err != nil`, nunca antes; verbo é `%w`
  (minúscula); mensagens de erro em minúscula, com `: ` antes do `%w`;
  não ignorar o erro de `os.WriteFile`; `gofmt -l .` e `go vet` antes de validar.
- `:=` pode esconder o `err` do resultado; só atribuir o erro do Close se `err == nil`.
- `OpenFile` só com `O_APPEND` abre em leitura: o `Open` não dá erro mas a
  escrita falha em silêncio se não verificares o `WriteString`. Com `O_CREATE`
  não é preciso ramo "se não existe" (e `os.Create` ainda trunca).
- `append` devolve uma slice nova: `s = append(s, x)`; o compilador recusa se
  deitares o resultado fora.
- `defer` sem `defer` à frente (`f.Close()` solto) fecha logo -> "file already closed".
- Pendente (estilo, não bloqueia): em `Acrescentar` tirar o ramo `os.IsNotExist`
  + `os.Create` e tratar o erro do `Close`; mensagens de erro em minúscula com `: `.

### Conceito fixado (posição do defer)

- O `defer` fica REGISTADO quando a execução passa pela linha; corre no fim da
  função, mas só se tiver sido registado. `return` antes da linha = nunca registado.
- Registo (na linha) vs execução (no fim). Argumentos/receptor avaliam-se no registo.
- Regra: abrir -> verificar erro -> `defer Close` -> usar. Nunca `defer` antes
  de verificar o erro (ex.: `resp.Body.Close()` com `resp == nil` dá panic).

### Exercicios7 — CONCLUÍDO (concorrência I)

- [x] 1. `SomarParalelo`: duas goroutines + channel com somas parciais
- [x] 2. `Gerar`: produtor, `close` do channel, leitura com `range`
- [x] 3. `QuadradosParalelo`: `sync.WaitGroup` (Add antes da goroutine, Done, Wait)
- [x] 4. `Dobrar`: pipeline, etapa que lê de `in` e fecha o seu channel
- [x] 5. `Trabalhadores`: worker pool com channels com buffer
- Validar com `cd exercicios7 && go test -race ./...` (os testes têm timeout
  de 2s: se bloquear, é deadlock ou channel por fechar).
- Validado em 2 nov: `go test -race ./...` passa. Falta só `gofmt -w main.go`
  (indentação/espaços em `Trabalhadores`).
- Pendente do exercicios6 (estilo): `Acrescentar` sem o ramo `os.Create`,
  tratar erro do `Close`, mensagens de erro em minúscula.

### Conceito fixado (JSON e ficheiros)

- `Marshal`: Go -> JSON (`[]byte`, usar `string(dados)` para ver texto);
  `Unmarshal(dados, &x)`: JSON -> Go (precisa de ponteiro).
- Só campos exportados (maiúscula) entram no JSON; minúscula é ignorada sem erro.
  A tag `` `json:"nome"` `` define o nome no JSON.
- Ficheiro inexistente: `errors.Is(err, os.ErrNotExist)`.
- Abrir: verificar erro, depois `defer f.Close()` logo a seguir. Ao escrever,
  o erro do Close conta.
- Testes com ficheiros: `t.TempDir()` + `filepath.Join`.

### Exercicios8 — CONCLUÍDO (concorrência II, 4 nov)

- [x] 1. `Contador`: Mutex + `defer Unlock` (validado 2 nov, `-race` ok)
- [x] 2. `Cache`: map + Mutex (`NovaCache` com `make`) (validado 4 nov, `-race` ok)
- [x] 3. `EsperarOuTimeout`: select + `time.After`
- [x] 4. `PrimeiroAResponder`: select sobre dois channels
- [x] 5. `BuscarTodos`: WaitGroup + Mutex
- Validar com `cd exercicios8 && go test -race ./...`.
- Lição (ex. 1): escreveu `defer c.mu.Lock()` em vez de `Unlock()` -> deadlock
  (timeout no teste). Regra: `Lock()` e logo a seguir `defer Unlock()`; reler
  a linha do `defer` antes de correr o teste. Cada `Lock` precisa do seu `Unlock`.
- Pendente (estilo): apagar `// TODO` soltos e, no fim, a linha `var _ = ...`.
- Lições (ex. 5): chave do map = URL, valor = conteúdo; falhas = lista de URLs (não erros);
  `wg.Wait()` FORA do ciclo; `buscar` fora do lock (só o map/slice protegidos);
  `select` executa um só case e não cria goroutines; `if err == nil` para o sucesso.
- Pendente (estilo): `var _ = ...`, `// TODO` solto, imports não usados.

## Próximos passos

1. (feito) `exercicios3`.
2. (feito) `defer` (exercicios4), `switch` + variádicas (exercicios5),
   JSON + ficheiros (exercicios6).
3. (feito) **Concorrência I** (exercicios7): goroutines, channels, `close`/`range`,
   `sync.WaitGroup`, channels com buffer.
4. **Concorrência II**: `select`, `sync.Mutex`, `-race` (exercício: descarregar
   vários URLs em paralelo).
4. Projeto final: CLI de tarefas (JSON + `os` + testes) ou API com
   `net/http` e interface `Repositorio`.

## Recursos

- A Tour of Go: https://go.dev/tour
- Go by Example: https://gobyexample.com
- Effective Go (mais tarde)
