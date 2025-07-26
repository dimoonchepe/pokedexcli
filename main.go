package main

import (
    "os"
    "bufio"
    "strings"
    "fmt"
)

type cliCommand struct {
    name        string
    description string
    callback    func() error
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

func cleanInput(text string) []string {
   return strings.Fields(strings.ToLower(text))
}
