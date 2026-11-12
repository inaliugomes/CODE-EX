# Conceitos para a API HTTP (lê por etapas, não tudo de uma vez)

Nunca fizeste uma API em Go: está tudo bem. Lê só a secção da etapa em que estás.

## Etapa 1 — Interface e Mutex (já conheces a base)
- **Interface**: conjunto de métodos. Quem os tiver, cumpre-a (implícito). Aqui
  `Repositorio` permite trocar memória por ficheiro sem mexer no servidor.
- **`var _ Repositorio = (*Memoria)(nil)`**: o compilador confirma que `Memoria`
  cumpre a interface. Se faltar um método, não compila.
- **Mutex**: como no exercicios8. Num servidor, cada pedido corre na sua
  própria goroutine, por isso o estado partilhado PRECISA de proteção.

## Etapa 2 — O modelo HTTP
- Cliente envia um **pedido** (método + caminho + cabeçalhos + corpo);
  servidor devolve uma **resposta** (código + cabeçalhos + corpo).
- Método diz a intenção: `GET` ler, `POST` criar, `PATCH` alterar parte,
  `DELETE` apagar.
- Um **handler** é uma função `func(w http.ResponseWriter, r *http.Request)`.
  - `r` = o pedido (ler método, caminho, corpo).
  - `w` = onde escreves a resposta.
- **`http.ServeMux`**: o "encaminhador": liga caminhos a handlers.
  `mux.HandleFunc("GET /tarefas", s.listar)` (Go 1.22+ aceita método no padrão).
- **Ordem importa em `w`**: primeiro cabeçalhos (`w.Header().Set(...)`), depois
  `w.WriteHeader(código)`, e só depois o corpo. Escrever o corpo primeiro fixa
  logo o código 200.
- **JSON**: `json.NewEncoder(w).Encode(x)` escreve; `json.NewDecoder(r.Body).Decode(&x)` lê.
  Uma slice `nil` vira `null` em JSON; uma slice vazia `[]Tarefa{}` vira `[]`.

## Etapa 3 — Códigos de estado
| Código | Significa | Quando |
|---|---|---|
| 200 | OK | listar |
| 201 | Created | criou uma tarefa |
| 204 | No Content | sucesso sem corpo (concluir, remover) |
| 400 | Bad Request | o cliente enviou lixo (JSON inválido, id `abc`, título vazio) |
| 404 | Not Found | id não existe |
| 500 | Internal Server Error | erro nosso, inesperado |
- Regra: 4xx = culpa do cliente, 5xx = culpa do servidor.
- **Traduzir erros**: o handler faz `errors.Is(err, ErrNaoEncontrada)` -> 404.
  O repositório não sabe de HTTP; o handler não sabe de ficheiros.

## Etapa 4 — Parâmetros no caminho
- `mux.HandleFunc("DELETE /tarefas/{id}", ...)` e dentro `r.PathValue("id")`
  (devolve texto; convertes com `strconv.Atoi`, como na CLI).

## Testar sem rede (`net/http/httptest`)
- `httptest.NewRequest(método, caminho, corpo)` cria um pedido falso.
- `httptest.NewRecorder()` é um `ResponseWriter` que regista a resposta
  (`rec.Code`, `rec.Body`).
- Chamas `handler.ServeHTTP(rec, req)` e verificas. Sem porta, sem rede, rápido.
- `httptest.NewServer(handler)`: servidor real numa porta aleatória (teste final).

## Etapa 6 — Arranque
- `http.ListenAndServe(":8080", handler)` bloqueia; só devolve em erro.
- Preferível `http.Server{Addr, Handler, ReadHeaderTimeout}`: sem timeouts, um
  cliente lento pode prender o servidor.
- **Middleware**: função que recebe um `http.Handler` e devolve outro (por
  exemplo, para registar cada pedido em log).

## Perguntas para pensares (respondemos juntos)
1. Porque é que cada pedido precisa de Mutex se o programa "só tem um main"?
2. Porque é que `Listar` devolve cópia?
3. O que acontece se escreveres no corpo antes de `WriteHeader(201)`?
4. `PATCH /tarefas/1/concluir` ou `PUT`? Qual o significado de cada um?
5. Porque é que o handler traduz erros e o repositório não devolve códigos HTTP?

## Fontes
- Go by Example: *HTTP Server*, *JSON*, *Context* (só depois).
- pkg.go.dev: `net/http`, `net/http/httptest`, `encoding/json`.
- Blog oficial: *Routing Enhancements for Go 1.22* (padrões do `ServeMux`).
- Tour of Go não cobre HTTP; a Go by Example e a documentação chegam.
