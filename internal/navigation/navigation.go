package navigation

import (
    "fmt"
    "errors"
    "io"
    "net/http"
    "encoding/json"
)

type config struct {
    nextURL string
    prevURL string
}
var conf config

func init() {
    conf = config {
        nextURL: "",
        prevURL: "",
    }
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

func CommandMap() error {
    url := "https://pokeapi.co/api/v2/location-area"
    if conf.nextURL != "" {
        url = conf.nextURL
    }
    return getLocations(url)
}

func CommandMapb() error {
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
