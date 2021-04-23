package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"

	"github.com/kyma-project/kyma/components/application-operator/pkg/apis/applicationconnector/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

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
	if err := initCtrlLog(log, level); err != nil {
		log.WithContext().Error("While initializing ctrl logger: %s", err.Error())
		os.Exit(5)
	}
	mgr, err := setupMgr(log, options)
	if err != nil {
		log.WithContext().Error("While initializing manager: %s", err.Error())
		os.Exit(5)
	}

	log.WithContext().With("options", options).Info("Starting Validation Proxy.")

	idCache := cache.New(
		time.Duration(options.cacheExpirationMinutes)*time.Minute,
		time.Duration(options.cacheCleanupMinutes)*time.Minute,
	)

	proxyHandler := validationproxy.NewProxyHandler(
		options.appNamePlaceholder,
		options.group,
		options.tenant,
		options.eventServicePathPrefixV1,
		options.eventServicePathPrefixV2,
		options.eventServiceHost,
		options.eventMeshPathPrefix,
		options.eventMeshHost,
		options.eventMeshDestinationPath,
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
	//TODO: stop handlers
	go func() {
		log.WithContext().With("server", "proxy").With("port", options.proxyPort).Error(proxyServer.ListenAndServe())
	}()

	go func() {
		log.WithContext().With("server", "external").With("port", options.externalAPIPort).Error(externalServer.ListenAndServe())
	}()

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.WithContext().Error("While starting manager: %s", err.Error())
		os.Exit(1)
	}
}

func initCtrlLog(log *logger.Logger, level logger.Level) error {
	zaprLogger := zapr.NewLogger(log.WithContext().Desugar())
	lvl, err := level.ToZapLevel()
	if err != nil {
		return errors.Wrap(err, "while getting zap log level")
	}
	zaprLogger.V((int)(lvl))
	ctrl.SetLogger(zaprLogger)
	return nil
}

func setupMgr(log *logger.Logger, opts *options) (manager.Manager, error) {
	scheme, err := setupScheme()
	if err != nil {
		return nil, err
	}
	//TODO: config resync period => options
	syncPeriod := time.Second * 2 * time.Second
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		SyncPeriod:         &syncPeriod,
		MetricsBindAddress: "0",
	})
	if err != nil {
		return nil, err
	}
	if err = (&controller.CacheReconciler{
		Client: mgr.GetClient(),
		Log:    log,
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		return nil, err
	}
	return mgr, nil
}

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}
