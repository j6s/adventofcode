/*
 * --- Day 6: Universal Orbit Map ---
 *
 * You've landed at the Universal Orbit Map facility on Mercury. Because navigation in space often involves transferring
 * between orbits, the orbit maps here are useful for finding efficient routes between, for example, you and Santa.
 * You download a map of the local orbits (your puzzle input).
 *
 * Except for the universal Center of Mass (COM), every object in space is in orbit around exactly one other object.
 * An orbit looks roughly like this:
 *
 *                   \
 *                    \
 *                     |
 *                     |
 * AAA--> o            o <--BBB
 *                     |
 *                     |
 *                    /
 *                   /
 *
 * In this diagram, the object BBB is in orbit around AAA. The path that BBB takes around AAA (drawn with lines) is only
 * partly shown. In the map data, this orbital relationship is written AAA)BBB, which means "BBB is in orbit around AAA".
 *
 * Before you use your map data to plot a course, you need to make sure it wasn't corrupted during the download.
 * To verify maps, the Universal Orbit Map facility uses orbit count checksums - the total number of direct orbits
 * (like the one shown above) and indirect orbits.
 *
 * Whenever A orbits B and B orbits C, then A indirectly orbits C. This chain can be any number of objects long:
 * if A orbits B, B orbits C, and C orbits D, then A indirectly orbits D.
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
 *
 * Visually, the above map of orbits looks like this:
 *
 *         G - H       J - K - L
 *        /           /
 * COM - B - C - D - E - F
 *                \
 *                 I
 *
 * In this visual representation, when two objects are connected by a line, the one on the right directly orbits the one on the left.
 *
 * Here, we can count the total number of orbits as follows:
 *
 *     D directly orbits C and indirectly orbits B and COM, a total of 3 orbits.
 *     L directly orbits K and indirectly orbits J, E, D, C, B, and COM, a total of 7 orbits.
 *     COM orbits nothing.
 *
 * The total number of direct and indirect orbits in this example is 42.
 *
 * What is the total number of direct and indirect orbits in your map data?
 *
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

func (object *SpaceObject) TotalNumberOfOrbitedObjects() int {
	orbitedObjects := 0
	if object.Orbits != nil {
		orbitedObjects += 1 + object.Orbits.TotalNumberOfOrbitedObjects()
	}
	return orbitedObjects
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

	orbits := 0
	for _, object := range objects {
		orbits += object.TotalNumberOfOrbitedObjects()
		fmt.Printf("%s: %d \n", object.Name, object.TotalNumberOfOrbitedObjects())
	}
	fmt.Print(orbits)
}
