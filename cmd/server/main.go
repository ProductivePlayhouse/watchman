// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/moov-io/base/admin"
	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/base/http/bind"
	"github.com/moov-io/base/log"
	"github.com/moov-io/watchman"
	"github.com/moov-io/watchman/internal/database"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	// "github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	// "github.com/awsdocs/aws-doc-sdk-examples/gov2/dynamodb/actions"
	// "github.com/awsdocs/aws-doc-sdk-examples/gov2/dynamodb/scenarios"	
)

var (
	httpAddr  = flag.String("http.addr", bind.HTTP("ofac"), "HTTP listen address")
	adminAddr = flag.String("admin.addr", bind.Admin("ofac"), "Admin HTTP listen address")

	flagBasePath  = flag.String("base-path", "/", "Base path to serve HTTP routes and webui from")
	flagLogFormat = flag.String("log.format", "", "Format for log lines (Options: json, plain")
	flagMaxProcs  = flag.Int("max-procs", runtime.NumCPU(), "Maximum number of CPUs used for search and endpoints")
	flagWorkers   = flag.Int("workers", 1024, "Maximum number of goroutines used for search")

	dataRefreshInterval = 12 * time.Hour
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*flagMaxProcs)

	var logger log.Logger
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		*flagLogFormat = v
	}
	if strings.ToLower(*flagLogFormat) == "json" {
		logger = log.NewJSONLogger()
	} else {
		logger = log.NewDefaultLogger()
	}

	logger.Logf(strings.Repeat("-", 88))
	logger.Logf("Watchman server attempting to connect to DynamoDB table...")

	scenarioMap := map[string]func(sdkConfig aws.Config){
		"movieTable": runMovieScenario,
	}
	choices := make([]string, len(scenarioMap))
	choiceIndex := 0
	for choice := range scenarioMap {
		choices[choiceIndex] = choice
		choiceIndex++
	}
	scenario := flag.String(
		"scenario", "",
		fmt.Sprintf("The scenario to run. Must be one of %v.", choices))
	flag.Parse()

	if runScenario, ok := scenarioMap[*scenario]; !ok {
		fmt.Printf("'%v' is not a valid scenario.\n", *scenario)
		flag.Usage()
	} else {
		sdkConfig, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			logger.Logf("unable to load SDK config, %v", err)
		}
		runScenario(sdkConfig)
	}
}

func addPingRoute(r *mux.Router) {
	r.Methods("GET").Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		moovhttp.SetAccessControlAllowHeaders(w, r.Header.Get("Origin"))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	})
}

// getDataRefreshInterval returns a time.Duration for how often OFAC should refresh data
//
// env is the value from an environmental variable
func getDataRefreshInterval(logger log.Logger, env string) time.Duration {
	if env != "" {
		if strings.EqualFold(env, "off") {
			return 0 * time.Second
		}
		if dur, _ := time.ParseDuration(env); dur > 0 {
			logger.Logf("Setting data refresh interval to %v", dur)
			return dur
		}
	}
	logger.Logf("Setting data refresh interval to %v (default)", dataRefreshInterval)
	return dataRefreshInterval
}

func setupWebui(logger log.Logger, r *mux.Router, basePath string) {
	dir := os.Getenv("WEB_ROOT")
	if dir == "" {
		dir = filepath.Join("webui", "build")
	}
	if _, err := os.Stat(dir); err != nil {
		logger.Logf("problem with webui=%s: %v", dir, err)
		os.Exit(1)
	}
	r.PathPrefix("/").Handler(http.StripPrefix(basePath, http.FileServer(http.Dir(dir))))
}

func handleDownloadStats(updates chan *DownloadStats, handle func(stats *DownloadStats)) {
	for {
		stats := <-updates
		if stats != nil {
			handle(stats)
		}
	}
}

func runMovieScenario(sdkConfig aws.Config) {
	scenarios.RunMovieScenario(
		sdkConfig,
		demotools.NewQuestioner(),
		"doc-example-movie-table",
		actions.MovieSampler{URL: "https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/samples/moviedata.zip"},
	)
}