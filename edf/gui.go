package edf

import (
	"errors"
	"fmt"
	"log"
	"os"

	//"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xwidget "fyne.io/x/fyne/widget"
)

var (
	win fyne.Window

	downloadEntry *widget.Label
	downloadFile  *widget.Button

	EnedisBegin *widget.Label
	EnedisEnd   *widget.Label

	calButton     *widget.Button
	calPopup      *widget.PopUp
	calBeginEntry *widget.Label
	calEndEntry   *widget.Label

	tempoProgress *widget.ProgressBar
	HcEntry       *widget.Entry
	analButton    *widget.Button
	largeText     *widget.Entry
)

/*
 */
func dialogDisplay() fyne.CanvasObject {
	largeText = widget.NewMultiLineEntry()
	largeText.SetMinRowsVisible(20)

	return largeText
}

/*
Mise à jour créneau HC par défaut
*/
func confirmCallback(response bool) {
	if !response {
		HcEntry.SetText("22:00-06:00")
		//slotHC, _ := buildHeuresCreuses(HcEntry.Text)

		//updateHoursDialog(slotHC)
	}
}

/*
 */
func analyze(win fyne.Window) {
	slotHC, err := buildHeuresCreuses(HcEntry.Text)
	if err != nil {
		cnf := dialog.NewConfirm("Erreur de Saisie  -  Format incorrect:",
			"hh:mn-hh:mn\n hh:mn-hh:mn,hh:mn-hh:mn\nhh:mn-hh:mn,hh:mn-hh:mn,hh:mn-hh:mn ",
			confirmCallback, win)
		cnf.SetDismissText("Défaut")
		cnf.SetConfirmText("OK")
		cnf.Show()
	} else {
		//updateHoursDialog(slotHC)
		consosYear := conso_totale(data, datesTempo, slotHC)

		prices := findPriceByPower(powerSelection)
		largeText.SetText("")

		result := displayYear(consosYear, prices)
		largeText.Append(result)

	}
}

/*
Affichage de l'aide pour l'application
*/
func displayHelp() {
	helpByte, err := os.ReadFile("./aide.txt")
	if err == nil {
		largeText.SetText(string(helpByte))
	}
}

/*
 */
func dialogAnalyze(app fyne.App, win fyne.Window) fyne.CanvasObject {
	labelVide := widget.NewLabel("   ")
	aidButton := widget.NewButton("  Aide  ", func() {
		displayHelp()

	})
	quitButton := widget.NewButton(" Quitter", func() {
		app.Quit()
	})
	analButton = widget.NewButton("Calculer", func() {
		analyze(win)
	})
	analButton.Disable()
	contQC := container.NewHBox(
		aidButton,
		labelVide,
		quitButton,
		labelVide,
		analButton,
	)
	return container.NewCenter(contQC)
}

/*
 */
func dialogEdfTarif(win fyne.Window) fyne.CanvasObject {
	var err error

	edfTarifLabel := widget.NewLabel("Tarifs EDF au: ")
	// edfTarifButton := widget.NewButton("Tarifs EDF au:", func() {
	// 	NewHoursDialog(win)
	// })
	edfDateLabel := widget.NewLabel("--/--/----         ")
	edfTarifs, err = loadFare()
	if err != nil {
		err = fetchFare()
		if err != nil {
			edfDateLabel.SetText("inconnus!")
		} else {
			edfDateLabel.SetText(edfTarifs.DateFare.Format(DMYLAYOUT))
		}
	} else {
		edfDateLabel.SetText(edfTarifs.DateFare.Format(DMYLAYOUT))
	}
	edfMajButton := widget.NewButton("M. à jour", func() {
		err = fetchFare()
		if err != nil {
			dialog.ShowError(err, win)
		} else {
			edfDateLabel.SetText(edfTarifs.DateFare.Format(DMYLAYOUT))
		}
	})

	return container.NewHBox(
		edfTarifLabel, edfDateLabel, edfMajButton,
	)
}

/*
Container définition créneaux HC
*/
func dialogHeuresCreuses() fyne.CanvasObject {
	//hcLabel := widget.NewLabel("Créneaux H. Creuses ")
	hcButton := widget.NewButton(" Créneaux H. Creuses ", func() {
		NewHoursDialog(win)
	})
	hcButton.Resize(fyne.NewSize(160, 36))
	hcButton.Move(fyne.NewPos(10, 0))
	// hcLabel.Resize(fyne.NewSize(140, 36))
	// hcLabel.Move(fyne.NewPos(0, 0))
	HcEntry = widget.NewEntry()
	HcEntry.SetPlaceHolder("hh:mn-hh:mn,hh:mn-hh:mn")

	// hh:mn
	regHm := `((0[0-9]|1[0-9]|2[0-3]):[0-5]\d)`
	// hh:mn-hh:mn
	regSlot := regHm + `-` + regHm
	// Formats des créneau HC acceptés
	//	hh:mn-hh:mn
	//	hh:mn-hh:mn,hh:mn-hh:mn
	//	hh:mn-hh:mn,hh:mn-hh:mn,hh:mn-hh:mn
	regSlots := `^` +
		regSlot + `$` +
		`|` + regSlot + `,` + regSlot + `$` +
		`|` + regSlot + `,` + regSlot + `,` + regSlot + `$`

	HcEntry.Validator = validation.NewRegexp(`^`+regSlots, "")

	HcEntry.Resize(fyne.NewSize(360, 36))
	HcEntry.Move(fyne.NewPos(180, 0))

	return container.NewWithoutLayout(
		//hcLabel, HcEntry,
		hcButton, HcEntry,
	)
}

/*
Container sélection de la puissance installée
*/
func dialogPower() fyne.CanvasObject {
	powerLabel := widget.NewLabel("Puissance Installée")
	powerItems := []string{"3 kVA", "6 kVA", "9 kVA", "12 kVA", "15 kVA",
		"18 kVA", "24 kVA", "30 kVA", "36 kVA"}
	powerSelect := widget.NewSelect(powerItems, func(value string) {
		powerSelection = value
	})
	powerSelect.SetSelected(powerItems[2])
	return container.NewHBox(powerLabel, powerSelect)
}

/*
Widget Progression chargement heures creuses
*/
func dialogTempoProgress() fyne.CanvasObject {
	tempoProgress = widget.NewProgressBar()
	tempoProgress.TextFormatter = func() string {
		return fmt.Sprintf("Récupération jours Tempo:   %.0f sur %.0f", tempoProgress.Value*365, tempoProgress.Max*365)
	}
	tempoProgress.SetValue(0.0)
	//tempoProgress.Resize(fyne.NewSize(200, 20))
	//cont := container.NewVBox(tempoProgress)
	//cont.Resize(fyne.NewSize(40, 140))

	return tempoProgress
}

/*
Sélection jour sur le calendrier
*/
func onSelectedCal(t time.Time) {
	var err error

	data, err = readEnedisFile(lines, t.Format(DMYLAYOUT))
	calPopup.Hide()
	if err == nil {
		to := t.Format(DMYLAYOUT)
		calEndEntry.SetText(to)
		from := t.AddDate(-1, 0, 0).Format(DMYLAYOUT)
		calBeginEntry.SetText(from)

		datesTempo, err = updateTempo(from, to)
		if err != nil {
			dialog.ShowError(errors.New("lecture données edf impossible"), win)
		}

		analButton.Enable()
	} else {
		calEndEntry.SetText("")
		calBeginEntry.SetText("")
		dialog.ShowError(err, win)
		analButton.Disable()
	}

}

/*
Container sélection dates d'analyse (durée 1 an)
*/
func dialogCalendar() fyne.CanvasObject {

	startingDate := time.Now()

	calBeginLabel := widget.NewLabel("Date Début   ")
	calBeginEntry = widget.NewLabel("--/--/----")
	containerBegin := container.NewHBox(calBeginLabel, calBeginEntry)

	calEndLabel := widget.NewLabel("Date Fin  ")
	calEndEntry = widget.NewLabel(startingDate.Format(DMYLAYOUT))
	calEndLabel.Alignment = fyne.TextAlignCenter
	calEndEntry.Alignment = fyne.TextAlignCenter

	calendar := xwidget.NewCalendar(startingDate, onSelectedCal)

	calButton = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {

		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(calButton)
		position.Y += calButton.Size().Height
		position.X += calButton.Size().Height

		calPopup = widget.NewPopUp(calendar, fyne.CurrentApp().Driver().CanvasForObject(calButton))
		calPopup.ShowAtPosition(position)
	})
	calButton.Disable()

	containerEnd := container.NewHBox(calEndLabel, calEndEntry, calButton)

	return container.NewGridWithColumns(2,
		containerBegin, containerEnd,
	)
}

/*
Container de sélection du fichier Enedis.cvs
*/
func dialogFileOpen(win fyne.Window, label string, filter []string) fyne.CanvasObject {
	empty := "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - "
	downloadLabel := widget.NewLabel(label)
	downloadEntry = widget.NewLabel(empty)

	downloadFile = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			loadFile := reader.URI().Path()
			downloadEntry.SetText(truncate(loadFile, 38))

			err = checkEnedisFile(loadFile)
			if err == nil {
				EnedisBegin.SetText(enedisStart.Format(DMYLAYOUT))
				EnedisEnd.SetText(enedisEnd.Format(DMYLAYOUT))
				calButton.Enable()
			} else {
				dialog.ShowError(err, win)
				calButton.Disable()
				analButton.Disable()
			}
		}, win)

		fd.SetFilter(storage.NewExtensionFileFilter(filter))
		fd.Show()
	})

	return container.NewHBox(
		downloadLabel, downloadEntry, downloadFile,
	)

}

/*
Container affichage dates début et fin du fichier Enedis
*/
func dialogDateSlot() fyne.CanvasObject {
	EnedisEndLabel := widget.NewLabel("Fin Enedis")
	EnedisEnd = widget.NewLabel("--/--/----")
	contEnedisEnd := container.NewHBox(
		EnedisEndLabel, EnedisEnd,
	)
	return container.NewGridWithColumns(2,
		dialogDateEnedis(), contEnedisEnd)
}

/*
Container date début du fichier Enedis
*/
func dialogDateEnedis() fyne.CanvasObject {
	EnedisBeginLabel := widget.NewLabel("Début Enedis")
	EnedisBegin = widget.NewLabel("--/--/----")

	return container.NewHBox(
		EnedisBeginLabel, EnedisBegin,
	)
}

/*
Container dialogue
*/
func buildWinEDF(app fyne.App, win fyne.Window) fyne.CanvasObject {
	contAnal := dialogAnalyze(app, win)

	return container.NewVBox(
		dialogFileOpen(win, "Fich Enedis.cvs", []string{".csv"}),
		dialogDateSlot(),
		dialogCalendar(),
		dialogTempoProgress(),
		dialogPower(),
		dialogHeuresCreuses(),
		dialogEdfTarif(win),
		dialogDisplay(),
		contAnal,
	)
}

/*
Fenêtre principale - caractéristiques
*/
func NewWinEDF(app fyne.App) fyne.Window {
	win = app.NewWindow("Comparaison Tarifs EDF V1.0")
	ctn := container.NewVBox(buildWinEDF(app, win))
	win.SetContent(ctn)
	win.Resize(fyne.NewSize(550, 600))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetMaster()
	win.CenterOnScreen()

	win.ShowAndRun()
	return win
}
