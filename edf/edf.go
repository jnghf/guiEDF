package edf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	FILEFARE = "./Work_Files/fare.json"
)

type Conso struct {
	Dte   string
	Hm    int
	Wh    int
	Color ColorTempo
}
type Consos []Conso

type Jour struct {
	Col      int
	Hp       int
	Hc       int
	HtempoHp [3]int
	HtempoHc [3]int
}
type JoursConso map[string]Jour

type Hconsos struct {
	Hbase    int
	Hp       int
	Hc       int
	HtempoHp [3]int
	HtempoHc [3]int
}

type EdfFares struct {
	BaseFare  map[string]BasePrice
	HchpFare  map[string]HcHpPrice
	TempoFare map[string]TempoPrice
	DateFare  time.Time
}

var (
	edfTarifs EdfFares

	data []string

	powerSelection string
)

/*
Trouver les prix en fonction de la puissance souscrite
*/
func findPriceByPower(power string) (prices Price) {
	prices.Base = edfTarifs.BaseFare[power]
	prices.HcHp = edfTarifs.HchpFare[power]
	prices.Tempo = edfTarifs.TempoFare[power]
	return prices
}

/*
Retourne mat de consommations journalières
*/
func consosJour(consos Consos, creneauxhc []HCreuses) JoursConso {
	i := 0
	jconso := make(JoursConso)
	dte := consos[0].Dte

	for i < len(consos) {
		var jc Jour

		for dte == consos[i].Dte {

			if isHeureCreuse(consos[i].Hm, creneauxhc) {
				jc.Hc += consos[i].Wh
				jc.HtempoHc[consos[i].Color] += consos[i].Wh
			} else {
				jc.Hp += consos[i].Wh
				jc.HtempoHp[consos[i].Color] += consos[i].Wh
			}
			i++
			if i >= len(consos) {
				return jconso
			}
		}
		if isHeureCreuse(consos[i].Hm, creneauxhc) {
			jc.Hc += consos[i].Wh
			jc.HtempoHc[consos[i].Color] += consos[i].Wh
		} else {
			jc.Hp += consos[i].Wh
			jc.HtempoHp[consos[i].Color] += consos[i].Wh
		}
		jconso[dte] = jc
		dte = consos[i].Dte
		i++
	}
	return jconso
}

/*
Retourne tableau des consommations par slot horaire (0.5h)
*/
func lireConsos(data []string, datesTempo TdTempo) Consos {
	var consos Consos
	var color ColorTempo

	firstTimeSlot := true
	var timeSlotStart time.Time
	timeSlot := 0.5
	for _, line := range data {
		dte, _ := time.Parse(time.RFC3339, line[0:25])
		if firstTimeSlot {
			firstTimeSlot = false
			timeSlotStart = dte
		} else {
			timeSlot = dte.Sub(timeSlotStart).Hours()
		}
		Wh, _ := strconv.Atoi(line[26:])
		htm := hmTominute(line[11:16])

		//  Heures Tempo de 06H00 à 06H00
		if htm > (6 * 60) {
			// Après 06h00 :  Couleur du jou
			df := dte.Format("2006-01-02")
			color = datesTempo[df]
		} else {
			// Avant 06h00 :  Couleur du jour précédent
			dyest := dte.AddDate(0, 0, -1)
			df := dyest.Format("2006-01-02")
			color = datesTempo[df]
		}

		cons := Conso{line[0:10], htm, int(float64(Wh) * timeSlot), color}
		consos = append(consos, cons)
		timeSlotStart = dte
	}
	return consos
}

/*
Sauvegarde des taris EDF
*/
func saveFare() {
	file, _ := json.MarshalIndent(edfTarifs, "", "")
	_ = os.WriteFile(FILEFARE, file, 0644)
}

/*
Chargement des tarifs EDF
*/
func loadFare() (EdfTarifs EdfFares, err error) {
	file, err := os.ReadFile(FILEFARE)
	if err != nil {
		return EdfTarifs, err
	}
	json.Unmarshal(file, &EdfTarifs)

	return EdfTarifs, nil
}

/*
Lecture tarifs de Base EDF
depuis données Web EDF
*/
func fetchBase(doc *goquery.Document) {

	tb := doc.Find(".table--responsive :contains('par EDF en option base')").Find("tbody").First()
	e := BasePrice{} ////
	c := make(map[string]BasePrice)

	tb.Find("tr").Each(func(iy int, tr *goquery.Selection) {
		hh := tr.Find("th strong")
		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
			case 0:
				e.Subscription = stringEuroToFloat(td.Text())
			case 1:
				e.Base = stringEuroToFloat(td.Text())
			}
		})

		c[hh.Text()] = e
		edfTarifs.BaseFare = c
	})
}

/*
Lecture tarifs Heures Pleines/Heures Creuses EDF
*/
func fetchHcHp(doc *goquery.Document) {

	tb := doc.Find(".table--responsive :contains('par EDF en option HP/HC')").Find("tbody").First()
	e := HcHpPrice{}
	c := make(map[string]HcHpPrice)

	tb.Find("tr").Each(func(iy int, tr *goquery.Selection) {
		hh := tr.Find("th strong")
		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
			case 0:
				e.Subscription = stringEuroToFloat(td.Text())
			case 1:
				e.Hc = stringEuroToFloat(td.Text())
			case 2:
				e.Hp = stringEuroToFloat(td.Text())
			}
		})

		c[hh.Text()] = e
		edfTarifs.HchpFare = c
	})
}

/*
Lecture tarifs Tempo EDF
*/
func fetchTempo(doc *goquery.Document) {

	tb := doc.Find(".table--responsive :contains('Grille tarifaire de l\\'offre Tempo par EDF')").Find("tbody").First()
	e := TempoPrice{}
	c := make(map[string]TempoPrice)

	tb.Find("tr").Each(func(iy int, tr *goquery.Selection) {
		hh := tr.Find("th strong")
		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
			case 0:
				e.Subscription = stringEuroToFloat(td.Text())
			case 1:
				e.BlueHC = stringEuroToFloat(td.Text())
			case 2:
				e.BlueHP = stringEuroToFloat(td.Text())
			case 3:
				e.WhiteHC = stringEuroToFloat(td.Text())
			case 4:
				e.WhiteHP = stringEuroToFloat(td.Text())
			case 5:
				e.RedHC = stringEuroToFloat(td.Text())
			case 6:
				e.RedHP = stringEuroToFloat(td.Text())
			}
		})

		c[hh.Text()] = e
		edfTarifs.TempoFare = c
	})
}

/*
Lecture Body document Web
*/
func fetchWeb(webPage string) (doc *goquery.Document, err error) {
	resp, err := http.Get(webPage)
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return doc, fmt.Errorf("taris edf inaccessibles:  %s", resp.Status)
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, nil
}

/*
Lecture des tarifs EDF depuis site Web
*/
func fetchFare() (err error) {
	doc, err := fetchWeb(EDF_WEB)
	if err != nil {
		return err
	}
	fetchBase(doc)
	fetchHcHp(doc)
	fetchTempo(doc)
	edfTarifs.DateFare = time.Now()
	saveFare()

	return nil
}
