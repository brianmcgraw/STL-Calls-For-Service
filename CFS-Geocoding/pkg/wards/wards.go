package wards

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

func IsPointInsidePolygon(fc []*geojson.Feature, point orb.Point) (int, bool) {

	for _, feature := range fc {
		polygon, _ := feature.Geometry.(orb.Polygon)

		if planar.PolygonContains(polygon, point) {
			actualWard := int(feature.Properties["WARD"].(float64))
			return actualWard, true
		}
	}
	return 0, false
}

func IsPointInsideMultiPolygon(fc []*geojson.Feature, point orb.Point) (string, bool) {

	var actualNeighborhood string
	var outcome bool
	for _, feature := range fc {
		multiPolygon, _ := feature.Geometry.(orb.MultiPolygon)

		if planar.MultiPolygonContains(multiPolygon, point) {
			actualNeighborhood = feature.Properties["NHD_NAME"].(string)
			outcome = true
		}
	}

	return actualNeighborhood, outcome
}
