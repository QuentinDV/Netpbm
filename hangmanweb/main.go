package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hangman", hangman)
	http.HandleFunc("/hangman/choose_difficulty", ChooseDifficulty)
	http.HandleFunc("/hangman/easy", Easy)
	http.HandleFunc("/hangman/medium", Medium)
	http.HandleFunc("/hangman/hard", Hard)
	http.HandleFunc("/hangman/defeated", Defeated)
	http.HandleFunc("/hangman/victory", Victory)

	// Ajoutez une nouvelle route pour gérer le formulaire
	http.HandleFunc("/process_form", ProcessForm)

	// Servir les fichiers statiques dans le répertoire "assets"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html/"))))

	fmt.Println("\n(http://localhost:8080/hangman) - Server started on port")
	fmt.Println("\n(http://localhost:8080/hangman/choose_difficulty) - Server started on port")
	http.ListenAndServe(":8080", nil)
}

func hangman(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/hangman.html")
}

func ChooseDifficulty(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/choose_difficulty.html")
}

func Easy(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/easy.html")
}

func Medium(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/medium.html")
}

func Hard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/hard.html")
}

func Defeated(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/defeated.html")
}

func Victory(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/victory.html")
}

func ProcessForm(w http.ResponseWriter, r *http.Request) {
	// Analyser les données du formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	pseudo := r.Form.Get("pseudo")

	// Vous pouvez faire ce que vous voulez avec les données du formulaire ici
	// Par exemple, imprimer le pseudo dans la console
	fmt.Println("Name:", pseudo)

	// Rediriger ou renvoyer une réponse au client
	// (vous pouvez personnaliser cela en fonction de votre logique d'application)
	http.Redirect(w, r, "/hangman/choose_difficulty", http.StatusSeeOther)
}
