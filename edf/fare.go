package edf

import (
	"fmt"
)

type BasePrice struct {
	Subscription float64
	Base         float64
}

type HcHpPrice struct {
	Subscription float64
	Hc           float64
	Hp           float64
}

type TempoPrice struct {
	Subscription float64
	BlueHC       float64
	BlueHP       float64
	WhiteHC      float64
	WhiteHP      float64
	RedHC        float64
	RedHP        float64
}

type Price struct {
	Base  BasePrice
	HcHp  HcHpPrice
	Tempo TempoPrice
}

/*
//Affichage des consommations totales
*/
func conso_totale(data []string, datesTempo TdTempo, creneauxhc []HCreuses) (hc Hconsos) {
	consos := lireConsos(data, datesTempo)
	jconso := consosJour(consos, creneauxhc)

	for _, jr := range jconso {
		hc.Hp += jr.Hp
		hc.Hc += jr.Hc
		hc.Hbase += (jr.Hc + jr.Hp)
		// Heures Tempo
		for i := 0; i < 3; i++ {
			hc.HtempoHc[i] += jr.HtempoHc[i]
			hc.HtempoHp[i] += jr.HtempoHp[i]
		}

	}

	return hc
}

/*
Préparation des résultats pour affichage
Heure de Base
Heure pleines / Heures Creuses
Tempo
*/
func displayYear(hc Hconsos, prices Price) (result string) {
	//  Affichage Tarif de Base
	cb := float64(hc.Hbase/1000) * (prices.Base.Base)
	result = fmt.Sprintf("\nBase: \t\t\t%2d kW  \tConso: \t\t%6.0f €\n", hc.Hbase/1000, cb)
	result += fmt.Sprintf("\t\t\t\t\tAbonnement: \t%6.0f €\n", prices.Base.Subscription)
	result += fmt.Sprintf("\t\t\t\t\tTotal: \t\t\t%6.0f €\n\n", cb+prices.Base.Subscription)
	//  Affichage Heures Pleines/Heures Creuses
	chp := float64(hc.Hp/1000) * (prices.HcHp.Hp)
	result += fmt.Sprintf("H. Pleines: \t\t%2d kW  \tConso: \t\t%6.0f €\n", hc.Hp/1000, chp)
	chc := float64(hc.Hc/1000) * (prices.HcHp.Hc)
	result += fmt.Sprintf("H. Creuses: \t\t%2d kW  \tConso: \t\t%6.0f €\n", hc.Hc/1000, chc)
	result += fmt.Sprintf("\t\t\t\t\tAbonnement: \t%6.0f €\n", prices.HcHp.Subscription)
	result += fmt.Sprintf("\t\t\t\t\tTotal: \t\t\t%6.0f €\n\n", chp+chc+prices.HcHp.Subscription)
	//  Affichage Tarif Tempo
	cbhp := float64(hc.HtempoHp[BLEU]/1000) * (prices.Tempo.BlueHP)
	result += fmt.Sprintf("T. Bleu HP: \t\t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHp[BLEU]/1000, cbhp)
	cbhc := float64(hc.HtempoHc[BLEU]/1000) * (prices.Tempo.BlueHC)
	result += fmt.Sprintf("T. Bleu HC: \t\t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHc[BLEU]/1000, cbhc)
	cwhp := float64(hc.HtempoHp[BLANC]/1000) * (prices.Tempo.WhiteHP)
	result += fmt.Sprintf("T. Blanc HP: \t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHp[BLANC]/1000, cwhp)
	cwhc := float64(hc.HtempoHc[BLANC]/1000) * (prices.Tempo.WhiteHC)
	result += fmt.Sprintf("T. Blanc HC: \t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHc[BLANC]/1000, cwhc)
	crhp := float64(hc.HtempoHp[ROUGE]/1000) * (prices.Tempo.RedHP)
	result += fmt.Sprintf("T. Rouge HP: \t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHp[ROUGE]/1000, crhp)
	crhc := float64(hc.HtempoHc[ROUGE]/1000) * (prices.Tempo.RedHC)
	result += fmt.Sprintf("T. Rouge HC: \t%2d kW  \tConso: \t\t%6.0f €\n", hc.HtempoHc[ROUGE]/1000, crhc)
	result += fmt.Sprintf("\t\t\t\t\tAbonnement: \t%6.0f €\n", prices.Tempo.Subscription)
	result += fmt.Sprintf("\t\t\t\t\tTotal: \t\t\t%6.0f €\n\n", cbhp+cbhc+cwhp+cwhc+crhp+crhc+prices.Tempo.Subscription)

	return result
}
