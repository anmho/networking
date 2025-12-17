package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/netutil"
)




var filename string

func getEmployeesFile(filename string) (*os.File, error) {
	var employeeFile *os.File

	if filename == "" {
		employeeFile = os.Stdin
	}

	return employeeFile, nil
}

type Server struct {
	employees []*Employee
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		log.Printf("error reading client message: %s", err)
	}

	fmt.Printf("read %d bytes: %s\n", n, string(buf))
}

func main() {
	flag.StringVar(&filename, "file", "", "employee file name")

	empFile, err := getEmployeesFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	employees, err := loadEmployees(empFile)
	if err != nil {
		log.Fatalln(err)
	}

	s := &Server{
		employees: employees,
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to start tcp socket", err)
	}


	// blocks others
	l = netutil.LimitListener(l, 1000)
	// right now a client could decide not to close a connection and it would hold it forever
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s", err)
		}
		go s.Handle(conn)
	}
}