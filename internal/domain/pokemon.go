package domain

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
