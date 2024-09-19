package edf

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const (
	DMYLAYOUT = "02/01/2006"
	YMDLAYOUT = "2006-01-02"
)

/*
réduction de la taille d'une chaîne à max caractères
*/
func truncate(s string, max int) string {
	if len(s) > max {
		return s[len(s)-max:]
	} else {
		return s
	}
}

/*
Conversion Euro string  -->  float
*/
func stringEuroToFloat(euro string) float64 {
	//  élimination \n \t €
	str := strings.Trim(euro, "\n\t\u00a0€")
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

/*
Conversion hh:mn --> entier
*/
func hmTominute(hstring string) int {
	//   h * 60 + min
	//  "10:32"  --> 10*60 + 32 = 632
	hsplit := strings.Split(hstring, ":")
	h, _ := strconv.Atoi(hsplit[0])
	min, _ := strconv.Atoi(hsplit[1])
	return h*60 + min
}

/*
Heure (hhmn) entre hStart et hEnd ?
*/
func inSlot(hStart, hEnd, h string) bool {
	hmstart := hmTominute(hStart)
	hmend := hmTominute(hEnd)
	hm := hmTominute(h)
	if hmend > hmstart {
		if hm > hmstart && hm <= hmend {
			return true
		}
	} else {
		if hm <= hmend || hm > hmstart {
			return true
		}
	}
	return false
}

/*
Chevauchement de 2 slots ?
*/
func overlap(h1, h2 string) bool {
	ha := strings.Split(h1, "-")
	hb := strings.Split(h2, "-")
	if inSlot(ha[0], ha[1], hb[0]) {
		return true
	}
	if inSlot(ha[0], ha[1], hb[1]) {
		return true
	}
	return false
}

/*
chevauchement de n slots ?
*/
func overlapSlots(entry string) bool {
	hs := strings.Split(entry, ",")
	l := len(hs)
	for i := 0; i < l-1; i++ {
		for j := i + 1; j < l; j++ {
			if overlap(hs[i], hs[j]) {
				return true
			}
		}
	}
	return false
}

/*
Vérification du format de l'heure (hh:mn)
*/
func checkHM(sHm string) (err error) {
	err1 := errors.New("err_checkHM")
	hm := strings.Split(sHm, ":")
	if len(hm) != 2 {
		return err1
	} else {
		h, err := strconv.Atoi(hm[0])
		if err != nil || h < 0 || h > 23 {
			return err1
		}
		m, err := strconv.Atoi(hm[1])
		if err != nil || m < 0 || m > 59 {
			return err1
		}
	}
	return nil
}
