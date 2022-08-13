package colormapping

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"image/color"
	"strconv"
	"strings"
)

//go:embed colors/*
var Files embed.FS

type ColorMapping struct {
	colors               map[string]*color.RGBA
	extendedpaletteblock map[string]bool
	extendedpalette      *Palette
}

func (m *ColorMapping) GetColor(name string, param2 int) *color.RGBA {
	//TODO: list of node->palette
	if m.extendedpaletteblock[name] {
		// param2 coloring
		return m.extendedpalette.GetColor(param2)
	}

	return m.colors[name]
}

func (m *ColorMapping) GetColors() map[string]*color.RGBA {
	return m.colors
}

func (m *ColorMapping) LoadBytes(buffer []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	count := 0
	line := 0

	for scanner.Scan() {
		line++

		txt := strings.Trim(scanner.Text(), " ")

		if len(txt) == 0 {
			//empty
			continue
		}

		if strings.HasPrefix(txt, "#") {
			//comment
			continue
		}

		parts := strings.Fields(txt)

		if len(parts) < 4 {
			return 0, errors.New("invalid line: #" + strconv.Itoa(line))
		}

		if len(parts) >= 4 {
			r, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return 0, err
			}

			g, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return 0, err
			}

			b, err := strconv.ParseInt(parts[3], 10, 32)
			if err != nil {
				return 0, err
			}

			a := int64(255)

			c := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			m.colors[parts[0]] = &c
			count++
		}
	}

	return count, nil
}

func (m *ColorMapping) LoadVFSColors(filename string) (int, error) {
	buffer, err := Files.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	return m.LoadBytes(buffer)
}

func (m *ColorMapping) LoadDefaults() error {
	list := []string{
		"advtrains.txt",
		"custom.txt",
		"mc2.txt",
		"miles.txt",
		"mtg.txt",
		"nodecore.txt",
		"scifi_nodes.txt",
		"vanessa.txt",
	}
	for _, name := range list {
		_, err := m.LoadVFSColors("colors/" + name)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewColorMapping() *ColorMapping {
	data, err := Files.ReadFile("colors/unifieddyes_palette_extended.png")
	if err != nil {
		panic(err)
	}

	extendedpalette, err := NewPalette(data)
	if err != nil {
		panic(err)
	}

	data, err = Files.ReadFile("colors/extended_palette.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	extendedpaletteblock := make(map[string]bool)

	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		txt := strings.Trim(scanner.Text(), " ")

		if len(txt) == 0 {
			//empty
			continue
		}

		extendedpaletteblock[txt] = true
	}

	return &ColorMapping{
		colors:               make(map[string]*color.RGBA),
		extendedpaletteblock: extendedpaletteblock,
		extendedpalette:      extendedpalette,
	}
}
