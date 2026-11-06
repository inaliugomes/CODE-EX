package main

import (
	"fmt"
	"runtime"
	"time"
)

// Armadilha: goroutine presa (leak). O channel SEM buffer obriga a goroutine
// a esperar por um receptor que nunca vem, porque o chamador desistiu.
func comLeak() {
	ch := make(chan string) // sem buffer
	go func() {
		time.Sleep(20 * time.Millisecond)
		ch <- "resultado" // fica presa aqui para sempre
	}()
	select {
	case v := <-ch:
		fmt.Println(v)
	case <-time.After(5 * time.Millisecond):
		fmt.Println("desisti")
	}
}

// Correção: buffer de 1. O envio cabe no buffer e a goroutine termina.
func semLeak() {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		ch <- "resultado"
	}()
	select {
	case v := <-ch:
		fmt.Println(v)
	case <-time.After(5 * time.Millisecond):
		fmt.Println("desisti")
	}
}

func presas(f func()) int {
	antes := runtime.NumGoroutine()
	f()
	time.Sleep(100 * time.Millisecond)
	return runtime.NumGoroutine() - antes
}

func main() {
	fmt.Println("goroutines presas (sem buffer):", presas(comLeak))
	fmt.Println("goroutines presas (com buffer):", presas(semLeak))
}
