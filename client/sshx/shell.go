package sshx

import (
	"context"
	"framework/errors"
	"io"

	"golang.org/x/crypto/ssh"
)

type ShellClient struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer

	c *ssh.Client
	s *ssh.Session
}

func (this *ShellClient) Init(ctx context.Context) error {
	s, err := this.c.NewSession()
	if err != nil {
		return errors.Error("failed to create session: ", err)
	}

	s.Stdin = this.in
	s.Stdout = this.out
	s.Stderr = this.errOut

	this.s = s
	return nil
}

func (this *ShellClient) Close() error {
	if this.s != nil {
		return this.s.Close()
	}
	return nil
}

func (this *ShellClient) work() error {
	var modes = ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := this.s.RequestPty("xterm", 80, 40, modes); err != nil {
		return errors.Error("request for pseudo terminal failed: ", err)
	}

	if err := this.s.Shell(); err != nil {
		return errors.Error("failed to start shell: ", err)
	}

	if err := this.s.Wait(); err != nil {
		return errors.Error("session ended with error: ", err)
	}

	return nil
}

func (this *ShellClient) Work() error {
	return this.work()
}
