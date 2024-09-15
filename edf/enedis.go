package edf

import (
	"errors"
	"os"

	"strings"
	"time"
)

type HCreuses struct {
	hstart string
	hend   string
}

var (
	enedisStart time.Time
	enedisEnd   time.Time
	lines       []string
)

/*
Recherche si hhmm(minutes) est dans un créneau d'heures creuses
*/
func isHeureCreuse(hmin int, timeSlots []HCreuses) bool {
	for _, slot := range timeSlots {
		hmstart := hmTominute(slot.hstart)
		hmend := hmTominute(slot.hend)
		if hmend > hmstart {
			if hmin > hmstart && hmin <= hmend {
				return true
			}
		} else {
			if hmin <= hmend || hmin > hmstart {
				return true
			}
		}
	}
	return false
}

/*
Création slice d'HeuresCreuses
hh:mm-hh:mm   				-->   [{hh:mm hh:mm}]
hh:mm-hh:mm,hh:mm-hh:mm   	-->   [{hh:mm hh:mm} {hh:mm hh:mm}]
*/
func buildHeuresCreuses(entry string) (slotHC []HCreuses, err error) {
	str := strings.Split(strings.ReplaceAll(entry, " ", ""), ",")

	for _, hphc := range str {
		h := strings.Split(hphc, "-")
		if len(h) != 2 {
			return slotHC, errors.New("format créneau incorrect")
		}
		if checkHM(h[0]) != nil || checkHM(h[1]) != nil {
			return slotHC, errors.New("format heure:minute incorrect")
		}
		slotHC = append(slotHC, HCreuses{h[0], h[1]})
	}

	if overlapSlots(entry) {
		return slotHC, errors.New("recouvrement de créneaux horaires")
	}

	return slotHC, nil
}

/*
Verification du fichier Enedis.
Ce fichier doit de type horaire
et mémoriser au moins un an
*/

func checkEnedisFile(filename string) (err error) {

	dataByte, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("fichier inconnu")
	}
	//  découpage en lignes
	lines = strings.Split(string(dataByte), "\n")
	//  Contrôle si le fichier est du type horaire
	fl := strings.Split(string(lines[0]), ";")
	if fl[len(fl)-1] != "Pas en minutes" {
		return errors.New("fichier enedis non compatible")
	}
	//  Contrôle si le fichier a au moins une durée d'un an
	enedisStart, err = time.Parse(DMYLAYOUT, strings.Split(string(lines[1]), ";")[2])
	if err != nil {
		return err
	}
	enedisEnd, err = time.Parse(DMYLAYOUT, strings.Split(string(lines[1]), ";")[3])
	if err != nil {
		return err
	}
	enedisDuration := enedisEnd.Sub(enedisStart)
	if enedisDuration.Hours()/24 < 365 {
		return errors.New("durée fichier enedis inférieure à 1 an")
	}
	return nil
}

/*
 */
func readEnedisFile(lines []string, endAnalyze string) (data []string, err error) {
	//  Contrôle si un an de données sont disponibles avant la date choisie
	end, err := time.Parse(DMYLAYOUT, endAnalyze)
	if err != nil {
		return nil, err
	}
	if end.After(enedisEnd) {
		return nil, errors.New("date ultérieure au fichier")
	}
	analyzeDuration := end.Sub(enedisStart)
	if analyzeDuration.Hours()/24 < 365 {
		return nil, errors.New("créneau d'analyse inférieur à 1 an")
	}
	//  Recherche les lignes de début de fin du créneau annuel à analyser
	//  Retourne les datas du créneau annuel
	oneYear := end.AddDate(-1, 0, 0).Format(YMDLAYOUT)

	for s := range lines {
		if lines[s][0:10] == oneYear {
			for e := s + 1; e < len(lines); e++ {
				if lines[e][0:10] == end.Format(YMDLAYOUT) {
					//  retourne datas d'une durée d'un an
					return lines[s+1 : e+1], nil
				}
			}
			return nil, errors.New("date de fin d'analyse non trouvée")
		}
	}
	return nil, errors.New("date de fin d'analyse -1 an non trouvée")
}
