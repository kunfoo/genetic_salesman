# genetic_salesman
A go implementation of the Travelling Salesman Problem using a Genetic
Algorithm.

I wrote this for the Parallel and Distributed Computing class at Hochschule
Darmstadt. You can use graphs from the TSPLIB project
(http://elib.zib.de/pub/Packages/mp-testdata/tsp/tsplib/tsp/index.html), more
precisely symmetric graphs of the "EUC_2D" (euclidian distance) type.

You can tune the genetic algorithm a little at runtime. Here are the available
command line options:
&nbsp;&nbsp;-e    use elitism (default true)
&nbsp;&nbsp;-i    initialize population with nearest neighbor algorithm
&nbsp;&nbsp;-l uint
&nbsp;&nbsp;&nbsp;&nbsp;pathlen lower bound
&nbsp;&nbsp;-m uint
&nbsp;&nbsp;&nbsp;&nbsp;mutation rate in percent (default 2)
&nbsp;&nbsp;-n uint
&nbsp;&nbsp;&nbsp;&nbsp;number of goroutines (default 1)
&nbsp;&nbsp;-p uint
&nbsp;&nbsp;&nbsp;&nbsp;population size (default 20)
&nbsp;&nbsp;-s uint
&nbsp;&nbsp;&nbsp;&nbsp;number of individuals in natural selection (default 4)
&nbsp;&nbsp;-t duration
&nbsp;&nbsp;&nbsp;&nbsp;timeout in seconds
