package ftpd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
)

type FtpServerOption struct {
	S3Endpoint       string
	Bucket           string
	PublicIP         string
	IpBind           string
	Port             int
	PassivePortStart int
	PassivePortStop  int
}

type Server struct {
	option *FtpServerOption
}

// NewServer returns a new FTP server driver
func NewFtpServer(option *FtpServerOption) (*Server, error) {
	var err error
	return &Server{
		option: option,
	}, err
}

// GetSettings returns some general settings around the server setup
func (s *Server) GetSettings() (*ftpserver.Settings, error) {
	var portRange *ftpserver.PortRange
	if s.option.PassivePortStart > 0 && s.option.PassivePortStop > s.option.PassivePortStart {
		portRange = &ftpserver.PortRange{
			Start: s.option.PassivePortStart,
			End:   s.option.PassivePortStop,
		}
	}

	return &ftpserver.Settings{
		ListenAddr:               fmt.Sprintf("%s:%d", s.option.IpBind, s.option.Port),
		PublicHost:               s.option.PublicIP,
		PassiveTransferPortRange: portRange,
		ActiveTransferPortNon20:  true,
		IdleTimeout:              -1,
		ConnectionTimeout:        20,
	}, nil
}

// ClientConnected is called to send the very first welcome message
func (s *Server) ClientConnected(cc ftpserver.ClientContext) (string, error) {
	return "Welcome to SeaweedFS FTP Server", nil
}

// ClientDisconnected is called when the user disconnects, even if he never authenticated
func (s *Server) ClientDisconnected(cc ftpserver.ClientContext) {
}

// AuthUser authenticates the user and selects an handling driver
func (s *Server) AuthUser(cc ftpserver.ClientContext, username, password string) (ftpserver.ClientDriver, error) {
	accFs, _ := LoadFs(s, username, password)

	return &ClientDriver{
		Fs: accFs,
	}, nil
}

// GetTLSConfig returns a TLS Certificate to use
// The certificate could frequently change if we use something like "let's encrypt"
func (s *Server) GetTLSConfig() (*tls.Config, error) {
	return nil, errors.New("no TLS certificate configured")
}

type ClientDriver struct {
	afero.Fs
}

// LoadFs loads a file system from an access description
func LoadFs(s *Server, username string, password string) (afero.Fs, error) {
	endpoint := s.option.S3Endpoint
	region := "default"
	bucket := s.option.Bucket

	sess, errSession := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(username, password, ""),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	if errSession != nil {
		return nil, errSession
	}

	s3Fs := NewFs(bucket, sess)

	return s3Fs, nil
}
