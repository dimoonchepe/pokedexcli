package main

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "errors"
    "io"
    "net/http"
    "encoding/json"
)

type cliCommand struct {
    name        string
    description string
    callback    func() error
}

type config struct {
    nextURL string
    prevURL string
}

type location struct {
    name string
    URL  string
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

var commands map[string]cliCommand
var conf config


func init() {
    commands = map[string]cliCommand {
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
	"map": {
	    name:        "map",
	    description: "List next 20 locations",
	    callback:    commandMap,
	},
	"mapb": {
	    name:        "mapb",
	    description: "List previous 20 locations",
	    callback:    commandMapb,
	},
    }
    conf = config {
        nextURL: "",
	prevURL: "",
    }
}


func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for true {
        fmt.Print("Pokedex > ") 
        scanner.Scan()
        command := scanner.Text()

        comm, exists := commands[command]
        if exists {
            comm.callback()
        } else {
            fmt.Println("Unknown command")
        }
    }
    
}

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp() error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:\n")
    for key, value := range commands {
        fmt.Printf("%s: %s\n", key, value.description)
    }
    return nil
}

func commandMap() error {
    url := "https://pokeapi.co/api/v2/location-area"
    if conf.nextURL != "" {
        url = conf.nextURL
    }
    return getLocations(url)
}

func commandMapb() error {
    if conf.prevURL == "" {
	fmt.Println("Already at the beginning of the locations list")	
	return nil
    }
    return getLocations(conf.prevURL)
}

func getLocations(url string) error {
    res, err := http.Get(url)
    if err != nil {
	return err
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
	return errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
    }
    if err != nil {
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

func cleanInput(text string) []string {
   return strings.Fields(strings.ToLower(text))
}
