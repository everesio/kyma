package controller

import (
	"context"
	"time"

	"github.com/kyma-project/kyma/common/logging/logger"
	gocache "github.com/patrickmn/go-cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type cacheWithLoader struct {
	cacheSync CacheSync
	appCache  *gocache.Cache
	log       *logger.Logger
}

func NewCacheWithLoader(log *logger.Logger, client client.Client, appCache *gocache.Cache) *cacheWithLoader {
	return &cacheWithLoader{
		cacheSync: NewCacheSync(log, client, appCache, "cache_loader"),
		appCache:  appCache,
		log:       log,
	}
}

func (c *cacheWithLoader) Get(parent context.Context, key string) ([]string, bool) {
	clientIDs, found := c.appCache.Get(key)
	if !found {
		ctx, cancel := context.WithTimeout(parent, 2*time.Second)
		defer cancel()
		//TODO: add negative caching & mutex loader for a key with wait
		c.log.WithContext().With("name", key).Warnf("Sync not found Application.")
		err := c.cacheSync.Sync(ctx, key)
		if err != nil {
			return []string{}, false
		}
		clientIDs, found = c.appCache.Get(key)
		if !found {
			return []string{}, false
		}
		c.log.WithContext().
			With("name", key).
			Warnf("Application not sync but loaded to the cache.")
	}
	return clientIDs.([]string), found
}
