package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dudiko2/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, conf *config) error
}

type commandMap map[string]cliCommand

func newCommands() commandMap {
	m := commandMap{
		"help": {
			name:        "help",
			description: "Displays usage info",
			callback:    runHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    runExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    runMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    runMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area",
			callback:    runExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    runCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect caught Pokemon",
			callback:    runInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List caughtr Pokemon",
			callback:    runPokedex,
		},
	}
	return m
}

type config struct {
	pokeapiClient           pokeapi.Client
	commands                commandMap
	prevLocationAreasOffset int
	nextLocationAreasOffset int
}

func newConfig() *config {
	c := config{
		pokeapiClient:           *pokeapi.NewClient(),
		commands:                newCommands(),
		prevLocationAreasOffset: -1,
		nextLocationAreasOffset: -1,
	}
	return &c
}

var caughtPokemon = map[string]pokeapi.PokemonRes{}

const locationAreasLimit = 20

func main() {
	conf := newConfig()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printPrompt()
		input := prepareInput(scanner)
		execCommand(input, conf)
	}
}

func printPrompt() {
	fmt.Print("pokedex > ")
}

func printUnknownCmd(cmd string) {
	fmt.Println("Unknown command: " + cmd)
}

func execCommand(inp parsedInput, conf *config) {
	if inp.command == "" {
		return
	}
	defer fmt.Println("")
	cmd, ok := conf.commands[inp.command]
	if !ok {
		printUnknownCmd(inp.command)
		runHelp(inp.arguments, conf)
		return
	}
	err := cmd.callback(inp.arguments, conf)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func runHelp(args []string, conf *config) error {
	res := "Usage:\n\n"
	for name, cmd := range conf.commands {
		res += fmt.Sprintf("%s: %s\n", name, cmd.description)
	}
	fmt.Print(res)
	return nil
}

func runExit(args []string, c *config) error {
	os.Exit(0)
	return nil
}

func forEach[T any](list []T, callback func(T, int)) {
	for i, item := range list {
		callback(item, i)
	}
}

func printLocationAreasNames(la []pokeapi.LocationAreasEntry) {
	forEach(la, func(area pokeapi.LocationAreasEntry, i int) {
		fmt.Println(area.Name)
	})
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseOffsetParam(urlString *string) int {
	if urlString == nil {
		return -1
	}
	// XXX handle errs
	u, _ := url.Parse(*urlString)
	offsetParam, _ := strconv.Atoi(u.Query().Get("offset"))
	return offsetParam
}

func getLocations(offset int, conf *config) error {
	pl := pokeapi.GetLocationAreasPayload{
		Offset: offset,
		Limit:  locationAreasLimit,
	}
	d, err := conf.pokeapiClient.GetLocationAreas(pl)
	if err != nil {
		return err
	}
	locationAreas := d.Results
	printLocationAreasNames(locationAreas)
	conf.nextLocationAreasOffset = parseOffsetParam(d.Next)
	conf.prevLocationAreasOffset = parseOffsetParam(d.Previous)
	return nil
}

func runMapNext(args []string, conf *config) error {
	offset := maxInt(0, conf.nextLocationAreasOffset)
	return getLocations(offset, conf)
}

func runMapBack(args []string, conf *config) error {
	offset := maxInt(0, conf.prevLocationAreasOffset)
	return getLocations(offset, conf)
}

func printLocationExplore(l pokeapi.LocationRes) {
	encountersLen := len(l.PokemonEncounters)
	if encountersLen < 1 {
		fmt.Println("No Pokemon found!")
		return
	}
	fmt.Printf("Found Pokemon in %s:\n", l.Name)
	for _, e := range l.PokemonEncounters {
		fmt.Printf("- %s\n", e.Pokemon.Name)
	}
}

func runExplore(args []string, conf *config) error {
	argsLen := len(args)
	if argsLen < 1 {
		return errors.New("missing argument: id")
	}
	locationID := args[0]
	fmt.Printf("Exploring %s...\n", locationID)
	d, err := conf.pokeapiClient.GetLocationArea(locationID)
	if err != nil {
		return err
	}
	printLocationExplore(d)
	return nil
}

func calculateChance(value int) bool {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	num := random.Intn(value)
	return num < 20
}

func runCatch(args []string, conf *config) error {
	argsLen := len(args)
	if argsLen < 1 {
		return errors.New("missing argument: id")
	}
	pokemonID := args[0]
	pokemon, err := conf.pokeapiClient.GetPokemonData(pokemonID)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	caught := calculateChance(pokemon.BaseExperience)
	if caught {
		fmt.Printf("%s was caught!", pokemon.Name)
		caughtPokemon[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!", pokemon.Name)
	}
	return nil
}

func printPokemon(p pokeapi.PokemonRes) {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %v\n", p.Height)
}

func runInspect(args []string, conf *config) error {
	argsLen := len(args)
	if argsLen < 1 {
		return errors.New("missing argument: id")
	}
	pokemonName := args[0]
	p, caught := caughtPokemon[pokemonName]
	if !caught {
		fmt.Printf("you have not caught %s\n", pokemonName)
		return nil
	}
	printPokemon(p)
	return nil
}

func runPokedex(args []string, conf *config) error {
	pokedexSize := len(caughtPokemon)
	if pokedexSize < 1 {
		fmt.Println("Your Pokedex is empty")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for pokemonName := range caughtPokemon {
		fmt.Printf("- %s\n", pokemonName)
	}
	return nil
}
