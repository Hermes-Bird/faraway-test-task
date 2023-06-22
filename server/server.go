package server

import (
	"context"
	"fmt"
	"github.com/Hermes-Bird/faraway-test-task.git/config"
	"github.com/Hermes-Bird/faraway-test-task.git/internal/challenge"
	"github.com/Hermes-Bird/faraway-test-task.git/internal/qoutes"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	port        string
	connTimeout time.Duration
	listener    net.Listener
	maxConn     int32
	mapLock     sync.Mutex
	connMap     map[string]net.Conn
	connAmount  atomic.Int32
}

func NewServer(cfg config.ServerConfig) *Server {
	return &Server{
		port:        cfg.Port,
		maxConn:     cfg.MaxConn,
		mapLock:     sync.Mutex{},
		connMap:     map[string]net.Conn{},
		connTimeout: time.Duration(cfg.ConnTimeout) * time.Millisecond,
	}
}

func (s *Server) HandleCancel(ctx context.Context, quitCh chan struct{}) {
	<-ctx.Done()
	s.listener.Close()

	s.mapLock.Lock()
	for id, con := range s.connMap {
		log.Printf("closing connection %s", id)
		delete(s.connMap, id)
		con.Close()
	}
	s.mapLock.Unlock()

	quitCh <- struct{}{}
}

func (s *Server) Start(ctx context.Context, quitCh chan struct{}) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))

	if err != nil {
		return err
	}
	s.listener = listener
	log.Printf("start listening on port %s\n", s.port)

	go s.HandleCancel(ctx, quitCh)

	for {
		if s.connAmount.Load() > s.maxConn {
			continue
		}

		con, err := listener.Accept()
		if err != nil {
			continue
		}

		err = con.SetDeadline(time.Now().Add(s.connTimeout))
		if err != nil {
			log.Printf("failed to set deadline for connection %s", con.RemoteAddr().String())
			con.Close()
			continue
		}

		s.connAmount.Add(1)
		conId := GenerateConId(con)

		s.mapLock.Lock()
		s.connMap[conId] = con
		s.mapLock.Unlock()

		go s.handleConnection(con, conId)
	}
}
func (s *Server) closeConnection(id string) {
	log.Printf("closing connection %s\n", id)
	s.mapLock.Lock()

	s.connMap[id].Close()
	delete(s.connMap, id)
	s.connAmount.Add(-1)

	s.mapLock.Unlock()

}

func GenerateConId(con net.Conn) string {
	return con.RemoteAddr().String()
}

func (s *Server) handleConnection(con net.Conn, id string) {
	defer s.closeConnection(id)
	log.Printf("start handling new connection %s\n", id)
	ch := challenge.GetChallenge()

	_, err := con.Write(ch)
	if err != nil {
		log.Printf("error while writing challenge to connection %s: %s\n", id, err.Error())
		return
	}
	log.Printf("challenge generated and sended to connection %s\n", id)

	bs := make([]byte, 4)
	n, err := con.Read(bs)
	if err != nil && err != io.EOF {
		log.Printf("error while reading solution from connection %s: %s\n", id, err.Error())
		return
	}

	log.Printf("geting solution from connection %s - %d\n", id, bs[:n])

	if !challenge.CheckChallenge(ch, bs[:n]) {
		log.Printf("solution is incorrect")
		return
	}

	log.Printf("solution are right for connection %s sending a qoute\n", id)

	_, err = con.Write(qoutes.GetQuoteBytes())
	if err != nil {
		return
	}
}
