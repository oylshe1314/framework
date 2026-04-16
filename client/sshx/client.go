package sshx

import (
	"bytes"
	"context"
	"framework/errors"
	"framework/option"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	option.Optional[Option]

	c *ssh.Client
}

func NewSshClient(opt *Option) *SshClient {
	var c = &SshClient{}
	c.SetOption(opt)
	return c
}

func (this *SshClient) Init(ctx context.Context) error {
	if this.GetOption() == nil {
		return errors.New("'SshClient' option is nil")
	}

	var opt = this.GetOption()
	if opt.Address == "" {
		return errors.New("'SshClient' option.address is empty")
	}

	if opt.User == "" {
		return errors.New("'SshClient' option.user is empty")
	}

	return nil
}

func (this *SshClient) Close() error {
	if this.c != nil {
		return this.c.Close()
	}
	return nil
}

func (this *SshClient) Dial() error {
	var opt = this.GetOption()

	var config = &ssh.ClientConfig{
		User:            opt.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if opt.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(opt.Password))
	}

	var key []byte
	if opt.Key != "" {
		key = []byte(opt.Key)
	} else if opt.KeyPath != "" {
		f, err := os.Open(opt.KeyPath)
		if err != nil {
			return err
		}

		key, err = io.ReadAll(f)
		if err != nil {
			return err
		}
	}

	if key != nil {
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}

		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	c, err := ssh.Dial("tcp", opt.Address, config)
	if err != nil {
		return err
	}

	this.c = c
	return nil
}

func (this *SshClient) Exec(args ...string) (string, error) {
	s, err := this.c.NewSession()
	if err != nil {
		return "", err
	}

	defer func() {
		_ = s.Close()
	}()

	var out = new(bytes.Buffer)
	var errOut = new(bytes.Buffer)

	s.Stdout = out
	s.Stderr = errOut

	err = s.Run(strings.Join(args, " "))
	if err != nil {
		return errOut.String(), err
	}

	return out.String(), nil
}

func (this *SshClient) NewShellClient(in io.Reader, out, errOut io.Writer) (*ShellClient, error) {
	var sc = &ShellClient{
		in:     in,
		out:    out,
		errOut: errOut,
		c:      this.c,
	}

	var err = sc.Init(context.Background())
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (this *SshClient) NewSftpClient() (*SftpClient, error) {
	var sc = &SftpClient{
		c: this.c,
	}

	var err = sc.Init(context.Background())
	if err != nil {
		return nil, err
	}

	return sc, nil
}
