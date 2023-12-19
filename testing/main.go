package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ColorData struct {
	Special struct {
		Background string `json:"background"`
		Foreground string `json:"foreground"`
	} `json:"special"`
	Colors struct {
		Color0  string `json:"color0"`
		Color1  string `json:"color1"`
		Color2  string `json:"color2"`
		Color3  string `json:"color3"`
		Color4  string `json:"color4"`
		Color5  string `json:"color5"`
		Color6  string `json:"color6"`
		Color7  string `json:"color7"`
		Color8  string `json:"color8"`
		Colot10 string `json:"color10"`
		Color11 string `json:"color11"`
		Color12 string `json:"color12"`
		Color13 string `json:"color13"`
		Color14 string `json:"color14"`
		Color15 string `json:"color15"`
	} `json:"colors"`
}

func main() {
	file, err := os.ReadFile("./colortemplate.json")
	if err != nil {
		panic(err)
	}

	destination := ColorData{}

	err = json.Unmarshal(file, &destination)
	if err != nil {
		panic(err)
	}

	fmt.Println(destination)
}
