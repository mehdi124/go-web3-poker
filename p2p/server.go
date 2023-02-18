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
	msgChan chan *Message
}

func NewServer(cfg ServerConfig) *Server {

	return &Server{
		handler:      *NewHandler(),
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
	}
}

func (s *Server) Start() {

	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server is running on port %s\n", s.ServerConfig.ListenAddr)
	s.acceptLoop()
}

func (s *Server) handleConn(conn net.Conn) {

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		s.msgChan <- &Message{
			From:    conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}
	}

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

		peer.Send([]byte("mehdi124 v0.1-alpha"))
		go s.handleConn(conn)
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
