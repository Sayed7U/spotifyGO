package main

import (
	"fmt"
	"log"
	"strings"
	"spotify/gospotify"
	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/glib"
)
var tracks gospotify.Tracks
var dirChosen bool = false
//go build -ldflags -H=windowsgui    flags to hide cmd in exe

type Dir struct {
	s string
}
func (d *Dir) SetString(ss string) {
	d.s = ss
}
func main() {
	fmt.Println("Welcome to the spotify data analyser.")
	tracks.TotalTimePlayed("Days")
	tracks.FindArtist("Ariana Grande")
	tracks.FindArtist("Tame Impala")
	tracks.FindArtist("Justin Bieber")
	tracks.FindArtist("Jhené Aiko")
	tracks.FindArtist("Drake")
	tracks.FindTrackName("Borderline")
	tracks.FindTrackName("Bad Day")
	tracks.FindTrackName("Lost in Yesterday")
	tracks.FindTrackName("Imagination")
	tracks.FindTrackName("God is a Woman")
	tracks.FindArtistTracks("Future")
	tracks.FindArtistPlayed()
	tracks.AverageTimePlayed("Seconds")

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win := setup_window("Spotify Analyser")

	box := setup_box(gtk.ORIENTATION_VERTICAL)

	sw := setup_sw()

	sw.SetHExpand(true)
	sw.SetVExpand(true)

	labelsGrid := setup_grid(gtk.ORIENTATION_VERTICAL)
	sw.Add(labelsGrid)
	labelsGrid.SetHExpand(true)

	// grid, err := gtk.GridNew()
	// if err != nil {
	// 	log.Fatal("Unable to create grid:", err)
	// }
	// grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	//Create a new label widget to show in the window.
	welcome := createLabel("Welcome to the Spotify Analyser")

	// // Add the label to the window.
	// win.Add(l)
	box.PackStart(welcome, false, false, 0)
	buttons := createButtons(win, labelsGrid)

	for i := range buttons {
		box.PackStart(buttons[i], false, false, 0)
	}

	box.PackStart(sw,true, true, 0)
	win.Add(box)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func setup_window(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func setup_box(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

func setup_grid(orient gtk.Orientation) *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(orient)

	return grid
}

func setup_tview() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv
}

func get_buffer_from_tview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func get_text_from_tview(tv *gtk.TextView) string {
	buffer := get_buffer_from_tview(tv)
	start, end := buffer.GetBounds()

	text, err := buffer.GetText(start, end, true)
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	return text
}

func set_text_in_tview(tv *gtk.TextView, text string) {
	buffer := get_buffer_from_tview(tv)
	buffer.SetText(text)
}

func createEmptyButton(label string) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}
func createLabel(label string) *gtk.Label {
	l, err := gtk.LabelNew(label)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	return l
}
func setup_sw() *gtk.ScrolledWindow {
	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	return sw
}

func setup_dialog(win *gtk.Window, message string) *gtk.MessageDialog {
	dialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, 
		gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", message)
	return dialog
}

func createButtons(win *gtk.Window, labelsGrid *gtk.Grid) []*gtk.Button {
	var buttons []*gtk.Button
	var dir string

	dirBtn := createEmptyButton("Choose Data Folder")
	totalTimeBtn := createEmptyButton("Total Time played")

	dirBtn.Connect("clicked", func() {
		fileDialogue, _ := gtk.FileChooserDialogNewWith2Buttons("Select folder",
			win,gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Cancel", gtk.RESPONSE_CANCEL, 
   			"OK", gtk.RESPONSE_OK)


   		res := fileDialogue.Run()
   		if (res == gtk.RESPONSE_OK) {
   			dir = fileDialogue.FileChooser.GetURI()
   			dir = strings.Trim(dir, "file:///")
   			dir = strings.Replace(dir,"/","\\",-1)
   			dirChosen = true
   			fmt.Println(dir)
   			tracks = gospotify.OpenJsonTracks("StreamingHistory",dir)
   			fileDialogue.Destroy()

   		}
		if (res == gtk.RESPONSE_CANCEL) {
			fileDialogue.Destroy()
		}
	})

	totalTimeBtn.Connect("clicked", func() {
		if (dirChosen == false) {
			l := createLabel("Please specify the data folder location")
			labelsGrid.Add(l)
			l.SetHExpand(true)
			labelsGrid.ShowAll()
		} else {
			// dlg := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, 
			// 	gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", "Pleas")
			value := tracks.TotalTimePlayed("Days")
			s := fmt.Sprintf("Total time played: %.3f.", value)
			l := createLabel(s)
			labelsGrid.Add(l)
			l.SetHExpand(true)
			labelsGrid.ShowAll()
		}
	})

	buttons = append(buttons,dirBtn)
	buttons = append(buttons,totalTimeBtn)

	return buttons

}