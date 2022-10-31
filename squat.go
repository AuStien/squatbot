package main

import (
	"fmt"
)

var squatMap = map[int]int{
	1:  10,
	2:  15,
	3:  20,
	4:  25,
	5:  0,
	6:  30,
	7:  35,
	8:  40,
	9:  45,
	10: 0,
	11: 50,
	12: 55,
	13: 60,
	14: 65,
	15: 0,
	16: 70,
	17: 75,
	18: 80,
	19: 85,
	20: 0,
	21: 90,
	22: 95,
	23: 100,
	24: 105,
	25: 0,
	26: 110,
	27: 115,
	28: 120,
	29: 125,
	30: 130,
	31: 135,
}

func Squats(day int) (int, error) {
	for k, v := range squatMap {
		if k == day {
			return v, nil
		}
	}
	return 0, fmt.Errorf("day %q not in map", day)
}
