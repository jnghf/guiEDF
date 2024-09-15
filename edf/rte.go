package edf

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	EDF_WEB   = "https://www.jechange.fr/energie/edf/tarifs#tarifs-reglementes"
	FILETEMPO = "./Work_Files/tempo.json"
)

const (
	BLEU ColorTempo = iota
	BLANC
	ROUGE
)

type ColorTempo int

type TdTempo map[string]ColorTempo

var datesTempo TdTempo

/*
Lecture Jour Tempo sur le site EDF
*/
func curlTempo(stempo string) (res map[string]interface{}, err error) {
	const url = "https://www.services-rte.com/cms/open_data/v1/tempo?season="

	var dt map[string](map[string]interface{})

	curl := exec.Command("curl", url+stempo)
	out, err := curl.Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(out), &dt)
	if err != nil {
		return nil, err
	}

	return dt["values"], nil
}

/*
Sauvegarde map tempo dans fichier Json
*/
func saveDateTempo(datestempo TdTempo) {
	file, _ := json.MarshalIndent(datestempo, "", "")
	_ = os.WriteFile(FILETEMPO, file, 0644)
}

/*
Lecteur fichier Json tempo --> map
*/
func loadTempo() (datestempo TdTempo, err error) {
	file, err := os.ReadFile(FILETEMPO)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(file, &datestempo)
	return datestempo, nil
}

/*
Trouver année du fichier RTE (yyyy-yyyy+1)
du 01/09/yyyy au 01/09/yyyy+1
*/
func findYearRte(fromTempo time.Time) string {
	year := fromTempo.Year()
	sept := time.Date(year, 9, 1, 0, 0, 0, 0, time.UTC) /// UTC ?
	if fromTempo.Before(sept) {
		return strconv.Itoa(year-1) + "-" + strconv.Itoa(year)
	} else {
		return strconv.Itoa(year) + "-" + strconv.Itoa(year+1)
	}
}

/*
Mise à jour fichier tempo de fromtempo à totempo
Création map si le fichier n'existe pas
Exec CURL site EDF sur les dates tempo non connues seulement
*/
func updateTempo(from, to string) (datesTempo TdTempo, err error) {
	var icol ColorTempo
	var nbJour float64 = 0

	fromTempo, _ := time.Parse(DMYLAYOUT, from)
	toTempo, _ := time.Parse(DMYLAYOUT, to)

	datesTempo, err = loadTempo()
	if err != nil {
		//  Nouvelle map
		datesTempo = make(TdTempo)
	}

	for fromTempo.Before(toTempo) {
		sTempo := fromTempo.Format(YMDLAYOUT)
		_, ok := datesTempo[sTempo] //  Lecture jour tempo
		if ok {                     //  Jour déjà connu
			fromTempo = fromTempo.Add(time.Hour * 24) //  Jour suivant
			nbJour++
		} else {

			sTempo = findYearRte(fromTempo)
			dTempo, err := curlTempo(sTempo)
			if err != nil {
				return datesTempo, err
			}
			for dt, color := range dTempo {
				switch color {
				case "BLUE":
					icol = BLEU
				case "WHITE":
					icol = BLANC
				case "RED":
					icol = ROUGE
				default:
					continue
				}
				datesTempo[dt] = icol

				nbJour++

				tempoProgress.SetValue(nbJour / 365)
				time.Sleep(time.Millisecond * 15)
			}
			fromTempo = fromTempo.Add(time.Hour * 24)
		}
	}
	//  Sauvegarde fichier Tempo
	saveDateTempo(datesTempo)

	return datesTempo, nil
}
