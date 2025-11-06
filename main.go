package main

import (
	"log"
	"os"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtklayershell"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/diamondburned/gotk4/pkg/pango"
)

func main() {
	app := gtk.NewApplication("com.github.animainmi.VoBriGo", 0)
	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

var attrs = []*pango.Attribute{
	pango.NewAttrScale(10),
	pango.NewAttrWeight(pango.WeightBold),
}

func activate(app *gtk.Application) {
	if !gtklayershell.IsSupported() {
		log.Fatalln("gtk-layer-shell not supported")
	}

	labelAttrs := pango.NewAttrList()
	for _, attr := range attrs {
		labelAttrs.Insert(attr)
	}

	label := gtk.NewLabel("This window floats in the middle.")
	label.SetAttributes(labelAttrs)
	label.Show()

	appwin := gtk.NewApplicationWindow(app)
	appwin.Add(label)
	appwin.SetTitle("Layer Shell Example")

	window := &appwin.Window

	gtklayershell.InitForWindow(window)
	gtklayershell.SetLayer(window, gtklayershell.LayerShellLayerTop)

	for edge := gtklayershell.Edge(0); edge < gtklayershell.LayerShellEdgeEntryNumber; edge++ {
		gtklayershell.SetAnchor(window, edge, false)
	}

	appwin.Show()
}
