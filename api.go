package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// PokemonSprites represents the sprite URLs
type PokemonSprites struct {
	FrontDefault string `json:"front_default"`
}

// PokemonSpecies represents the species information
type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonStatInfo represents individual stat details
type PokemonStatInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonStat represents a stat entry
type PokemonStat struct {
	BaseStat int             `json:"base_stat"`
	Effort   int             `json:"effort"`
	Stat     PokemonStatInfo `json:"stat"`
}

// PokemonTypeInfo represents individual type details
type PokemonTypeInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonType represents a type entry
type PokemonType struct {
	Slot int             `json:"slot"`
	Type PokemonTypeInfo `json:"type"`
}

// PokemonAbilityInfo represents individual ability details
type PokemonAbilityInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonAbility represents an ability entry
type PokemonAbility struct {
	Ability  PokemonAbilityInfo `json:"ability"`
	IsHidden bool               `json:"is_hidden"`
	Slot     int                `json:"slot"`
}

// Pokemon represents the full structure for a Pokémon
type Pokemon struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Height    int              `json:"height"`
	Weight    int              `json:"weight"`
	Sprites   PokemonSprites   `json:"sprites"`
	Species   PokemonSpecies   `json:"species"`
	Stats     []PokemonStat    `json:"stats"`
	Types     []PokemonType    `json:"types"`
	Abilities []PokemonAbility `json:"abilities"`
}

// ParsePokemonData unmarshals JSON data into a Pokemon struct
func ParsePokemonData(jsonData string) (Pokemon, error) {
	var pokemon Pokemon
	err := json.Unmarshal([]byte(jsonData), &pokemon)
	if err != nil {
		// Keep identifier in the error message if possible, though it's not passed here.
		// This function is usually called by GetPokemon which has the identifier.
		return pokemon, fmt.Errorf("error parsing pokemon JSON data: %w", err)
	}
	return pokemon, nil
}

// GetPokemon fetches Pokémon data from the PokéAPI by name or ID.
func GetPokemon(identifier string) (Pokemon, error) {
	var pokemon Pokemon
	normalizedIdentifier := strings.ToLower(identifier)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", normalizedIdentifier)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return pokemon, fmt.Errorf("network error fetching pokemon '%s': %w", identifier, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try to read body for more info from API, but don't let it mask the status code error
		// bodyBytes, _ := ioutil.ReadAll(resp.Body)
		// return pokemon, fmt.Errorf("API error for pokemon '%s' (status: %s): %s", identifier, resp.Status, string(bodyBytes))
		return pokemon, fmt.Errorf("pokemon '%s' not found (API status: %s)", identifier, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pokemon, fmt.Errorf("error reading response body for pokemon '%s': %w", identifier, err)
	}

	pokemon, err = ParsePokemonData(string(body))
	if err != nil {
		// The ParsePokemonData error is already quite generic; GetPokemon adds context.
		return pokemon, fmt.Errorf("error parsing pokemon data for '%s': %w", identifier, err)
	}

	return pokemon, nil
}
