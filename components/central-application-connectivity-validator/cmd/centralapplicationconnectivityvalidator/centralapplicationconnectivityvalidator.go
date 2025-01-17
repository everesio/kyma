package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kyma-project/kyma/common/logging/logger"
	"github.com/kyma-project/kyma/common/logging/tracing"
	"github.com/kyma-project/kyma/components/central-application-connectivity-validator/internal/controller"
	"github.com/kyma-project/kyma/components/central-application-connectivity-validator/internal/externalapi"
	"github.com/kyma-project/kyma/components/central-application-connectivity-validator/internal/validationproxy"
	"github.com/patrickmn/go-cache"
)

func main() {
	options, err := parseOptions()
	if err != nil {
		if logErr := logger.LogFatalError("Failed to parse options: %s", err.Error()); logErr != nil {
			fmt.Printf("Failed to initializie default fatal error logger: %s,Failed to parse options: %s", logErr, err)
		}
		os.Exit(1)
	}
	if err = options.validate(); err != nil {
		if logErr := logger.LogFatalError("Failed to validate options: %s", err.Error()); logErr != nil {
			fmt.Printf("Failed to initializie default fatal error logger: %s,Failed to validate options: %s", logErr, err)
		}
		os.Exit(1)
	}
	level, err := logger.MapLevel(options.LogLevel)
	if err != nil {
		if logErr := logger.LogFatalError("Failed to map log level from options: %s", err.Error()); logErr != nil {
			fmt.Printf("Failed to initializie default fatal error logger: %s, Failed to map log level from options: %s", logErr, err)
		}

		os.Exit(2)
	}
	format, err := logger.MapFormat(options.LogFormat)
	if err != nil {
		if logErr := logger.LogFatalError("Failed to map log format from options: %s", err.Error()); logErr != nil {
			fmt.Printf("Failed to initializie default fatal error logger: %s, Failed to map log format from options: %s", logErr, err)
		}
		os.Exit(3)
	}
	log, err := logger.New(format, level)
	if err != nil {
		if logErr := logger.LogFatalError("Failed to initialize logger: %s", err.Error()); logErr != nil {
			fmt.Printf("Failed to initializie default fatal error logger: %s, Failed to initialize logger: %s", logErr, err)
		}
		os.Exit(4)
	}
	if err := logger.InitKlog(log, level); err != nil {
		log.WithContext().Error("While initializing klog logger: %s", err.Error())
		os.Exit(5)
	}

	log.WithContext().With("options", options).Info("Starting Validation Proxy.")

	idCache := cache.New(
		time.Duration(options.cacheExpirationSeconds)*time.Second,
		time.Duration(options.cacheCleanupIntervalSeconds)*time.Second,
	)
	idCache.OnEvicted(func(key string, i interface{}) {
		log.WithContext().
			With("controller", "cache_janitor").
			With("name", key).
			Warnf("Deleted the application from the cache on cache eviction.")
	})

	proxyHandler := validationproxy.NewProxyHandler(
		options.appNamePlaceholder,
		options.eventingPathPrefixV1,
		options.eventingPathPrefixV2,
		options.eventingPathPrefixEvents,
		options.eventingPublisherHost,
		options.eventingDestinationPath,
		options.appRegistryPathPrefix,
		options.appRegistryHost,
		idCache,
		log)

	tracingMiddleware := tracing.NewTracingMiddleware(proxyHandler.ProxyAppConnectorRequests)

	proxyServer := http.Server{
		Handler: validationproxy.NewHandler(tracingMiddleware),
		Addr:    fmt.Sprintf(":%d", options.proxyPort),
	}

	externalServer := http.Server{
		Handler: externalapi.NewHandler(),
		Addr:    fmt.Sprintf(":%d", options.externalAPIPort),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		controller.Start(log, options.kubeConfig, options.apiServerURL, options.syncPeriod, idCache)
	}()

	go func() {
		log.WithContext().With("server", "proxy").With("port", options.proxyPort).Fatal(proxyServer.ListenAndServe())
	}()

	go func() {
		log.WithContext().With("server", "external").With("port", options.externalAPIPort).Fatal(externalServer.ListenAndServe())
	}()

	wg.Wait()
}
