package main

import (
	_ "embed"
	"log"
	"os"
	"strings"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var styleCSS string

func main() {
	app := gtk.NewApplication("com.github.animainmi.VoBriGo", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	if !gtk4layershell.IsSupported() {
		log.Fatalln("gtk-layer-shell not supported")
	}

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), loadCSS(styleCSS),
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	popover := gtk.NewPopover()
	appwin := gtk.NewApplicationWindow(app)
	window := &appwin.Window
	window.SetTitle("gotk4 Example")
	window.SetChild(popover)
	window.SetDefaultSize(700, 25)
	window.SetVisible(true)

	gtk4layershell.InitForWindow(window)
	gtk4layershell.SetLayer(window, gtk4layershell.LayerShellLayerTop)
	gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeTop, true)
}

func loadCSS(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		// Optional line parsing routine.
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromString(content)
	return prov
}
