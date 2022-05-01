package main

import (
	"github.com/urfave/cli"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	msgr "simpleMessenger/pkg"
	"syscall"
)

const (
	DefaultPort    = "13005"
	DefaultLogFile = "log.txt"
)

type Server struct {
	port    string
	logfile *os.File
}

var server Server

func ParseArgs() error {
	app := cli.NewApp()
	app.Name = "simpleMessenger"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: DefaultPort,
			Usage: "port to listen on",
		},
		&cli.StringFlag{
			Name:  "logfile",
			Value: DefaultLogFile,
			Usage: "file logs are written to",
		},
	}
	app.Action = func(c *cli.Context) error {
		server.port = c.String("port")
		logfile := c.String("logfile")

		var err error
		server.logfile, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(server.logfile)
		return nil
	}
	return app.Run(os.Args)
}

func Setup() (*net.TCPListener, error) {
	msgrRPC := new(Msgr)
	err := rpc.Register(msgrRPC)
	if err != nil {
		return nil, err
	}

	t, err := net.ResolveTCPAddr("tcp", "localhost:"+server.port)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", t)
	if err != nil {
		return nil, err
	}

	err = msgr.Start()
	if err != nil {
		log.Fatalf("messenger connection failed to be established: %s", err)
	}

	log.Printf("listen: %s", server.port)
	return l, nil
}

func main() {
	err := ParseArgs()
	if err != nil {
		return
	}

	l, err := Setup()
	if err != nil {
		log.Println(err.Error())
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		err = msgr.ShutDown()
		if err != nil {
			log.Fatalf("couldnot properly shut down messenger connection: %s", err)
		}
	}()

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				continue
			}
			log.Printf("client: %s", c.RemoteAddr())
			go jsonrpc.ServeConn(c)
		}
	}()

	<-done
}
