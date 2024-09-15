package main

import (
	"fyne.io/fyne/v2/app"

	"EDF/guiEDF/edf"
)

/*
CALCUL DES BUDGETS EDF SUR UNE ANNEE DE CONSOMMATION
*/
func main() {

	myApp := app.NewWithID("123")
	//  Création fenêtre de l'application
	edf.NewWinEDF(myApp)

}
