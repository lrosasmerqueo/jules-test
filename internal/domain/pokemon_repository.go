package domain

// PokemonRepository defines the interface for accessing Pokemon data.
type PokemonRepository interface {
	GetByID(id int) (Pokemon, error)
	GetByName(name string) (Pokemon, error)
	// GetByIdentifier could be an alternative if we want to keep a single method
	// GetByIdentifier(identifier string) (Pokemon, error)
}
