package framework

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"framework/errors"
	"framework/option"
	"framework/util"
	"framework/util/jsonx"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	ProgramHash = "EMPTY"
	ConfigHash  = "EMPTY"
	DataHash    = "EMPTY"
)

func Start(svr Server) {
	os.Exit(start(svr))
}

func start(svr Server) int {

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

	return run(svr, serverOption(opt))
}

type serverComponent struct {
	name   string
	server Server
}

func run(svr Server, opt serverOption) int {

	fmt.Println("Server set option")
	var components, err = opt.setOption(svr)
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
	for _, component := range components {
		time.Sleep(time.Millisecond * 100)
		err = component.server.Init(ctx)
		if err != nil {
			return 1
		}

		ctx = ContextWithServer(ctx, component.name, component.server)
	}

	fmt.Println("Program-Hash: ", ProgramHash)
	fmt.Println("Config-Hash: ", ConfigHash)
	fmt.Println("Data-Hash: ", DataHash)
	fmt.Println("Profile-Active: ", Active)

	fmt.Println("Server startup")
	for _, component := range components {
		time.Sleep(time.Millisecond * 100)
		err = component.server.Start()
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
		for i := len(components) - 1; i >= 0; i-- {
			time.Sleep(time.Millisecond * 100)
			_ = components[i].server.Close()
		}
		cancel()
	}()

	// Waiting servers shutdown
	<-ctx.Done()
	return 0
}

type serverOption option.Option

func (opt serverOption) setOption(svr Server) ([]*serverComponent, error) {
	var st = reflect.TypeOf(svr)
	if st.Kind() != reflect.Pointer || st.Elem().Kind() != reflect.Struct {
		return nil, errors.New("server should be a struct pointer")
	}

	return opt.reflectSetOption(st.Elem().Name(), reflect.ValueOf(svr), st)
}

func (opt serverOption) reflectSetOption(name string, sv reflect.Value, st reflect.Type) ([]*serverComponent, error) {
	var err error
	var components []*serverComponent

	var sit = reflect.TypeOf((*Server)(nil)).Elem()

	st = st.Elem()
	sv = sv.Elem()

	for i := range st.NumField() {
		var sf = st.Field(i)
		if !sf.IsExported() {
			continue
		}

		var ft = sf.Type
		if ft.Kind() != reflect.Pointer {
			ft = reflect.PointerTo(ft)
		}

		if ft.Implements(sit) {
			var tags = strings.Split(sf.Tag.Get("server"), ",")

			var svrName = tags[0]
			if svrName == "" {
				svrName = util.LowerCamelCase(sf.Name)
			}

			var val = option.Option(opt).Get(svrName)
			if val == nil {
				continue
			}

			tmp, ok := val.(map[string]any)
			if !ok {
				return nil, errors.Errorf("can not set option '%s' of %s to server '%s' ", name, reflect.TypeOf(tmp).String(), sf.Name)
			}

			var fv = sv.Field(i)
			if sf.Type.Kind() == reflect.Pointer {
				if fv.IsNil() {
					return nil, errors.Errorf("the pointer of the server '%s' is nil", sf.Name)
				}
			} else {
				fv = fv.Addr()
			}

			var subComponents []*serverComponent
			subComponents, err = (serverOption(tmp)).reflectSetOption(svrName, fv, ft)
			if err != nil {
				return nil, err
			}

			components = append(components, subComponents...)
		}
	}

	for i := range st.NumMethod() {
		var sm = st.Method(i)
		if !sm.IsExported() {
			continue
		}

		if sm.Name == "SetOption" {
			var mt = sm.Type
			if mt.NumIn() != 1 {
				return nil, errors.Errorf("invalid parameter of the server '%s' option set funcation", st.Name())
			}

			var at = mt.In(0)
			if at.Kind() != reflect.Struct || at.Elem().Kind() != reflect.Struct {
				return nil, errors.Errorf("the parameter of the server '%s' option set funcation should be a struct or struct pointer", st.Name())
			}

			if at.Kind() == reflect.Pointer {
				var av = reflect.New(at.Elem())
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}

				sv.Method(i).Call([]reflect.Value{av})
			} else {
				var av = reflect.New(at)
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}
				sv.Method(i).Call([]reflect.Value{av.Elem()})
			}
		}
	}

	for st.Kind() == reflect.Pointer {
		st = st.Elem()
	}

	components = append(components, &serverComponent{name: name, server: sv.Interface().(Server)})
	return components, nil
}
