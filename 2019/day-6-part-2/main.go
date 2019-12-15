/*
 * --- Part Two ---
 *
 * Now, you just need to figure out how many orbital transfers you (YOU) need to take to get to Santa (SAN).
 *
 * You start at the object YOU are orbiting; your destination is the object SAN is orbiting.
 * An orbital transfer lets you move from any object to an object orbiting or orbited by that object.
 *
 * For example, suppose you have the following map:
 *
 * COM)B
 * B)C
 * C)D
 * D)E
 * E)F
 * B)G
 * G)H
 * D)I
 * E)J
 * J)K
 * K)L
 * K)YOU
 * I)SAN
 *
 * Visually, the above map of orbits looks like this:
 *
 *                           YOU
 *                          /
 *         G - H       J - K - L
 *        /           /
 * COM - B - C - D - E - F
 *                \
 *                 I - SAN
 *
 * In this example, YOU are in orbit around K, and SAN is in orbit around I.
 * To move from K to I, a minimum of 4 orbital transfers are required:
 *
 *     K to J
 *     J to E
 *     E to D
 *     D to I
 *
 * Afterward, the map of orbits looks like this:
 *
 *         G - H       J - K - L
 *        /           /
 * COM - B - C - D - E - F
 *                \
 *                 I - SAN
 *                  \
 *                   YOU
 *
 * What is the minimum number of orbital transfers required to move from the object YOU are orbiting to the object SAN is orbiting?
 * (Between the objects they are orbiting - not between YOU and SAN.)
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type SpaceObject struct {
	Name     string
	Orbits   *SpaceObject
	Orbiting []*SpaceObject
}

type Path struct {
	VisitedObjects []*SpaceObject
}

func (object *SpaceObject) TotalNumberOfOrbitedObjects() int {
	orbitedObjects := 0
	if object.Orbits != nil {
		orbitedObjects += 1 + object.Orbits.TotalNumberOfOrbitedObjects()
	}
	return orbitedObjects
}

// Gather all possible paths to all possible other nodes
func (object *SpaceObject) allPathsToAllOtherObjects() [][]*SpaceObject {
	paths := [][]*SpaceObject{
		[]*SpaceObject{object},
	}

	i := 0
	for {
		if i >= len(paths) {
			log.Printf("HEEELLLOOO??? %d \n", i)
			return paths
		}

		currentPath := paths[i]
		lastNode := currentPath[len(currentPath)-1]

		possibleNextVisitedNodes := lastNode.Orbiting
		if lastNode.Orbits != nil {
			possibleNextVisitedNodes = append(possibleNextVisitedNodes, lastNode.Orbits)
		}

		for _, possibleNextNode := range possibleNextVisitedNodes {
			// Check that one node has not been visited twice
			nodeWasAlreadyVisited := false
			for _, node := range currentPath {
				if node == possibleNextNode {
					nodeWasAlreadyVisited = true
				}
			}

			fmt.Printf("%d [%s] %v\n", i, possibleNextNode.Name, nodeWasAlreadyVisited)

			if !nodeWasAlreadyVisited {
				nextPathToCheck := append(currentPath, possibleNextNode)
				if possibleNextNode.Name == "I" {
					names := make([]string, len(nextPathToCheck))
					for i, el := range nextPathToCheck {
						names[i] = el.Name
					}
					fmt.Println(strings.Join(names, " > "))
				}
				paths = append(paths, nextPathToCheck)
			}
		}

		i++
	}

	return paths
}

// Depth-first search for a path to the destination.
// Returns all visited nodes if a path was found or an empty slice if not
func (object *SpaceObject) findPathTo(destination *SpaceObject, visited []*SpaceObject) []*SpaceObject {
	visited = append(visited, object)

	if object == destination {
		return visited
	}

	objectsToCheck := object.Orbiting
	if object.Orbits != nil {
		objectsToCheck = append(objectsToCheck, object.Orbits)
	}

	for _, objectToCheck := range objectsToCheck {
		alreadyVisited := false
		for _, v := range visited {
			if v == objectToCheck {
				alreadyVisited = true
				break
			}
		}

		if alreadyVisited {
			continue
		}

		path := objectToCheck.findPathTo(destination, visited)
		if len(path) > 0 {
			return path
		}
	}

	return make([]*SpaceObject, 0)
}

func (object *SpaceObject) PathTo(destination *SpaceObject) []*SpaceObject {
	return object.findPathTo(destination, []*SpaceObject{})
	// paths := object.allPathsToAllOtherObjects()

	// for i, path := range paths {
	// 	names := make([]string, len(path))
	// 	for i, p := range path {
	// 		names[i] = p.Name
	// 	}
	// 	fmt.Printf("[%d] %s\n", i, strings.Join(names, " > "))
	// }
	// fmt.Println("=============================================")

	// // Filter out the possible paths to destination
	// pathsEndingInDesiredDestination := make([][]*SpaceObject, 0)
	// for _, path := range paths {
	// 	fmt.Printf("%v <> %v\n", destination.Name, path[len(path) - 1].Name)
	// 	if path[len(path) - 1] == destination {
	// 		pathsEndingInDesiredDestination = append(pathsEndingInDesiredDestination, path)
	// 	}
	// }

	// fmt.Printf("%v\n", pathsEndingInDesiredDestination)

	// // Get shortest path to destination from that
	// return make([]*SpaceObject, 1)
}

func main() {
	com := SpaceObject{Name: "COM"}
	objects := map[string]*SpaceObject{"COM": &com}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ")")

		// Ensure that both exist
		if _, exists := objects[parts[0]]; !exists {
			objects[parts[0]] = &SpaceObject{Name: parts[0], Orbits: nil}
		}
		if _, exists := objects[parts[1]]; !exists {
			objects[parts[1]] = &SpaceObject{Name: parts[1], Orbits: nil}
		}

		// Add relationship
		orbited := objects[parts[0]]
		orbiting := objects[parts[1]]
		orbited.Orbiting = append(orbited.Orbiting, orbiting)
		orbiting.Orbits = orbited
		objects[parts[1]] = orbiting
	}

	you := objects["YOU"]
	santa := objects["SAN"]
	path := you.PathTo(santa)

	// Path is end-to-end but we only need to change from our orbit to
	// the same orbit as santa. Therefor 2 are subtracted for source and
	// destination. Also: We are already orbiting the first planet, we must
	// not transfer to it. Therefor another 1 is subtracted.
	orbitalTransfers := len(path) - 3
	fmt.Print(orbitalTransfers)
}
