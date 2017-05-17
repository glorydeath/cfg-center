package scheduler

import (
	"github.com/4paradigm/cfg-center/src/cfg-server/cfgLoader"
	log "github.com/auxten/logrus"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/signal"
	"syscall"
)

type Scheduler struct {
	reload_ch   chan uint8
	sig_ch      chan os.Signal
	cfgDataPath string
}

func New(cfgDataPath string) *Scheduler {
	return &Scheduler{
		reload_ch:   make(chan uint8, 1),
		sig_ch:      make(chan os.Signal, 1),
		cfgDataPath: cfgDataPath,
	}
}

func (s Scheduler) setDirWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = watcher.Add(s.cfgDataPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debug("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename ||
					event.Op&fsnotify.Remove == fsnotify.Remove {
					s.reload_ch <- uint8(0)
					log.Info("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Error("error:", err)
			}
		}
	}()
}

func (s Scheduler) setSigHandler() {
	// Reload signal
	signal.Notify(s.sig_ch, syscall.SIGUSR1)
}

func (s Scheduler) Run(cfgm *cfgLoader.CfgManager) {
	//todo timer
	s.setSigHandler()
	//todo inotify
	s.setDirWatcher()
	//todo even git
	for {
		select {
		case <-s.reload_ch:
			go cfgm.LoadCfg()
		case <-s.sig_ch:
			go cfgm.LoadCfg()
		}
	}
}
