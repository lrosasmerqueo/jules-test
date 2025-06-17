package cache

// Package cache provides an in-memory cache for Pokemon data.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"pokedex-app/internal/domain"
)

// InMemoryCache implements the domain.PokemonRepository interface using hardcoded data.
type InMemoryCache struct {
	data map[string]string // Stores raw JSON strings
}

const bulbasaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":64,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/1.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/1.ogg"},"forms":[{"name":"bulbasaur","url":"https://pokeapi.co/api/v2/pokemon-form/1/"}],"game_indices":[],"height":7,"held_items":[],"id":1,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/1/encounters","moves":[],"name":"bulbasaur","order":1,"past_abilities":[],"past_types":[],"species":{"name":"bulbasaur","url":"https://pokeapi.co/api/v2/pokemon-species/1/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png"},"stats":[{"base_stat":45,"effort":0,"stat":{"name":"hp"}},{"base_stat":49,"effort":0,"stat":{"name":"attack"}},{"base_stat":49,"effort":0,"stat":{"name":"defense"}},{"base_stat":65,"effort":1,"stat":{"name":"special-attack"}},{"base_stat":65,"effort":0,"stat":{"name":"special-defense"}},{"base_stat":45,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":69}`
const ivysaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":142,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/2.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/2.ogg"},"forms":[{"name":"ivysaur","url":"https://pokeapi.co/api/v2/pokemon-form/2/"}],"game_indices":[],"height":10,"held_items":[],"id":2,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/2/encounters","moves":[],"name":"ivysaur","order":2,"past_abilities":[],"past_types":[],"species":{"name":"ivysaur","url":"https://pokeapi.co/api/v2/pokemon-species/2/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/2.png"},"stats":[{"base_stat":60,"effort":0,"stat":{"name":"hp"}},{"base_stat":62,"effort":0,"stat":{"name":"attack"}},{"base_stat":63,"effort":0,"stat":{"name":"defense"}},{"base_stat":80,"effort":1,"stat":{"name":"special-attack"}},{"base_stat":80,"effort":1,"stat":{"name":"special-defense"}},{"base_stat":60,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":130}`
const venusaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":236,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/3.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/3.ogg"},"forms":[{"name":"venusaur","url":"https://pokeapi.co/api/v2/pokemon-form/3/"}],"game_indices":[],"height":20,"held_items":[],"id":3,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/3/encounters","moves":[],"name":"venusaur","order":3,"past_abilities":[],"past_types":[],"species":{"name":"venusaur","url":"https://pokeapi.co/api/v2/pokemon-species/3/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/3.png"},"stats":[{"base_stat":80,"effort":0,"stat":{"name":"hp"}},{"base_stat":82,"effort":0,"stat":{"name":"attack"}},{"base_stat":83,"effort":0,"stat":{"name":"defense"}},{"base_stat":100,"effort":2,"stat":{"name":"special-attack"}},{"base_stat":100,"effort":1,"stat":{"name":"special-defense"}},{"base_stat":80,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":1000}`

// NewInMemoryCache creates a new InMemoryCache and populates it with hardcoded data.
func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		data: make(map[string]string),
	}
	cache.data["1"] = bulbasaurJSONConst
	cache.data["bulbasaur"] = bulbasaurJSONConst
	cache.data["2"] = ivysaurJSONConst
	cache.data["ivysaur"] = ivysaurJSONConst
	cache.data["3"] = venusaurJSONConst
	cache.data["venusaur"] = venusaurJSONConst
	return cache
}

// parsePokemonData unmarshals JSON data from the cache into a domain.Pokemon struct.
func parsePokemonData(jsonData []byte) (domain.Pokemon, error) {
	var pokemon domain.Pokemon
	err := json.Unmarshal(jsonData, &pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("error parsing cached pokemon JSON data: %w", err)
	}
	return pokemon, nil
}

// GetByName implements the domain.PokemonRepository interface.
func (c *InMemoryCache) GetByName(name string) (domain.Pokemon, error) {
	jsonData, ok := c.data[strings.ToLower(name)]
	if !ok {
		return domain.Pokemon{}, fmt.Errorf("pokemon '%s' not found in cache", name)
	}
	return parsePokemonData([]byte(jsonData))
}

// GetByID implements the domain.PokemonRepository interface.
func (c *InMemoryCache) GetByID(id int) (domain.Pokemon, error) {
	jsonData, ok := c.data[strconv.Itoa(id)]
	if !ok {
		return domain.Pokemon{}, fmt.Errorf("pokemon with ID %d not found in cache", id)
	}
	return parsePokemonData([]byte(jsonData))
}
