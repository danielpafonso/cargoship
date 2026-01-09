package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"cargoship/cmd/shipper/transport"
	"cargoship/internal/configurations"
	"cargoship/internal/logging"
	"cargoship/internal/manifests"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	scriptLogger logging.Logger
	filesLogger  logging.Logger
)

// ftpServer process configured services for a ftp server
func ftpServer(server configurations.FtpConfig, times *[]manifests.ShipperManifest) {
	// create connection to ftp
	ftpUrl := fmt.Sprintf("%s:%d", server.Hostname, server.Port)
	conn, err := ftp.Dial(ftpUrl, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	// login
	err = conn.Login(server.User, server.Pass)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	// service loop
	for _, service := range server.Services {
		if !service.Enable {
			continue
		}
		switch service.Mode {
		case "import":
			transport.DownloadFiles(server.Name, conn, service, times, scriptLogger, filesLogger)
		case "export":
			transport.UploadFiles(server.Name, conn, service, times, scriptLogger, filesLogger)
		default:
			scriptLogger.LogWarn(
				fmt.Sprintf("ERROR Unknown mode, %s, on service %s.\n", service.Mode, service.Name),
			)
		}
	}
}

// sftpServer process configured services for a sftp server
func sftpServer(server configurations.FtpConfig, times *[]manifests.ShipperManifest) {
	// create connection to sftp
	config := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Hostname, server.Port), config)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	// create sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	defer client.Close()
	// service loop
	for _, service := range server.Services {
		if !service.Enable {
			continue
		}
		if service.Mode == "import" {
			transport.SftpDownloadFiles(server.Name, client, service, times, scriptLogger, filesLogger)
		} else if service.Mode == "export" {
			transport.SftpUploadFiles(server.Name, client, service, times, scriptLogger, filesLogger)
		} else {
			scriptLogger.LogWarn(
				fmt.Sprintf("ERROR Unknown mode, %s, on service %s.\n", service.Mode, service.Name),
			)
		}
	}
}

func main() {
	// command line flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "shipper_config.yaml", "Path to configuration yaml")
	flag.StringVar(&configFilepath, "c", "shipper_config.yaml", "")
	flag.Usage = func() {
		fmt.Print(`Usage of shipper:
  -c, --config  path to configuation yaml 
  -h, --help    display this help message
`)
	}
	flag.Parse()

	// read script configuration
	configs, err := configurations.ShipperReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}

	// start loggers
	scriptLogger.Init(configs.Log.Script, configs.Log2Console)
	filesLogger.Init(configs.Log.Files, configs.Log2Console)
	defer scriptLogger.Close()
	defer filesLogger.Close()

	// read ftp times state
	times, err := manifests.ShipperReadTimes(configs.TimesPath)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}

	// defer update/write ftp time state file
	defer manifests.ShipperWriteTimes(&times, configs.TimesPath)

	for _, server := range configs.Ftps {
		scriptLogger.LogInfo(fmt.Sprintf("Connect to server: %s\n", server.Name))
		if server.Protocol == "ftp" {
			ftpServer(server, &times)
		} else if server.Protocol == "sftp" {
			sftpServer(server, &times)
		} else {
			scriptLogger.LogWarn(
				fmt.Sprintf("ERROR Unknown protocol, %s, on server %s.\n", server.Protocol, server.Name),
			)
		}
	}
}
