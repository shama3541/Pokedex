package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shama3541/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, string) error
}

type config struct {
	Next     string
	Previous string
}

type Pokemon struct {
	Name        string `json:"name"`
	Pokemondata interface{}
}

var pokedex = make(map[string]Pokemon)
var commands map[string]cliCommand

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the map of the Pokedex",
			callback:    CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous map of the Pokedex",
			callback:    commandPrevious,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Pokedex",
			callback:    exploreArea,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    inspect,
		},
	}
	cfg := config{}
	cache := pokecache.NewCache(5 * time.Second)
	Location := ""
	for {
		fmt.Print("Pokedex > ")

		// Wait for user input
		scanned := scanner.Scan()
		if !scanned {
			break // EOF or error
		}

		line := scanner.Text()
		words := cleanInput(line)

		if len(words) == 0 {
			continue // Empty input, prompt again
		}
		commandName := words[0]
		command, exists := commands[commandName]
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}
		if commandName == "explore" || commandName == "catch" || commandName == "inspect" {
			Location = words[1]
			command.callback(&cfg, cache, Location)
			continue
		}

		command.callback(&cfg, cache, Location)
	}
}

func commandExit(cfg *config, cache *pokecache.Cache, Location string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache, Location string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range commands {
		fmt.Printf("%s - %s\n", cmd.name, cmd.description)
	}
	return nil
}

func CommandMap(cfg *config, cache *pokecache.Cache, Location string) error {
	var body []byte

	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
		cfg.Next = url
	}

	if cached, found := cache.Get(url); found {
		fmt.Println("🔁 Using cached data")
		body = cached
	} else {
		fmt.Println("🌐 Fetching from API:", url)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching data: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error: received status code %d", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}

		cache.Add(url, body)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	if nexturl, ok := data["next"].(string); ok {
		cfg.Next = nexturl
	} else {
		cfg.Next = ""
	}

	if prevurl, ok := data["previous"].(string); ok {
		cfg.Previous = prevurl
	} else {
		cfg.Previous = ""
	}

	results, ok := data["results"].([]interface{})
	if !ok {
		return fmt.Errorf("error: 'results' is not an array")
	}

	for _, location := range results {
		locMap, ok := location.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error: location is not a map")
		}

		name, ok := locMap["name"].(string)
		if !ok {
			return fmt.Errorf("error: 'name' key not found or not a string")
		}

		fmt.Println(name)
	}

	return nil
}

func commandPrevious(cfg *config, cache *pokecache.Cache, Location string) error {
	if cfg.Previous == "nil" || cfg.Previous == "" {
		fmt.Println("No previous page available.")
		return nil
	}

	cfg.Next = cfg.Previous
	cfg.Previous = ""

	return CommandMap(cfg, cache, Location)

}

func exploreArea(cfg *config, cache *pokecache.Cache, Location string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", Location)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()
	var mydata map[string]interface{}
	jsondata, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if err := json.Unmarshal(jsondata, &mydata); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	fmt.Printf("Exploring %s....\n", Location)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range mydata["pokemon_encounters"].([]interface{}) {
		pokemonMap, ok := pokemon.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error: pokemon encounter is not a map")
		}
		pokemonName, ok := pokemonMap["pokemon"].(map[string]interface{})["name"].(string)
		if !ok {
			return fmt.Errorf("error: 'name' key not found or not a string")
		}
		fmt.Println("-", pokemonName)
	}
	return nil
}

func catch(cfg *config, cache *pokecache.Cache, pokemon string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %v", err)
	}
	defer res.Body.Close()
	var mypokemon map[string]interface{}
	jsondata, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if err := json.Unmarshal(jsondata, &mypokemon); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	chance := rand.Float64() * 1000
	if _, ok := pokedex[pokemon]; !ok {

		if chance >= mypokemon["base_experience"].(float64) {
			fmt.Printf("Caught %s!\n", mypokemon["name"])
			pokedex[mypokemon["name"].(string)] = Pokemon{
				Name:        mypokemon["name"].(string),
				Pokemondata: mypokemon,
			}
		} else {
			fmt.Printf("%s escaped!\n", mypokemon["name"])
		}
	}
	return nil
}

func inspect(cfg *config, cache *pokecache.Cache, pokemon string) error {
	if pkmn, ok := pokedex[pokemon]; ok {
		fmt.Printf("Inspecting %s...\n", pkmn.Name)

		// Type assert Pokémon data
		dataMap, ok := pkmn.Pokemondata.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid pokemon data format")
		}

		// Print basic info
		fmt.Println("Name:", pkmn.Name)
		if height, ok := dataMap["height"].(float64); ok {
			fmt.Println("Height:", height)
		}
		if weight, ok := dataMap["weight"].(float64); ok {
			fmt.Println("Weight:", weight)
		}

		// Print stats
		fmt.Println("Stats:")
		if statsArr, ok := dataMap["stats"].([]interface{}); ok {
			for _, statItem := range statsArr {
				statMap, ok := statItem.(map[string]interface{})
				if !ok {
					continue
				}
				baseStat := statMap["base_stat"].(float64)
				statInfo := statMap["stat"].(map[string]interface{})
				statName := statInfo["name"].(string)
				fmt.Printf("  -%s: %.0f\n", statName, baseStat)
			}
		}

		// Print types
		fmt.Println("Types:")
		if typesArr, ok := dataMap["types"].([]interface{}); ok {
			for _, t := range typesArr {
				typeMap, ok := t.(map[string]interface{})
				if !ok {
					continue
				}
				typeInfo := typeMap["type"].(map[string]interface{})
				typeName := typeInfo["name"].(string)
				fmt.Printf("  - %s\n", typeName)
			}
		}

	} else {
		fmt.Println("You have not caught that Pokémon")
	}
	return nil
}
