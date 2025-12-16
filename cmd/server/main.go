package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Employee struct {
	Name string // 
	Age int // 8 bytes
	Salary float64 // 8 bytes
}

func (e Employee) String() string {
	return fmt.Sprintf("Employee {%s, %d, %0.2f}", e.Name, e.Age, e.Salary)

}

func csvRowToEmployee(row []string) (*Employee, error) {
	name := row[0]
	ageStr := row[2]

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse age %s: %w", ageStr, err)
	}


	salaryStr := row[2]
	salary, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse salary %s: %w", salaryStr, err)
	}

	return &Employee{
		Name: name,
		Age: age,
		Salary: salary,
	}, nil
}

func loadEmployees(employeeFile *os.File) ([]*Employee, error) {
	reader := csv.NewReader(employeeFile)

	var employees []*Employee

	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read columns: %w", err)
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("reading row: %w", err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("expeted %d rows. got %d", 3, len(row))
		}

		employee, err := csvRowToEmployee(row)
		if err != nil {
			return nil, fmt.Errorf("failed to parse employee row: %v", row)
		}

		employees = append(employees, employee)
	}

	return employees, nil
}


var filename string

func main() {
	flag.StringVar(&filename, "file", "", "employee file name")

	var employeeFile *os.File

	if filename == "" {
		employeeFile = os.Stdin
	}

	employees, err := loadEmployees(employeeFile)
	if err != nil {
		log.Fatalln(err)
	}

	_ = employees

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to start tcp socket", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error accepting connection: %w", err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if err != nil {
		log.Printf("error reading client message: %s", err)
	}

	fmt.Printf("read %d bytes: %s\n", n, string(buf))
}