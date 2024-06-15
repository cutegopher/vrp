package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"vrp/vrp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Unexpected number of arguments\nCorrect usage: programe-name input-file-name")
	}
	loads, distanceFromStart, err := vrp.ParseFile(os.Args[1])

	if err != nil {
		log.Fatalf("Error parsing the file %s %v", os.Args[1], err)
	}
	done := []int{}
	xdistance := vrp.FillDistanceMatrix(loads)
	routesPlanned := 0
	allRoutes := [][]int{}

	for routesPlanned < len(loads) {
		// first node
		newID := vrp.NearestFromStart(distanceFromStart, done)
		routesPlanned++
		done = append(done, newID)
		route := []int{newID}

		var currentLocation vrp.Coordinates = vrp.Start
		var ReturnPath float64 = vrp.CalcDistance(loads[newID-1].DropOff, vrp.Start)
		var disTravelled float64 = vrp.CalcDistance(currentLocation, loads[newID-1].PickUp)
		disTravelled = disTravelled + vrp.CalcDistance(loads[newID-1].PickUp, loads[newID-1].DropOff)
		currentLocation = loads[newID-1].DropOff

		for {
			var newnearestNeighID int
			var newDistance, newReturnPath float64
			iterate := true

            // second node
			newnearestNeighID = vrp.NearestNeighbor(loads, xdistance, newID, done)
			if newnearestNeighID == -1 {
				iterate = false
			} else {
				newDistance = vrp.CalcDistance(currentLocation, loads[newnearestNeighID-1].PickUp)
				newDistance = newDistance + vrp.CalcDistance(loads[newnearestNeighID-1].PickUp, loads[newnearestNeighID-1].DropOff)
				newReturnPath = vrp.CalcDistance(loads[newnearestNeighID-1].DropOff, vrp.Start)
			}
			if (disTravelled+newDistance+newReturnPath) >= 720 || !iterate {
				disTravelled = disTravelled + ReturnPath
				allRoutes = append(allRoutes, route)
				break
			}
			currentLocation = loads[newnearestNeighID-1].DropOff
			disTravelled = disTravelled + newDistance
			ReturnPath = newReturnPath
			routesPlanned++
			done = append(done, newnearestNeighID)
			route = append(route, newnearestNeighID)
		}
	}

	for i := 0; i < len(allRoutes); i++ {
		strs := []string{}
		for j := 0; j < len(allRoutes[i]); j++ {
			strs = append(strs, strconv.Itoa(allRoutes[i][j]))
		}
		fmt.Printf("[%s]\n", strings.Join(strs, ","))
	}
}
