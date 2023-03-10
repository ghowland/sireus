package extdata

import (
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"time"
)

// StoreQueryResult will store a QueryResult in the cache
func StoreQueryResult(session *data.InteractiveSession, site *data.Site, query data.BotQuery, startTime time.Time, queryResult data.QueryResult) {
	// Create and store the QueryResult pool item
	newCacheItem := data.QueryResultPoolItem{
		QueryServer:     query.QueryServer,
		Query:           query.Query,
		InteractiveUUID: session.UUID,
		TimeRequested:   startTime,
		TimeReceived:    util.GetTimeNow(),
		Result:          queryResult,
		IsValid:         true, //TODO(ghowland): Check instead of force set.  If it's not valid, we need a way to tell them about the problem, and show them for the BotGroup and Bots so they arent confused as to why it's not working.  Can tell them why it's malformed and show them the result so they can troubleshoot it.
		QueryStartTime:  session.QueryStartTime,
		QueryDuration:   session.QueryDuration,
	}

	// Save this result to the cache
	QueryCacheSet(session, site, query, newCacheItem)
}

// GetCachedQueryResult returns a cached query result.  Web App requests should set errorOverInterval=false, which
// is used by the background query system to test missing or expired query results as equivalent.
func GetCachedQueryResult(session *data.InteractiveSession, site *data.Site, query data.BotQuery) (data.QueryResult, error) {
	queryKey := GetQueryKey(session, query)

	result, ok := app.GetQueryResultByQueryKey(site, queryKey)
	if !ok {
		return data.QueryResult{}, errors.New(fmt.Sprintf("Could not find Query Result: Server: %s  Name: %s", query.QueryServer, query.Name))
	}

	if session.IgnoreCacheQueryMismatch && (result.QueryStartTime != session.QueryStartTime || result.QueryDuration != session.QueryDuration) {
		return data.QueryResult{}, errors.New(fmt.Sprintf("Does not match start and duration: Start: %v  Duration: %v", session.QueryStartTime, session.QueryDuration))
	}

	// Test if it is older than the Interval refresh, this
	since := util.GetTimeNow().Sub(result.TimeReceived)

	// If we don't want to return values if they are over the Interval, then mark them
	if since.Seconds() > time.Duration(query.Interval).Seconds() {
		if session.IgnoreCacheOverInterval {
			return data.QueryResult{}, errors.New(fmt.Sprintf("Query Result found, but over interval: Server: %s  Name: %s", query.QueryServer, query.Name))
		} else {
			//TODO(ghowland): For specific BotGroup queries, we have an additional check for Staleness, it doesnt use Interval above...  Deal with that later.
		}
	}

	// Returning the cached result
	return result.Result, nil
}

func QueryCacheSet(session *data.InteractiveSession, site *data.Site, query data.BotQuery, newCacheItem data.QueryResultPoolItem) {
	queryKey := GetQueryKey(session, query)

	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryPoolSyncLock.Lock()
	defer site.QueryResultCache.QueryPoolSyncLock.Unlock()

	site.QueryResultCache.PoolItems[queryKey] = newCacheItem
}

// QueryLockClear will clear the Query Lock, so we can make this Query again after the Interval
func QueryLockClear(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	delete(site.QueryResultCache.QueryLocks, queryKey)
}

// QueryLockSet will set the Query Lock for a QueryKey, so we won't request this Query again until it finishes or the
// AppConfig.QueryLockTimeout expires
func QueryLockSet(site *data.Site, queryKey string) {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	site.QueryResultCache.QueryLocks[queryKey] = util.GetTimeNow()
}

// GetQueryKey returns "(QueryServer).(Query)", so it can be shared by any BotGroup
func GetQueryKey(session *data.InteractiveSession, query data.BotQuery) string {
	// Key on the Query itself, so if different BotGroups share the same query from the same QueryServer, it's shared
	output := fmt.Sprintf("%d.%s.%s", session.UUID, query.QueryServer, query.Query)
	return output
}

// IsQueryLocked returned whether this Query currently being requested.  Don't want to request more than once at a time
func IsQueryLocked(session *data.InteractiveSession, site *data.Site, query data.BotQuery) bool {
	// Block until we can lock, for goroutine safety
	site.QueryResultCache.QueryLocksSyncLock.Lock()
	defer site.QueryResultCache.QueryLocksSyncLock.Unlock()

	queryKey := GetQueryKey(session, query)

	queryLockTime, ok := site.QueryResultCache.QueryLocks[queryKey]
	if !ok {
		return false
	}

	since := util.GetTimeNow().Sub(queryLockTime)
	if since.Seconds() < time.Duration(data.SireusData.AppConfig.QueryLockTimeout).Seconds() {
		return true
	}

	return false
}
