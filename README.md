# genetic_salesman
A go implementation of the Travelling Salesman Problem using a Genetic
Algorithm.

I wrote this for the Parallel and Distributed Computing Class at Hochschule
Darmstadt. You can use graphs from the TSPLIB project
(http://elib.zib.de/pub/Packages/mp-testdata/tsp/tsplib/tsp/index.html), more
precisely symmetric graphs of the "EUC_2D" (euclidian distance) type.

You can tune the genetic algorithm a little at runtime. Here are the available
command line options:
  -e    use elitism (default true)
  -i    initialize population with nearest neighbor algorithm
  -l uint
        pathlen lower bound
  -m uint
        mutation rate in percent (default 2)
  -n uint
        number of goroutines (default 1)
  -p uint
        population size (default 20)
  -s uint
        number of individuals in natural selection (default 4)
  -t duration
        timeout in seconds
