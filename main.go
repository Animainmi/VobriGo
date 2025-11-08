package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
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

	menu := gio.NewMenu()
	menu.Append("Item 1", "app.item1")
	menu.Append("Item 2", "app.item2")
	menu.Append("Item 3", "app.item3")

	action1 := gio.NewSimpleAction("item1", nil)
	action2 := gio.NewSimpleAction("item2", nil)
	action3 := gio.NewSimpleAction("item3", nil)

	app.AddAction(action1)
	app.AddAction(action2)
	app.AddAction(action3)

	popoverMenu := gtk.NewPopoverMenuFromModel(menu)

	menuButton := gtk.NewMenuButton()

	menuButton.SetPopover(popoverMenu)

	BrightnessAdjustment := gtk.NewAdjustment(50, 0, 100, 1, 0, 0)
	VolumeAdjustment := gtk.NewAdjustment(50, 0, 100, 1, 0, 0)

	BriScale := gtk.NewScale(gtk.OrientationHorizontal, BrightnessAdjustment)
	VoScale := gtk.NewScale(gtk.OrientationHorizontal, VolumeAdjustment)

	BriScale.SetHExpand(true)
	BriScale.ConnectValueChanged(func() {
		brightnessFormatted := fmt.Sprintf("%d%%", int(BrightnessAdjustment.Value()))
		cmd := exec.Command("brightnessctl", "set", brightnessFormatted)
		cmd.Run()
		realBri := exec.Command("brightnessctl", "get")
		out, err := realBri.Output()
		if err != nil {
			log.Fatalf("Failed to execute command: %v", err)
		}
		fmt.Println(string(out))
	})

	VoScale.SetHExpand(true)
	VoScale.ConnectValueChanged(func() {
		volumeFormatted := fmt.Sprintf("%d%%", int(VolumeAdjustment.Value()))
		cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", volumeFormatted)
		cmd.Run()
	})

	mainBox.SetEndWidget(BriScale)
	mainBox.SetStartWidget(menuButton)

	window.SetTitle("gotk4 Example")
	window.SetChild(mainBox)

	window.SetDefaultSize(700, 25)
	window.SetVisible(true)

	gtk4layershell.InitForWindow(window)
	gtk4layershell.SetLayer(window, gtk4layershell.LayerShellLayerTop)
	gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeTop, true)
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
