package fyneui // Changed package name to fyneui to avoid conflict

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app" // Added for app.New()
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"pokedex-app/internal/application"
	"pokedex-app/internal/domain"
)

// FyneUI holds all UI elements and the application service.
type FyneUI struct {
	app                      fyne.App
	mainWindow               fyne.Window
	pokemonService           *application.PokemonService
	pokemonImage             *canvas.Image
	nameLabel                *widget.Label
	idLabel                  *widget.Label
	typesLabel               *widget.Label
	heightLabel              *widget.Label
	weightLabel              *widget.Label
	speciesLabel             *widget.Label
	hpLabel                  *widget.Label
	attackLabel              *widget.Label
	defenseLabel             *widget.Label
	spAttackLabel            *widget.Label
	spDefenseLabel           *widget.Label
	speedLabel               *widget.Label
	abilitiesListLabel       *widget.Label
	errorLabel               *widget.Label
	prevButton               *widget.Button
	nextButton               *widget.Button
	searchEntry              *widget.Entry
	currentPokemonID         int
	placeholderImageResource fyne.Resource
}

// NewFyneUI creates and initializes a new Fyne-based UI.
func NewFyneUI(appService *application.PokemonService) *FyneUI {
	fyneApp := fyne.CurrentApp()
	if fyneApp == nil {
		fyneApp = app.New()
	} // Ensure app exists
	mainWindow := fyneApp.NewWindow("Go Pokedex - Refactored")

	ui := &FyneUI{
		app:              fyneApp,
		mainWindow:       mainWindow,
		pokemonService:   appService,
		currentPokemonID: 1, // Default starting Pokemon
	}

	// Load placeholder image
	var errPholder error
	ui.placeholderImageResource, errPholder = fyne.LoadResourceFromPath("placeholder.png")
	if errPholder != nil {
		log.Println("CRITICAL: Could not load placeholder image:", errPholder)
		ui.pokemonImage = canvas.NewImageFromImage(nil) // Fallback
	} else {
		ui.pokemonImage = canvas.NewImageFromResource(ui.placeholderImageResource)
	}
	ui.pokemonImage.FillMode = canvas.ImageFillContain
	ui.pokemonImage.SetMinSize(fyne.NewSize(150, 150))

	// Initialize labels and other widgets
	ui.nameLabel = widget.NewLabelWithStyle("Name: N/A", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	ui.idLabel = widget.NewLabelWithStyle("#---", fyne.TextAlignCenter, fyne.TextStyle{})
	ui.typesLabel = widget.NewLabel("Type(s): N/A")
	ui.heightLabel = widget.NewLabel("Height: N/A")
	ui.weightLabel = widget.NewLabel("Weight: N/A")
	ui.speciesLabel = widget.NewLabel("Species: N/A")
	ui.hpLabel = widget.NewLabel("HP: N/A")
	ui.attackLabel = widget.NewLabel("Attack: N/A")
	ui.defenseLabel = widget.NewLabel("Defense: N/A")
	ui.spAttackLabel = widget.NewLabel("Sp. Attack: N/A")
	ui.spDefenseLabel = widget.NewLabel("Sp. Defense: N/A")
	ui.speedLabel = widget.NewLabel("Speed: N/A")
	ui.abilitiesListLabel = widget.NewLabel("Abilities: N/A")
	ui.errorLabel = widget.NewLabel("Enter Pokémon Name/ID and press Enter, or use Next/Previous.")
	ui.errorLabel.Wrapping = fyne.TextWrapWord

	// Setup layout
	leftPanel := container.NewVBox(ui.pokemonImage, ui.nameLabel, ui.idLabel, layout.NewSpacer())
	aboutTabContent := container.NewVBox(ui.typesLabel, ui.heightLabel, ui.weightLabel, ui.speciesLabel)
	statsTabContent := container.NewVBox(ui.hpLabel, ui.attackLabel, ui.defenseLabel, ui.spAttackLabel, ui.spDefenseLabel, ui.speedLabel)
	abilitiesTabContent := container.NewVBox(ui.abilitiesListLabel)
	tabs := container.NewAppTabs(
		container.NewTabItem("About", aboutTabContent),
		container.NewTabItem("Stats", statsTabContent),
		container.NewTabItem("Abilities", abilitiesTabContent),
	)
	rightPanel := tabs
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.3

	// Navigation buttons and search entry
	ui.prevButton = widget.NewButton("Previous", func() {
		if ui.currentPokemonID > 1 {
			ui.fetchAndDisplayPokemon(strconv.Itoa(ui.currentPokemonID - 1))
		}
	})
	ui.nextButton = widget.NewButton("Next", func() {
		ui.fetchAndDisplayPokemon(strconv.Itoa(ui.currentPokemonID + 1))
	})
	ui.searchEntry = widget.NewEntry()
	ui.searchEntry.SetPlaceHolder("Search by name or ID...")
	ui.searchEntry.OnSubmitted = func(searchText string) {
		trimmedText := strings.TrimSpace(searchText)
		if trimmedText != "" {
			ui.fetchAndDisplayPokemon(trimmedText)
			ui.searchEntry.SetText("") // Clear search entry
		}
	}

	bottomNav := container.NewBorder(nil, nil, ui.prevButton, ui.nextButton, ui.searchEntry)
	uiContent := container.NewBorder(nil, container.NewVBox(bottomNav, ui.errorLabel), nil, nil, split)

	ui.mainWindow.SetContent(uiContent)
	ui.mainWindow.Resize(fyne.NewSize(800, 650))

	return ui
}

// Start begins the Fyne application UI.
func (ui *FyneUI) Start() {
	ui.clearPokemonDisplay()                                     // Clear display initially
	ui.fetchAndDisplayPokemon(strconv.Itoa(ui.currentPokemonID)) // Initial load
	if ui.nameLabel.Text == "Name: N/A" {                        // Check if initial load really failed
		ui.errorLabel.SetText("Failed to load initial Pokedex data. Check connection or search.")
		ui.nameLabel.SetText("Pokedex")
		ui.idLabel.SetText("#")
	}
	ui.mainWindow.ShowAndRun()
}

func (ui *FyneUI) concatenateTypeNames(pokemonTypes []domain.PokemonType) string {
	var typeNames []string
	for _, t := range pokemonTypes {
		typeNames = append(typeNames, strings.Title(t.Type.Name))
	}
	return strings.Join(typeNames, ", ")
}

func (ui *FyneUI) concatenateAbilityNames(pokemonAbilities []domain.PokemonAbility) string {
	var abilityNames []string
	for _, a := range pokemonAbilities {
		abilityNames = append(abilityNames, strings.Title(a.Ability.Name))
	}
	return strings.Join(abilityNames, ", ")
}

func (ui *FyneUI) clearPokemonDisplay() {
	ui.nameLabel.SetText("Name: N/A")
	ui.idLabel.SetText("#---")
	if ui.placeholderImageResource != nil {
		ui.pokemonImage.Resource = ui.placeholderImageResource
	} else {
		ui.pokemonImage.Resource = nil
	}
	ui.pokemonImage.Refresh()
	ui.speciesLabel.SetText("Species: N/A")
	ui.heightLabel.SetText("Height: N/A")
	ui.weightLabel.SetText("Weight: N/A")
	ui.typesLabel.SetText("Type(s): N/A")
	ui.hpLabel.SetText("HP: N/A")
	ui.attackLabel.SetText("Attack: N/A")
	ui.defenseLabel.SetText("Defense: N/A")
	ui.spAttackLabel.SetText("Sp. Attack: N/A")
	ui.spDefenseLabel.SetText("Sp. Defense: N/A")
	ui.speedLabel.SetText("Speed: N/A")
	ui.abilitiesListLabel.SetText("Abilities: N/A")
	if ui.errorLabel != nil {
		ui.errorLabel.SetText("Enter Pokémon Name/ID and press Enter, or use Next/Previous.")
	}
	ui.nameLabel.Refresh()
	ui.idLabel.Refresh()
	ui.speciesLabel.Refresh()
	ui.heightLabel.Refresh()
	ui.weightLabel.Refresh()
	ui.typesLabel.Refresh()
	ui.hpLabel.Refresh()
	ui.attackLabel.Refresh()
	ui.defenseLabel.Refresh()
	ui.spAttackLabel.Refresh()
	ui.spDefenseLabel.Refresh()
	ui.speedLabel.Refresh()
	ui.abilitiesListLabel.Refresh()
	if ui.errorLabel != nil {
		ui.errorLabel.Refresh()
	}
}

func (ui *FyneUI) displayPokemonData(pokemon domain.Pokemon) {
	ui.nameLabel.SetText(strings.Title(pokemon.Name))
	ui.idLabel.SetText(fmt.Sprintf("#%d", pokemon.ID))
	ui.typesLabel.SetText(fmt.Sprintf("Type(s): %s", ui.concatenateTypeNames(pokemon.Types)))
	ui.heightLabel.SetText(fmt.Sprintf("Height: %.1f m", float64(pokemon.Height)/10.0))
	ui.weightLabel.SetText(fmt.Sprintf("Weight: %.1f kg", float64(pokemon.Weight)/10.0))
	if pokemon.Species.Name != "" {
		ui.speciesLabel.SetText(fmt.Sprintf("Species: %s", strings.Title(pokemon.Species.Name)))
	} else {
		ui.speciesLabel.SetText("Species: N/A")
	}
	defaultStatText := func(s string) string { return fmt.Sprintf("%s: N/A", s) }
	ui.hpLabel.SetText(defaultStatText("HP"))
	ui.attackLabel.SetText(defaultStatText("Attack"))
	ui.defenseLabel.SetText(defaultStatText("Defense"))
	ui.spAttackLabel.SetText(defaultStatText("Sp. Attack"))
	ui.spDefenseLabel.SetText(defaultStatText("Sp. Defense"))
	ui.speedLabel.SetText(defaultStatText("Speed"))
	for _, s := range pokemon.Stats {
		statValue := fmt.Sprintf("%s: %d", strings.Title(s.Stat.Name), s.BaseStat)
		switch s.Stat.Name {
		case "hp":
			ui.hpLabel.SetText(statValue)
		case "attack":
			ui.attackLabel.SetText(statValue)
		case "defense":
			ui.defenseLabel.SetText(statValue)
		case "special-attack":
			ui.spAttackLabel.SetText(statValue)
		case "special-defense":
			ui.spDefenseLabel.SetText(statValue)
		case "speed":
			ui.speedLabel.SetText(statValue)
		}
	}
	ui.abilitiesListLabel.SetText(fmt.Sprintf("Abilities: %s", ui.concatenateAbilityNames(pokemon.Abilities)))
	if pokemon.Sprites.FrontDefault != "" {
		imgURL, err := url.Parse(pokemon.Sprites.FrontDefault)
		if err != nil {
			log.Printf("Error parsing image URL %s: %v\n", pokemon.Sprites.FrontDefault, err)
			ui.errorLabel.SetText(fmt.Sprintf("Error (image URL): %v", err))
			ui.pokemonImage.Resource = ui.placeholderImageResource
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
				ui.pokemonImage.Resource = res
			case err := <-errorChan:
				log.Printf("Error loading image resource from %s: %v. Using placeholder.", pokemon.Sprites.FrontDefault, err)
				ui.errorLabel.SetText(fmt.Sprintf("Image load failed: %v. Using placeholder.", err))
				ui.pokemonImage.Resource = ui.placeholderImageResource
			case <-time.After(15 * time.Second): // Increased timeout
				log.Printf("Timeout loading image from %s. Using placeholder.", pokemon.Sprites.FrontDefault)
				ui.errorLabel.SetText("Image load timed out. Using placeholder.")
				ui.pokemonImage.Resource = ui.placeholderImageResource
			}
		}
	} else {
		log.Println("No front_default sprite URL found. Using placeholder.")
		ui.pokemonImage.Resource = ui.placeholderImageResource
	}
	ui.pokemonImage.Refresh()
	ui.nameLabel.Refresh()
	ui.idLabel.Refresh()
	ui.typesLabel.Refresh()
	ui.heightLabel.Refresh()
	ui.weightLabel.Refresh()
	ui.speciesLabel.Refresh()
	ui.hpLabel.Refresh()
	ui.attackLabel.Refresh()
	ui.defenseLabel.Refresh()
	ui.spAttackLabel.Refresh()
	ui.spDefenseLabel.Refresh()
	ui.speedLabel.Refresh()
	ui.abilitiesListLabel.Refresh()
	ui.errorLabel.Refresh()
}

func (ui *FyneUI) fetchAndDisplayPokemon(identifier string) {
	ui.clearPokemonDisplay()
	ui.errorLabel.SetText(fmt.Sprintf("Fetching %s...", identifier))
	if ui.prevButton != nil {
		ui.prevButton.Disable()
	}
	if ui.nextButton != nil {
		ui.nextButton.Disable()
	}
	if ui.searchEntry != nil {
		ui.searchEntry.Disable()
	}
	defer func() {
		if ui.prevButton != nil {
			ui.prevButton.Enable()
		}
		if ui.nextButton != nil {
			ui.nextButton.Enable()
		}
		if ui.searchEntry != nil {
			ui.searchEntry.Enable()
		}
		if ui.currentPokemonID <= 1 {
			if ui.prevButton != nil {
				ui.prevButton.Disable()
			}
		}
	}()
	log.Printf("UI: Attempting to fetch: %s", identifier)
	pokemon, err := ui.pokemonService.FetchPokemon(identifier)
	if err != nil {
		log.Printf("UI: Error fetching pokemon '%s': %v", identifier, err)
		ui.clearPokemonDisplay() // Clear again to ensure N/A state
		ui.errorLabel.SetText(fmt.Sprintf("Error: %v", err))
		ui.nameLabel.SetText("Not Found")
		ui.idLabel.SetText("#---")
		ui.nameLabel.Refresh()
		ui.idLabel.Refresh()
		return
	}
	log.Printf("UI: Successfully fetched %s.", pokemon.Name)
	ui.errorLabel.SetText(fmt.Sprintf("Displaying %s", strings.Title(pokemon.Name)))
	ui.currentPokemonID = pokemon.ID
	ui.displayPokemonData(pokemon)
}
