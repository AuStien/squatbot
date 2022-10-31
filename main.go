package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	squats, err := Squats(time.Now().Day())
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().Unix())
	if squats == 0 {
		fmt.Printf(restMessages[rand.Intn(len(restMessages))])
	} else {
		fmt.Printf(squatMessages[rand.Intn(len(squatMessages))], squats)
	}
}
