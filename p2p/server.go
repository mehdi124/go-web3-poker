package p2p

import (
	"bytes"
	"encoding/gob"
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
	transport *TCPTransport
	peers     map[net.Addr]*Peer
	ServerConfig
	addPeer chan *Peer
	delPeer chan *Peer
	msgChan chan *Message
}

func NewServer(cfg ServerConfig) *Server {

	s := &Server{
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

	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("game server started with info :")

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
			}).Info("player disconnected")

			fmt.Printf("player disconnected %s\n", peer.conn.RemoteAddr())
			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			s.SendHandshake(peer)

			if err := s.handshake(peer); err != nil {
				logrus.Errorf("Handshake with incoming player failed: %s", err)
				continue
			}

			go peer.ReadLoop(s.msgChan)

			logrus.WithFields(logrus.Fields{
				"addr":            peer.conn.RemoteAddr(),
				"listener_server": s.ServerConfig.ListenAddr,
			}).Info("handshake successfull: player added")

			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}

		}
	}
}

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		Version:     s.Version,
		GameVariant: s.GameVariant,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())

}

type Handshake struct {
	Version     string
	GameVariant GameVariant
}

func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("invalid game variant %s", hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("invalid version %s", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Logger.Info("received handshake")

	return nil
}

func (s *Server) handleMessage(msg *Message) error {

	fmt.Printf("%+v\n", msg)
	return nil
}
