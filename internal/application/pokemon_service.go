package application

// Package application contains the application services (use cases).

import (
	"fmt"
	"strconv"
	"strings"

	"pokedex-app/internal/domain"
)

// PokemonService orchestrates fetching of Pokemon data.
// It uses one primary repository and an optional fallback repository.
type PokemonService struct {
	primaryRepo  domain.PokemonRepository
	fallbackRepo domain.PokemonRepository // Optional
}

// NewPokemonService creates a new PokemonService.
// fallbackRepo can be nil if no fallback is desired.
func NewPokemonService(primary domain.PokemonRepository, fallback domain.PokemonRepository) *PokemonService {
	return &PokemonService{
		primaryRepo:  primary,
		fallbackRepo: fallback,
	}
}

// FetchPokemon fetches a Pokemon by its identifier (name or ID).
// It tries the primary repository first, then the fallback if available and primary fails.
func (s *PokemonService) FetchPokemon(identifier string) (domain.Pokemon, error) {
	var pokemon domain.Pokemon
	var err error

	// Try to parse identifier as ID first, then as name
	id, errConv := strconv.Atoi(identifier)

	// Attempt with primary repository
	if s.primaryRepo != nil {
		if errConv == nil {
			// Identifier is likely an ID
			pokemon, err = s.primaryRepo.GetByID(id)
		} else {
			// Identifier is likely a name
			pokemon, err = s.primaryRepo.GetByName(strings.ToLower(identifier))
		}
		if err == nil {
			return pokemon, nil // Successfully fetched from primary
		}
		fmt.Printf("Primary repository failed for '%s': %v\n", identifier, err) // Log or handle error more gracefully
	}

	// Attempt with fallback repository if primary failed and fallback exists
	if s.fallbackRepo != nil {
		fmt.Printf("Attempting fallback repository for '%s'\n", identifier)
		if errConv == nil {
			pokemon, err = s.fallbackRepo.GetByID(id)
		} else {
			pokemon, err = s.fallbackRepo.GetByName(strings.ToLower(identifier))
		}
		if err == nil {
			return pokemon, nil // Successfully fetched from fallback
		}
		fmt.Printf("Fallback repository also failed for '%s': %v\n", identifier, err) // Log or handle error
	}

	// If both failed or no repositories configured properly
	if err != nil {
		return domain.Pokemon{}, fmt.Errorf("failed to fetch pokemon '%s' from all sources: %w", identifier, err)
	}
	return domain.Pokemon{}, fmt.Errorf("failed to fetch pokemon '%s': no repositories configured or other error", identifier)
}
