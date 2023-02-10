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

	for poolIndex, result := range site.QueryResultCache.PoolItems {
		// Is this result a match?
		if result.QueryServer == query.QueryServer && result.Query == query.Query {
			// Mark found and set it back over the current location, we're done
			QueryCacheSet(site, poolIndex, newCacheItem)
			return
		}
	}

	// If we didn't find this result and return already, append it
	QueryCacheAppend(site, newCacheItem)
}

// Returns a cached query result.  Web App requests should set errorOverIntervall=false, which is used by the
// background query system to test missing or expired query results as equivolent.
func GetCachedQueryResult(site *data.Site, query data.BotQuery, errorOverInterval bool) (data.QueryResult, error) {

	for _, result := range site.QueryResultCache.PoolItems {
		// Is this result a match?
		if result.QueryServer == query.QueryServer && result.Query == query.Query {
			since := time.Now().Sub(result.TimeReceived)

			// If we don't want to return values if they are over the interval, then mark them
			if since.Seconds() > time.Duration(query.Interval).Seconds() {
				if errorOverInterval {
					return data.QueryResult{}, errors.New(fmt.Sprintf("Query Result found, but over interval: Server: %s  Name: %s", query.QueryServer, query.Name))
				} else {
					// This is an expired result.  Any Bots that use this are now Stale and can't be IsAvailable, so can't execute Actions
					result.Result.IsExpired = true
				}
			}

			// Returning the cached result.  May have set IsExpired above.
			return result.Result, nil
		}
	}

	return data.QueryResult{}, errors.New(fmt.Sprintf("Could not find Query Result: Server: %s  Name: %s", query.QueryServer, query.Name))
}

func QueryCacheSet(site *data.Site, poolIndex int, newCacheItem data.QueryResultPoolItem) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	// Mark found and set it back over the current location, we're done
	site.QueryResultCache.PoolItems[poolIndex] = newCacheItem
}

func QueryCacheAppend(site *data.Site, newCacheItem data.QueryResultPoolItem) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	// If we didn't find this result and return already, append it
	site.QueryResultCache.PoolItems = append(site.QueryResultCache.PoolItems, newCacheItem)
}

// Clear the Query Lock, so we can make this Query again after the Interval
func QueryLockClear(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	site.QueryResultCache.QueryLocks[queryKey] = time.UnixMilli(0)
}

// Set the Query Lock, so we won't request this Query again until it finishes or the AppConfig.QueryLockTimeout expires
func QueryLockSet(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	site.QueryResultCache.QueryLocks[queryKey] = time.Now()
}
