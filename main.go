package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	wait := rand.Intn(int(time.Hour) * 8)

	time.Sleep(time.Duration(wait))

	squats, err := Squats(time.Now().Day())
	if err != nil {
		panic(err)
	}

	if squats == 0 {
		fmt.Printf(restMessages[rand.Intn(len(restMessages))])
	} else {
		fmt.Printf(squatMessages[rand.Intn(len(squatMessages))], squats)
	}
}
