package command

import (
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/ftpd"
	ftpserver "github.com/fclairamb/ftpserverlib"
)

var (
	ftpServer            *ftpserver.FtpServer
	ftpStandaloneOptions FtpOptions
)

type FtpOptions struct {
	s3Endpoint *string
	bucket     *string
	ip         *string
	publicIP   *string
	port       *int
	portStart  *int
	portStop   *int
}

func init() {
	cmdFtp.Run = runFTP
	ftpStandaloneOptions.s3Endpoint = cmdFtp.Flag.String("s3.endpoint", "localhost:8333", "s3 server address")
	ftpStandaloneOptions.bucket = cmdFtp.Flag.String("s3.bucket", "test", "s3 bucket")
	ftpStandaloneOptions.ip = cmdFtp.Flag.String("ip", "localhost", "ftp server ip address")
	ftpStandaloneOptions.publicIP = cmdFtp.Flag.String("public.ip", "localhost", "ftp server public ip address")
	ftpStandaloneOptions.port = cmdFtp.Flag.Int("port", 2121, "ftp server listen port")
	ftpStandaloneOptions.portStart = cmdFtp.Flag.Int("port.start", 2121, "ftp server listen start port")
	ftpStandaloneOptions.portStop = cmdFtp.Flag.Int("port.stop", 2130, "ftp server listen stop port")
}

var cmdFtp = &Command{
	UsageLine: "ftp [-port=2121] [-s3.endpoint=<ip:port>] [-s3.bucket=test]",
	Short:     "start a FTP API compatible server that is backed by a s3",
	Long:      "start a FTP API compatible server that is backed by a s3.",
}

func runFTP(cmd *Command, args []string) bool {
	return ftpStandaloneOptions.startS3Server()
}

func (ftpOpt *FtpOptions) startS3Server() bool {

	driver, _ := ftpd.NewFtpServer(&ftpd.FtpServerOption{
		S3Endpoint:       *ftpOpt.s3Endpoint,
		Bucket:           *ftpOpt.bucket,
		IpBind:           *ftpOpt.ip,
		PublicIP:         *ftpOpt.publicIP,
		Port:             *ftpOpt.port,
		PassivePortStart: *ftpOpt.portStart,
		PassivePortStop:  *ftpOpt.portStop,
	})
	ftpServer = ftpserver.NewFtpServer(driver)

	fmt.Println("start ftp server success")

	if err := ftpServer.ListenAndServe(); err != nil {
		fmt.Println("start ftp server error")
	}

	return true

}
