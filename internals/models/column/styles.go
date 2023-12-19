package column

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/lipgloss"
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

const (
	margin = 4
	APPEND = -1
)

func (m *Model) getSize(width, height int) {
	m.width = width / margin
}

func Palette() (*ColorData, error) {
	file, err := os.ReadFile("./colortemplate.json")
	if err != nil {
		return nil, err
	}

	destination := &ColorData{}

	err = json.Unmarshal(file, &destination)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

func (m *Model) GetColors() error {
	data, err := Palette()
	if err != nil {
		return err
	}

	m.colors = *data

	return nil
}

func (m *Model) SetStyle() {
	err := m.GetColors()
	if err != nil {
		panic(err)
	}

	t_color := lipgloss.AdaptiveColor{Dark: m.colors.Colors.Color4, Light: m.colors.Colors.Color12}
	b_color := lipgloss.AdaptiveColor{Dark: m.colors.Colors.Color5, Light: m.colors.Colors.Color13}

	m.Contents.Styles.Title.
		Background(t_color)
	m.Contents.Styles.FilterPrompt.
		Background(t_color)

	m.Contents.Styles.FilterCursor.
		Foreground(b_color)
}

func (m *Model) getStyle() lipgloss.Style {

	if m.Focused() {
		color := lipgloss.AdaptiveColor{Dark: m.colors.Colors.Color3, Light: m.colors.Colors.Color11}
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(color).
			Height(m.height).
			Width(m.width + 2*margin)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(m.height).
		Width(m.width + margin)
}
