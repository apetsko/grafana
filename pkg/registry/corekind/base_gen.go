// THIS FILE IS GENERATED. EDITING IS FUTILE.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     BaseCoreRegistryJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package corekind

import (
	"fmt"

	"github.com/grafana/grafana/pkg/kinds/dashboard"
	"github.com/grafana/grafana/pkg/kinds/playlist"
	"github.com/grafana/grafana/pkg/kinds/stars"
	"github.com/grafana/grafana/pkg/kinds/svg"
	"github.com/grafana/grafana/pkg/kindsys"
	"github.com/grafana/thema"
)

// Base is a registry of kindsys.Interface. It provides two modes for accessing
// kinds: individually via literal named methods, or as a slice returned from
// an All*() method.
//
// Prefer the individual named methods for use cases where the particular kind(s) that
// are needed are known to the caller. For example, a dashboard linter can know that it
// specifically wants the dashboard kind.
//
// Prefer All*() methods when performing operations generically across all kinds.
// For example, a validation HTTP middleware for any kind-schematized object type.
type Base struct {
	all                   []kindsys.Interface
	numRaw, numStructured int
	dashboard             *dashboard.Kind
	playlist              *playlist.Kind
	svg                   *svg.Kind
	stars                 *stars.Kind
}

// type guards
var (
	_ kindsys.Structured = &dashboard.Kind{}
	_ kindsys.Structured = &playlist.Kind{}
	_ kindsys.Raw        = &svg.Kind{}
	_ kindsys.Structured = &stars.Kind{}
)

// Dashboard returns the [kindsys.Interface] implementation for the dashboard kind.
func (b *Base) Dashboard() *dashboard.Kind {
	return b.dashboard
}

// Playlist returns the [kindsys.Interface] implementation for the playlist kind.
func (b *Base) Playlist() *playlist.Kind {
	return b.playlist
}

// SVG returns the [kindsys.Interface] implementation for the svg kind.
func (b *Base) SVG() *svg.Kind {
	return b.svg
}

// Stars returns the [kindsys.Interface] implementation for the stars kind.
func (b *Base) Stars() *stars.Kind {
	return b.stars
}

func doNewBase(rt *thema.Runtime) *Base {
	var err error
	reg := &Base{
		numRaw:        1,
		numStructured: 3,
	}

	reg.dashboard, err = dashboard.NewKind(rt)
	if err != nil {
		panic(fmt.Sprintf("error while initializing the dashboard Kind: %s", err))
	}
	reg.all = append(reg.all, reg.dashboard)

	reg.playlist, err = playlist.NewKind(rt)
	if err != nil {
		panic(fmt.Sprintf("error while initializing the playlist Kind: %s", err))
	}
	reg.all = append(reg.all, reg.playlist)

	reg.svg, err = svg.NewKind()
	if err != nil {
		panic(fmt.Sprintf("error while initializing the svg Kind: %s", err))
	}
	reg.all = append(reg.all, reg.svg)

	reg.stars, err = stars.NewKind(rt)
	if err != nil {
		panic(fmt.Sprintf("error while initializing the stars Kind: %s", err))
	}
	reg.all = append(reg.all, reg.stars)

	return reg
}
