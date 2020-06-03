package main

import (
	"fmt"
	"log"
	"strings"
	"spotify/gospotify"
	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/glib"
)

//go build -ldflags -H=windowsgui    flags to hide cmd in exe

type Dir struct {
	s string
}
func (d *Dir) SetString(ss string) {
	d.s = ss
}
func main() {
	fmt.Println("Welcome to the spotify data analyser.")
	tracks := gospotify.OpenJsonTracks("StreamingHistory",
		"C:\\Users\\Sayed\\Desktop\\my_spotify_data\\MyData")
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
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Spotify Analyser")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	//Create a new label widget to show in the window.
	l, err := gtk.LabelNew("Welcome!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// // Add the label to the window.
	// win.Add(l)
	buttons := createButtons(win)

	for i := range buttons {
		grid.Attach(buttons[i],i,i,1,1)
	}
	grid.Attach(l,2,2,1,1)
	win.Add(&grid.Container.Widget)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func createButtons(win *gtk.Window) []*gtk.Button {
	var buttons []*gtk.Button
	dirClicked := false
	var dir string
	
	dirBtn, err := gtk.ButtonNewWithLabel("Data DIR")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	totalTimeBtn, err := gtk.ButtonNewWithLabel("Total Time Played")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}

	dirBtn.Connect("clicked", func() {
		dirClicked = true
		fileDialogue, _ := gtk.FileChooserDialogNewWith2Buttons("Select folder",
			win,gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Cancel", gtk.RESPONSE_CANCEL, 
   			"OK", gtk.RESPONSE_OK)


   		res := fileDialogue.Run()
   		if (res == gtk.RESPONSE_OK) {
   			dir = fileDialogue.FileChooser.GetURI()
   			dir = strings.Trim(dir, "file:///")
   			dir = strings.Replace(dir,"/","\\",-1)
   			fmt.Println(dir)
   			fileDialogue.Destroy()

   		}
		if (res == gtk.RESPONSE_CANCEL) {
			fileDialogue.Destroy()
		}
	})

	buttons = append(buttons,dirBtn)
	buttons = append(buttons,totalTimeBtn)

	return buttons

}