package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/register-traefik/etcd"
	"github.com/octoblu/register-traefik/healthchecker"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("register-traefik:main")

func main() {
	app := cli.NewApp()
	app.Name = "register-traefik"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "etcd-uri, e",
			EnvVar: "REGISTER_TRAEFIK_ETCD_URI",
			Usage:  "Etcd URI to register traefik server to",
		},
		cli.StringFlag{
			Name:   "server-key, s",
			EnvVar: "REGISTER_TRAEFIK_SERVER_KEY",
			Usage:  "Etcd key to register traefik server to",
		},
		cli.StringFlag{
			Name:   "uri, u",
			EnvVar: "REGISTER_TRAEFIK_URI",
			Usage:  "URI to healthcheck, must return status 200",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	etcdURI, serverKey, uri := getOpts(context)

	sigTerm := make(chan os.Signal)
	signal.Notify(sigTerm, syscall.SIGTERM)

	sigTermReceived := false

	go func() {
		<-sigTerm
		fmt.Println("SIGTERM received, waiting to exit")
		sigTermReceived = true
	}()

	for {
		if sigTermReceived {
			fmt.Println("Cleaning up server.")
			onNotHealthy(etcdURI, serverKey)
			fmt.Println("I'll be back.")
			os.Exit(0)
		}

		loop(etcdURI, serverKey, uri)
	}
}

func loop(etcdURI, serverKey, uri string) {
	if healthchecker.Healthy(fmt.Sprintf("%v/healthcheck", uri)) {
		onHealthy(etcdURI, serverKey, uri)
	} else {
		onNotHealthy(etcdURI, serverKey)
	}

	time.Sleep(5 * time.Second)
}

func onHealthy(etcdURI, serverKey, uri string) {
	debug("onHealthy")
	var err error

	urlKey := fmt.Sprintf("%v/url", serverKey)
	weightKey := fmt.Sprintf("%v/weight", serverKey)

	err = etcd.Set(etcdURI, urlKey, uri)
	FatalIfError("etcdclient.Set urlKey", err)
	err = etcd.Set(etcdURI, weightKey, "10")
	FatalIfError("etcdclient.Set weightKey", err)
}

func onNotHealthy(etcdURI, serverKey string) {
	debug("onNotHealthy")

	err := etcd.DelDir(etcdURI, serverKey)
	FatalIfError("etcd.DelDir", err)
}

func getOpts(context *cli.Context) (string, string, string) {
	etcdURI := context.String("etcd-uri")
	serverKey := context.String("server-key")
	uri := context.String("uri")

	if etcdURI == "" || serverKey == "" || uri == "" {
		cli.ShowAppHelp(context)

		if etcdURI == "" {
			color.Red("  Missing required flag --etcd-uri or REGISTER_TRAEFIK_ETCD_URI")
		}
		if serverKey == "" {
			color.Red("  Missing required flag --server-key or REGISTER_TRAEFIK_SERVER_KEY")
		}
		if uri == "" {
			color.Red("  Missing required flag --uri or REGISTER_TRAEFIK_URI")
		}
		os.Exit(1)
	}

	return etcdURI, serverKey, uri
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Fatalln(errorMessage, err.Error())
	}

	return version.String()
}

// FatalIfError prints error and dies if error is non nil
func FatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalf("ERROR(%v):\n\n%v", msg, err)
}
