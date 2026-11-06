# Projeto 1 — CLI de tarefas (`tarefas`)

Objetivo: um programa de terminal para gerir uma lista de tarefas, guardada
num ficheiro JSON. Usa o que já aprendeste: structs, métodos, erros, maps,
JSON, ficheiros, testes em tabela e (opcional) Mutex.

```
$ tarefas add "Estudar Go"
Tarefa 1 criada.
$ tarefas list
[ ] 1  Estudar Go
$ tarefas done 1
Tarefa 1 concluída.
```

## Como se planeia um projeto (o método)

1. **Definir o âmbito**: o que faz e, sobretudo, o que NÃO faz (evita ficar
   infinito).
2. **Desenhar os dados**: que structs existem e o que guardam.
3. **Dividir em camadas**: separar a lógica do resto, para testar sem terminal
   nem ficheiros.
4. **Etapas pequenas**: cada etapa acaba com testes a passar.
5. **Testes primeiro**: escreves o teste, vês falhar, depois implementas.
6. **Definir "pronto"**: critérios claros para fechar o projeto.

## Âmbito

Dentro: `add`, `list`, `done`, `remove`, persistência em JSON, erros claros,
testes.
Fora (por agora): prioridades, datas, categorias, interface gráfica, API.
Extra opcional no fim: flag `--all`, ficheiro configurável, Mutex.

## Desenho

### Dados
```go
type Tarefa struct {
    ID        int    `json:"id"`
    Titulo    string `json:"titulo"`
    Concluida bool   `json:"concluida"`
}
```

### Camadas (porquê: cada uma testa-se separada)

| Camada | Ficheiro | Responsabilidade | Sabe de ficheiros? |
|---|---|---|---|
| Lógica | `tarefa.go` | `Lista` e métodos `Adicionar`, `Concluir`, `Remover` | Não |
| Persistência | `armazem.go` | `Carregar(caminho)` e `Guardar(caminho, lista)` | Sim |
| Interface | `main.go` | ler `os.Args`, chamar as outras, imprimir | Só via armazem |

Regra: `main.go` fica fino. Quanto mais lógica fora dele, mais testável.

### Erros próprios (usar `errors.Is`)
- `ErrTituloVazio` — `add ""`
- `ErrNaoEncontrada` — `done 99`, `remove 99`

### Decisões a tomar tu (etapa 0)
1. Como se gera o ID? (sugestão: maior ID existente + 1; porquê não o tamanho
   da lista?)
2. `Concluir` duas vezes a mesma tarefa é erro ou não faz nada?
3. O que acontece se o ficheiro JSON não existir? (já fizeste isto no
   exercicios6)

## Etapas

Cada etapa: testes a passar com `go test -race ./...`, `gofmt -l .` e `go vet`
sem avisos.

### Etapa 0 — Preparar (15 min)
- Criar `tarefa.go`, `armazem.go`, `main.go` e os `_test.go`.
- Responder às 3 decisões acima, por escrito, num comentário.
- Pronto quando: `go build ./...` compila.

### Etapa 1 — Lógica (1–1,5h)
- Struct `Tarefa` e tipo `Lista` (`[]Tarefa`).
- `Adicionar(titulo) (Tarefa, error)`, `Concluir(id) error`, `Remover(id) error`.
- Testes em tabela, com casos-limite: título vazio, ID inexistente, lista
  vazia, concluir duas vezes, remover o último.
- Pronto quando: lógica 100% testada, sem tocar em ficheiros.
- Atenção: métodos que alteram a lista precisam de receptor por ponteiro.

### Etapa 2 — Persistência (1h)
- `Guardar` e `Carregar` com JSON (`MarshalIndent`, `os.WriteFile`).
- Ficheiro inexistente -> lista vazia, sem erro.
- Testes com `t.TempDir()` e `filepath.Join`. Caso: guardar e carregar
  devolvem o mesmo; JSON corrompido dá erro.
- Pronto quando: ida e volta igual e erros com `%w`.

### Etapa 3 — CLI (1–1,5h)
- `main.go`: ler `os.Args`, `switch` por comando (`add`, `list`, `done`,
  `remove`), carregar, aplicar, guardar.
- Mensagens ao utilizador no stdout, erros no stderr, `os.Exit(1)` em erro.
- Comando desconhecido ou sem argumentos: mostrar ajuda.
- Dica de desenho: uma função `executar(args []string) error` testável, e o
  `main` só chama e imprime o erro.
- Pronto quando: usar a CLI à mão e todos os comandos funcionam.

### Etapa 4 — Robustez (1h)
- `ID` inválido (`done abc`): erro claro, sem panic.
- Escrita segura: gravar num ficheiro temporário e renomear (para não
  corromper se falhar a meio).
- Pronto quando: nenhum input razoável faz panic.

### Etapa 5 — Extras (opcional)
- `list --all` / só pendentes.
- Caminho do ficheiro por variável de ambiente.
- Mutex na `Lista` (reaproveita o exercicios8).
- Escrever um `README.md` curto com exemplos de uso.

## Definição de pronto
- [ ] Todos os comandos funcionam à mão.
- [ ] `go test -race ./...` passa; `gofmt -l .` e `go vet` limpos.
- [ ] Sem lógica no `main.go` além de ler argumentos e imprimir.
- [ ] Casos-limite testados (vazio, inexistente, inválido).
- [ ] Erros comparáveis com `errors.Is`.

## Como trabalhamos
- Tu escreves o código; eu escrevo plano, esqueleto de testes e valido.
- Primeiro pistas, solução só se pedires.
- Fim de cada etapa: validamos e atualizo o PROGRESSO.md.
