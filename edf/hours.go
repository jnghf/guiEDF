package edf

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type hoursSlot struct {
	begin int
	end   int
}

/*
Création fenêtre de dialogue
pour la saisie des créneaux Heures Creuses
pas de 30mn (0.5h)
*/
func NewHoursDialog(win fyne.Window) {
	//48 boutons controlables à la souris (24h*2)
	buttons := []fyne.CanvasObject{}
	for i := 0; i < 24*2; i++ {
		buttons = append(buttons, NewHoldableButton(""))
	}

	contAm := container.NewHBox(buttons[0:24]...)
	contPm := container.NewHBox(buttons[24:]...)
	// 	Labels AM et PM (12h)
	labelsAm := []fyne.CanvasObject{}
	labelsPm := []fyne.CanvasObject{}
	for i := 0; i < 12; i++ {
		labelsAm = append(labelsAm, widget.NewLabel(fmt.Sprintf("%02d ", i)))
		labelsPm = append(labelsPm, widget.NewLabel(fmt.Sprintf("%02d ", i+12)))
	}
	contLabelAm := container.NewHBox(labelsAm...)
	contLabelPm := container.NewHBox(labelsPm...)

	//  Création form dialog
	items := []*widget.FormItem{
		widget.NewFormItem("AM", contLabelAm),
		widget.NewFormItem("", contAm),

		widget.NewFormItem("", contPm),
		widget.NewFormItem("PM", contLabelPm),
		widget.NewFormItem("", widget.NewLabel("")), // pour espace
	}

	slotHC, err := buildHeuresCreuses(HcEntry.Text)
	if err == nil {
		updateHoursDialog(buttons, slotHC)
	}

	dialog.ShowForm("Créneaux Heures Creuses", "Valider", "Annuler", items, func(b bool) {
		if !b {
			//  Annuler
			return
		} else {
			// Valider
			HcEntry.SetText(find_slots(buttons))
		}
	}, win)

}

/*
Conversion "hh:mn" --> int  (0 .. 47 ; pas de 30mn)
*/
func hmToInt(str string) int {
	sp := strings.Split(str, ":")
	h, _ := strconv.Atoi(sp[0])
	m, _ := strconv.Atoi(sp[1])
	return h*2 + m/30
}

/*
Mise à jour des boutons Heures Creuses en fonction d'une entrée manuelle
"hh:mn-hh:mn,hh:mn-hh:mn,   ..........."
*/

func updateHoursDialog(buttons []fyne.CanvasObject, slotHC []HCreuses) {
	//  Raz boutons HC
	for i := 0; i < 48; i++ {
		buttons[i].(*HoldableButton).FocusLost()
		buttons[i].(*HoldableButton).On = false
	}

	for _, slot := range slotHC {
		hs := hmToInt(slot.hstart)
		he := hmToInt(slot.hend)

		if he > hs {
			for i := hs; i < he; i++ {
				buttons[i].(*HoldableButton).FocusGained()
				buttons[i].(*HoldableButton).On = true
			}
		} else {
			//  chevauchement 00h00
			for i := hs; i < 48; i++ {
				buttons[i].(*HoldableButton).FocusGained()
				buttons[i].(*HoldableButton).On = true
			}
			for i := 0; i < he; i++ {
				buttons[i].(*HoldableButton).FocusGained()
				buttons[i].(*HoldableButton).On = true
			}
		}

	}
}

/*
Conversion entier = 48 * 0.5 h en string "hh:mn"
hh: 0 --> 23
mn: 00 ou 30
*/
func intToHm(i int) (str string) {
	var s string
	if i%2 == 0 {
		s = "00"
	} else {
		s = "30"
	}
	return fmt.Sprintf("%02d:%s", i/2, s)
}

/*
Conversion slots H Début-H Fin en chaîne
Ex:
[{3 8} {30 36} {40 45}]  --> "01:30-04:00,15:00-18:00,20:00-22:30"
*/
func slotsToString(hourSlots []hoursSlot) (str string) {
	str = ""
	for i, slot := range hourSlots {
		str += fmt.Sprintf("%s-%s", intToHm(slot.begin), intToHm(slot.end))
		if i < len(hourSlots)-1 {
			str += ","
		}
	}
	return str
}

/*
Récupération des slots Heures Creuses sélectionnés
*/
func find_slots(hourButtons []fyne.CanvasObject) string {

	var bS bool = false // mémo début slot
	var hourSlots []hoursSlot
	var hS hoursSlot

	for i, h := range hourButtons {

		if h.(*HoldableButton).On {
			if !bS {
				hS.begin = i
				bS = true
			}
			if i == 47 {
				hS.end = 0
				hourSlots = append(hourSlots, hS)
				bS = false
			}
		} else {
			if bS {
				hS.end = i
				hourSlots = append(hourSlots, hS)
				bS = false
			}

		}
	}

	// chevauchement 00h00  -->  concaténation de 2 slots
	if len(hourSlots) > 1 && (hourSlots[0].begin == 0 && hourSlots[len(hourSlots)-1].end == 0) {
		hourSlots[0].begin = hourSlots[len(hourSlots)-1].begin
		//effacer dernier slot
		hourSlots = hourSlots[:len(hourSlots)-1]
	}

	return slotsToString(hourSlots)
}
