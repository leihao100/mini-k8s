package listwatch

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
)

// Lister is any object that knows how to perform an initial list.
type Lister interface {
	// List should return a list type object; the Items field will be extracted, and the
	// ResourceVersion field will be used to start the watch in the right place.
	List(options config.ListOptions) (config.ApiObjectList, error)
}

// Watcher is any object that knows how to start a watch on a resource.
type Watcher interface {
	// Watch should begin a watch at the specified version.
	Watch(options config.ListOptions) (watch.Interface, error)
}

// ListerWatcher is any object that knows how to perform an initial list and start a watch on a resource.
type ListerWatcher interface {
	Lister
	Watcher
}

// ListFunc knows how to list resources
type ListFunc func(options config.ListOptions) (config.ApiObjectList, error)

// WatchFunc knows how to watch resources
type WatchFunc func(options config.ListOptions) (watch.Interface, error)

// ListWatch knows how to list and watch a set of apiserver resources.  It satisfies the ListerWatcher interface.
// It is a convenience function for users of NewReflector, etc.
// ListFunc and WatchFunc must not be nil
type ListWatch struct {
	ListFunc  ListFunc
	WatchFunc WatchFunc
	// DisableChunking requests no chunking for this list watcher.
	DisableChunking bool
}

// List a set of apiserver resources
func (lw *ListWatch) List() (config.ApiObjectList, error) {
	// ListWatch is used in Reflector, which already supports pagination.
	// Don't paginate here to avoid duplication.
	var options config.ListOptions
	return lw.ListFunc(options)
}

// Watch a set of apiserver resources
func (lw *ListWatch) Watch() (watch.Interface, error) {
	var options config.ListOptions
	return lw.WatchFunc(options)
}

// NewListWatchFromClient creates a new ListWatch from the specified client, resource, namespace and field selector.
func NewListWatchFromClient(c apiClient.Interface) *ListWatch {
	optionsModifier := func(options *config.ListOptions) {
		// options.FieldSelector = fieldSelector.String()
	}
	return NewFilteredListWatchFromClient(c, optionsModifier)
}

// NewFilteredListWatchFromClient creates a new ListWatch from the specified client, resource, namespace, and option modifier.
// Option modifier is a function takes a ListOptions and modifies the consumed ListOptions. Provide customized modifier function
// to apply modification to ListOptions with a field selector, a label selector, or any other desired options.
func NewFilteredListWatchFromClient(c apiClient.Interface, optionsModifier func(options *config.ListOptions)) *ListWatch {
	listFunc := func(options config.ListOptions) (config.ApiObjectList, error) {
		optionsModifier(&options)
		return c.GetAll()
	}
	watchFunc := func(options config.ListOptions) (watch.Interface, error) {
		options.Watch = true
		optionsModifier(&options)
		return c.WatchAll()
	}
	return &ListWatch{ListFunc: listFunc, WatchFunc: watchFunc}
}
