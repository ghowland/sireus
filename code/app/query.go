package app

import "github.com/ghowland/sireus/code/data"

func GetQueryResultByQueryKey(site *data.Site, queryKey string) (data.QueryResultPoolItem, bool) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryPoolSyncLock.Lock()
	defer site.QueryResultCache.QueryPoolSyncLock.Unlock()

	result, ok := site.QueryResultCache.PoolItems[queryKey]
	return result, ok
}
