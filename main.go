package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Question struct {
	Question string
	Answer   string
}

type Quiz struct {
	Questions chan Question
	Counter   int
	Done      chan bool
	Timer     *time.Timer
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

	fmt.Printf("Quiz will take %v seconds\n", *timeLimit)
	quiz.Timer = time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemloop:
	for {
		select {
		case <-quiz.Timer.C:
			fmt.Println("\nTime's up!")
			fmt.Println("Total Score:", quiz.Counter)
			break problemloop
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
		Questions: make(chan Question, 3),
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
				break
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
	close(quiz.Questions)
}

func (quiz *Quiz) askQuestion(q Question) bool {
	fmt.Printf("Q >> %v\n", q.Question)
	fmt.Print("Your Answer: ")

	answerReceived := make(chan string, 1)
	go func() {
		var input string
		fmt.Scanln(&input)
		answerReceived <- input
	}()

	select {
	case <-quiz.Timer.C:
		quiz.Done <- true
		fmt.Println("Interupt!")
		return false
	case input := <-answerReceived:
		if input == q.Answer {
			fmt.Println("Your Answer is correct")
			return true
		}
		fmt.Println("Your Answer is wrong")
		return false
	}
}
