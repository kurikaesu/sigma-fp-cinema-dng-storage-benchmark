package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	location := os.Args[1]
	fps, err := strconv.Atoi(os.Args[2])

	if err != nil {
		log.Fatal("Second argument needs to be a whole number for FPS")
	}

	minDelay := int64(1000.0 / fps)
	dngSize := 12630528

	dngObj := make([]byte, dngSize)
	count, err := rand.Read(dngObj)

	if err != nil || count != dngSize {
		log.Fatal("Could not create a 12 bit DNG buffer")
	}

	fileCounter := 0
	filePrefix := "benchmark_"
	fileSuffix := ".dng"

	log.Printf("Starting to write files to (%s)", location)

	for {
		startTime := time.Now()
		fileName := fmt.Sprintf("%s/%s%d%s", location, filePrefix, fileCounter, fileSuffix)
		err = os.WriteFile(fileName, dngObj, 0644)
		if err != nil {
			log.Printf("Could not write (%s)", fileName)
			break
		}

		fileCounter += 1

		duration := time.Since(startTime)
		timeToDelay := minDelay - duration.Milliseconds()
		if timeToDelay > 0 {
			time.Sleep(time.Duration(timeToDelay * int64(time.Millisecond)))
		} else {
			log.Printf("Could not write files fast enough. Stopping at (%s)", fileName)
			break
		}
	}

	log.Printf("Cleaning up")

	for num := 0; num < fileCounter; num++ {
		fileName := fmt.Sprintf("%s/%s%d%s", location, filePrefix, num, fileSuffix)
		err = os.Remove(fileName)

		if err != nil {
			log.Printf("Could not clean up (%s). Skipping", fileName)
		}
	}

	writtenMiB := float32(dngSize*fileCounter) / 1024
	log.Printf("Benchmark stopped after writing at most (%f) MiB", writtenMiB)
	log.Printf("In GiB that would be (%f)", writtenMiB/1024)
}
