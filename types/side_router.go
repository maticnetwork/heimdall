package types

import (
	"fmt"
	"regexp"
)

var (
	_ SideRouter = (*router)(nil)

	isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
)

// SideHandlers handler for side-tx processing
type SideHandlers struct {
	SideTxHandler SideTxHandler
	PostTxHandler PostTxHandler
}

// SideRouter implements router.
type SideRouter interface {
	AddRoute(r string, h *SideHandlers) (rtr SideRouter)
	HasRoute(r string) bool
	GetRoute(path string) (h *SideHandlers)
	Seal()
}

type router struct {
	routes map[string]*SideHandlers
	sealed bool
}

// NewSideRouter new router
func NewSideRouter() SideRouter {
	return &router{
		routes: make(map[string]*SideHandlers),
	}
}

// Seal seals the router which prohibits any subsequent route handlers to be
// added. Seal will panic if called more than once.
func (rtr *router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}

	rtr.sealed = true
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
func (rtr *router) AddRoute(path string, h *SideHandlers) SideRouter {
	if rtr.sealed {
		panic("router sealed; cannot add route handler")
	}

	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}

	if rtr.HasRoute(path) {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h

	return rtr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (rtr *router) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (rtr *router) GetRoute(path string) *SideHandlers {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
