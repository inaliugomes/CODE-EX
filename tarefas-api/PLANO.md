# Projeto 2 — API HTTP de tarefas (`tarefasapi`)

Objetivo: o mesmo domínio da CLI, agora como servidor HTTP com JSON.
Novo: `net/http`, handlers, códigos de estado, interface `Repositorio`,
`httptest`. Reaproveita: JSON, erros com `errors.Is`, Mutex, `-race`.

```
$ curl -X POST localhost:8080/tarefas -d '{"titulo":"Estudar Go"}'
{"id":1,"titulo":"Estudar Go","concluida":false}
$ curl localhost:8080/tarefas
[{"id":1,"titulo":"Estudar Go","concluida":false}]
$ curl -X PATCH localhost:8080/tarefas/1/concluir     # 204
$ curl -X DELETE localhost:8080/tarefas/1             # 204
```

## Âmbito
Dentro: listar, criar, concluir, remover; erros em JSON; testes com `httptest`;
repositório em memória e em ficheiro.
Fora (por agora): autenticação, base de dados, paginação, frameworks.

## Rotas e códigos de estado

| Método e caminho | Sucesso | Erros |
|---|---|---|
| `GET /tarefas` | 200 + lista (`[]`, nunca `null`) | — |
| `POST /tarefas` (`{"titulo":"..."}`) | 201 + tarefa criada | 400 JSON inválido ou título vazio |
| `PATCH /tarefas/{id}/concluir` | 204 | 400 id inválido, 404 não encontrada |
| `DELETE /tarefas/{id}` | 204 | 400 id inválido, 404 não encontrada |

Erros: corpo `{"erro":"mensagem"}` com `Content-Type: application/json`.
O handler traduz erros do repositório com `errors.Is`
(`ErrTituloVazio` -> 400, `ErrNaoEncontrada` -> 404; o resto -> 500).

## Camadas

| Camada | Ficheiro | Responsabilidade |
|---|---|---|
| Domínio | `tarefa.go` | `Tarefa` e erros |
| Repositório | `repositorio.go` | interface + `Memoria` (Mutex) |
| Ficheiro | `armazem.go` (etapa 5) | `Repositorio` em JSON, reutiliza o `Guardar` seguro |
| HTTP | `servidor.go` | rotas e handlers, só conhece a interface |
| Arranque | `main.go` | escolhe o repositório e faz `ListenAndServe` |

## Etapas (cada uma acaba com `go test -race ./...`, `gofmt -l .`, `go vet`)

### Etapa 0 — Esqueleto (feito)
Estrutura criada; `go build` compila e os testes de `Memoria` falham (vermelho).

### Etapa 1 — `Memoria` (1h)
- Implementar os 4 métodos com `sync.Mutex`; `Listar` devolve cópia.
- Pronto quando `repositorio_test.go` passa, incluindo `-race`.
- Pistas: ID = maior + 1; `defer Unlock()` logo depois do `Lock()`;
  `slices.Clone` ou `append([]Tarefa(nil), ...)` para copiar.

### Etapa 2 — `GET /tarefas` (1h)
- Servidor com `http.ServeMux`; handler que lista em JSON.
- Conceitos: `http.ResponseWriter`, `*http.Request`, `json.NewEncoder(w)`,
  `Content-Type`, `httptest.NewRecorder` e `httptest.NewRequest`.
- Testes (escreve-os eu quando chegares aqui): lista vazia dá `[]`.

### Etapa 3 — `POST /tarefas` (1–1,5h)
- Ler o corpo com `json.NewDecoder(r.Body).Decode`, 201, 400 em erros.
- Mapear erros do repositório para códigos com `errors.Is`.

### Etapa 4 — `PATCH` e `DELETE` (1h)
- `r.PathValue("id")` (padrões do `ServeMux` do Go 1.22+), `strconv.Atoi`.
- 400 para id inválido, 404 para inexistente, 204 sem corpo.

### Etapa 5 — Repositório em ficheiro (1h)
- Novo tipo que cumpre `Repositorio`, usando `Carregar`/`Guardar` da CLI.
- Prova da interface: o servidor não muda uma linha.

### Etapa 6 — Arranque e robustez (1h)
- Porta por variável de ambiente, `http.Server` com timeouts,
  teste final com `httptest.NewServer` (pedido real por HTTP).
- Extras: `log` de pedidos (middleware), paragem limpa com `context`.

## Definição de pronto
- [ ] Todas as rotas funcionam com `curl`.
- [ ] `go test -race ./...` passa; `gofmt -l .` e `go vet` limpos.
- [ ] Handlers só conhecem a interface `Repositorio`.
- [ ] Códigos de estado corretos (200, 201, 204, 400, 404).
- [ ] Depois: comparar com projetos de outras pessoas (ver abaixo).

## Depois de pronto: leitura comparada
Com o nosso projeto feito, escolhemos 2–3 projetos pequenos no GitHub e
comparamos: organização de pacotes, tratamento de erros, testes, handlers.
Código descarregado fica FORA deste repositório e só se lê antes de correr.

## Como trabalhamos
Tu escreves o código; eu escrevo plano, testes e valido. Pistas primeiro,
solução só se pedires. Nunca faço commit.
