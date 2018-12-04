package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

/**
Problem: Have to close the goroutines from the previous questions, not just leave them hanging!
- Pretty sure the done goroutine is not firing
 */

func main(){
	fileName := flag.String("test-file",  "problems.csv", "The name of the test file")
	timer := flag.Int("question-timer", 10, "How long to allow for each question")
	flag.Parse()

	// get file name from flags
	csv, err := openTestFile(*fileName)

	if err != nil {
		fmt.Printf("error opening file %s: %v", *fileName, err)
		return
	}

	runQuiz(csv, *timer)
}

func runQuiz(csv *csv.Reader, timer int) {

	var totalCorrect, totalQuestions int

	// if we receive on the timer channel, cancel
	timerChan := make(chan struct{})
	answerChan := make(chan string)
	done := make(chan struct{})

	// skip 0th line
	csv.Read()
	for {
		question, answer, err := readLine(csv)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				// skip badly formed questions
				fmt.Println(err)
				continue
			}
		}
		totalQuestions ++

		fmt.Printf("Q%d: %s\n", totalQuestions, question)

		go countdown(timer, timerChan, done)
		go getInput(answerChan, done)

		// select between the timer countdown ending or the user entering input
		select {
			case <- timerChan:
				continue
			case input := <- answerChan:
				isCorrect := compareAnswers(answer, input)
				showAnswer(isCorrect)

				if isCorrect {
					totalCorrect ++
				}
		}

	}

	fmt.Printf("\n\nQuiz Finished\n")
	fmt.Printf("Your score: %d/%d\n", totalCorrect, totalQuestions)
}

func countdown(timer int, timerChan chan struct{}, doneChan chan struct{}){
	ticker := time.NewTicker(1 * time.Second)
	for i:=timer; i >= 0; i-- {
		select {
			case <- ticker.C:
				continue
			case <- doneChan:
				ticker.Stop()
				return
		}
	}

	timerChan <- struct{}{}
	ticker.Stop()
}

func readLine(csv *csv.Reader) (string, string, error) {
	line, err := csv.Read()
	if err != nil {
		return "", "", err
	}

	return line[0], line[1], nil
}

func showAnswer(isCorrect bool) {
	if isCorrect {
		fmt.Println("Correct")
	} else {
		fmt.Println("Incorrect")
	}
}

// read in csv file from default problems.csv
func openTestFile(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(file))
	return reader, nil
}

func getInput(answerChan chan string, doneChan chan struct{}) {
	for {
		select {
		case <-doneChan:
			return
		default:
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			answerChan <- input
		}
	}
}

func compareAnswers(input string, answer string) bool {
	inputNormalised := strings.ToLower(strings.TrimSpace(input))
	answerNormalised := strings.ToLower(strings.TrimSpace(answer))
	return inputNormalised == answerNormalised
}