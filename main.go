package main

import (
	_ "embed"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed resources/problems.csv
var defaultCsv string

type score struct {
	total   int
	correct int
}

func (s *score) show() {
	fmt.Printf("You answered %d right out of %d\n", s.correct, s.total)
}

var quiz_file = flag.String("f", "", "the quiz file")
var timer_flag = flag.String("t", "30s", "quiz timer")
var shuffle = flag.Bool("s", false, "Shuffle questions")

func quiz(ch chan bool, records [][]string, score *score) {

	for _, record := range records {
		// A record is a slice of values in a row
		ques, ans := record[0], record[1]
		//ans_int, _ := strconv.Atoi(ans)
		//fmt.Println(ques, " ", ans)
		fmt.Printf("%s : ", ques)
		score.total += 1
		var got string
		fmt.Scanln(&got)
		got = strings.TrimSpace(got)
		if got == ans {
			score.correct += 1
		}
	}
	ch <- true
}

func main() {

	// The user can give the command line flag -f to give a custom csv quiz file, if not, then use the
	// default problems.csv file.
	flag.Parse() // parse the flags in command line
	//fmt.Println("The quiz file is", *quiz_file)

	timer, err := time.ParseDuration(*timer_flag)
	if err != nil {
		log.Fatal(err)
	}
	// Open the csv file
	var source io.Reader
	if *quiz_file == "" {
		//fmt.Println("Using the default problems file")
		source = strings.NewReader(defaultCsv)
	} else {
		source, err = os.Open(*quiz_file)
		if err != nil {
			log.Fatal(err)
		}
	}
	//fmt.Println("Opened file", file.Name())

	// Parse the csv file
	csv_file := csv.NewReader(source)
	records, err := csv_file.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Shuffle if flag is enabled
	if *shuffle == true {
		//fmt.Println("Gotta shuffle!")
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })
	}
	done := make(chan bool)
	score := score{}
	// take consent from user to start the quiz
	fmt.Println("Press Enter to start the quiz. The timer will start as soon as press Enter.")
	var start int
	fmt.Scanln(&start)
	go quiz(done, records, &score)
	select {
	case <-done:
		score.show()
	case <-time.After(timer):
		fmt.Printf("\nTime over!\n")
		score.show()
	}
}
