package sshx

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpClient struct {
	c  *ssh.Client
	fc *sftp.Client
}

func (this *SftpClient) Init(ctx context.Context) error {
	fc, err := sftp.NewClient(this.c)
	if err != nil {
		return err
	}

	this.fc = fc
	return nil
}

func (this *SftpClient) Close() error {
	if this.fc != nil {
		return this.fc.Close()
	}
	return nil
}

// Upload uploads a local file to a remote location.
// localPath specifies the path to the local file.
// remotePath specifies the remote destination:
//   - If it ends with "/", it is treated as a remote directory path,
//     and the remote file will keep the same base name as the local file.
//   - If it does NOT end with "/", it is treated as the full remote file path,
//     allowing the caller to customize the remote file name.
func (this *SftpClient) Upload(localPath, remotePath string) error {
	var filename = filepath.Base(localPath)

	lf, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer lf.Close()

	if remotePath[len(remotePath)-1] == '/' {
		err = this.fc.MkdirAll(remotePath)
		if err != nil {
			return err
		}
		remotePath = filepath.Join(remotePath, filename)
	} else {
		err = this.fc.MkdirAll(filepath.Dir(remotePath))
		if err != nil {
			return err
		}
	}

	rf, err := this.fc.Create(remotePath)
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
