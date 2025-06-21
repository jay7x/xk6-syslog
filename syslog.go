package syslog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/syslog", new(Syslog))
}

// Config defines how to connect/send
type Config struct {
	Transport string     `json:"transport"` // "udp", "tcp", or "tls"
	Timeout   int        `json:"timeout"`   // timeout in seconds
	TLS       *TLSConfig `json:"tls,omitempty"`
}

type TLSConfig struct {
	CA                 string `json:"ca,omitempty"`   // PEM-encoded CA cert
	ClientCert         string `json:"cert,omitempty"` // PEM-encoded client cert
	ClientKey          string `json:"key,omitempty"`  // PEM-encoded private key
	ServerName         string `json:"serverName,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty"`
}

// Syslog is the main module struct
type Syslog struct{}

// Connect opens a new connection to a syslog server
func (s *Syslog) Connect(address string, config Config) (*Connection, error) {
	var conn net.Conn
	var err error

	switch config.Transport {
	case "tcp":
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}
		conn, err = net.DialTCP("tcp", nil, addr)
		if err != nil {
			return nil, err
		}

	case "tls":
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		tcpConn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return nil, err
		}

		tlsCfg := &tls.Config{
			// #nosec G402 - it's only true if it's explicitly set by user
			InsecureSkipVerify: config.TLS.InsecureSkipVerify,
			ServerName:         config.TLS.ServerName,
		}

		if config.TLS.ClientCert != "" && config.TLS.ClientKey != "" {
			cert, err := tls.X509KeyPair([]byte(config.TLS.ClientCert), []byte(config.TLS.ClientKey))
			if err != nil {
				return nil, fmt.Errorf("failed to load client cert/key: %w", err)
			}
			tlsCfg.Certificates = []tls.Certificate{cert}
		}

		if config.TLS.CA != "" {
			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM([]byte(config.TLS.CA)); !ok {
				return nil, fmt.Errorf("failed to parse CA certificate")
			}
			tlsCfg.RootCAs = pool
		}

		tlsConn := tls.Client(tcpConn, tlsCfg)
		err = tlsConn.Handshake()
		if err != nil {
			return nil, err
		}

		conn = tlsConn

	default: // "udp"
		addr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			return nil, err
		}
		conn, err = net.DialUDP("udp", nil, addr)
		if err != nil {
			return nil, err
		}
	}

	if config.Timeout > 0 {
		err = conn.SetDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		if err != nil {
			return nil, err
		}
	}

	return &Connection{
		conn:      conn,
		transport: config.Transport,
	}, nil
}

// Connection holds an open network connection and transport type
type Connection struct {
	conn      net.Conn
	transport string // "udp", "tcp", or "tls"
}

// Send sends raw bytes
func (c *Connection) Send(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Exports returns the module's JS API
func (s *Syslog) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]any{
			"connect": s.Connect,
		},
	}
}
