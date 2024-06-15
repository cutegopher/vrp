package vrp

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Coordinates struct {
	X float64
	Y float64
}

type Load struct {
	ID                int
	PickUp            Coordinates
	DropOff           Coordinates
	Distance          float64
	DistanceFromStart float64
	DistanceToStart   float64
}

type idDistance struct {
	id       int
	distance float64
}

var Start = Coordinates{X: 0.0, Y: 0.0}

func isPlanned(id int, done []int) bool {
	for i := range done {
		if done[i] == id {
			return true
		}
	}
	return false
}

func getIdFromDistance(neighbors []float64, distance float64) int {
	for i, _ := range neighbors {
		if distance == neighbors[i] {
			return i
		}
	}
	return 0
}

// FillDistanceMatrix calculates the distance between the PickUp points
// and drop off points.
func FillDistanceMatrix(loads []Load) [200][200]float64 {
	distanceMatrix := [200][200]float64{}
	for i := range loads {
		for j := range loads {
			if i == j {
				continue
			}
			distance := CalcDistance(loads[i].DropOff, loads[j].PickUp)
			distanceMatrix[i][j] = distance
		}
	}
	return distanceMatrix
}

// NearestFromStart returns the nearest PickUp location from the
// Start that is available.
func NearestFromStart(idDistances []idDistance, done []int) int {
	for _, d := range idDistances {
		if isPlanned(d.id, done) {
			continue
		}
		return d.id
	}
	return 0
}

// NearestNeighbor returns the nearest load number to the loadiID id
func NearestNeighbor(loads []Load, matrix [200][200]float64, id int, done []int) int {
	distances := matrix[id-1][:]
	iDistances := []idDistance{}

	for i, d := range distances {
		if d == 0 {
			continue
		}
		iDistances = append(iDistances, idDistance{
			id:       i + 1,
			distance: d,
		})
	}
	slices.SortFunc(iDistances,
		func(a, b idDistance) int {
			return cmp.Compare(a.distance, b.distance)
		})

	for _, d := range iDistances[1:] {

		if isPlanned(d.id, done) {
			continue
		}
		return d.id
	}
	return -1
}

// CalcDistance calculates the distance between two Cartesian points (X,Y). 
func CalcDistance(a, b Coordinates) float64 {
	var xDiff float64 = b.X - a.X
	var yDiff float64 = b.Y - a.Y
	return math.Sqrt((xDiff * xDiff) + (yDiff * yDiff))
}

// parseEntry parses the string s and matches it against the regular expression
// to extract load#, load-PickUp location and load-dropOff location
func parseEntry(s string) (Load, error) {
	var loadNo int
	var err error

	re := regexp.MustCompile(`(\d+)\s+\((.*),(.*)\)\s+\((.*),(.*)\)`)
	matchString := re.FindAllStringSubmatch(s, -1)

	if len(matchString) != 1 {
		return Load{}, fmt.Errorf("no matching load entry")
	}
	if len(matchString[0]) != 6 {
		return Load{}, fmt.Errorf("invalid number of records in the load entry")
	}
	if loadNo, err = strconv.Atoi(matchString[0][1]); err != nil {
		return Load{}, fmt.Errorf("falied to parse int value %v", err)
	}
	floats := []float64{}
	for _, flt := range matchString[0][2:] {
		var fl float64
		if fl, err = strconv.ParseFloat(flt, 64); err != nil {
			return Load{}, fmt.Errorf("error in parsing float value %v", err)
		}
		floats = append(floats, fl)
	}
	ii := Load{
		ID: loadNo,
		PickUp: Coordinates{
			X: floats[0],
			Y: floats[1],
		},
		DropOff: Coordinates{
			X: floats[2],
			Y: floats[3],
		},
	}
	return ii, nil
}

// ParseFile parses the input to extract load#, load-PickUp, location
// and load-dropOff location
func ParseFile(fileName string) ([]Load, []idDistance, error) {
	loads := []Load{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Load{}, []idDistance{}, fmt.Errorf("error reading the file %s", fileName)
	}
	entries := strings.Split(string(data), "\n")
	if len(entries) == 0 {
		return []Load{}, []idDistance{}, fmt.Errorf("file %s is empty", fileName)
	}

	iDistances := make([]idDistance, len(loads))

	for _, entry := range entries {
		record, err := parseEntry(entry)
		if err != nil {
			continue
		}
		record.Distance = CalcDistance(record.PickUp, record.DropOff)
		record.DistanceFromStart = CalcDistance(Start, record.PickUp)
		record.DistanceToStart = CalcDistance(record.DropOff, Start)

		loads = append(loads, record)
		iDistances = append(iDistances, idDistance{
			id:       record.ID,
			distance: record.DistanceFromStart,
		})
	}

	if len(loads) == 0 {
		return []Load{}, []idDistance{}, fmt.Errorf("no matching records found in the file %s", fileName)
	}
	slices.SortFunc(iDistances,
		func(a, b idDistance) int {
			return cmp.Compare(a.distance, b.distance)
		})
	return loads, iDistances, nil
}
