package rlog

import (
	"fmt"
	"io"
	"os"
	"sync"

	apex "github.com/apex/log"

	"github.com/fatih/color"
	colorable "github.com/mattn/go-colorable"
)

var Default = NewHandler(os.Stderr)

var bold = color.New(color.Bold)

var colors = [...]*color.Color{
	apex.DebugLevel: color.New(color.FgWhite),
	apex.InfoLevel:  color.New(color.FgBlue),
	apex.WarnLevel:  color.New(color.FgYellow),
	apex.ErrorLevel: color.New(color.FgRed),
	apex.FatalLevel: color.New(color.FgMagenta),
}

// Strings mapping.
var prefixes = [...]string{
	apex.DebugLevel: "DBUG",
	apex.InfoLevel:  "INFO",
	apex.WarnLevel:  "WARN",
	apex.ErrorLevel: "EROR",
	apex.FatalLevel: "CRIT",
}

type LvlMap map[string]apex.Level

type Handler struct {
	mu     sync.Mutex
	Lvls   LvlMap
	Writer io.Writer
}

func NewHandler(w io.Writer) apex.Handler {
	lvlmap := parseEnv()

	if f, ok := w.(*os.File); ok {
		return &Handler{
			Writer: colorable.NewColorable(f),
			Lvls:   lvlmap,
		}
	}

	return &Handler{
		Writer: w,
		Lvls:   LvlMap{"*": apex.DebugLevel},
	}
}

func NewHandlerWithPkgMap(w io.Writer, lvlmap LvlMap) apex.Handler {
	if f, ok := w.(*os.File); ok {
		return &Handler{
			Writer: colorable.NewColorable(f),
			Lvls:   lvlmap,
		}
	}

	return &Handler{
		Writer: w,
		Lvls:   lvlmap,
	}
}

func (h *Handler) shouldLog(pkg string, e *apex.Entry) bool {
	if maxLvl, ok := h.Lvls[pkg]; ok {
		return e.Level >= maxLvl
	} else {
		return e.Level >= h.Lvls["*"]
	}
}

func (h *Handler) HandleLog(e *apex.Entry) error {
	pkg := getPackage()
	if !h.shouldLog(pkg, e) {
		return nil
	}

	var (
		color = colors[e.Level]
		level = prefixes[e.Level]
		names = e.Fields.Names()
		msg   = e.Message
	)

	h.mu.Lock()
	defer h.mu.Unlock()

	if e.Level == apex.FatalLevel {
		msg = bold.Sprint(msg)
	}
	color.Fprintf(h.Writer, "%s %-25s", bold.Sprintf("%*s|", 4, level), msg)

	for _, name := range names {
		if name == "source" {
			continue
		}
		fmt.Fprintf(h.Writer, " %s=%#v", bold.Sprint(color.Sprint(name)), e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)
	return nil
}
