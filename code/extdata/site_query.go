package extdata

import (
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"time"
)

// Store this QueryResult in the cache
func StoreQueryResult(interactiveUUID int64, site *data.Site, query data.BotQuery, startTime time.Time, queryResult data.QueryResult) {
	// Create and store the QueryResult pool item
	newCacheItem := data.QueryResultPoolItem{
		QueryServer:     query.QueryServer,
		Query:           query.Query,
		InteractiveUUID: interactiveUUID,
		TimeRequested:   startTime,
		TimeReceived:    time.Now(),
		Result:          queryResult,
		IsValid:         true, //TODO(ghowland): Check instead of force set.  If it's not valid, we need a way to tell them about the problem, and show them for the BotGroup and Bots so they arent confused as to why it's not working.  Can tell them why it's malformed and show them the result so they can troubleshoot it.
	}

	// Save this result to the cache
	QueryCacheSet(site, query, newCacheItem)
}

// Returns a cached query result.  Web App requests should set errorOverIntervall=false, which is used by the
// background query system to test missing or expired query results as equivolent.
func GetCachedQueryResult(site *data.Site, query data.BotQuery, errorOverInterval bool) (data.QueryResult, error) {
	queryKey := GetQueryKey(query)

	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryPoolSyncLock.Lock()
	defer site.QueryResultCache.QueryPoolSyncLock.Unlock()

	result, ok := site.QueryResultCache.PoolItems[queryKey]
	if !ok {
		return data.QueryResult{}, errors.New(fmt.Sprintf("Could not find Query Result: Server: %s  Name: %s", query.QueryServer, query.Name))
	}

	// Test if it is older than the Interval refresh, this
	since := time.Now().Sub(result.TimeReceived)

	// If we don't want to return values if they are over the Interval, then mark them
	if since.Seconds() > time.Duration(query.Interval).Seconds() {
		if errorOverInterval {
			return data.QueryResult{}, errors.New(fmt.Sprintf("Query Result found, but over interval: Server: %s  Name: %s", query.QueryServer, query.Name))
		} else {
			//TODO(ghowland): For specific BotGroup queries, we have an additional check for Staleness, it doesnt use Interval above...  Deal with that later.
		}
	}

	// Returning the cached result
	return result.Result, nil
}

func QueryCacheSet(site *data.Site, query data.BotQuery, newCacheItem data.QueryResultPoolItem) {
	queryKey := GetQueryKey(query)

	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryPoolSyncLock.Lock()
	defer site.QueryResultCache.QueryPoolSyncLock.Unlock()

	site.QueryResultCache.PoolItems[queryKey] = newCacheItem
}

// Clear the Query Lock, so we can make this Query again after the Interval
func QueryLockClear(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	delete(site.QueryResultCache.QueryLocks, queryKey)
}

// Set the Query Lock, so we won't request this Query again until it finishes or the AppConfig.QueryLockTimeout expires
func QueryLockSet(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	site.QueryResultCache.QueryLocks[queryKey] = time.Now()
}

// QueryKey is "(QueryServer).(Query)", so it can be shared by any BotGroup
func GetQueryKey(query data.BotQuery) string {
	// Key on the Query itself, so if different BotGroups share the same query from the same QueryServer, it's shared
	output := fmt.Sprintf("%s.%s", query.QueryServer, query.Query)
	return output
}

// Is this Query currently being requested?  We dont want to request more than once at a time
func IsQueryLocked(site *data.Site, botGroup data.BotGroup, query data.BotQuery) bool {
	queryKey := GetQueryKey(query)

	queryLockTime, ok := site.QueryResultCache.QueryLocks[queryKey]
	if !ok {
		return false
	}

	since := time.Now().Sub(queryLockTime)
	if since.Seconds() < time.Duration(data.SireusData.AppConfig.QueryLockTimeout).Seconds() {
		return true
	}

	return false
}
