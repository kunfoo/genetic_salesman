package main

import (
    "os"
    "fmt"
    "math"
    "bufio"
    "strconv"
    "errors"
    "strings"
    "math/rand"
)

/*
 * simple parser for data from the TSPLIB project and returns an array of City
 * only works with graphs of type TSP (symmetrical) and euclidean coordinates
 */
func readTSPLIBData(filename string) ([]City, error) {
    var cities []City

    file, err := os.Open(filename)
    if nil != err {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    SCAN_HEADER:
    for scanner.Scan() {
        tokens := strings.Split(scanner.Text(), ":")
        for i := range(tokens) {
            tokens[i] = strings.TrimSpace(tokens[i])
        }
        // sanity check some header fields
        switch tokens[0] {
        case "TYPE":
            if "TSP" != tokens[1] {
                s := fmt.Sprintf("Error: unsupported tsp TYPE %s", tokens[1])
                return nil, errors.New(s)
            }
        case "EDGE_WEIGHT_TYPE":
            if "EUC_2D" != tokens[1] {
                s := fmt.Sprintf("Error: unsupported tsp EDGE_WEIGHT_TYPE %s", tokens[1])
                return nil, errors.New(s)
            }
        case "DIMENSION":
            dim, err := strconv.Atoi(tokens[1])
            if nil != err {
                s := fmt.Sprintf("Error parsing DIMENSION %s", tokens[1])
                return nil, errors.New(s)
            }
            // allocate memory for all cities
            cities = make([]City, dim)
        case "NODE_COORD_SECTION":
            break SCAN_HEADER
        }
    }

    for scanner.Scan() {
        tokens := strings.Fields(scanner.Text())
        if "EOF" == tokens[0] {
            break
        }
        id, err := strconv.ParseUint(tokens[0], 10, 64)
        if nil != err {
            s := fmt.Sprintf("Error parsing line %v", tokens)
            return nil, errors.New(s)
        }
        x, err := strconv.ParseFloat(tokens[1], 64)
        if nil != err {
            s := fmt.Sprintf("Error parsing x coordinate from line %v", tokens)
            return nil, errors.New(s)
        }
        y, err := strconv.ParseFloat(tokens[2], 64)
        if nil != err {
            s := fmt.Sprintf("Error parsing y coordinate from line %v", tokens)
            return nil, errors.New(s)
        }
        cities[id - 1] = City{id, x, y}
    }

    return cities, nil
}

// creates a randomly ordered Tour
func createRandomTour(cities []City) Tour {
    randomCities := make([]City, len(cities) + 1)
    randomOrder := rand.Perm(len(cities) - 1)

    first := cities[0]
    randomCities[0] = first
    var length uint64
    idx := 1
    for _, randomCityIdx := range randomOrder {
        nextCity := cities[randomCityIdx + 1]
        randomCities[idx] = nextCity
        length += nextCity.Distance(randomCities[idx - 1])
        idx++
    }
    randomCities[idx] = first
    length += first.Distance(randomCities[idx - 1])

    return Tour{randomCities, length}
}

// creates a tour according to the nearest neighbor algorithm
func createNearestNeighborTour(cities []City) Tour {
    cityMap := make(map[uint64]City)
    for _, city := range cities[1:] {
        cityMap[city.Id] = city
    }

    var t Tour
    firstCity := cities[0]
    curCity := firstCity
    t.Cities = append(t.Cities, curCity)
    nearestNeighbor := cities[1]
    minDistance := curCity.Distance(nearestNeighbor)
    for len(cityMap) > 0 {
        for _, city := range cityMap {
            curDistance := curCity.Distance(city)
            if curDistance == minDistance {
                /*
                 * use city with highest id as next city, else this function
                 * behaves differently on every call, because range on a
                 * map iterates differently
                 */
                if city.Id > nearestNeighbor.Id {
                    nearestNeighbor = city
                }
            } else if curDistance < minDistance {
                nearestNeighbor = city
                minDistance = curDistance
            }
        }
        curCity = nearestNeighbor
        delete(cityMap, curCity.Id)
        t.Cities = append(t.Cities, curCity)
        t.Length += minDistance
        minDistance = math.MaxUint64
    }
    t.Cities = append(t.Cities, firstCity)
    t.Length += nearestNeighbor.Distance(firstCity)

    return t
}
