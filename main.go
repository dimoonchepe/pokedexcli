package main

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "github.com/dimoonchepe/pokedexcli/internal/navigation"
)

type cliCommand struct {
    name        string
    description string
    callback    func(arg string) error
}

var commands map[string]cliCommand


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
            callback:    navigation.CommandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "List previous 20 locations",
            callback:    navigation.CommandMapb,
        },
        "explore": {
            name:        "explore",
            description: "Explore selected location",
            callback:    navigation.CommandExplore,
        },
        "catch": {
            name:        "catch",
            description: "Try catching a specified pokemon",
            callback:    navigation.CommandCatch,
        },
        "inspect": {
            name:        "inspect",
            description: "Get stats of caught pokemon",
            callback:    navigation.CommandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "List all caught pokemon",
            callback:    navigation.CommandPokedex,
        },
    }
}


func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for true {
        fmt.Print("Pokedex > ") 
        scanner.Scan()
        fields := strings.Fields(scanner.Text())
        if len(fields) > 2 {
            fmt.Println("Unknown command")
        }
        command := fields[0]
        arg := ""
        if len(fields) == 2 {
            arg = fields[1]
        }

        comm, exists := commands[command]
        if exists {
            comm.callback(arg)
        } else {
            fmt.Println("Unknown command")
        }
    }
    
}

func commandExit(_ string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(_ string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println("")
    for key, value := range commands {
        fmt.Printf("%s: %s\n", key, value.description)
    }
    return nil
}


func cleanInput(text string) []string {
   return strings.Fields(strings.ToLower(text))
}
