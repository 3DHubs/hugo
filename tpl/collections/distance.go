// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collections

import (
	"errors"
	"math"
	"reflect"
	"strings"
)

// Where returns a filtered subset of a given data type.
func (ns *Namespace) DistanceSort(seq interface{}, fieldName interface{}, lat interface{}, lon interface{}) (interface{}, error) {
	if seq == nil {
		return nil, errors.New("sequence must be provided")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		// ok
	default:
		return nil, errors.New("can't sort " + reflect.ValueOf(seq).Type().String())
	}

	sortByField, ok := fieldName.(string)
	if !ok {
		return nil, errors.New("fieldName should be a string.")
	}

	centerLat, ok := lat.(float64)
	if !ok {
		return nil, errors.New("centerLat should be a float.")
	}

	centerLon, ok := lon.(float64)
	if !ok {
		return nil, errors.New("centerLon should be a float.")
	}

	// Create a list of pairs that will be used to do the sort
	p := pairList{SortAsc: true, SliceType: reflect.SliceOf(seqv.Type().Elem())}
	p.Pairs = make([]pair, seqv.Len())

	path := strings.Split(strings.Trim(sortByField, "."), ".")

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < seqv.Len(); i++ {
			p.Pairs[i].Value = seqv.Index(i)
			v := p.Pairs[i].Value
			var err error
			for _, elemName := range path {
				v, err = evaluateSubElem(v, elemName)
				if err != nil {
					return nil, err
				}
			}

			var location map[string]interface{}
			location = v.Interface().(map[string]interface{})
			p.Pairs[i].Key = reflect.ValueOf(Distance(centerLat, centerLon, location["lat"].(float64), location["lon"].(float64)))
		}
	case reflect.Map:
		keys := seqv.MapKeys()
		for i := 0; i < seqv.Len(); i++ {
			p.Pairs[i].Value = seqv.MapIndex(keys[i])
			v := p.Pairs[i].Value
			var err error
			for _, elemName := range path {
				v, err = evaluateSubElem(v, elemName)
				if err != nil {
					return nil, err
				}
			}
			var location map[string]interface{}
			location = v.Interface().(map[string]interface{})
			p.Pairs[i].Key = reflect.ValueOf(Distance(centerLat, centerLon, location["lat"].(float64), location["lon"].(float64)))
		}
	}

	return p.sort(), nil
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
