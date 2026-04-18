package framework

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/option"
	"github.com/oylshe1314/framework/profile"
	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/util"
)

var (
	ProgramHash = "EMPTY"
	DataHash    = "EMPTY"
	ConfigHash  = "EMPTY"
)

func Start(svr server.Server) {
	os.Exit(start(svr))
}

func start(svr server.Server) int {

	var showVersion bool
	var configFile string
	var flagOption option.FlagOption

	flag.BoolVar(&showVersion, "v", false, "# Show version information and exit.")
	flag.BoolVar(&showVersion, "version", false, "# Show version information and exit.")
	flag.StringVar(&configFile, "c", "config.json", "# Start the server with a config file.")
	flag.Var(flagOption, "D", "# Used to define configuration options, format: option.subOption=value.")
	flag.Parse()

	PrintVersion(os.Stdout)

	if showVersion {
		return 0
	}

	CheckExpired()

	var err error
	var allHash []string

	allHash, _, err = util.HashAll(md5.New(), true, nil, nil, []string{filepath.Dir(os.Args[0])})
	if err != nil {
		fmt.Println("Calculate program hash failed, ", err)
		return 1
	}

	ProgramHash = allHash[0]

	opt, err := option.ReadJson(configFile)
	if err != nil {
		fmt.Println("Read config file failed, ", err)
		return 1
	}

	opt.Merge(option.Option(flagOption))

	var optStr = opt.String()
	fmt.Println(optStr)

	allHash, _, err = util.HashAll(md5.New(), true, []string{optStr}, nil, nil)
	if err != nil {
		fmt.Println("Calculate config hash failed, ", err)
		return 1
	}

	ConfigHash = allHash[0]

	return run(svr, opt)
}

func run(svr server.Server, opt option.Option) int {

	fmt.Println("Server set option")
	var cs, err = Option(opt).setOption(svr)
	if err != nil {
		fmt.Println("Server init failed, ", err)
		return 1
	}

	var stopped atomic.Bool

	stopped.Store(false)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if !stopped.Load() {
			cancel()
		}
	}()

	fmt.Println("Server initialization")
	for _, cc := range cs {
		time.Sleep(time.Millisecond * 100)
		err = cc.Init(ctx)
		if err != nil {
			return 1
		}

		switch tc := cc.(type) {
		case server.Component[server.Server]:
			ctx = ContextWithServer(ctx, tc.Name(), tc.Server())
		case client.Component[client.Client]:
			ctx = ContextWithClient(ctx, tc.Name(), tc.Client())
		}
	}

	fmt.Println("Program-Hash: ", ProgramHash)
	fmt.Println("Data-Hash: ", DataHash)
	fmt.Println("Config-Hash: ", ConfigHash)
	fmt.Println("Profile-Active: ", profile.Active)

	fmt.Println("Server startup")
	for _, cc := range cs {
		time.Sleep(time.Millisecond * 100)
		err = cc.Start()
		if err != nil {
			fmt.Println("Server start failed, ", err)
			return 1
		}
	}

	// Start the close monitor
	go func() {
		var sigChan = make(chan os.Signal)

		runtime.GOMAXPROCS(runtime.NumCPU())
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
		defer signal.Stop(sigChan)

		<-sigChan

		stopped.Store(true)

		// Server shutdown order: reverse of server startup
		for i := len(cs) - 1; i >= 0; i-- {
			time.Sleep(time.Millisecond * 100)
			_ = cs[i].Close()
		}
		cancel()
	}()

	// Waiting servers shutdown
	<-ctx.Done()
	return 0
}
