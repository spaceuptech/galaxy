package runner

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/runner/driver"
	"github.com/spaceuptech/galaxy/utils"
	"github.com/spaceuptech/galaxy/utils/auth"
)

// Runner is the module responsible to manage the runner
type Runner struct {
	// For storing config
	config *Config

	// For handling http related stuff
	router *mux.Router

	// For internal use
	auth     *auth.Module
	driver   driver.Driver
	debounce *utils.Debounce

	// For autoscaler
	db       *badger.DB
	chAppend chan *model.ProxyMessage
}

// New creates a new instance of the runner
func New(c *Config) (*Runner, error) {
	// Add the proxy port to the driver config
	proxyPort, err := strconv.Atoi(c.ProxyPort)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy port (%s) provided", c.ProxyPort)
	}
	c.Driver.ProxyPort = uint32(proxyPort)

	// Initialise all modules
	a, err := auth.New(c.Auth)
	if err != nil {
		return nil, err
	}

	d, err := driver.New(a, c.Driver)
	if err != nil {
		return nil, err
	}

	debounce := utils.NewDebounce()

	opts := badger.DefaultOptions("/tmp/galaxy.db")
	opts.Logger = &logrus.Logger{Out: ioutil.Discard}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	// Periodically run the garbage collector
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	// Return a new runner instance
	return &Runner{
		config: c,
		router: mux.NewRouter(),

		// For internal use
		auth:     a,
		driver:   d,
		debounce: debounce,

		// For autoscaler
		db:       db,
		chAppend: make(chan *model.ProxyMessage, 10),
	}, nil
}

// Start begins the runner
func (runner *Runner) Start() error {
	// Initialise the various routes of the runner
	runner.routes()

	// Start necessary routines for autoscaler
	go runner.routineAdjustScale()
	for i := 0; i < 10; i++ {
		go runner.routineDumpDetails()
	}

	// Start proxy server
	go func() {
		// Create a new router
		router := mux.NewRouter()
		router.PathPrefix("/").HandlerFunc(runner.handleProxy())

		// Start http server
		corsObj := utils.CreateCorsObject()
		logrus.Infof("Starting runner proxy on port %s", runner.config.ProxyPort)
		if err := http.ListenAndServe(":"+runner.config.ProxyPort, corsObj.Handler(router)); err != nil {
			logrus.Fatalln("Proxy server failed:", err)
		}
	}()

	// Start the http server
	corsObj := utils.CreateCorsObject()
	logrus.Infof("Starting runner on port %s", runner.config.Port)
	return http.ListenAndServe(":"+runner.config.Port, corsObj.Handler(runner.router))
}
