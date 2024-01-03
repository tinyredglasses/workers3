package sockets

import (
	"context"
	"net"
	"time"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type SecureTransport string

const (
	// SecureTransportOn indicates "Use TLS".
	SecureTransportOn SecureTransport = "on"
	// SecureTransportOff indicates "Do not use TLS".
	SecureTransportOff SecureTransport = "off"
	// SecureTransportStartTLS indicates "Do not use TLS initially, but allow the socket to be upgraded
	// to use TLS by calling *Socket.StartTLS()".
	SecureTransportStartTLS SecureTransport = "starttls"
)

type SocketOptions struct {
	SecureTransport SecureTransport `json:"secureTransport"`
	AllowHalfOpen   bool            `json:"allowHalfOpen"`
}

const defaultDeadline = 999999 * time.Hour

func Connect(ctx context.Context, addr string, opts *SocketOptions) (net.Conn, error) {
	connect, err := cfruntimecontext.GetRuntimeContextValue(ctx, "connect")
	if err != nil {
		return nil, err
	}
	optionsObj := jsutil.NewObject()
	if opts != nil {
		if opts.AllowHalfOpen {
			optionsObj.Set("allowHalfOpen", true)
		}
		if opts.SecureTransport != "" {
			optionsObj.Set("secureTransport", opts.SecureTransport)
		}
	}
	sockVal := connect.Invoke(addr, optionsObj)
	deadline := time.Now().Add(defaultDeadline)
	return newSocket(ctx, sockVal, deadline, deadline), nil
}
