package server

import "sync"

type ClientManager struct {
	clients map[*Client]struct{}
	mx *sync.Mutex
}


func NewClientManager() *ClientManager{
	return &ClientManager{
		clients: make(map[*Client]struct{}),
		mx: &sync.Mutex{},
	}	
}


func (m *ClientManager) AddClient(c *Client){

	m.mx.Lock()
	defer m.mx.Unlock()

	if _, ok := m.clients[c]; ok{
		return
	}

	m.clients[c]=struct{}{}
}


func (m *ClientManager) RemoveClient(c *Client){
	m.mx.Lock()
	defer m.mx.Unlock()

	
}