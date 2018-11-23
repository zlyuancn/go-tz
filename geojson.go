package gotz

import (
	"encoding/json"
)

type FeatureCollection struct {
	featureCollection
}

type featureCollection struct {
	Features []*Feature `json:"features"`
}

type Feature struct {
	feature
}

type feature struct {
	Geometry   Geometry `json:"geometry"`
	Properties struct {
		Tzid string `json:"tzid"`
	} `json:"properties"`
}

type Geometry struct {
	geometry
}

type geometry struct {
	Type        string    `json:"type"`
	Coordinates [][]Point `json:"coordinates"`
}

var jPolyType struct {
	Type       string      `json:"type"`
	Geometries []*Geometry `json:"geometries"`
}

var jPolygon struct {
	Coordinates [][][]float64 `json:"coordinates"`
}

var jMultiPolygon struct {
	Coordinates [][][][]float64 `json:"coordinates"`
}

func (g *Geometry) UnmarshalJSON(data []byte) (err error) {
	if err := json.Unmarshal(data, &jPolyType); err != nil {
		return err
	}
	g.Type = "MultiPolygon"

	if jPolyType.Type == "Polygon" {
		if err := json.Unmarshal(data, &jPolygon); err != nil {
			return err
		}
		//Create a bounding box
		pol := make([]Point, len(jPolygon.Coordinates[0]))
		for i, v := range jPolygon.Coordinates[0] {
			pol[i].Lon = v[0]
			pol[i].Lat = v[1]
		}
		b := getBoundingBox(pol)
		g.Coordinates = append(g.Coordinates, b, pol)
		return nil
	}

	if jPolyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &jMultiPolygon); err != nil {
			return err
		}
		for _, poly := range jMultiPolygon.Coordinates {
			pol := make([]Point, len(poly[0]))
			for i, v := range poly[0] {
				pol[i].Lon = v[0]
				pol[i].Lat = v[1]
			}
			b := getBoundingBox(pol)
			g.Coordinates = append(g.Coordinates, b, pol)
		}
		return nil
	}
	return nil
}
