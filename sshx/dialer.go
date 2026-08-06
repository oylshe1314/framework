package sshx

import (
	"context"
	"net"

	"golang.org/x/crypto/ssh"
)

type SshDialer struct {
	c *ssh.Client
}

func (this *SshDialer) Dial(network, address string) (c net.Conn, err error) {
	return this.c.Dial(network, address)
}

func (this *SshDialer) DialContext(ctx context.Context, network, address string) (c net.Conn, err error) {
	return this.c.DialContext(ctx, network, address)
}
