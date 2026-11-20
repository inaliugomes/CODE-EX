package main

import "sync"

type Armazem interface {
	Criar(l Link) error
	Procurar(codigo string) (Link, error)
	RegistarVisita(codigo string) (Link, error)
}

type memoria struct {
	mu    sync.Mutex
	links map[string]Link
}

func novaMemoria() *memoria {
	return &memoria{links: make(map[string]Link)}
}

func (m *memoria) Criar(l Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.links[l.Codigo]; ok {
		return ErrCodigoExiste
	}
	m.links[l.Codigo] = l
	return nil
}

func (m *memoria) Procurar(codigo string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[codigo]
	if !ok {
		return Link{}, ErrNaoExiste
	}
	return l, nil
}

func (m *memoria) RegistarVisita(codigo string) (Link, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.links[codigo]
	if !ok {
		return Link{}, ErrNaoExiste
	}
	l.Visitas++
	m.links[codigo] = l
	return l, nil
}
