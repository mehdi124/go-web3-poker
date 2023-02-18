package p2p

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

func (gv GameVariant) String() string {

	switch gv {
	case TexasHoldem:
		return "TEXAS_HOLDEM"
	case Other:
		return "other"
	default:
		return "unknown"
	}
}

const (
	TexasHoldem GameVariant = iota
	Other
)

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	handler   Handler
	transport *TCPTransport
	listener  net.Listener
	peers     map[net.Addr]*Peer
	ServerConfig
	addPeer chan *Peer
	delPeer chan *Peer
	msgChan chan *Message
}

func NewServer(cfg ServerConfig) *Server {

	s := &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer

	return s
}

func (s *Server) Start() {

	go s.loop()

	fmt.Printf("game server is running on port %s\n", s.ListenAddr)
	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("started new game server")

	s.transport.ListenAndAccept()
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

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player disconnected")

			fmt.Printf("player disconnected %s\n", peer.conn.RemoteAddr())
			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:

			go peer.ReadLoop(s.msgChan)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")

			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())
			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgChan:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}

		}
	}
}
