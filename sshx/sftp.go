package sshx

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/oylshe1314/framework/errors"
	"github.com/pkg/sftp"
)

type SftpClient struct {
	c *sftp.Client
}

func (this *SftpClient) Init(ctx context.Context) error {
	return nil
}

func (this *SftpClient) Close() error {
	if this.c != nil {
		return this.c.Close()
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
	var filename = path.Base(localPath)

	lf, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer lf.Close()

	if remotePath[len(remotePath)-1] == '/' {
		err = this.c.MkdirAll(remotePath)
		if err != nil {
			return err
		}
		remotePath = path.Join(remotePath, filename)
	} else {
		err = this.c.MkdirAll(path.Dir(remotePath))
		if err != nil {
			return err
		}
	}

	rf, err := this.c.Create(remotePath)
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

func (this *SftpClient) Download(localPath, remotePath string) error {
	var filename = path.Base(localPath)

	rf, err := this.c.Open(remotePath)
	if err != nil {
		return err
	}

	defer rf.Close()

	if localPath[len(remotePath)-1] == '/' {
		err = os.MkdirAll(localPath, os.ModePerm)
		if err != nil {
			return err
		}
		localPath = path.Join(localPath, filename)
	} else {
		err = os.MkdirAll(path.Dir(localPath), os.ModePerm)
		if err != nil {
			return err
		}
	}

	lf, err := os.Create(localPath)
	if err != nil {
		return err
	}

	defer lf.Close()

	_, err = io.Copy(lf, rf)
	if err != nil {
		return err
	}

	return nil
}

func (this *SftpClient) Write(buff []byte, remoteFile string) error {
	if remoteFile[len(remoteFile)-1] == '/' {
		return errors.New("invalid remote filename")
	}

	var err = this.c.MkdirAll(path.Dir(remoteFile))
	if err != nil {
		return err
	}

	rf, err := this.c.Create(remoteFile)
	if err != nil {
		return err
	}

	defer rf.Close()

	_, err = rf.Write(buff)
	if err != nil {
		return err
	}

	return nil
}
