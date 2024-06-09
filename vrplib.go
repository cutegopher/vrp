package vrp

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
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
	Scheduled         bool
}

var start = Coordinates{X: 0.0, Y: 0.0}

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
// of all 200 x 200 load#.
func FillDistanceMatrix(loads []Load) [200][200]float64 {
	distanceMatrix := [200][200]float64{}
	for i := range loads {
		for j := range loads {
			if i == j {
				continue
			}
			distance := CalcDistance(loads[i].PickUp, loads[j].PickUp)
			distanceMatrix[i][j] = distance
			distanceMatrix[j][i] = distance
		}
	}
	return distanceMatrix
}

// GetNearestNeighbor finds the nearest neighbor for a given PickUp point.
func GetNearestNeighbor(distanceMatrix [200][200]float64, id int, done []int) (int, float64) {
	sortDist := make([]float64, len(distanceMatrix[id]))
	copy(sortDist, distanceMatrix[id][:])
	sort.Float64s(sortDist)
	var idNeigh int
	var distance float64

	for _, distance = range sortDist[1:] {
		idNeigh = getIdFromDistance(distanceMatrix[id][:], distance)
		if isPlanned(idNeigh+1, done) {
			continue
		}
		return idNeigh + 1, distance
	}
	return idNeigh, distance
}

// NearestToStart calculates nearest node to the Starting point (0.0, 0.0)
func NearestToStart(loads []Load, done []int) (Load, error) {
	var shortestDistance float64
	var nearestLoc Load

	if len(loads) == 0 {
		return nearestLoc, fmt.Errorf("the list is empty")
	}
	if len(loads) == 1 {
		return loads[0], fmt.Errorf("just one element found in the list")
	}

	shortestDistance = 9999999.00
	found := false

	for _, load := range loads {
		if isPlanned(load.ID, done) {
			continue
		}
		if load.DistanceFromStart < shortestDistance {
			nearestLoc = load
			shortestDistance = load.DistanceFromStart
			found = true
		}
	}
	if !found {
		return Load{}, fmt.Errorf("no matching neighbor found")
	}
	return nearestLoc, nil
}

func CalcDistance(a, b Coordinates) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

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
	return Load{
		ID: loadNo,
		PickUp: Coordinates{
			X: floats[0],
			Y: floats[1],
		},
		DropOff: Coordinates{
			X: floats[2],
			Y: floats[3],
		},
	}, nil
}

func ParseFile(fileName string) ([]Load, error) {
	loads := []Load{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Load{}, fmt.Errorf("error reading the file %s", fileName)
	}
	entries := strings.Split(string(data), "\n")
	if len(entries) == 0 {
		return []Load{}, fmt.Errorf("file %s is empty", fileName)
	}
	for _, entry := range entries {
		record, err := parseEntry(entry)
		if err != nil {
			continue
		}
		record.Distance = CalcDistance(record.PickUp, record.DropOff)
		record.DistanceFromStart = CalcDistance(start, record.PickUp)
		record.DistanceToStart = CalcDistance(record.DropOff, start)
		loads = append(loads, record)
	}

	if len(loads) == 0 {
		return []Load{}, fmt.Errorf("no matching records found in the file %s", fileName)
	}
	return loads, nil
}
