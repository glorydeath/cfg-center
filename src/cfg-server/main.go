package main

import (
	"flag"
	"fmt"
	"github.com/4paradigm/cfg-center/src/cfg-server/cfgLoader"
	"github.com/4paradigm/cfg-center/src/cfg-server/httpApi"
	"github.com/4paradigm/cfg-center/src/cfg-server/scheduler"
	log "github.com/auxten/logrus"
	"os"
)

var (
	versionStr = "unknown"
)

type Config struct {
	IsDebug     bool
	ListenPort  int
	CfgDataPath string
}

func init_args(cfg *Config) error {
	flag.BoolVar(&cfg.IsDebug, "d", false, "Debug Switch")
	flag.IntVar(&cfg.ListenPort, "port", 2120, "The TCP Port to Listen")
	flag.StringVar(&cfg.CfgDataPath, "data", "../../conf_data", "Config Data Dir Path")

	version := flag.Bool("v", false, "Show Version")

	flag.Parse()
	if *version {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if cfg.IsDebug {
		log.SetLevel(log.DebugLevel)
	}
	//log.SetLevel(log.DebugLevel)
	log.Debug("initing")

	return nil
}

func main() {
	var cfg Config
	err := init_args(&cfg)

	// Debug log
	//log.SetOutput(os.Stderr)

	if err != nil {
		log.Fatal(err.Error())
	}

	cfg_mgr := cfgLoader.New(cfg.CfgDataPath)
	cfg_mgr.LoadCfg()

	scheduler := scheduler.New(cfg.CfgDataPath)
	go scheduler.Run(cfg_mgr)

	log.Error(httpApi.HTTPServerStart(cfg.ListenPort, cfg_mgr))
}
