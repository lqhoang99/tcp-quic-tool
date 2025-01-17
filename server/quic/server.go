package quic

import (
	"crypto/tls"
	"github.com/lqhoang99/tcp-quic-tools/server/util"
	"github.com/lqhoang99/tcp-quic-tools/util/cli"
	"github.com/lqhoang99/tcp-quic-tools/util/connection_type"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"sync"
)

// QUIC Server implementation.
type Server struct {
	// TLS configuration to use
	tlsConfig tls.Config
}

// Create new QUIC server.
func NewServer(options *cli.Options) (*Server, error) {
	server := Server{
		tlsConfig: options.TlsConfiguration,
	}

	return &server, nil
}

func (s *Server) GetType() connection_type.ConnectionType {
	return connection_type.QUIC
}

func (s *Server) Listen(addr *string) (*sync.WaitGroup, error) {
	log.Println("Setting up QUIC listener")
	listener, err := quic.ListenAddr(*addr, &s.tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go s.listen(listener, &wg)

	return &wg, nil
}

// Start listening to incoming connections.
func (s *Server) listen(listener quic.Listener, wg *sync.WaitGroup) {
	for {
		sess, err := listener.Accept()
		if err != nil {
			log.Println("QUIC server failed while listening to incoming connection requests. Cancelling listening.")
			break
		}

		wg.Add(1)
		go s.inSession(&sess, wg)
	}
}

func (s *Server) inSession(session *quic.Session, wg *sync.WaitGroup) {
	log.Println("Accepted new connection")

	stream, err := (*session).AcceptStream()
	if err != nil {
		log.Println("QUIC session failed while trying to accept stream. Cancelling session.")
		wg.Done()
		return
	}

	_, err = io.Copy(util.LoggingWriter{
		Writer: stream,
	}, stream)

	wg.Done()
}
