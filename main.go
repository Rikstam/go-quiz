package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the providedCSV file")
	}
	problems := parseLines(lines)
	// fmt.Println(problems)

	// initialise the timer after all the parsing is done
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	printQuestions(problems, &correct, *timer)
}

func printQuestions(problems []problem, correct *int, timer time.Timer) {
	for i, p := range problems {
		fmt.Printf("Problem #%d %s = \n", i+1, p.q)
		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d, out of %d.\n", *correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == p.a {
				*correct++
			}
		}
	}
	fmt.Printf("You scored %d, out of %d.\n", *correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
