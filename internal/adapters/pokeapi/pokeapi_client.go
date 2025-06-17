package pokeapi

// Package pokeapi provides an adapter to fetch Pokemon data from the public PokeAPI.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"pokedex-app/internal/domain" // Assuming this is the correct module path
)

// PokeAPIClient implements the domain.PokemonRepository interface using the public PokeAPI.
type PokeAPIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewPokeAPIClient creates a new PokeAPIClient.
func NewPokeAPIClient(baseURL string, timeoutSeconds int) *PokeAPIClient {
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2" // Default base URL
	}
	return &PokeAPIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// parsePokemonData unmarshals JSON data into a domain.Pokemon struct
// It's kept as an unexported helper function within the adapter.
func parsePokemonData(jsonData []byte) (domain.Pokemon, error) {
	var pokemon domain.Pokemon
	err := json.Unmarshal(jsonData, &pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("error parsing pokemon JSON data: %w", err)
	}
	return pokemon, nil
}

// GetByIdentifier fetches Pokémon data from the PokéAPI by name or ID string.
// This is a helper that can be used by GetByName and GetByID.
func (c *PokeAPIClient) GetByIdentifier(identifier string) (domain.Pokemon, error) {
	var pokemon domain.Pokemon
	normalizedIdentifier := strings.ToLower(identifier)

	url := fmt.Sprintf("%s/pokemon/%s", c.BaseURL, normalizedIdentifier)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return pokemon, fmt.Errorf("network error fetching pokemon '%s': %w", identifier, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pokemon, fmt.Errorf("pokemon '%s' not found (API status: %s)", identifier, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pokemon, fmt.Errorf("error reading response body for pokemon '%s': %w", identifier, err)
	}

	pokemon, err = parsePokemonData(body) // Use the local parsePokemonData
	if err != nil {
		return pokemon, fmt.Errorf("error parsing pokemon data for '%s': %w", identifier, err)
	}

	return pokemon, nil
}

// GetByName implements the domain.PokemonRepository interface.
func (c *PokeAPIClient) GetByName(name string) (domain.Pokemon, error) {
	return c.GetByIdentifier(name)
}

// GetByID implements the domain.PokemonRepository interface.
func (c *PokeAPIClient) GetByID(id int) (domain.Pokemon, error) {
	return c.GetByIdentifier(fmt.Sprintf("%d", id))
}
