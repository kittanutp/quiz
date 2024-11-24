package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Question struct {
	Question string
	Answer   string
}

type Quiz struct {
	Questions   chan Question
	Counter     int
	ProcessDone chan struct{}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	quiz := NewQuiz()

	go quiz.createQuizFromCSV(*csvFilename)

	var i string
	fmt.Println("Press Any Key to start quiz")
	fmt.Scanln(&i)

	quiz.start(*timeLimit)

	for {
		select {
		case <-quiz.ProcessDone:
			fmt.Println("\nTime's up!")
			fmt.Println("Total Score:", quiz.Counter)
			return
		case q, ok := <-quiz.Questions:

			if !ok {
				fmt.Println("Quiz completed!")
				fmt.Println("Total Score:", quiz.Counter)
				return
			}
			if quiz.askQuestion(q) {
				quiz.Counter++
			}
		}
	}
}

func NewQuiz() *Quiz {
	return &Quiz{
		Questions:   make(chan Question, 5),
		ProcessDone: make(chan struct{}),
	}
}

func (quiz *Quiz) createQuizFromCSV(csvFilename string) {
	file, err := os.Open(csvFilename)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				close(quiz.Questions)
				return
			}
			log.Fatal("Error while reading the file", err)
		}
		if len(record) != 2 {
			continue
		}
		quiz.Questions <- Question{
			Question: record[0],
			Answer:   record[1],
		}
	}
}

func (quiz *Quiz) start(defaultTime int) {
	fmt.Printf("Quiz will take %v seconds\n", defaultTime)
	go func() {
		time.Sleep(time.Duration(defaultTime) * time.Second)
		close(quiz.ProcessDone)
	}()
}

func (quiz *Quiz) askQuestion(q Question) bool {
	fmt.Printf("Problem: %s = ", q.Question)
	answerCh := make(chan string)

	go func() {
		var answer string
		fmt.Scanf("%s\n", &answer)
		answerCh <- answer
	}()

	select {
	case <-quiz.ProcessDone:
		return false
	case input := <-answerCh:
		// Compare sanitized input with the correct answer
		input = strings.TrimSpace(strings.ToLower(input))
		return input == strings.ToLower(q.Answer)
	}
}
