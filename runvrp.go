package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vrp/vrp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Unexpected number of arguments\nCorrect usage: programe-name input-file-name")
	}
	loads, err := vrp.ParseFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing the file %s %v", os.Args[1], err)
	}
	done := []int{}
	xdistance := vrp.FillDistanceMatrix(loads)
	routesPlanned := 0
	allRoutes := [][]int{}

	for routesPlanned < len(loads) {
		near, err := vrp.NearestToStart(loads, done)
		if err != nil {
			fmt.Println("Could not find the node nearest to the Starting point")
			continue
		}

		// first node
		routesPlanned++
		done = append(done, near.ID)
		route := []int{near.ID}

		travelTime := near.DistanceFromStart           // #1
		travelTime = travelTime + near.DistanceToStart // #6
		middle1 := near.Distance                       // #2 (conditional, when second node cannot be added.
		var middle float64 = middle1

		// second node
		for {
			nearestNeighID, neighDistance := vrp.GetNearestNeighbor(xdistance, near.ID-1, done)             // second node
			newtravelTime := neighDistance                                                                  // #3
			newtravelTime = newtravelTime + vrp.CalcDistance(loads[nearestNeighID-1].DropOff, near.DropOff) // 5
			middle2 := loads[nearestNeighID-1].Distance                                                     // #4.

			if (travelTime + newtravelTime + middle2) <= 720 {
				middle = middle2
				travelTime = travelTime + newtravelTime

				done = append(done, nearestNeighID)
				route = append(route, nearestNeighID)

				routesPlanned++
				_ = travelTime + newtravelTime // time taken by the load to be delivered
				continue
			}
			_ = travelTime + middle // time taken by the load to be delivered
			travelTime = 0.0
			allRoutes = append(allRoutes, route)
			break
		}
	}

	for i := range allRoutes {
		fmt.Print("[")
		subStr := ""
		for _, j := range allRoutes[i] {
			subStr = subStr + fmt.Sprintf("%d,", j)
		}
		fmt.Printf("%s]\n", strings.Trim(subStr, ","))
	}
}
