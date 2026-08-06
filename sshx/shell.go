package sshx

import (
	"context"

	"github.com/oylshe1314/framework/errors"
	"golang.org/x/crypto/ssh"
)

type ShellClient struct {
	s *ssh.Session
}

func (this *ShellClient) Init(ctx context.Context) error {
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

	var err error
	err = this.s.RequestPty("xterm-256color", 640, 400, modes)
	if err != nil {
		return err
	}

	err = this.s.Shell()
	if err != nil {
		return errors.Error("failed to start shell: ", err)
	}

	err = this.s.Wait()
	if err != nil {
		return errors.Error("session ended with error: ", err)
	}

	return nil
}

func (this *ShellClient) Work() error {
	return this.work()
}
