// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oldham123/go-dex/internal/storage/csv"
)

func main() {
	log.Printf("Starting PokeDex API server...")

	// Initialize store
	store := csv.NewStore()

	// Load data - using relative path from where the binary runs
	if err := store.LoadFromCSV("data/pokemon.csv"); err != nil {
		log.Fatalf("Failed to load Pokemon data: %v", err)
	}

	// Demo: Get total count
	allPokemon := store.List()
	log.Printf("Successfully loaded %d Pokemon", len(allPokemon))

	// Demo: Retrieve specific Pokemon
	if pokemon, err := store.GetByID(1); err == nil {
		log.Printf("Found Pokemon #1: %s (Height: %d, Weight: %d)",
			pokemon.Identifier,
			pokemon.Height,
			pokemon.Weight,
		)
	}

	// Demo: Print first 5 Pokemon
	fmt.Println("\nFirst 5 Pokemon loaded:")
	for i, p := range allPokemon {
		if i >= 5 {
			break
		}
		fmt.Printf("#%d: %s (Species ID: %d)\n",
			p.Id,
			p.Identifier,
			p.SpeciesId,
		)
	}

	// Set up HTTP server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("\nServer starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
