package main

import (
	"fmt"
	"container/list"
	"log"
	"strings"
	"spotify/gospotify"
	"github.com/gotk3/gotk3/gtk"
	//"strconv"
	//"github.com/gotk3/gotk3/glib"
)
var tracks gospotify.Tracks
var dirChosen bool = false
var labelList *list.List = list.New()
//go build -ldflags -H=windowsgui    flags to hide cmd in exe

func main() {
	fmt.Println("Welcome to the spotify data analyser.")
	// tracks.TotalTimePlayed("Days")
	// tracks.FindArtistTracksNo("Ariana Grande")
	// tracks.FindArtistTracksNo("Tame Impala")
	// tracks.FindArtistTracksNo("Justin Bieber")
	// tracks.FindArtistTracksNo("Jhené Aiko")
	// tracks.FindArtistTracksNo("Drake")
	// tracks.FindTrackName("Borderline")
	// tracks.FindTrackName("Bad Day")
	// tracks.FindTrackName("Lost in Yesterday")
	// tracks.FindTrackName("Imagination")
	// tracks.FindTrackName("God is a Woman")
	// tracks.FindArtistTracks("Future")
	// tracks.FindArtistPlayed()
	// tracks.AverageTimePlayed("Seconds")

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
	welcome := createLabelObj("")
	welcome.SetMarkup("<big><b>Welcome to the Spotify Analyser</b></big>")

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

func createEmptyButton(label string) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}
func createLabelObj(label string) *gtk.Label {
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
func setup_entry() *gtk.Entry {
	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create new entry ", err)
	}
	entry.SetVisibility(true)
	return entry
}

func getDialogContentArea(dialog *gtk.MessageDialog) *gtk.Box {
	box, err := dialog.GetContentArea()
	if err != nil {
		log.Fatal("Unable to get content area of dialog: ", err)
	}
	return box
}

func getEntryText(entry *gtk.Entry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text of entry: ", err)
	}
	return text
}

func createLabel(grid *gtk.Grid, s string) {
	l := createLabelObj("")
	l.SetMarkup(s)
	labelList.PushBack(l)
	grid.Add(l)
	l.SetHExpand(true)
	grid.ShowAll()
}

func createLabelNoMarkup(grid *gtk.Grid, s string) {
	l := createLabelObj(s)
	labelList.PushBack(l)
	grid.Add(l)
	l.SetHExpand(true)
	grid.ShowAll()
}

func removeAllLabels() {
	for (labelList.Len() != 0) {
		e := labelList.Back()
		lab, ok := labelList.Remove(e).(*gtk.Label)
		if !ok {
			log.Print("Element to remove is not a *gtk.Label")
			return
		}
		lab.Destroy()
	}
}

func createButtons(win *gtk.Window, labelsGrid *gtk.Grid) []*gtk.Button {
	var buttons []*gtk.Button
	var dir string

	dirBtn := createEmptyButton("Choose Data Folder")
	totalTimeBtn := createEmptyButton("Time Played")
	FindArtistBtn := createEmptyButton("Find Artist Tracks")
	clearBtn := createEmptyButton("Clear Output")

	dirBtn.Connect("clicked", func() {

		fileDialogue, _ := gtk.FileChooserDialogNewWith2Buttons("Select folder",
			win,gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Cancel", gtk.RESPONSE_CANCEL, 
   			"OK", gtk.RESPONSE_OK)

   		res := fileDialogue.Run()
   		if (res == gtk.RESPONSE_OK) {
   			removeAllLabels()
   			dir = fileDialogue.FileChooser.GetURI()
   			dir = strings.Trim(dir, "file:///")
   			dir = strings.Replace(dir,"/","\\",-1)
   			dirChosen = true
   			fmt.Println(dir)
   			var fileURLs string
   			tracks, fileURLs = gospotify.OpenJsonTracks("StreamingHistory",dir)
   			if len(tracks) == 0 {
   				createLabel(labelsGrid, 
   					"\nNo StreamingHistory.json " +
   					"files found in this directory: \n" + dir + "\n")
   			} else {
   				createLabel(labelsGrid, "\n" + fileURLs + "\n")
   			}
   			fileDialogue.Destroy()

   		}
		if (res == gtk.RESPONSE_CANCEL) {
			fileDialogue.Destroy()
		}
	})

	totalTimeBtn.Connect("clicked", func() {
		if (dirChosen == false) {
			createLabel(labelsGrid, 
				"<b>Please specify the data folder location</b>\n")
		} else {
			btnDlg := setup_dialog(win, 
				"Select the time format (days, hours, minutes), default: seconds")
			dlgBox := getDialogContentArea(btnDlg)

			userEntry := setup_entry()
			dlgBox.PackEnd(userEntry, false, false, 0)
			btnDlg.ShowAll()

			res := btnDlg.Run()
			input := getEntryText(userEntry)
			if (res == gtk.RESPONSE_OK && len(input) != 0) {
				input = strings.ToLower(input)
				value, format := tracks.TotalTimePlayed(input)
				value2, _:= tracks.AverageTimePlayed("seconds")
				s := fmt.Sprintf("<b>Total time played:</b> %.3f %s.\n" +
					"<b> Average time played:</b> %.3f seconds.\n", value, format, 
					value2)
				createLabel(labelsGrid, s)
			} else {
				createLabel(labelsGrid, "<b>Please enter a value.</b>\n")
			}
			btnDlg.Destroy()
		}
	})

	FindArtistBtn.Connect("clicked", func() {
		if (dirChosen == false) {
			createLabel(labelsGrid, 
				"<b>Please specify the data folder location</b>")
		} else {
			btnDlg := setup_dialog(win, 
				"Select the Artist to search")
			dlgBox := getDialogContentArea(btnDlg)
			userEntry := setup_entry()
			dlgBox.PackEnd(userEntry, false, false, 0)
			btnDlg.ShowAll()

			res := btnDlg.Run()
			input := getEntryText(userEntry)
			if (res == gtk.RESPONSE_OK && len(input) != 0) {
				noOfTracks := tracks.FindArtistTracksNo(input)
				artistTracks := tracks.FindArtistTracks(input).ToString()
				if len(artistTracks) == 0 {
					createLabel(labelsGrid, "<b>No Tracks found for " + input + "</b>\n")
				} else {
					s := fmt.Sprintf("<b>No. of times %s was played:</b> %v", input,
					 noOfTracks)
					createLabel(labelsGrid, s)
					createLabelNoMarkup(labelsGrid, artistTracks + "\n")
				}
			} else {
				createLabel(labelsGrid, "<b>Please enter an artist.</b>\n")
			}
			btnDlg.Destroy()
		}
	})

	clearBtn.Connect("clicked", func() {
		removeAllLabels()
	})

	buttons = append(buttons,dirBtn)
	buttons = append(buttons,totalTimeBtn)
	buttons = append(buttons,FindArtistBtn)
	buttons = append(buttons, clearBtn)

	return buttons

}