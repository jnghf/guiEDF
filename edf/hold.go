package edf

import (
	//"log"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type HoldableButton struct {
	widget.Button
	On bool
}

var (
	downMemo = false
	onMemo   = false
)

/*
Bouton avec contrôle souris
Down - Up - Event
*/
func NewHoldableButton(label string) *HoldableButton {
	button := &HoldableButton{}
	button.ExtendBaseWidget(button)
	button.Text = label
	button.On = false
	return button
}

/*
Clic souris sur le bouton
*/
func (h *HoldableButton) MouseDown(*desktop.MouseEvent) {
	//log.Println("down ", h.On, h.Text)
	if h.On {
		h.FocusLost()
		h.On = false
		onMemo = false
	} else {
		h.FocusGained()
		h.On = true
		onMemo = true
	}
	downMemo = true
}

/*
Relachement souris sur le bouton
*/
func (h *HoldableButton) MouseUp(*desktop.MouseEvent) {
	//log.Println("up ", h.On)
	downMemo = false
}

/*
Mouvement souris sur le bouton
*/
func (h *HoldableButton) MouseIn(*desktop.MouseEvent) {
	//log.Println("event ", downMemo, onMemo, h.Text)
	if downMemo {
		if onMemo {
			h.FocusGained()
			h.On = true
		} else {
			h.FocusLost()
			h.On = false
		}
	}
}
