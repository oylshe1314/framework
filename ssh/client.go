package ssh

import (
	"bytes"
	"github.com/pkg/sftp"
	std "golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"strings"
)

type Client struct {
	c *std.Client
}

func Open(address, user, pass string, key []byte) (*Client, error) {
	var config = &std.ClientConfig{
		User:            user,
		HostKeyCallback: std.InsecureIgnoreHostKey(),
	}

	if pass != "" {
		config.Auth = append(config.Auth, std.Password(pass))
	}

	if key != nil {
		signer, err := std.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}

		config.Auth = append(config.Auth, std.PublicKeys(signer))
	}

	c, err := std.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}

	return &Client{c: c}, nil
}

func (this *Client) Close() error {
	return this.c.Close()
}

func (this *Client) Exec(args ...string) (string, error) {
	s, err := this.c.NewSession()
	if err != nil {
		return "", err
	}

	defer s.Close()

	var o, e = new(bytes.Buffer), new(bytes.Buffer)

	s.Stdout = o
	s.Stderr = e

	err = s.Run(strings.Join(args, " "))
	if err != nil {
		return e.String(), err
	}

	return o.String(), nil
}

func (this *Client) Upload(localFile, remotePath string) error {
	fc, err := sftp.NewClient(this.c)
	if err != nil {
		return err
	}

	defer fc.Close()

	lf, err := os.Open(localFile)
	if err != nil {
		return err
	}

	defer lf.Close()

	if remotePath[len(remotePath)-1] == '/' {
		err = fc.MkdirAll(remotePath)
		if err != nil {
			return err
		}
		remotePath = path.Join(remotePath, path.Base(localFile))
	} else {
		err = fc.MkdirAll(path.Dir(remotePath))
		if err != nil {
			return err
		}
	}

	rf, err := fc.Create(remotePath)
	if err != nil {
		return err
	}

	defer rf.Close()

	_, err = io.Copy(rf, lf)
	if err != nil {
		return err
	}

	return nil
}
