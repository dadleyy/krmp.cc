package krmp

import "regexp"

type Terminal func(*RequestRuntime)
type Middleware func(Terminal) Terminal

type Route struct {
	Method  string
	Path    *regexp.Regexp
	Handler Terminal
}
