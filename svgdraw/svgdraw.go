// svgdraw
package svgdraw

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Svg struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/svg svg"`
	Height  int      `xml:"height,attr"`
	Width   int      `xml:"width,attr"`
	GList   []*SvgG  `xml:"g"`
}

//fill="black" font="sans-serif" font-size="30" render-order="30" stroke="none" text-anchor="middle"

type SvgG struct {
	XMLName     xml.Name    `xml:"g"`
	Title       string      `xml:"title"`
	Rects       []*SvgRect  `xml:"rect"`
	Images      []*SvgImage `xml:"image"`
	Texts       []*SvgText  `xml:"text"`
	RenderOrder string      `xml:"render-order,attr,omitempty"`
	Fill        string      `xml:"fill,attr,omitempty"`
	Font        string      `xml:"font,attr,omitempty"`
	FontSize    string      `xml:"font-size,attr,omitempty"`
	Stroke      string      `xml:"stroke,attr,omitempty"`
	TextAnchor  string      `xml:"text-anchor,attr,omitempty"`
}

type SvgRect struct {
	XMLName       xml.Name `xml:"rect"`
	Id            string   `xml:"id,attr,omitempty"`
	Fill          string   `xml:"fill,attr,omitempty"`
	FillOpacity   float32  `xml:"fill-opacity,attr,omitempty"`
	Height        float32  `xml:"height,attr,omitempty"`
	Width         float32  `xml:"width,attr,omitempty"`
	X             float32  `xml:"x,attr,omitempty"`
	Y             float32  `xml:"y,attr,omitempty"`
	Rx            float32  `xml:"rx,attr,omitempty"`
	Ry            float32  `xml:"ry,attr,omitempty"`
	Stroke        string   `xml:"stroke,attr,omitempty"`
	StrokeOpacity float32  `xml:"stroke-opacity,attr,omitempty"`
	StrokeWidth   float32  `xml:"stroke-width,attr,omitempty"`
	Transform     string   `xml:"transform,attr,omitempty"`
}

//<ns0:text dx="58.23433899673569" dy="61.77261527154374" id="text_key_E" style="alignment-baseline:middle" transform="rotate(6.587978123034058,903.6142227704469,563.1951975245701)" x="845.3798837737112" y="501.4225822530263">E</ns0:text>
type SvgText struct {
	XMLName       xml.Name `xml:"text"`
	Id            string   `xml:"id,attr,omitempty"`
	Fill          string   `xml:"fill,attr"`
	FillOpacity   float32  `xml:"fill-opacity,attr,omitempty"`
	Height        float32  `xml:"height,attr,omitempty"`
	Width         float32  `xml:"width,attr,omitempty"`
	X             float32  `xml:"x,attr,omitempty"`
	Y             float32  `xml:"y,attr,omitempty"`
	Dx            float32  `xml:"dx,attr,omitempty"`
	Dy            float32  `xml:"dy,attr,omitempty"`
	Stroke        string   `xml:"stroke,attr,omitempty"`
	StrokeOpacity float32  `xml:"stroke-opacity,attr,omitempty"`
	StrokeWidth   float32  `xml:"stroke-width,attr,omitempty"`
	Transform     string   `xml:"transform,attr,omitempty"`
	Style         string   `xml:"style,attr,omitempty"`
	Text          string   `xml:",innerxml"`
}
type SvgImage struct {
	XMLName   xml.Name `xml:"image"`
	Id        string   `xml:"id,attr,omitempty"`
	Height    float32  `xml:"height,attr,omitempty"`
	Width     float32  `xml:"width,attr,omitempty"`
	X         float32  `xml:"x,attr,omitempty"`
	Y         float32  `xml:"y,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
	Style     string   `xml:"style,attr,omitempty"`
	Href      string   `xml:"http://www.w3.org/1999/xlink href,attr"`
}

type Configuration struct {
	CmdIcons map[string]string
}

var config Configuration

func KeyMap(KeyCommand map[string]string, hotkeys []string, modifier string) {

}

func SvgFileParse(filename string) (*Svg, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer xmlFile.Close()

	svg := &Svg{}
	decoder := xml.NewDecoder(xmlFile)
	if err := decoder.Decode(svg); err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return svg, nil
}

func SvgFileSave(svg *Svg, filename string) error {

	xmlFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer xmlFile.Close()

	xmlFile.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\n")

	enc := xml.NewEncoder(xmlFile)
	enc.Indent("  ", "    ")

	if err := enc.Encode(svg); err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	return nil
}

func DrawKeys(keyAction map[string]string, filename string) {
	svg, _ := SvgFileParse("configs/keyboard.svg")

	//var textLayer *SvgG
	var iconLayer *SvgG
	var btnsLayer *SvgG

	for _, g := range svg.GList {
		if g.Title == "Text Layer" {
			//textLayer = g
		}
		if g.Title == "Icons Layer" {
			iconLayer = g
		}
		if g.Title == "Buttons Layer" {
			btnsLayer = g
		}
	}

	//for _, r := range btnsLayer.Rects {
	//	st := SvgText{}
	//	st.XMLName.Local = "text"
	//	st.Id = "text_" + r.Id
	//	st.Dx = r.Width / 2.0
	//	st.Dy = r.Height / 2.0
	//	st.Style = "alignment-baseline:middle"
	//	st.Transform = r.Transform
	//	st.X = r.X
	//	st.Y = r.Y
	//	split_id := strings.Split(r.Id, "_")
	//	key := split_id[len(split_id)-1]

	//	st.Text = key
	//	textLayer.Texts = append(textLayer.Texts, &st)
	//}

	for _, r := range btnsLayer.Rects {
		split_id := strings.Split(r.Id, "_")
		key := split_id[len(split_id)-1]

		if action, ok := keyAction[key]; ok {
			icon, okk := config.CmdIcons[action]
			if !okk {
				icon = action
			}
			st := SvgImage{}
			st.XMLName.Local = "image"
			st.Id = "icon_" + r.Id
			st.Width = r.Width
			st.Height = r.Height
			st.Style = "alignment-baseline:middle"
			st.Transform = r.Transform
			st.X = r.X
			st.Y = r.Y
			st.Href = "icons/" + icon + ".png"

			iconLayer.Images = append(iconLayer.Images, &st)
		}

	}

	SvgFileSave(svg, filename)
}

func init() {
	fmt.Println("Init svgdraw")

	file, _ := os.Open("configs/svg.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}
