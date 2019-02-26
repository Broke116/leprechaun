package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/cron"
	"github.com/kilgaloon/leprechaun/server"
)

// VERSION of application
const (
	VERSION = "1.0.0-rc"
	RELEASE = "Calimero"
)

func main() {
	shutdownSignal := make(chan os.Signal, 1)

	iniPath := flag.String("ini", "/etc/leprechaun/configs/config.ini", "Path to .ini configuration")
	pidPath := flag.String("pid", "/var/run/leprechaun/.pid", "PID file of process")
	cmd := flag.String("cmd", "run", "Command for app to run")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	c := strings.Fields(*cmd)[0]

	if !api.IsAPIRunning() {
		configs := config.NewConfigs()
		client.New("client", configs.New("client", *iniPath), *debug)
		server.New("server", configs.New("server", *iniPath), *debug)
		cron.New("cron", configs.New("cron", *iniPath), *debug)

		a := api.New("")
		a.Register(client.Agent)
		a.Register(server.Agent)
		a.Register(cron.Agent)
		go a.Start()

		d.New(
			pidPath, 
			os.Args, 
			api,
			[client.Agent, server.Agent, cron.Agent]
		)
	}

	switch c {
	case "run":
		go client.Agent.Start()
		go server.Agent.Start()
		go cron.Agent.Start()
	case "client:start":
		go client.Agent.Start()
	case "server:start":
		go server.Agent.Start()
	case "cron:start":
		go cron.Agent.Start()
	case "client:stop":
		*cmd = "client stop"
		fallthrough
	case "server:stop":
		*cmd = "server stop"
		fallthrough
	case "cron:stop":
		*cmd = "cron stop"
		fallthrough
	default:
		api.Resolver(api.Cmd(*cmd))
		os.Exit(0)
	}

	signal.Notify(shutdownSignal,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGSTOP,
		syscall.SIGTERM)

	<-shutdownSignal

	os.Exit(0)
}
