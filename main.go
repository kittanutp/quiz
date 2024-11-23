package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type Question struct {
	Question string
	Answer   string
}

type Quiz struct {
	Questions chan Question
	Counter   int
}

const path = "problems.csv"

func main() {
	quiz := NewQuiz()
	go quiz.createQuizFromCSV()
	fmt.Println("Start Quiz")
	var i string
	var resp string
	for q := range quiz.Questions {
		fmt.Printf("Q >> %v\n", q.Question)
		fmt.Print("Your Answer: ")
		fmt.Scanln(&i)
		if i == q.Answer {
			resp = "Correct"
			quiz.Counter += 1

		} else {
			resp = "Wrong"
		}
		fmt.Println("Your Answer is", resp)
	}
	fmt.Println("Total Score", quiz.Counter)
}

func NewQuiz() *Quiz {
	return &Quiz{
		Questions: make(chan Question, 3),
	}
}

func (quiz *Quiz) createQuizFromCSV() {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
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
	defer file.Close()
	defer close(quiz.Questions)
}
