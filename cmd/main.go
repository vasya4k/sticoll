package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"time"

	auth_pb "sticoll/auth"
	"sticoll/rest"

	bolt "github.com/coreos/bbolt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	debug bool
)

type device struct {
	cfg        *rest.GRPCCfg
	Stats      gRPCStats
	ifxPointCh chan ifxPoint
	Opts       []grpc.DialOption
}

func init() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	logrus.WithFields(logrus.Fields{
		"nCPU": nuCPU,
	}).Info("Running with number of CPUs")
}

func appCLISetup() *cli.App {
	app := cli.NewApp()
	app.Name = "Godjun - example of interaction with BGP via NETCONF in Junos"
	app.Usage = "Provides a REST API to show bgp neighbor on a Junos router with some helpfult additions. More info https://netopscasts.com/first/"
	app.Email = "egor.krv@gmail.com"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "debug mode currently just prety prints JSON",
			Destination: &debug,
		},
	}
	return app
}

func readCfg(db *bolt.DB) ([]*rest.GRPCCfg, error) {
	gCfgs := make([]*rest.GRPCCfg, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				var g rest.GRPCCfg
				err := json.Unmarshal(v, &g)
				if err != nil {
					return err
				}
				gCfgs = append(gCfgs, &g)
				// fmt.Println(string(k), string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {

		return nil, err
	}
	return gCfgs, nil
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	app := appCLISetup()
	app.Action = func(c *cli.Context) {
		db, err := bolt.Open("my.db", 0600, nil)
		if err != nil {
			logFatal("dbopen", "failure to open db file", err)
		}
		defer func() {
			err = db.Close()
			if err != nil {
				logFatal("dbclose", "failure to close db file", err)
			}
		}()
		// creating new influx structure and initialising
		var ifx influxDB
		err = ifx.NewClientAndPoints()
		if err != nil {
			logFatal(ifxLogTopic, eventNewClientErr, err)
		}
		ifx.dataCh = make(chan ifxPoint)
		go ifx.sendToInflux()
		cfgs, err := readCfg(db)
		if err != nil {
			logErrEvent(cfgErrTopic, cfgReadErrEv, err)
		}
		cfgCh := make(chan *rest.GRPCCfg)
		go func() {
			err = rest.StartHTTPSrv(db, &cfgs, cfgCh)
			if err != nil {
				logErrEvent("http", "failure to start http server", err)
				os.Exit(1)
			}
		}()
		// creating gorutines for each device and passing influx channel
		// many device rutines pass data to a single influx rutine which writes data into the DB
		// works on startup only
		for _, cfg := range cfgs {
			d := device{
				cfg:        cfg,
				ifxPointCh: ifx.dataCh,
			}
			fmt.Println(cfg)
			go d.prepConAndSubscribe()
		}
		// waiting for new devices and connecting to them
		for newCfg := range cfgCh {
			d := device{
				cfg:        newCfg,
				ifxPointCh: ifx.dataCh,
			}
			cfgs = append(cfgs, newCfg)
			go d.prepConAndSubscribe()
			fmt.Println("Adeed new device", newCfg)
		}

	}
	app.Run(os.Args)
}

func (d *device) prepConAndSubscribe() {
	defer d.cfg.RUnlock()
	var conn *grpc.ClientConn
	for {
		fmt.Println("NEw conn")
		err := addDialOptions(d)
		if err != nil {
			logErrEvent(grpcTopic, grpcDialOptsErrEv, err)
			return
		}
		hostname := d.cfg.Host + ":" + strconv.Itoa(d.cfg.Port)
		fmt.Println("Dial")
		conn, err = grpc.Dial(hostname, d.Opts...)
		if err != nil {
			logFatalEvent(grpcTopic, grpcConnErrEv, err)
		}
		defer conn.Close()
		if d.cfg.User != "" && d.cfg.Password != "" {
			fmt.Println("Auth")
			user := d.cfg.User
			pass := d.cfg.Password
			if d.cfg.Meta == false {
				fmt.Println("Meta")
				d.cfg.RLock()
				if d.cfg.Removed {
					logInfoEvent(grpcTopic, "device removed", "exiting gorutine")
					return
				}
				d.cfg.RUnlock()
				fmt.Println("Login")
				dat, err := auth_pb.NewLoginClient(conn).LoginCheck(context.Background(), &auth_pb.LoginRequest{
					UserName: user,
					Password: pass,
					ClientId: d.cfg.CID,
				})
				fmt.Println("After Login")
				if err != nil {
					logErrEvent(grpcTopic, grpcLoginErrEv, err)
					time.Sleep(10 * time.Second)
					continue
				}
				if dat.Result == false {
					logErrEvent(grpcTopic, grpcAuthErrEv, errors.New("login failure"))
				} else {
					break
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	d.subscribe(conn)
}

func addDialOptions(d *device) error {
	if d.cfg.TLS.CA != "" {
		cert, err := tls.LoadX509KeyPair(d.cfg.TLS.ClientCrt, d.cfg.TLS.ClientKey)
		if err != nil {
			logFatalEvent(tlsLogTopic, tlsLogTopic, err)
			return err
		}
		certPool := x509.NewCertPool()
		caInBytes, err := ioutil.ReadFile(d.cfg.TLS.CA)
		if err != nil {
			logFatalEvent(tlsLogTopic, tlsCAReadEv, err)
			return err
		}
		ok := certPool.AppendCertsFromPEM(caInBytes)
		if !ok {
			logFatalEvent(tlsLogTopic, tlsCertAppendEv, err)
			return errors.New("AppendCertsFromPEM")
		}
		transportCreds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   d.cfg.TLS.ServerName,
			RootCAs:      certPool,
		})
		d.Opts = append(d.Opts, grpc.WithTransportCredentials(transportCreds))
	} else {
		d.Opts = append(d.Opts, grpc.WithInsecure())
	}
	d.Opts = append(d.Opts, grpc.WithStatsHandler(&statsHandler{cfg: d}))
	if d.cfg.Compression != "" {
		var dc grpc.Decompressor
		if d.cfg.Compression == "gzip" {
			dc = grpc.NewGZIPDecompressor()
		} else if d.cfg.Compression == "deflate" {
			dc = newDEFLATEDecompressor()
		}
		compressionOpts := grpc.Decompressor(dc)
		d.Opts = append(d.Opts, grpc.WithDecompressor(compressionOpts))
	}
	if d.cfg.WS != 0 {
		d.Opts = append(d.Opts, grpc.WithInitialWindowSize(d.cfg.WS))
	}
	return nil
}
