package main

import (
    "os"
    "fmt"
    "math/rand"
    "reflect"
    "flag"
    "time"
    "runtime"
)

var elitism bool = true
var numRoutines uint = 8
var selectionSize uint = 7
var mutationRate uint = 3   // in percent
var initializeNN = false    // initialize population with nearest neighbour algorithm

type Population struct {
    Individuals []Tour
    Fittest Tour
}

/*
 * create an initial population of populationSize, either randomly or with the
 * nearest neighbor algorithm
 */
func newPopulation(cities []City, populationSize uint) Population {
    individuals := make([]Tour, populationSize)

    if initializeNN {
        individuals[0] = createNearestNeighborTour(cities)
        for idx := uint(1); idx < populationSize; idx++ {
            individuals[idx] = createRandomTour(cities)
        }
    } else {
        for idx := range individuals {
            individuals[idx] = createRandomTour(cities)
        }
    }

    return Population{individuals, getFittest(individuals)}
}

// find the fittest individual
func getFittest(individuals []Tour) Tour {
    fittest := individuals[0]
    for _, ind := range individuals[1:] {
        if fittest.Length > ind.Length {
            fittest = ind
        }
    }

    return fittest
}

// Choose some Tours out of Population p randomly, and return the fittest one
func naturalSelection(p Population) Tour {
    t := make([]Tour, selectionSize)

    for idx := range t {
        r := rand.Intn(len(p.Individuals))
        t[idx] = p.Individuals[r]
    }

    return getFittest(t)
}

// calculate tour length of cities
func getTourLength(cities []City) uint64 {
    curCity := cities[1]
    length := curCity.Distance(cities[0])
    for _, city := range cities[2:] {
        length += curCity.Distance(city)
        curCity = city
    }

    return length
}

// return a random index out of range cities, excluding the first and last one
func getRandomIndex(cities []City) int {
    return rand.Intn(len(cities) - 2) + 1
}

func mutate(t Tour) Tour {
    mutation := t.Cities[:]

    for idx1 := 1; idx1 < len(mutation) - 1; idx1++ {
        if r := uint(rand.Intn(100)); r < mutationRate {
            idx2 := getRandomIndex(mutation)
            mutation[idx1], mutation[idx2] = mutation[idx2], mutation[idx1]
        }
    }

    return Tour{mutation, getTourLength(mutation)}
}

/*
 * take a random subset of cities from Tour t1 and copy it into a new Tour, take
 * the remaining cities from Tour t2
 */
func crossover(t1, t2 Tour) Tour {
    cities := make([]City, len(t1.Cities))
    // the origin city must always stay the same
    cities[0] = t1.Cities[0]
    cities[len(cities) - 1] = t1.Cities[len(cities) - 1]

    startIdx := 0
    endIdx := 0
    for startIdx >= endIdx {
        startIdx = getRandomIndex(t1.Cities)
        endIdx = getRandomIndex(t2.Cities)
    }

    // remember which cities we already took from t1
    cityMap := make(map[uint64]bool)
    // copy cities in the range [startIdx:endIdx] from t1 into cities
    for idx := startIdx; idx <= endIdx; idx++ {
        nextCity := t1.Cities[idx]
        cities[idx] = nextCity
        cityMap[nextCity.Id] = true
    }

    // copy remaining cities from t2 to cities
    citiesIdx := 1
    t2Idx := 1
    maxIdx := len(cities) - 1
    for citiesIdx < maxIdx && t2Idx < maxIdx {
        // this is the range we already copied from t1
        if citiesIdx == startIdx {
            citiesIdx = endIdx + 1
        }
        nextCity := t2.Cities[t2Idx]
        if _, ok := cityMap[nextCity.Id]; ok {
            // city Id was already taken from t1
            t2Idx++
            continue
        }
        cities[citiesIdx] = nextCity
        citiesIdx++
        t2Idx++
    }

    return Tour{cities, getTourLength(cities)}
}

/* 
 * do natural selection from population, cross over fittest individuals and
 * mutate children
 */
func evolvePopulation(p Population) Population {
    populationSize := len(p.Individuals)
    individuals := make([]Tour, populationSize)

    offset := 0
    if elitism {
        // keep fittest individual from last generation
        individuals[0] = p.Fittest
        offset++
    }

    for idx := offset; idx < populationSize; idx++ {
        parent1 := naturalSelection(p)
        parent2 := naturalSelection(p)
        child := crossover(parent1, parent2)
        individuals[idx] = mutate(child)
    }

    return Population{individuals, getFittest(individuals)}
}

// this function is executed by the goroutines
func runner(p Population, quit chan bool, result chan<- Tour, fittest <-chan Tour) {
    for {
        select {
        case <- quit:
            close(result)
            quit <- true
            return
        case newFittest := <-fittest:
            if newFittest.Length < p.Fittest.Length {
                p.Fittest = newFittest
                p.Individuals[0] = newFittest
            }
        default:
            curFittest := p.Fittest
            p = evolvePopulation(p)
            if curFittest.Length > p.Fittest.Length {
                result <- p.Fittest
            }
        }
    }
}

func main() {
    _elitism        := flag.Bool("e", true, "use elitism")
    _selectionSize  := flag.Uint("s", 4, "number of individuals in natural selection")
    _mutationRate   := flag.Uint("m", 2, "mutation rate in percent")
    _initializeNN   := flag.Bool("i", false, "initialize population with nearest neighbor algorithm")
    populationSize  := flag.Uint("p", 20, "population size")
    numRoutines     := flag.Uint("n", 1, "number of goroutines")
    timeout         := flag.Duration("t", time.Duration(0), "timeout in seconds")
    lowerBound      := flag.Uint64("l", 0, "pathlen lower bound")
    flag.Parse()

    if flag.NArg() != 1 {
        fmt.Printf("Usage: %s [options] <tsp-file>\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }
    elitism       = *_elitism
    selectionSize = *_selectionSize
    mutationRate  = *_mutationRate
    initializeNN  = *_initializeNN

    cities, err := readTSPLIBData(flag.Arg(0))
    if nil != err {
        fmt.Println(err)
        os.Exit(1)
    }

    maxProcs := int(*numRoutines) + 1
    if maxProcs > runtime.NumCPU() {
        // do not start more more threads to run goroutines than logical CPUs on the system
        maxProcs = runtime.NumCPU()
    }
    runtime.GOMAXPROCS(maxProcs)
    rand.Seed(time.Now().Unix())
    p := newPopulation(cities, *populationSize)  // create an initial random population

    var quit []chan bool    // to signal finish
    var fittest []chan Tour // to send fittest individual to all runners
    var result []chan Tour  // to receive results from runners
    var selectSet []reflect.SelectCase  // this allows to select on an array of channels
    selectSet = append(selectSet, reflect.SelectCase{
            Dir: reflect.SelectDefault, // add a default case, so reflect.Select() does not block
    })
    for idx := uint(0); idx < *numRoutines; idx++ {
        quit = append(quit, make(chan bool, 1))
        fittest = append(fittest, make(chan Tour, 50))
        result = append(result, make(chan Tour, 50)) // buffered channel enhances performance
        selectSet = append(selectSet, reflect.SelectCase{
            Dir: reflect.SelectRecv,            // direction: receive
            Chan: reflect.ValueOf(result[idx]), // the channel to receive from
        })
        go runner(p, quit[idx], result[idx], fittest[idx])
    }

    if 0 == *lowerBound && timeout.Seconds() == 0 {
        nnt := createNearestNeighborTour(cities)
        *lowerBound = nnt.Length
    }

    startTime := time.Now()
    timeChan := time.Tick(*timeout)
    fmt.Printf("goal: %d (%g s deadline)\n", *lowerBound, timeout.Seconds())
    curFitness := p.Fittest.Length

    // main infinite loop
    MAIN:
    for {
        select {
        case <-timeChan:    // timeout
            break MAIN
        default:
            // select from array of result channels
            from, tmpValue, _ := reflect.Select(selectSet)
            if from != 0 { // default case
                value := tmpValue.Interface().(Tour)    // tmpValue is of type Value and not Tour
                if value.Length < curFitness {
                    curFitness = value.Length
                    if curFitness < *lowerBound {
                        break MAIN
                    }
                    fmt.Printf("routine %d: %d\n", from, curFitness)
                    for idx, ch := range fittest {
                        if idx != from {
                            ch <- value
                        }
                    }
                }
            }
        }
    }

    fmt.Println("timeout -> sending quit signal to all goroutines")
    for idx, ch := range quit {
        ch <- true      // signal quit to all runners
        <-result[idx]   // in case a runner is blocked on a full result channel
        <-ch            // wait for runner to return
    }

    fmt.Printf("final result: %d (worked %g s)\n", curFitness, time.Since(startTime).Seconds())
}
