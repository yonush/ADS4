package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"ADS4/internal/models"

	"github.com/gocarina/gocsv"
)

func ReadCsv() []*models.CoursesCSV {
	csvFile, csvFileError := os.OpenFile("./data/courses.csv", os.O_RDONLY, os.ModePerm)
	// If an error occurs during os.OpenFIle, panic and halt execution.
	if csvFileError != nil {
		panic(csvFileError)
	}
	// Ensure the file is closed once the function returns
	defer csvFile.Close()

	var courses []*models.CoursesCSV
	// Parse the CSV data into the articles slice. If an error occurs, panic.
	if unmarshalError := gocsv.UnmarshalFile(csvFile, &courses); unmarshalError != nil {
		panic(unmarshalError)
	}

	return courses
}

func main() {
	file, err := os.Open("./data/courses.csv")

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()
	//exampl for writing CSV files
	/*   records := [][]string{
	     {"1", "John Doe", "john@email.com"},
	     {"2", "Jane Smith", "jane@email.com"},
	  	}*/
	/*
		w := csv.NewWriter(f)
		// Override default delimiter rune
		w.Comma = ‘;‘

		// Configure quoting behavior
		w.ForceQuote = false
		w.AlwaysQuote = false

		// Override newline terminator
		w.UseCRLF = true
		err := w.Write(record)
	*/
	reader := csv.NewReader(file)
	reader.Comma = ','   // Set delimiter
	reader.Comment = '#' // Ignore comments
	// Trim spaces around fields
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	// Raise errors instead of skipping bad records
	reader.LazyQuotes = true
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record[0], record[1], record[2])
	}

	courses := ReadCsv()
	for _, course := range courses {
		fmt.Printf("CourseCode: %s, Description: %s, Level: %d, Status: %s\n",
			course.CourseCode, course.Description, course.Level, course.Status)
	}
	/*
		records, err := reader.ReadAll()

		if err != nil {
			fmt.Println("Error reading records")
		}

		for _, eachrecord := range records {
			//fmt.Println(eachrecord)

			fmt.Println(eachrecord[0], eachrecord[1], eachrecord[2])
		}
	*/
}
