package main

import (
	"fmt"
	"math/rand"
)

type cartaDaGioco struct {
	valore string
	seme   string
}

type mazzo struct {
	carte []cartaDaGioco
	n     int
}

func creaCarta(val string, s string) cartaDaGioco {
	return cartaDaGioco{valore: val, seme: s}
}

func creaMazzo() mazzo {
	var maz mazzo
	valori := [10]string{"asso", "due", "tre", "quattro", "cinque", "sei", "sette", "fante", "donna", "re"}
	semi := [4]string{"quadri", "picche", "cuori", "fiori"}

	for _, v := range valori {
		for _, s := range semi {
			maz.carte = append(maz.carte, creaCarta(v, s))
			maz.n++
		}
	}
	return maz
}

func mischia(mazzo *mazzo) {
	for i := range mazzo.carte {
		j := rand.Intn(i + 1)
		mazzo.carte[i], mazzo.carte[j] = mazzo.carte[j], mazzo.carte[i]
	}
}

func preleva(mazzo mazzo) (cartaDaGioco, mazzo, error) {
	if mazzo.n <= 0 {
		return cartaDaGioco{}, mazzo, fmt.Errorf("mazzo vuoto")
	}
	carta := mazzo.carte[mazzo.n-1]
	mazzo.carte = mazzo.carte[:mazzo.n-1]
	mazzo.n--
	return carta, mazzo, nil
}

func inizio() (bool, mazzo) {
	var gioca string
	mazzo := creaMazzo()
	mischia(&mazzo)
	fmt.Print("\n*+*+* Munizza's CASINO *+*+*\nVuoi giocare a BlackJack? (s/n)\n>>")
	fmt.Scan(&gioca)
	if gioca == "s" {
		return true, mazzo
	} else if gioca == "n" {
		return false, mazzo
	} else {
		fmt.Println("Input non valido, riprova")
		return inizio()
	}
}

func giocata(maz *mazzo, saldo *int) {
	var carteGiocatore []cartaDaGioco

	// Distribuisci due carte iniziali
	for i := 0; i < 2; i++ {
		carta, nuovoMazzo, err := preleva(*maz)
		if err != nil {
			fmt.Println("Errore nel prelevare le carte:", err)
			return
		}
		carteGiocatore = append(carteGiocatore, carta)
		*maz = nuovoMazzo
	}

	for {
		// Mostra le carte e il punteggio attuale
		fmt.Printf("\nLe tue carte:\n")
		stampaCarteGiocatore(carteGiocatore)
		punteggioGiocatore := calcolaPunteggio(carteGiocatore)
		fmt.Printf("\nPunteggio attuale: %d\n", punteggioGiocatore)

		// Controlla se il giocatore ha superato 21
		if punteggioGiocatore > 21 {
			fmt.Printf("\nHAI SBALLATO, HAI PERSO...\n")
			*saldo--
			return
		}

		// Chiedi al giocatore se vuole un'altra carta
		var risposta string
		fmt.Printf("\nVuoi un'altra carta? (s/n): ")
		fmt.Scan(&risposta)

		if risposta == "s" {
			carta, nuovoMazzo, err := preleva(*maz)
			if err != nil {
				fmt.Println("Errore nel prelevare la carta:", err)
				return
			}
			carteGiocatore = append(carteGiocatore, carta)
			*maz = nuovoMazzo
		} else if risposta == "n" {
			fmt.Printf("\nTi sei fermato.\n")
			break
		} else {
			fmt.Println("Input non valido, per favore inserisci 's' o 'n'.")
		}
	}

	// Fine turno: mostra il punteggio finale
	fmt.Printf("\nTurno terminato.")
	fmt.Printf("\nLe tue carte finali:\n")
	stampaCarteGiocatore(carteGiocatore)
	punteggioGiocatore := calcolaPunteggio(carteGiocatore)
	fmt.Printf("\nPunteggio finale: %d\n", punteggioGiocatore)

	// Giocata del banco
	fmt.Printf("\n---GIOCATA BANCO---\n\n")
	carteBanco := giocataBanco(*maz, punteggioGiocatore)
	punteggioBanco := calcolaPunteggio(carteBanco)
	stampaCarteGiocatore(carteBanco)
	fmt.Printf("\n\nPunteggio banco: %d\n", punteggioBanco)

	// Determina il risultato
	fmt.Print("\nVerdetto:\n\n")
	if punteggioBanco > 21 || punteggioGiocatore > punteggioBanco {
		fmt.Println("HAI VINTO!!!")
		*saldo++
	} else if punteggioGiocatore == punteggioBanco {
		fmt.Println("PAREGGIO")
	} else {
		fmt.Println("HAI PERSO...")
		*saldo--
	}
}

func giocataBanco(mazzo mazzo, puntiGiocatore int) []cartaDaGioco {
	var carteBanco []cartaDaGioco
	for i := 0; i < 2; i++ {
		carta, nuovoMazzo, err := preleva(mazzo)
		if err != nil {
			fmt.Println("Errore nel prelevare le carte:", err)
			return carteBanco
		}
		carteBanco = append(carteBanco, carta)
		mazzo = nuovoMazzo
	}

	// Logica del banco: deve fermarsi se il punteggio è almeno 17
	for punteggioBanco := calcolaPunteggio(carteBanco); punteggioBanco < 17; punteggioBanco = calcolaPunteggio(carteBanco) {
		carta, nuovoMazzo, err := preleva(mazzo)
		if err != nil {
			fmt.Println("Errore nel prelevare le carte:", err)
			return carteBanco
		}
		carteBanco = append(carteBanco, carta)
		mazzo = nuovoMazzo
	}
	return carteBanco
}

func calcolaPunteggio(carte []cartaDaGioco) int {
	valori := map[string]int{
		"asso":    11,
		"due":     2,
		"tre":     3,
		"quattro": 4,
		"cinque":  5,
		"sei":     6,
		"sette":   7,
		"fante":   10,
		"donna":   10,
		"re":      10,
	}

	punteggio := 0
	assoCount := 0

	for _, carta := range carte {
		punteggio += valori[carta.valore]
		if carta.valore == "asso" {
			assoCount++
		}
	}

	for punteggio > 21 && assoCount > 0 {
		punteggio -= 10
		assoCount--
	}

	return punteggio
}

func stampaCarteGiocatore(carte []cartaDaGioco) {
	for _, v := range carte {
		fmt.Printf("%v\t", v)
	}
	fmt.Println()
}

func main() {
	gioca := "s"
	var cash int = 10
	inizia, mazzo := inizio()
	if !inizia {
		return
	}

	for cash > 0 && gioca == "s" {
		fmt.Printf("--------------------\nSALDO disponibile: %d$\n--------------------\n", cash)
		giocata(&mazzo, &cash)
		for {
			fmt.Printf("Vuoi giocare ancora? (s/n): ")
			fmt.Scan(&gioca)
			if gioca == "s" || gioca == "n" {
				break
			}
			fmt.Println("Input non valido, riprova.")
		}
	}

	fmt.Printf("\nGrazie per aver giocato. SALDO: %d$", cash)
}
