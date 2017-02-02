package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

var inputFilePath string
var prettyPrint bool

func init() {
	flag.StringVar(&inputFilePath, "input-file-path", "", "The input file containing a GeoJSON feature with a polygon")
	flag.BoolVar(&prettyPrint, "pretty-print", true, "Wether the output JSON should be well formatted or not")
	flag.Parse()
}

func main() {
	r, err := os.Open(inputFilePath)
	if err != nil {
		exitWithError(fmt.Errorf("Failed to open the input file: %v", err))
	}
	defer r.Close()

	normalizedFeature, err := normalizeGeoJsonFeature(r)
	if err != nil {
		exitWithError(fmt.Errorf("Failed to to normalize feature: %v", err))
	}

	b, err := encodeFeatureIntoJSON(normalizedFeature)
	if err != nil {
		exitWithError(fmt.Errorf("Failed to encode the feature into JSON: %v", err))
	}

	fmt.Println(string(b[:]))
}

func normalizeGeoJsonFeature(r io.Reader) (Feature, error) {
	feature, err := readFeature(r)
	if err != nil {
		return feature, fmt.Errorf("Failed to read the feature: %v", err)
	}
	feature.Geometry.Coordinates = normalizeGeoJsonCoordinates(feature.Geometry.Coordinates)
	return feature, nil
}

func readFeature(r io.Reader) (Feature, error) {
	feature := Feature{}
	if err := json.NewDecoder(r).Decode(&feature); err != nil {
		return feature, fmt.Errorf("Failed to decode the geojson feature from reader: %v", err)
	}
	return feature, nil
}

func normalizeGeoJsonCoordinates(coordinates [][][]float64) [][][]float64 {
	boundaries := getCoordinatesBoundaries(coordinates[0])
	lngCrossesTheDateLine := ((boundaries.MaxLng - boundaries.MinLng) > 180)
	latCrossesTheDateLine := ((boundaries.MaxLat - boundaries.MinLat) > 90)
	var normalizedCoordinates [][]float64
	for _, lngLat := range coordinates[0] {
		normalizedLngLat := append([]float64(nil), lngLat...)
		if lngCrossesTheDateLine && (normalizedLngLat[0] < 0) {
			normalizedLngLat[0] += 360
		}
		if latCrossesTheDateLine && (normalizedLngLat[1] < 0) {
			normalizedLngLat[1] += 180
		}
		normalizedCoordinates = append(normalizedCoordinates, normalizedLngLat)
	}
	return [][][]float64{normalizedCoordinates}
}

func getCoordinatesBoundaries(coordinates [][]float64) Boundaries {
	boundaries := Boundaries{
		MinLng: math.MaxFloat64,
		MaxLng: -math.MaxFloat64,
		MinLat: math.MaxFloat64,
		MaxLat: -math.MaxFloat64,
	}
	for _, coord := range coordinates {
		if coord[0] < boundaries.MinLng {
			boundaries.MinLng = coord[0]
		}
		if coord[0] > boundaries.MaxLng {
			boundaries.MaxLng = coord[0]
		}
		if coord[1] < boundaries.MinLat {
			boundaries.MinLat = coord[1]
		}
		if coord[1] > boundaries.MaxLat {
			boundaries.MaxLat = coord[1]
		}
	}
	return boundaries
}

func encodeFeatureIntoJSON(feature Feature) ([]byte, error) {
	b, err := json.Marshal(feature)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal the feature: %v", err)
	}

	if prettyPrint {
		formattedJSON := new(bytes.Buffer)
		if err := json.Indent(formattedJSON, b, "", "    "); err != nil {
			return nil, fmt.Errorf("Failed to index json: %v", err)
		}
		b = formattedJSON.Bytes()
	}

	return b, nil
}

func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

type Feature struct {
	Id         string      `json:"id"`
	Properties interface{} `json:"properties"`
	Geometry   struct {
		Type        string        `json:"type"`
		Orientation string        `json:"orientation"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Type string `json:"type"`
}

type Boundaries struct {
	MinLng, MaxLng, MinLat, MaxLat float64
}
