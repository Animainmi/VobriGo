package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var styleCSS string

func main() {
	app := gtk.NewApplication("com.github.animainmi.VobriGo", gio.ApplicationFlagsNone)
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

	appwin := gtk.NewApplicationWindow(app)
	window := &appwin.Window
	mainBox := gtk.NewCenterBox()

	BrightnessAdjustment := gtk.NewAdjustment(50, 0, 100, 1, 0, 0)
	VolumeAdjustment := gtk.NewAdjustment(50, 0, 100, 1, 0, 0)

	BriScale := gtk.NewScale(gtk.OrientationHorizontal, BrightnessAdjustment)
	VoScale := gtk.NewScale(gtk.OrientationHorizontal, VolumeAdjustment)

	BriScale.SetHExpand(true)
	SetupBrightness(BrightnessAdjustment, BriScale)

	VoScale.SetHExpand(true)
	VoScale.ConnectValueChanged(func() {
		volumeFormatted := fmt.Sprintf("%d%%", int(VolumeAdjustment.Value()))
		cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", volumeFormatted)
		cmd.Run()
	})

	timeLabel := gtk.NewButtonWithLabel("time")

	mainBox.SetCenterWidget(BriScale)
	mainBox.SetEndWidget(timeLabel)

	window.SetTitle("gotk4 Example")
	window.SetChild(mainBox)

	window.SetDefaultSize(700, 25)
	window.SetVisible(true)

	gtk4layershell.InitForWindow(window)
	gtk4layershell.SetLayer(window, gtk4layershell.LayerShellLayerTop)
	gtk4layershell.AutoExclusiveZoneEnable(window)
	gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeTop, true)

	go func() {
		for t := range time.Tick(time.Second) {
			currentTime := t
			glib.IdleAdd(func() {
				timeLabel.SetLabel(fmt.Sprintf(
					"%s", currentTime.Format("15:04"),
				))
			})
		}
	}()
}

func SetupBrightness(adj *gtk.Adjustment, scale *gtk.Scale) {
	get := func() (current, max float64) {
		out, _ := exec.Command("brightnessctl", "g").Output()
		maxOut, _ := exec.Command("brightnessctl", "m").Output()
		current, _ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
		max, _ = strconv.ParseFloat(strings.TrimSpace(string(maxOut)), 64)
		return
	}

	// Initial update
	if cur, max := get(); max > 0 {
		adj.SetLower(0)
		adj.SetUpper(max)
		adj.SetValue(cur)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			if cur, max := get(); max > 0 {
				glib.IdleAdd(func() { adj.SetValue(cur) })
			}
		}
	}()
	scale.Connect("destroy", func() { ticker.Stop() })

	scale.ConnectValueChanged(func() {
		p := int(math.Max(adj.Value(), 1))
		exec.Command("brightnessctl", "s", fmt.Sprintf("%d", p), "-q").Run()
	})
}

func loadCSS(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromString(content)
	return prov
}
