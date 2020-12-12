package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

var traceName string
var cpuProfName string

var traceFile *os.File
var cpuProfFile *os.File

var loopStart time.Time
var loopTime time.Duration

func main() {
	var N int
	var seed int64
	var printSome int
	var runSerial bool
	var numCPU int

	flag.IntVar(&N, "N", 1000000, "number of work items")
	flag.Int64Var(&seed, "seed", 1, "random generator seed")
	flag.IntVar(&printSome, "printsome", 0, "print first num values of result to verify")
	flag.BoolVar(&runSerial, "serial", false, "run non-parallelized for loop instead")
	flag.IntVar(&numCPU, "numgr", runtime.NumCPU(), "number of goroutines to use in parallel loop")
	flag.StringVar(&traceName, "trace", "", "output trace of loop to file")
	flag.StringVar(&cpuProfName, "cpuprofile", "", "output CPU profile of loop to file")
	flag.Parse()

	// initialize input array of N values
	rand.Seed(seed)
	inputArray := make([]float64, N)
	for i := 0; i < N; i++ {
		inputArray[i] = 10 * (rand.Float64() - 0.5) // -5 to 5
	}

	// allocate output array
	outputArray := make([]float64, N)

	// execute loop
	preLoop()
	if runSerial {
		for i := 0; i < N; i++ {
			outputArray[i] = sinc(inputArray[i] * math.Pi)
		}
	} else {
		parallel.WithNumGoroutines(numCPU).For(N, func(i int) {
			outputArray[i] = sinc(inputArray[i] * math.Pi)
		})
	}
	postLoop()

	// print execution time
	fmt.Println(loopTime)

	// print some output values
	if printSome > 0 {
		fmt.Println("inputs:", inputArray[:printSome])
		fmt.Println("outputs:", outputArray[:printSome])
	}
}

func sinc(x float64) float64 {
	if x == 0.0 {
		return 1.0
	}
	return math.Sin(x) / x
}

func preLoop() {
	if traceName != "" {
		var err error
		traceFile, err = os.Create(traceName)

		if err != nil {
			log.Fatalln(err)
		}

		err = trace.Start(traceFile)

		if err != nil {
			log.Fatalln(err)
		}
	}

	if cpuProfName != "" {
		var err error
		cpuProfFile, err = os.Create(cpuProfName)

		if err != nil {
			log.Fatalln(err)
		}

		err = pprof.StartCPUProfile(cpuProfFile)

		if err != nil {
			log.Fatalln(err)
		}
	}

	loopStart = time.Now()
}

func postLoop() {
	loopStop := time.Now()
	loopTime = loopStop.Sub(loopStart)

	if traceName != "" {
		trace.Stop()
		err := traceFile.Close()

		if err != nil {
			log.Fatalln(err)
		}
	}

	if cpuProfName != "" {
		pprof.StopCPUProfile()
		err := cpuProfFile.Close()

		if err != nil {
			log.Fatalln(err)
		}
	}
}
