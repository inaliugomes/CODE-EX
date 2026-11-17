# Estudo de Go — Objetivo e Estado

Última atualização: 2026-11-15 (tarefas-api: Etapa 5 feita, `Armazem` em ficheiro JSON)

Próximo: tarefas-api Etapa 6 (arranque no `main.go` com o `Armazem`, porta por variável de ambiente, `http.Server` com timeouts, teste com `httptest.NewServer`). Depois, limpeza de estilo e leitura comparada.

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

### Rever amanhã — conceitos HTTP da Etapa 2 (o utilizador achou confuso)

Começar por pedir ao utilizador que explique, por palavras suas, o fluxo abaixo; só depois avançar.

- Fluxo: pedido -> `Handler()` (o mux) -> procura método + caminho -> chama a função
  do handler `(w http.ResponseWriter, r *http.Request)` -> escreve a resposta em `w`.
- Mux = tabela "padrão -> função". `HandleFunc("GET /tarefas", s.listar)` precisa de
  2 argumentos; o método dentro do padrão é do Go 1.22+.
- `s.listar` SEM parênteses = passar a função (o mux chama-a quando chegar um pedido).
- O mux é criado e preenchido UMA vez (`NovoServidor`), guardado no `Servidor`, e o
  `Handler()` devolve esse mesmo mux. Erro que cometeu: devolver `&Servidor{}` vazio e
  `http.NewServeMux()` novo no `Handler()` (rota perdida); variável `route` não usada.
- Ordem na resposta: `Header().Set(...)` ANTES do `Encode`/`Write` (a 1.ª escrita envia
  cabeçalhos + 200; depois já não se mudam). `http.Error` escreve texto + código.
- `nil` vs `[]`: `slices.Clone(nil)` dá `nil` e o JSON sai `null`; garantir `[]Tarefa{}`.
- Testes sem rede: `httptest.NewRequest` + `httptest.NewRecorder` + `Handler().ServeHTTP(rec, req)`.
  Um repositório falso (`repoComErro` embebe `Repositorio`) força o caminho do erro 500.
- Pendente (estilo): mensagem do 500 é `"...."` (pôr texto real, ideal JSON `{"erro":...}`),
  apagar o TODO do comentário de `Servidor`, tratar/aceitar o erro do `Encode`.

### Projeto 1 — tarefas-cli (`../tarefas-cli/PLANO.md`)

- [x] Etapa 0: esqueleto criado (`tarefa.go`, `armazem.go`, `main.go`, testes).
- [x] Etapa 1: lógica (`Adicionar`, `Concluir`, `Remover`) — `go test -race` passa.
- [x] Etapa 2: persistência (`Carregar`/`Guardar`), testes primeiro.
- [x] Etapa 3: CLI (`executar` em `main.go`, testes em `main_test.go`) — `go test -race` passa, testado à mão.
- [x] Etapa 4: robustez — concluída (10 nov). Validada: 16 testes com `-race`, `go vet` e `gofmt` limpos.
  - [x] ID inválido: `lerID` + erro `id inválido: ...` em `done`/`remove`; `TestExecutarIDInvalido` passa.
  - [x] Escrita segura no `Guardar`: `CreateTemp` + `Write` + `Close` (erro verificado) + `Rename(tmp.Name(), caminho)`; `Remove` do temporário em cada falha. `go test -race`, `go vet` e `gofmt` limpos.
  - [x] `git checkout -- tarefas.json`, apagado `tarefas.json.tmp-123`, removido código morto (`if len(args) == 0` dentro de `add`/`remove`).
  - Combinado: o utilizador escreve o código; o Claude escreve os testes e valida.
  - Lições: `CreateTemp` devolve o ficheiro JÁ aberto (usar `tmp.Write`, não `os.WriteFile(tmp.Name())`); nada de nomes fixos no `Rename`/`Remove` (ignoram o `caminho`); `Rename` na mesma pasta é atómico; verificar o `Close` antes do `Rename`; devolver o erro do `Rename`; apagar o temporário em toda a falha.
  - [x] Estilo: `tmp.Close()` antes do `os.Remove` no ramo do `Write`; "texto vazio" -> "comando em falta".
  - [x] `tarefas-cli/.gitignore` criado.
  - Pendente (git, decisão do utilizador): acrescentar `/tarefas` ao `.gitignore` (hoje só tem `tarefas.json`, sem newline final) e `git rm --cached tarefas-cli/tarefas tarefas-cli/tarefas.json`, porque ambos continuam versionados e o `.gitignore` não os afeta.
### Projeto 2 — tarefas-api (`../tarefas-api/PLANO.md`, conceitos em `CONTEXTO.md`)

API HTTP de tarefas com `net/http`. Primeira API do utilizador em Go.
Método igual ao da CLI: o utilizador escreve o código; o Claude escreve plano, testes e valida.
Decisão (10 nov): construir o nosso primeiro e só depois ler/comparar 2–3 projetos pequenos
de outras pessoas no GitHub (fora deste repositório; ler antes de correr; não copiar linha a linha).

- [x] Etapa 0: esqueleto (`tarefa.go`, `repositorio.go`, `servidor.go`, `main.go`), `PLANO.md`, `CONTEXTO.md`,
  `repositorio_test.go` (6 testes, a falhar de propósito). `go vet` e `gofmt` limpos.
- [x] Etapa 1: `Memoria` (12 nov) — `go test -race` passa (6 testes), `gofmt`/`go vet` limpos.
- [x] Etapa 2: `GET /tarefas` (13 nov) — `go test -race` passa (5 testes), `gofmt`/`go vet` limpos.
- [x] Etapa 3: `POST /tarefas` (14 nov) — 201, 400, 500 com `errors.Is`; `go test -race` passa.
- [x] Etapa 4: `PATCH .../concluir` e `DELETE` (14 nov) — `r.PathValue` + `Atoi` (função `letID`), 204, 400, 404, 500. 24 testes com `-race`, `gofmt`/`go vet` limpos.
- Lições (Etapas 3 e 4): `Decode(&x)` devolve só `error` e escreve na variável (como `Unmarshal`);
  o contrato do corpo é definido por NÓS (struct + tag); os dados vêm no corpo (`r.Body`), não na rota;
  variável de pedido GLOBAL = data race (a `entrada` tem de ser local ao handler);
  depois de `http.Error` SEMPRE `return` (faltou 3 vezes: 500 no POST, 400 no DELETE, e corpo duplicado);
  `errors.Is(nil, X)` é falso -> envolver em `if err != nil` antes de distinguir 404/500, senão o sucesso dá 500;
  `Encode` antes de `WriteHeader(201)` fixa 200; sem `WriteHeader` nem corpo o Go envia 200, por isso 204 é explícito;
  400 = pedido mal formado (`abc`), 404 = bem formado mas não existe; o mux dá 405 sozinho para método errado.
- Respondido (Etapa 2): `listar` com erro -> 500 + mensagem no corpo, não o nome da constante; o erro vai para o CLIENTE.
- Pendente (estilo, tratar no fim de todas as etapas): erros em JSON `{"erro":"..."}` em vez de texto;
  tratar o erro do `Encode`; `else` depois de `return`; renomear `letID` -> `lerID` e `delete` -> `remover`;
  mensagem do 400 diz "Erro do servidor" (devia ser "id inválido"); "Error", "nao", espaço antes da vírgula;
  corrigir o 500 do `listar` e apagar o TODO do comentário de `Servidor`.
- [x] Etapa 5: repositório em ficheiro (15 nov) — `Armazem` com `Carregar`/`Guardar` copiados da CLI (`[]Tarefa`); 11 testes novos, `go test -race` passa, `gofmt`/`go vet` limpos; `servidor.go` intacto.
  - Lições: o `Guardar` tem de ficar DENTRO do `if` em que se alterou (fora, só corria quando NÃO encontrava e o `return nil` perdia a alteração);
    o erro de domínio tem de ser o MESMO (`ErrTituloVazio`, não `errors.New` novo) senão o `errors.Is` falha e dá 500 em vez de 400;
    ficheiro inexistente -> `nil` (não `err`) no `Carregar`, senão o 1.º arranque dá 500;
    o Mutex cobre carregar+alterar+guardar (50 `Adicionar` em paralelo sem perdas nem IDs repetidos);
    não encontrado -> `ErrNaoEncontrada` sem gravar; uma slice `nil` viraria `null` no JSON.
  - Pendente (estilo): `return Guardar(...)` direto; `Carregar`/`Guardar` privados; `if l == nil` no `Guardar`; maiúsculas/`nao` nas mensagens.
- [ ] Etapa 6: arranque, timeouts, teste com `httptest.NewServer`.
- Depois: leitura comparada de projetos de outras pessoas.
- Pendente (git, decisão do utilizador): `tarefas-cli` ainda tem o binário e o `tarefas.json` versionados.

- Lições (Etapa 3): `||` vs `&&` numa validação (com `||` é sempre verdadeira);
  verificar `len(args)` antes de indexar (`args[0]`/`args[1]` dão panic);
  `Guardar` depois de TODA a operação que altera (add/done/remove), não só no add;
  devolver o erro real (`%w`) em vez de o esconder; os testes comparam o texto
  EXATO (maiúsculas, pontos, `\n`, dois espaços no `list`); escrever em `saida`
  (`fmt.Fprintf`), não em `fmt.Printf`; `main` trata o erro de `executar`
  (stderr + `os.Exit(1)`).
- Conceito: `go run main.go` só compila esse ficheiro; com vários ficheiros no
  mesmo `package` usa `go run .` (a pasta toda é um só programa). `os.Args`
  lê o que se escreve no terminal; `main()` é o ponto de arranque.
- Pendente (estilo): mensagens de erro (`"tem nao encontrado"`, `"texto vazio"`,
  erro cru do `Atoi`), função auxiliar `lerID` para `done`/`remove`, apagar o
  `TODO` do `main`, `.gitignore` para `tarefas.json` e o binário `tarefas`,
  respostas 2 e 3 no comentário de `tarefa.go`.
- Lições (Etapa 1): ID = maior existente + 1 (nunca `len`, senão repete após
  remover); lista vazia dá panic se se indexar `[len-1]`; receptor por ponteiro
  e `(*l)[i]` para alterar (o `range` dá cópia); `return` logo ao encontrar,
  em vez de variável auxiliar; não alterar a slice a meio de um `range`.
- Pendente (estilo): tirar `if id < 0` redundante, mensagens de erro mais
  limpas, preencher as respostas 2 e 3 no comentário de `tarefa.go`.

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
