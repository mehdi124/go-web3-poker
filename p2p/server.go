package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Server struct {
	handler  Handler
	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]*Peer
	ServerConfig
	addPeer chan *Peer
	delPeer chan *Peer
	msgChan chan *Message
}

func NewServer(cfg ServerConfig) *Server {

	return &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
	}
}

func (s *Server) Start() {

	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server is running on port %s\n", s.ListenAddr)
	s.acceptLoop()
}

func (s *Server) handleConn(p *Peer) {

	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		s.msgChan <- &Message{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}
	}

	s.delPeer <- p

}

func (s *Server) Connect(addr string) error {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer
	return peer.Send([]byte(s.Version))
}

func (s *Server) acceptLoop() {

	for {

		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer

		peer.Send([]byte(s.Version))
		go s.handleConn(peer)
	}

}

func (s *Server) listen() error {

	ln, err := net.Listen("tcp", s.ServerConfig.ListenAddr)
	if err != nil {
		return err
	}

	s.listener = ln

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			fmt.Printf("player disconnected %s\n", peer.conn.RemoteAddr())
			delete(s.peers, peer.conn.RemoteAddr())
		case peer := <-s.addPeer:
			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())
			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgChan:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}

		}
	}
}
