// =====================================================================
// RESUMO DO BÁSICO DE GO
// Para correr:  cd resumo  &&  go run main.go
// =====================================================================

// Todo o programa executável pertence ao pacote "main"
package main

// Importar pacotes (vários de uma vez, entre parênteses)
import (
	"errors"  // Criar erros
	"fmt"     // Escrever no ecrã e formatar texto
	"strings" // Funções para trabalhar com texto
)

// =====================================================================
// 1. CONSTANTES — valores que nunca mudam
// =====================================================================
const Pi = 3.14159

// Constantes agrupadas
const (
	DiasDaSemana = 7
	Saudacao     = "Olá"
)

// =====================================================================
// 2. STRUCTS — agrupar dados relacionados num tipo novo
// =====================================================================
type Pessoa struct {
	Nome  string
	Idade int
}

// =====================================================================
// 3. MÉTODOS — funções "ligadas" a um tipo
// =====================================================================

// Recetor por valor (p Pessoa): recebe uma CÓPIA, não altera o original
func (p Pessoa) Apresentar() string {
	return fmt.Sprintf("Sou %s e tenho %d anos", p.Nome, p.Idade)
}

// Recetor por ponteiro (p *Pessoa): consegue ALTERAR o original
func (p *Pessoa) FazerAnos() {
	p.Idade++
}

// =====================================================================
// 4. INTERFACES — definem um comportamento (que métodos um tipo tem)
// =====================================================================

// Qualquer tipo com um método "Area() float64" é uma Forma.
// Não é preciso escrever "implements": é automático.
type Forma interface {
	Area() float64
}

type Retangulo struct{ Largura, Altura float64 }
type Circulo struct{ Raio float64 }

func (r Retangulo) Area() float64 { return r.Largura * r.Altura }
func (c Circulo) Area() float64   { return Pi * c.Raio * c.Raio }

// =====================================================================
// 5. FUNÇÕES
// =====================================================================

// Função simples: parâmetros com tipo, e o tipo de retorno no fim
func somar(a int, b int) int {
	return a + b
}

// Parâmetros do mesmo tipo podem partilhar o tipo
func multiplicar(a, b int) int {
	return a * b
}

// Vários valores de retorno (muito usado em Go)
func dividirInteiros(a, b int) (int, int) {
	return a / b, a % b // quociente e resto
}

// Função variádica: aceita qualquer número de argumentos (recebe um slice)
func somarTodos(numeros ...int) int {
	total := 0
	for _, n := range numeros {
		total += n
	}
	return total
}

// =====================================================================
// 6. ERROS — em Go, os erros são valores devolvidos (não há exceções)
// =====================================================================

// Convenção: o erro é sempre o ÚLTIMO valor de retorno
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("não é possível dividir por zero")
	}
	return a / b, nil // nil = "sem erro"
}

func main() {
	// =================================================================
	// 7. VARIÁVEIS
	// =================================================================

	// Forma longa: var nome tipo = valor
	var nome string = "Ana"

	// O tipo pode ser deduzido pelo valor
	var cidade = "Lisboa"

	// Forma curta (só dentro de funções) — a mais usada
	idade := 25

	// Várias de uma vez
	x, y := 10, 20

	// Sem valor inicial → fica com o "valor zero" do tipo
	var contador int    // 0
	var ativo bool      // false
	var texto string    // "" (texto vazio)
	var decimal float64 // 0.0

	fmt.Println(nome, cidade, idade, x, y)
	fmt.Println("Valores zero:", contador, ativo, texto, decimal)

	// =================================================================
	// 8. TIPOS BÁSICOS
	// =================================================================
	var inteiro int = 42     // número inteiro
	var real float64 = 3.75  // número decimal
	var letra rune = 'A'     // um caráter (entre plicas)
	var frase string = "Go!" // texto (entre aspas)
	var verdade bool = true  // true ou false
	var b byte = 255         // um byte (0 a 255)

	// Conversão de tipos: é sempre explícita em Go
	convertido := float64(inteiro) + real
	fmt.Println(inteiro, real, string(letra), frase, verdade, b, convertido)

	// =================================================================
	// 9. OPERADORES
	// =================================================================
	fmt.Println(7+2, 7-2, 7*2, 7/2, 7%2)             // 9 5 14 3 1  (7/2 entre inteiros = 3)
	fmt.Println(7.0 / 2)                             // 3.5
	fmt.Println(5 > 3, 5 == 3, 5 != 3)               // true false true
	fmt.Println(true && false, true || false, !true) // E, OU, NÃO

	contador++    // soma 1
	contador += 5 // contador = contador + 5
	fmt.Println("Contador:", contador)

	// =================================================================
	// 10. IF / ELSE — sem parênteses, chavetas obrigatórias
	// =================================================================
	if idade >= 18 {
		fmt.Println("Maior de idade")
	} else if idade >= 13 {
		fmt.Println("Adolescente")
	} else {
		fmt.Println("Criança")
	}

	// If com declaração curta: a variável só existe dentro do if
	if resultado, err := dividir(10, 0); err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Println("Resultado:", resultado)
	}

	// =================================================================
	// 11. FOR — o ÚNICO ciclo em Go (não há while)
	// =================================================================

	// Clássico
	for i := 0; i < 3; i++ {
		fmt.Println("i =", i)
	}

	// Estilo "while"
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println("n =", n) // 128

	// Percorrer um intervalo de números (Go 1.22+)
	for i := range 3 {
		fmt.Print(i, " ") // 0 1 2
	}
	fmt.Println()

	// break sai do ciclo; continue salta para a próxima volta
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // ignora os pares
		}
		if i > 7 {
			break // pára quando passa de 7
		}
		fmt.Print(i, " ") // 1 3 5 7
	}
	fmt.Println()

	// Ciclo infinito: for { ... }  (sai-se com break ou return)

	// =================================================================
	// 12. SWITCH — não precisa de "break" (não "cai" para o seguinte)
	// =================================================================
	dia := 3
	switch dia {
	case 1:
		fmt.Println("Segunda")
	case 2, 3, 4: // vários valores no mesmo caso
		fmt.Println("Meio da semana")
	default:
		fmt.Println("Outro dia")
	}

	// Switch sem valor: funciona como uma cadeia de if/else
	switch {
	case idade < 18:
		fmt.Println("Menor")
	default:
		fmt.Println("Adulto")
	}

	// =================================================================
	// 13. ARRAYS — tamanho fixo
	// =================================================================
	var notas [3]int // [0 0 0]
	notas[0] = 15    // aceder por índice (começa em 0)
	cores := [3]string{"vermelho", "verde", "azul"}
	fmt.Println(notas, cores, len(cores))

	// =================================================================
	// 14. SLICES — tamanho variável (os mais usados)
	// =================================================================
	numeros := []int{1, 2, 3}
	numeros = append(numeros, 4, 5) // acrescentar elementos
	parte := numeros[1:3]           // do índice 1 ao 3 (exclusive) → [2 3]
	reservado := make([]int, 0, 10) // len 0, cap 10
	fmt.Println(numeros, parte, len(reservado), cap(reservado))

	// Percorrer com range: índice e valor
	for i, v := range numeros {
		fmt.Print(i, ":", v, " ")
	}
	fmt.Println()

	// =================================================================
	// 15. MAPS — pares chave → valor (como um dicionário)
	// =================================================================
	idades := map[string]int{
		"Ana": 25,
		"Rui": 30,
	}
	idades["Eva"] = 22    // acrescentar ou alterar
	delete(idades, "Rui") // remover

	// Verificar se uma chave existe ("comma ok")
	if valor, existe := idades["Ana"]; existe {
		fmt.Println("Ana tem", valor, "anos")
	}

	// Percorrer (atenção: a ordem NÃO é garantida)
	for chave, valor := range idades {
		fmt.Println(chave, "→", valor)
	}

	// =================================================================
	// 16. USAR FUNÇÕES
	// =================================================================
	fmt.Println(somar(2, 3), multiplicar(4, 5))

	quociente, resto := dividirInteiros(17, 5)
	fmt.Println("17 / 5 =", quociente, "resto", resto)

	_, soResto := dividirInteiros(9, 4) // "_" ignora um valor que não queremos
	fmt.Println("Resto:", soResto)

	fmt.Println(somarTodos(1, 2, 3, 4)) // vários argumentos
	fmt.Println(somarTodos(numeros...)) // passar um slice inteiro

	// Função anónima guardada numa variável
	dobro := func(v int) int { return v * 2 }
	fmt.Println(dobro(21))

	// =================================================================
	// 17. TRATAR ERROS — o padrão "if err != nil"
	// =================================================================
	resultado, err := dividir(10, 2)
	if err != nil {
		fmt.Println("Erro:", err)
		return // em programas reais: sair, ou devolver o erro para cima
	}
	fmt.Println("10 / 2 =", resultado)

	// =================================================================
	// 18. STRUCTS E MÉTODOS EM USO
	// =================================================================
	p := Pessoa{Nome: "Ana", Idade: 25}
	fmt.Println(p.Nome) // aceder a um campo
	fmt.Println(p.Apresentar())
	p.FazerAnos()               // altera p, porque o recetor é ponteiro
	fmt.Println(p.Apresentar()) // agora com 26

	// =================================================================
	// 19. PONTEIROS — guardam o ENDEREÇO de uma variável
	// =================================================================
	valor := 10
	ptr := &valor      // & = "endereço de"
	*ptr = 20          // * = "o valor que está nesse endereço"
	fmt.Println(valor) // 20 — alterámos valor através do ponteiro

	// =================================================================
	// 20. INTERFACES EM USO
	// =================================================================
	formas := []Forma{Retangulo{3, 4}, Circulo{1}}
	for _, f := range formas {
		fmt.Printf("%T tem área %.2f\n", f, f.Area()) // %T mostra o tipo
	}

	// =================================================================
	// 21. STRINGS E FORMATAÇÃO
	// =================================================================
	s := "Olá, Mundo"
	fmt.Println(len(s))                      // nº de BYTES (o "á" ocupa 2)
	fmt.Println(strings.ToUpper(s))          // OLÁ, MUNDO
	fmt.Println(strings.Contains(s, "Mun"))  // true
	fmt.Println(strings.Split("a,b,c", ",")) // [a b c]

	// Verbos de formatação mais usados com Printf / Sprintf
	fmt.Printf("%d inteiro | %.2f decimal | %s texto | %v qualquer | %t bool\n",
		42, 3.14159, "Go", p, true)
	fmt.Printf("%+v\n", p) // struct com os nomes dos campos

	// =================================================================
	// 22. DEFER — adia a execução até a função terminar
	// =================================================================
	// Muito usado para "arrumar": fechar ficheiros, ligações, etc.
	defer fmt.Println("Fim do programa (isto corre em último lugar)")

	// =================================================================
	// 23. GOROUTINES E CANAIS — concorrência (só uma amostra)
	// =================================================================
	canal := make(chan string) // canal para enviar strings entre goroutines

	go func() { // "go" corre a função em paralelo (numa goroutine)
		canal <- "mensagem da goroutine" // enviar para o canal
	}()

	mensagem := <-canal // receber do canal (espera até chegar)
	fmt.Println(mensagem)
}
