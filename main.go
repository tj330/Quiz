package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "quiz.csv", "a csv file which contains content in format 'question: answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exit(fmt.Sprintf("Error encountered while opening the file: %s\n", *csvFilename))
	}
	
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		exit("Error encountered while reading the file")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit)* time.Second)
	

	count := 0

	problemloop:
		for i, p := range problems {
			fmt.Printf("Problem %d: %s\n", i+1, p.q)
			answerCh := make(chan string)
			go func() {
				var answer string
				fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()
			select {
			case <-timer.C:
				fmt.Println()
				break problemloop
			case answer := <-answerCh:
				if answer == p.a {
					count++
				}
			}
		}

	fmt.Printf("Time is up!\nYou scored %d out of %d\n", count, len(problems))
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

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
