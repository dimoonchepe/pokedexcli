package navigation

import (
    "fmt"
    "errors"
    "io"
    "time"
    "net/http"
    "encoding/json"
    "math/rand"
    "github.com/dimoonchepe/pokedexcli/internal/pokecache"
)

type config struct {
    nextURL string
    prevURL string
}
var conf config
var cache pokecache.Cache
var pokedex map[string]pokemonDetails

func init() {
    conf = config {
        nextURL: "",
        prevURL: "",
    }
    cache = pokecache.NewCache(5 * time.Second)
    pokedex = make(map[string]pokemonDetails)
}

type locationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationDetails struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemonDetails struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
}


func CommandMap(_ string) error {
    url := "https://pokeapi.co/api/v2/location-area"
    if conf.nextURL != "" {
        url = conf.nextURL
    }
    return getLocations(url)
}

func CommandMapb(_ string) error {
    if conf.prevURL == "" {
	fmt.Println("Already at the beginning of the locations list")	
	return nil
    }
    return getLocations(conf.prevURL)
}

func getLocations(url string) error {
    body, err := makeCachedRequest(url)
    if (err != nil) {
        return err
    }
    response := locationResponse{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return err
    }
    conf.nextURL = ""
    if response.Next != nil {
        conf.nextURL = *response.Next
    }
    conf.prevURL = ""
    if response.Previous != nil {
    	conf.prevURL = *response.Previous
    }
    for i := 0; i < len(response.Results); i++ {
        fmt.Println(response.Results[i].Name)
    }
    return nil
}

func CommandExplore(locationName string) error {
    fmt.Println("Exploring", locationName)
    url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)

    body, err := makeCachedRequest(url)
    if (err != nil) {
        return err
    }
    response := locationDetails{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return err
    }
    pokemon := []string{}
    for _, enc := range response.PokemonEncounters {
        pokemon = append(pokemon, enc.Pokemon.Name)
    }
    if len(pokemon) == 0 {
        fmt.Println("No Pokemon Found")
        return nil
    }
    fmt.Println("Found Pokemon:")
    for _, pok := range pokemon {
        fmt.Println("-", pok)
    }
    return nil
}

func CommandCatch(name string) error {
    url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

    body, err := makeCachedRequest(url)
    if (err != nil) {
        return err
    }
    response := pokemonDetails{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return err
    }
    random := rand.Intn(response.BaseExperience)
    fmt.Printf("Throwing a Pokeball at %s...\n", name)
    if random < 30 {
        fmt.Println(name, "was caught!")
        pokedex[name] = response
    } else {
        fmt.Println(name, "escaped!")
    }
    return nil 
}

func makeCachedRequest(url string) ([]byte, error) {
    cached, exists := cache.Get(url)
    if exists {
        return cached, nil
    }
    res, err := http.Get(url)
    if err != nil {
        return []byte{}, err
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
        return []byte{}, errors.New(
            fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", 
                res.StatusCode, body),
        )
    }
    if err != nil {
        return []byte{}, err
    }
    return body, nil
}
