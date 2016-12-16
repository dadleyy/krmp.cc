package krmp

import "regexp"

type Terminal func(*RequestRuntime) (Result, error)
type Middleware func(Terminal) Terminal

type Route struct {
	Method  string
	Path    *regexp.Regexp
	Handler Terminal
}
