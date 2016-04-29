package dependency

import (
	"errors"
	"sort"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

// ErrStopped is a special error that is returned when a dependency is
// prematurely stopped, usually due to a configuration reload or a process
// interrupt.
var ErrStopped = errors.New("dependency stopped")

// Dependency is an interface for a dependency that Consul Template is capable
// of watching.
type Dependency interface {
	Fetch(*ClientSet, *QueryOptions) (interface{}, *ResponseMetadata, error)
	CanShare() bool
	HashCode() string
	Display() string
	Stop()
}

// ServiceTags is a slice of tags assigned to a Service
type ServiceTags []string

// Contains returns true if the tags exists in the ServiceTags slice.
func (t ServiceTags) Contains(s string) bool {
	for _, v := range t {
		if v == s {
			return true
		}
	}
	return false
}

// QueryOptions is a list of options to send with the query. These options are
// client-agnostic, and the dependency determines which, if any, of the options
// to use.
type QueryOptions struct {
	AllowStale bool
	WaitIndex  uint64
	WaitTime   time.Duration
}

// Converts the query options to Consul API ready query options.
func (r *QueryOptions) consulQueryOptions() *consulapi.QueryOptions {
	return &consulapi.QueryOptions{
		AllowStale: r.AllowStale,
		WaitIndex:  r.WaitIndex,
		WaitTime:   r.WaitTime,
	}
}

// ResponseMetadata is a struct that contains metadata about the response. This
// is returned from a Fetch function call.
type ResponseMetadata struct {
	LastIndex   uint64
	LastContact time.Duration
}

// deepCopyAndSortTags deep copies the tags in the given string slice and then
// sorts and returns the copied result.
func deepCopyAndSortTags(tags []string) []string {
	newTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		newTags = append(newTags, tag)
	}
	sort.Strings(newTags)
	return newTags
}

// respWithMetadata is a short wrapper to return the given interface with fake
// response metadata for non-Consul dependencies.
func respWithMetadata(i interface{}) (interface{}, *ResponseMetadata, error) {
	return i, &ResponseMetadata{
		LastContact: 0,
		LastIndex:   uint64(time.Now().Unix()),
	}, nil
}
