package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

const gradesFilename = "grades.csv"
const studentDataFieldCount = 7

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

/* Could also have used bufio.Scanner */
func parseCSV(filePath string) []student {
	students := []student{}

	/* Open file, grab file descriptor */
	fd, err := os.Open(gradesFilename)
	defer fd.Close()
	if err != nil {
		log.Fatalf("error opening grades csv file: %v", err)
	}

	/* Initialize a new reader */
	reader := bufio.NewReader(fd)

	/* Discard title line, assuming it always exists */
	_, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading student data title in grades csv file: %v", err)
	}

	/* Iterate line-by-line */
	rowCount := 0
	endOfFile := false
	for endOfFile == false {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			endOfFile = true
		}

		/* Removing newline from the end and splitting into fields */
		if s[len(s)-1] == '\n' {
			s = s[:len(s)-1]
		}
		fields := strings.Split(s, ",")

		/* Create student */
		parsedStudent, err := parseStudent(fields)
		if err != nil {
			log.Fatalf("error in row %d: %v", rowCount, err)
		}

		students = append(students, parsedStudent)
		rowCount++
	}

	return students
}

func calculateGrade(students []student) []studentStat {
	/* Initialize slice and a struct to easily match grades with scores */
	studentStats := []studentStat{}
	gradeBoundaries := []struct {
		grade    Grade
		boundary float32
	}{
		{A, 70},
		{B, 50},
		{C, 35},
		{F, 0},
	}

	/* Iterate over each student */
	for _, student := range students {
		var finalGrade Grade
		finalScore := float32(student.test1Score+student.test2Score+student.test3Score+student.test4Score) / 4
		for _, gradeBoundary := range gradeBoundaries {
			grade, boundary := gradeBoundary.grade, gradeBoundary.boundary
			if finalScore >= boundary {
				finalGrade = grade
				break
			}
		}
		studentStats = append(studentStats, studentStat{student: student, finalScore: finalScore, grade: finalGrade})
	}

	return studentStats
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	topper := studentStat{}
	for _, cur := range gradedStudents {
		if cur.finalScore > topper.finalScore {
			topper = cur
		}
	}
	return topper
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	universityToppers := map[string]studentStat{}

	for _, cur := range gs {
		university := cur.university
		universityTopper, exists := universityToppers[university]

		/* If curStudent has greater score than current university topper OR if this university has been seen for the first time */
		if (exists && cur.finalScore > universityTopper.finalScore) || !exists {
			universityToppers[university] = cur
		}
	}
	return universityToppers
}

/**** Helpers ****/

func validateStudentData(fields []string) error {
	/* Add list of validations */
	if len(fields) != studentDataFieldCount {
		return fmt.Errorf("invalid number of fields in student data count; wanted: %d, got: %d", studentDataFieldCount, len(fields))
	}

	return nil
}

func parseStudent(fields []string) (student, error) {
	/* Validation */
	err := validateStudentData(fields)
	if err != nil {
		return student{}, fmt.Errorf("%v: [data: %v]", err, fields)
	}

	/* Parse fields */
	firstName, lastName, university := fields[0], fields[1], fields[2]
	testScores := []int{}
	for testScoreIdx := 3; testScoreIdx < 7; testScoreIdx++ {
		testScore, err := strconv.Atoi(fields[testScoreIdx])
		if err != nil {
			return student{}, fmt.Errorf("error parsing test scores from csv %v", err)
		}
		testScores = append(testScores, testScore)
	}

	/* Create student */
	parsedStudent := student{firstName: firstName, lastName: lastName, university: university, test1Score: testScores[0], test2Score: testScores[1], test3Score: testScores[2], test4Score: testScores[3]}
	return parsedStudent, nil
}
