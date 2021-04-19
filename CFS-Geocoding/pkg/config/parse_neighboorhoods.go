package config

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/paulmach/orb/geojson"
)

// Relying on geojson files from https://github.com/slu-openGIS/STL_BOUNDARY_Nhood
func (c Config) BuildNeighborhoodFC() (fc *geojson.FeatureCollection, err error) {
	var features []*geojson.Feature

	files, _ := ioutil.ReadDir(c.Geofile_Loc)

	fc = geojson.NewFeatureCollection()

	for _, v := range files {

		if v.Name() == "wards.geojson" {
			continue
		}
		b, err := ioutil.ReadFile(c.Geofile_Loc + "/" + v.Name())

		if err != nil {
			log.Printf("Error reading file: %v %v", v.Name(), err)
		}

		featureCollection, err := geojson.UnmarshalFeatureCollection(b)

		if err != nil {
			log.Fatalf("Unable to decode geojson file: %v", err)
		}
		for _, f := range featureCollection.Features {
			ftwo := f
			features = append(features, ftwo)
		}

	}

	fc.Features = features

	if len(fc.Features) == 0 {
		return fc, errors.New("Error creating neighborhoods geojson config")
	}

	return fc, err
}
