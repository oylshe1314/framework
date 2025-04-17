package server

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/profile"
	"github.com/oylshe1314/framework/util"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	ProgramHash = "EMPTY"
	DataHash    = "EMPTY"
	ConfigHash  = "EMPTY"
)

func Start(svc Server) {
	os.Exit(start(svc))
}

func start(svr Server) int {

	var showVersion bool
	var configFile string
	var flagOptions FlagOptions

	flag.BoolVar(&showVersion, "v", false, "# Show version information and exit.")
	flag.BoolVar(&showVersion, "version", false, "# Show version information and exit.")
	flag.StringVar(&configFile, "conf", "config.json", "# Start the server with a config file.")
	flag.Var(&flagOptions, "D", "# Used to define configuration options, format: option.subOption=value.")
	flag.Parse()

	if showVersion {
		printVersion(os.Stdout)
		return 0
	}

	if util.Unix() >= expiration {
		fmt.Println("The server was expired")
		return 0
	}

	logVersion(log.DefaultLogger)

	log.DefaultLogger.Info("Config: ", configFile)

	hashAll, _, err := util.HashAll(md5.New(), true, nil, []string{filepath.Dir(os.Args[0])}, nil)
	if err != nil {
		log.DefaultLogger.Error("Calculate program hash failed, ", err)
		return 1
	}

	ProgramHash = hashAll[0]

	options, err := ReadOptions(configFile)
	if err != nil {
		log.DefaultLogger.Error("Read config file failed, ", err)
		return 1
	}

	additionalOptions, err := flagOptions.Parse()
	if err != nil {
		log.DefaultLogger.Error("Parse flag options failed, ", err)
		return 1
	}

	hashAll, _, err = util.HashAll(md5.New(), true, flagOptions, nil, []string{configFile})
	if err != nil {
		log.DefaultLogger.Error("Calculate program hash failed, ", err)
		return 1
	}

	ConfigHash = hashAll[0]

	options.Merge(additionalOptions)

	logOptions(log.DefaultLogger, options)

	log.DefaultLogger.Info("Server init")
	err = options.Init(svr)
	if err != nil {
		log.DefaultLogger.Error("Server init failed, ", err)
		return 1
	}

	return run(svr)
}

func run(svr Server) int {

	var code = 0
	var stopped = false
	var sigChan = make(chan os.Signal)

	runtime.GOMAXPROCS(runtime.NumCPU())
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL)

	var logger = svr.Logger()

	logger.Info("Server start")
	logger.Info("Program-Hash: ", ProgramHash)
	logger.Info("Data-Hash: ", DataHash)
	logger.Info("Config-Hash: ", ConfigHash)
	logger.Info("Profile-Active: ", profile.Active)
	go func(pCode *int) {
		err := svr.Serve()
		if err == nil {
			logger.Info("Server stopped")
		} else {
			if !stopped {
				*pCode = 1
			}
			logger.Error("Server start failed, ", err)
		}
		if !stopped {
			sigChan <- syscall.SIGQUIT
		}
	}(&code)

	<-sigChan
	stopped = true
	_ = svr.Close()
	close(sigChan)
	return code
}
