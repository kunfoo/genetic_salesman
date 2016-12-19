package main

import (
    "fmt"
    "math"
)

type City struct {
    Id uint64
    X, Y float64
}

// calculate distance of City c to City d
func (c City) Distance(d City) uint64 {
    x_square := math.Pow(c.X - d.X, 2)
    y_square := math.Pow(c.Y - d.Y, 2)
    return uint64(math.Sqrt(x_square + y_square))
}

type Tour struct {
    Cities []City
    Length uint64
}

// print a simple representation of Tour t and its length
func (t Tour) Print() {
    fmt.Printf("%d", t.Cities[0].Id)
    for _, city := range t.Cities[1:] {
        fmt.Printf(" => %d", city.Id)
    }
    fmt.Printf("\nTour length: %d\n", t.Length)
}

// verify that tour t contains all cities in city
func (t Tour) Verify(cities []City) bool {
    cityMap := make(map[uint64]City)
    for _, c := range cities[1:] {
        cityMap[c.Id] = c
    }

    length := uint64(0)
    first := t.Cities[0]
    last := t.Cities[len(t.Cities) - 1]
    if first != last {
        fmt.Println("Tour verification: first and last city in tour don't match")
        return false
    }
    if first != cities[0] {
        fmt.Printf("Tour verification: strange first city in tour (%d)\n", first.Id)
        return false
    }

    prevCity := first
    for _, c := range t.Cities[1:len(cities)] {
        if _, ok := cityMap[c.Id]; ok {
            delete(cityMap, c.Id)
            length += c.Distance(prevCity)
            prevCity = c
        } else {
            fmt.Printf("Tour verification: duplicate city %d in tour\n", c.Id)
            return false
        }
    }
    length += prevCity.Distance(first)

    if len(cityMap) > 0 {
        fmt.Printf("Tour verification: missing cities %v in tour\n", cityMap)
        return false
    }

    if length != t.Length {
        fmt.Printf("Tour verification: wrong length %d, should be %d\n", t.Length, length)
        return false
    }

    return true
}
