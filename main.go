package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Declare UI elements as global variables
var (
	pokemonImage *canvas.Image
	nameLabel    *widget.Label
	idLabel      *widget.Label

	// About Tab
	typesLabel   *widget.Label
	heightLabel  *widget.Label
	weightLabel  *widget.Label
	speciesLabel *widget.Label

	// Stats Tab
	hpLabel        *widget.Label
	attackLabel    *widget.Label
	defenseLabel   *widget.Label
	spAttackLabel  *widget.Label
	spDefenseLabel *widget.Label
	speedLabel     *widget.Label

	// Abilities Tab
	abilitiesListLabel *widget.Label

	// Navigation & Error
	errorLabel             *widget.Label
	prevButton             *widget.Button
	nextButton             *widget.Button
	searchEntry            *widget.Entry
	currentPokemonID       int = 1
	hardcodedPokemonData   map[string]string
	placeholderImageResource fyne.Resource
)

const bulbasaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":64,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/1.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/1.ogg"},"forms":[{"name":"bulbasaur","url":"https://pokeapi.co/api/v2/pokemon-form/1/"}],"game_indices":[],"height":7,"held_items":[],"id":1,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/1/encounters","moves":[],"name":"bulbasaur","order":1,"past_abilities":[],"past_types":[],"species":{"name":"bulbasaur","url":"https://pokeapi.co/api/v2/pokemon-species/1/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png"},"stats":[{"base_stat":45,"effort":0,"stat":{"name":"hp"}},{"base_stat":49,"effort":0,"stat":{"name":"attack"}},{"base_stat":49,"effort":0,"stat":{"name":"defense"}},{"base_stat":65,"effort":1,"stat":{"name":"special-attack"}},{"base_stat":65,"effort":0,"stat":{"name":"special-defense"}},{"base_stat":45,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":69}`
const ivysaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":142,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/2.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/2.ogg"},"forms":[{"name":"ivysaur","url":"https://pokeapi.co/api/v2/pokemon-form/2/"}],"game_indices":[],"height":10,"held_items":[],"id":2,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/2/encounters","moves":[],"name":"ivysaur","order":2,"past_abilities":[],"past_types":[],"species":{"name":"ivysaur","url":"https://pokeapi.co/api/v2/pokemon-species/2/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/2.png"},"stats":[{"base_stat":60,"effort":0,"stat":{"name":"hp"}},{"base_stat":62,"effort":0,"stat":{"name":"attack"}},{"base_stat":63,"effort":0,"stat":{"name":"defense"}},{"base_stat":80,"effort":1,"stat":{"name":"special-attack"}},{"base_stat":80,"effort":1,"stat":{"name":"special-defense"}},{"base_stat":60,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":130}`
const venusaurJSONConst = `{"abilities":[{"ability":{"name":"overgrow","url":"https://pokeapi.co/api/v2/ability/65/"},"is_hidden":false,"slot":1},{"ability":{"name":"chlorophyll","url":"https://pokeapi.co/api/v2/ability/34/"},"is_hidden":true,"slot":3}],"base_experience":236,"cries":{"latest":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/3.ogg","legacy":"https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/legacy/3.ogg"},"forms":[{"name":"venusaur","url":"https://pokeapi.co/api/v2/pokemon-form/3/"}],"game_indices":[],"height":20,"held_items":[],"id":3,"is_default":true,"location_area_encounters":"https://pokeapi.co/api/v2/pokemon/3/encounters","moves":[],"name":"venusaur","order":3,"past_abilities":[],"past_types":[],"species":{"name":"venusaur","url":"https://pokeapi.co/api/v2/pokemon-species/3/"},"sprites":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/3.png"},"stats":[{"base_stat":80,"effort":0,"stat":{"name":"hp"}},{"base_stat":82,"effort":0,"stat":{"name":"attack"}},{"base_stat":83,"effort":0,"stat":{"name":"defense"}},{"base_stat":100,"effort":2,"stat":{"name":"special-attack"}},{"base_stat":100,"effort":1,"stat":{"name":"special-defense"}},{"base_stat":80,"effort":0,"stat":{"name":"speed"}}],"types":[{"slot":1,"type":{"name":"grass"}},{"slot":2,"type":{"name":"poison"}}],"weight":1000}`

func init() {
	hardcodedPokemonData = map[string]string{
		"1":           bulbasaurJSONConst,
		"bulbasaur":   bulbasaurJSONConst,
		"2":           ivysaurJSONConst,
		"ivysaur":     ivysaurJSONConst,
		"3":           venusaurJSONConst,
		"venusaur":    venusaurJSONConst,
	}
}

func concatenateTypeNames(pokemonTypes []PokemonType) string {
	var typeNames []string
	for _, t := range pokemonTypes {
		typeNames = append(typeNames, strings.Title(t.Type.Name))
	}
	return strings.Join(typeNames, ", ")
}

func concatenateAbilityNames(pokemonAbilities []PokemonAbility) string {
	var abilityNames []string
	for _, a := range pokemonAbilities {
		abilityNames = append(abilityNames, strings.Title(a.Ability.Name))
	}
	return strings.Join(abilityNames, ", ")
}

func clearPokemonDisplay() {
	nameLabel.SetText("Name: N/A")
	idLabel.SetText("#---")
	if placeholderImageResource != nil {
		pokemonImage.Resource = placeholderImageResource
	} else {
		// Fallback if placeholder itself failed to load, though unlikely if file is present
		pokemonImage.Resource = nil
	}
	pokemonImage.Refresh()

	speciesLabel.SetText("Species: N/A")
	heightLabel.SetText("Height: N/A")
	weightLabel.SetText("Weight: N/A")
	typesLabel.SetText("Type(s): N/A")

	hpLabel.SetText("HP: N/A")
	attackLabel.SetText("Attack: N/A")
	defenseLabel.SetText("Defense: N/A")
	spAttackLabel.SetText("Sp. Attack: N/A")
	spDefenseLabel.SetText("Sp. Defense: N/A")
	speedLabel.SetText("Speed: N/A")
	abilitiesListLabel.SetText("Abilities: N/A")

	if errorLabel != nil { // errorLabel might not be initialized when this is first called by main
		errorLabel.SetText("Enter Pokémon Name/ID and press Enter, or use Next/Previous.")
	}

	// Refresh all labels
	nameLabel.Refresh(); idLabel.Refresh(); speciesLabel.Refresh(); heightLabel.Refresh()
	weightLabel.Refresh(); typesLabel.Refresh(); hpLabel.Refresh(); attackLabel.Refresh()
	defenseLabel.Refresh(); spAttackLabel.Refresh(); spDefenseLabel.Refresh(); speedLabel.Refresh()
	abilitiesListLabel.Refresh()
	if errorLabel != nil { errorLabel.Refresh() }
}


func displayPokemonData(pokemon Pokemon) {
	nameLabel.SetText(strings.Title(pokemon.Name))
	idLabel.SetText(fmt.Sprintf("#%d", pokemon.ID))

	typesLabel.SetText(fmt.Sprintf("Type(s): %s", concatenateTypeNames(pokemon.Types)))
	heightLabel.SetText(fmt.Sprintf("Height: %.1f m", float64(pokemon.Height)/10.0))
	weightLabel.SetText(fmt.Sprintf("Weight: %.1f kg", float64(pokemon.Weight)/10.0))
	if pokemon.Species.Name != "" {
		speciesLabel.SetText(fmt.Sprintf("Species: %s", strings.Title(pokemon.Species.Name)))
	} else {
		speciesLabel.SetText("Species: N/A")
	}

	// Reset stats before setting, useful if a Pokemon has fewer stats than previous
	defaultStatText := func(s string) string { return fmt.Sprintf("%s: N/A", s) }
	hpLabel.SetText(defaultStatText("HP"))
	attackLabel.SetText(defaultStatText("Attack"))
	defenseLabel.SetText(defaultStatText("Defense"))
	spAttackLabel.SetText(defaultStatText("Sp. Attack"))
	spDefenseLabel.SetText(defaultStatText("Sp. Defense"))
	speedLabel.SetText(defaultStatText("Speed"))

	for _, s := range pokemon.Stats {
		statValue := fmt.Sprintf("%s: %d", strings.Title(s.Stat.Name), s.BaseStat)
		switch s.Stat.Name {
		case "hp":
			hpLabel.SetText(statValue)
		case "attack":
			attackLabel.SetText(statValue)
		case "defense":
			defenseLabel.SetText(statValue)
		case "special-attack":
			spAttackLabel.SetText(statValue)
		case "special-defense":
			spDefenseLabel.SetText(statValue)
		case "speed":
			speedLabel.SetText(statValue)
		}
	}

	abilitiesListLabel.SetText(fmt.Sprintf("Abilities: %s", concatenateAbilityNames(pokemon.Abilities)))

	if pokemon.Sprites.FrontDefault != "" {
		imgURL, err := url.Parse(pokemon.Sprites.FrontDefault)
		if err != nil {
			log.Printf("Error parsing image URL %s: %v\n", pokemon.Sprites.FrontDefault, err)
			errorLabel.SetText(fmt.Sprintf("Error (image URL): %v", err))
			pokemonImage.Resource = placeholderImageResource // Fallback to placeholder
		} else {
			resourceChan := make(chan fyne.Resource)
			errorChan := make(chan error)
			go func() {
				res, e := fyne.LoadResourceFromURL(imgURL)
				if e != nil {
					errorChan <- e
					return
				}
				resourceChan <- res
			}()

			select {
			case res := <-resourceChan:
				log.Println("Successfully loaded image from URL:", pokemon.Sprites.FrontDefault)
				pokemonImage.Resource = res
			case err := <-errorChan:
				log.Printf("Error loading image resource from %s: %v. Using placeholder.", pokemon.Sprites.FrontDefault, err)
				errorLabel.SetText(fmt.Sprintf("Image load failed: %v. Using placeholder.", err))
				pokemonImage.Resource = placeholderImageResource
			case <-time.After(15 * time.Second):
				log.Printf("Timeout loading image from %s. Using placeholder.", pokemon.Sprites.FrontDefault)
				errorLabel.SetText("Image load timed out. Using placeholder.")
				pokemonImage.Resource = placeholderImageResource
			}
		}
	} else {
		log.Println("No front_default sprite URL found. Using placeholder.")
		pokemonImage.Resource = placeholderImageResource
	}
	// Refresh all UI elements that were updated
	pokemonImage.Refresh(); nameLabel.Refresh(); idLabel.Refresh(); typesLabel.Refresh()
	heightLabel.Refresh(); weightLabel.Refresh(); speciesLabel.Refresh(); hpLabel.Refresh()
	attackLabel.Refresh(); defenseLabel.Refresh(); spAttackLabel.Refresh(); spDefenseLabel.Refresh()
	speedLabel.Refresh(); abilitiesListLabel.Refresh(); errorLabel.Refresh()
}

func fetchAndDisplayPokemon(identifier string) {
	clearPokemonDisplay() // Clear display before fetching new data
	errorLabel.SetText(fmt.Sprintf("Fetching %s...", identifier))

	if prevButton != nil { prevButton.Disable() }
	if nextButton != nil { nextButton.Disable() }
	if searchEntry != nil { searchEntry.Disable() }

	defer func() {
		if prevButton != nil { prevButton.Enable() }
		if nextButton != nil { nextButton.Enable() }
		if searchEntry != nil { searchEntry.Enable() }
		if currentPokemonID <= 1 {
			if prevButton != nil { prevButton.Disable() }
		}
		// Max ID check could be added here if we had a reliable max_pokemon_id
		// For now, allow next to always be enabled unless at an extreme.
		// if currentPokemonID >= MAX_POKEMON_ID_KNOWN { nextButton.Disable() }
	}()

	log.Printf("Attempting to fetch: %s", identifier)
	pokemon, err := GetPokemon(identifier)
	if err == nil {
		log.Printf("Successfully fetched %s from API.", identifier)
		errorLabel.SetText(fmt.Sprintf("Displaying %s", strings.Title(pokemon.Name)))
		currentPokemonID = pokemon.ID
		displayPokemonData(pokemon)
		return
	}

	log.Printf("Live API call for '%s' failed: %v. Trying hardcoded data.", identifier, err)
	errorLabel.SetText(fmt.Sprintf("%v. Trying cache...", err)) // Display the refined error from GetPokemon

	jsonDataString, ok := hardcodedPokemonData[strings.ToLower(identifier)]
	if !ok {
		log.Printf("Pokemon '%s' not found in hardcoded data.", identifier)
		clearPokemonDisplay() // Clear again to ensure N/A state
		errorLabel.SetText(fmt.Sprintf("Network unavailable, and Pokémon '%s' not found in local cache.", identifier))
		nameLabel.SetText("Not Found"); idLabel.SetText("#---") // Explicitly set "Not Found"
		nameLabel.Refresh(); idLabel.Refresh()
		return
	}

	parsedPokemon, parseErr := ParsePokemonData(jsonDataString)
	if parseErr != nil {
		log.Printf("Error parsing hardcoded data for '%s': %v", identifier, parseErr)
		clearPokemonDisplay() // Clear again
		errorLabel.SetText(fmt.Sprintf("Error parsing local data for '%s': %v", identifier, parseErr))
		return
	}

	log.Printf("Displaying %s from hardcoded data.", parsedPokemon.Name)
	errorLabel.SetText(fmt.Sprintf("Displaying %s (Cached - Network Down)", strings.Title(parsedPokemon.Name)))
	currentPokemonID = parsedPokemon.ID
	displayPokemonData(parsedPokemon)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Go Pokedex")

	var errPholder error // Define error variable for placeholder loading
	placeholderImageResource, errPholder = fyne.LoadResourceFromPath("placeholder.png")
	if errPholder != nil {
		log.Fatalln("CRITICAL: Could not load placeholder image:", errPholder)
		// In a real app, you might use a built-in icon or draw something basic
		pokemonImage = canvas.NewImageFromImage(nil)
	} else {
		pokemonImage = canvas.NewImageFromResource(placeholderImageResource)
	}
	pokemonImage.FillMode = canvas.ImageFillContain
	pokemonImage.SetMinSize(fyne.NewSize(150, 150))

	nameLabel = widget.NewLabelWithStyle("Name: N/A", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	idLabel = widget.NewLabelWithStyle("#---", fyne.TextAlignCenter, fyne.TextStyle{})
	leftPanel := container.NewVBox(pokemonImage, nameLabel, idLabel, layout.NewSpacer())

	typesLabel = widget.NewLabel("Type(s): N/A")
	heightLabel = widget.NewLabel("Height: N/A")
	weightLabel = widget.NewLabel("Weight: N/A")
	speciesLabel = widget.NewLabel("Species: N/A")
	aboutTabContent := container.NewVBox(typesLabel, heightLabel, weightLabel, speciesLabel)

	hpLabel = widget.NewLabel("HP: N/A")
	attackLabel = widget.NewLabel("Attack: N/A")
	defenseLabel = widget.NewLabel("Defense: N/A")
	spAttackLabel = widget.NewLabel("Sp. Attack: N/A")
	spDefenseLabel = widget.NewLabel("Sp. Defense: N/A")
	speedLabel = widget.NewLabel("Speed: N/A")
	statsTabContent := container.NewVBox(hpLabel, attackLabel, defenseLabel, spAttackLabel, spDefenseLabel, speedLabel)

	abilitiesListLabel = widget.NewLabel("Abilities: N/A")
	abilitiesTabContent := container.NewVBox(abilitiesListLabel)

	tabs := container.NewAppTabs(
		container.NewTabItem("About", aboutTabContent),
		container.NewTabItem("Stats", statsTabContent),
		container.NewTabItem("Abilities", abilitiesTabContent),
	)
	rightPanel := tabs
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.3

	errorLabel = widget.NewLabel("Enter Pokémon Name/ID and press Enter, or use Next/Previous.")
	errorLabel.Wrapping = fyne.TextWrapWord

	prevButton = widget.NewButton("Previous", func() {
		if currentPokemonID > 1 {
			fetchAndDisplayPokemon(strconv.Itoa(currentPokemonID - 1))
		}
	})
	nextButton = widget.NewButton("Next", func() {
		fetchAndDisplayPokemon(strconv.Itoa(currentPokemonID + 1))
	})
	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by name or ID...")
	searchEntry.OnSubmitted = func(searchText string) {
		trimmedText := strings.TrimSpace(searchText)
		if trimmedText != "" {
			fetchAndDisplayPokemon(trimmedText)
			searchEntry.SetText("") // Clear search entry after submission
		}
	}

	bottomNav := container.NewBorder(nil, nil, prevButton, nextButton, searchEntry)
	uiContent := container.NewBorder(nil, container.NewVBox(bottomNav, errorLabel), nil, nil, split)

	myWindow.SetContent(uiContent)
	myWindow.Resize(fyne.NewSize(800, 650))

	clearPokemonDisplay() // Clear display initially
	fetchAndDisplayPokemon(strconv.Itoa(currentPokemonID)) // Initial load

	if nameLabel.Text == "Name: N/A" { // Check if initial load really failed
		errorLabel.SetText("Failed to load initial Pokedex data. Check connection or search.")
		nameLabel.SetText("Pokedex") // General title if no data
		idLabel.SetText("#")
	}

	myWindow.ShowAndRun()
}
