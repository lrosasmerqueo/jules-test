package main

import (
	"log"

	"fyne.io/fyne/v2/app" // Import for app.New()

	"pokedex-app/internal/adapters/cache"
	"pokedex-app/internal/adapters/gui/fyneui"
	"pokedex-app/internal/adapters/pokeapi"
	"pokedex-app/internal/application"
)

func main() {
	// Initialize Fyne application
	fyneApp := app.New()
	// It's important that fyneApp.SetUniqueID() is called here if needed,
	// but for this app, it might not be strictly necessary unless using preferences.

	// 1. Initialize Adapters
	// PokeAPI client - primary data source
	pokeapiClient := pokeapi.NewPokeAPIClient("https://pokeapi.co/api/v2", 10) // 10 seconds timeout

	// In-memory cache - fallback data source
	inMemoryCache := cache.NewInMemoryCache()

	// 2. Initialize Application Service
	// PokemonService will use pokeapiClient first, then inMemoryCache if pokeapiClient fails
	pokemonAppService := application.NewPokemonService(pokeapiClient, inMemoryCache)

	// 3. Initialize UI Adapter
	// The FyneUI needs the application service to interact with
	gui := fyneui.NewFyneUI(pokemonAppService) // Pass the app service

	// 4. Start the application
	log.Println("Starting Pokedex application...")
	gui.Start() // This will call mainWindow.ShowAndRun()
}
