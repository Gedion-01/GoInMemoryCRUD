package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Gedion-01/Go-Crud-Challenge/db"
	"github.com/Gedion-01/Go-Crud-Challenge/types"
)

type Server struct {
	listener         net.Listener
	quit             chan struct{}
	exited           chan struct{}
	db               db.PersonStore
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
}

func NewServer(personStore db.PersonStore) *Server {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("failed to create listener", err.Error())
	}
	srv := &Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               personStore,
		connections:      map[int]net.Conn{},
		connCloseTimeout: 10 * time.Second,
	}
	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	var id int
	fmt.Println("listening for clients")
	for {
		select {
		case <-srv.quit:
			fmt.Println("shutting down the server")
			err := srv.listener.Close()
			if err != nil {
				fmt.Println("could not close listener", err.Error())
			}
			if len(srv.connections) > 0 {
				srv.warnConnections(srv.connCloseTimeout)
				<-time.After(srv.connCloseTimeout)
				srv.closeConnections()
			}
			close(srv.exited)
			return
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Println("failed to set listener deadline", err.Error())
			}

			conn, err := tcpListener.Accept()
			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("failed to accept connection", err.Error())
			}

			write(conn, "Welcome to MemoryDB server")
			srv.connections[id] = conn

			go func(connID int) {
				fmt.Println("client with id", connID, "joined")
				srv.handleConn(conn)
				delete(srv.connections, connID)

				fmt.Println("client with id", connID, "left")
			}(id)
			id++
		}
	}
}

func write(conn net.Conn, format string, args ...interface{}) {
	_, err := fmt.Fprintf(conn, format+"\n-> ", args...)
	if err != nil {
		log.Fatal(err)
	}
}

func writeAll(conn net.Conn, persons []types.Person) {
	for _, person := range persons {
		_, err := fmt.Fprintf(conn, "ID: %s\nName: %s\nAge: %s\nHobbies: %v\n\n", person.ID, person.Name, person.Age, person.Hobbies)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err := fmt.Fprintf(conn, "-> ")
	if err != nil {
		log.Fatal(err)
	}
}

func writePerson(conn net.Conn, person types.Person) {
	_, err := fmt.Fprintf(conn, "ID: %s\nName: %s\nAge: %s\nHobbies: %v\n-> ", person.ID, person.Name, person.Age, person.Hobbies)
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		switch {
		case len(values) == 4 && values[0] == "set":
			age := values[2]
			personParams := types.CreatePersonParams{
				Name:    values[1],
				Age:     age,
				Hobbies: strings.Split(values[3], ","),
			}
			validationErrors := personParams.Validate()
			if len(validationErrors) > 0 {
				var errorMessages []string
				for field, err := range validationErrors {
					errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, err))
				}
				write(conn, strings.Join(errorMessages, "\n"))
				continue
			}
			person := types.NewPersonFromParams(personParams)

			srv.db.Set(person)
			write(conn, "1 row inserted: \nID: %s, \nName: %s, \nAge: %s, \nHobbies: %v", person.ID, person.Name, person.Age, person.Hobbies)
		case len(values) == 2 && values[0] == "get":
			id := values[1]
			person, found := srv.db.Get(id)
			if !found {
				write(conn, fmt.Sprintf("Person with ID %s not found", id))
			} else {
				writePerson(conn, *person)
			}
		case len(values) >= 2 && values[0] == "update":
			id := values[1]
			updatePersonParams := &types.CreatePersonParams{}
			if len(values) > 2 {
				updatePersonParams.Name = values[2]
			}
			if len(values) > 3 {
				updatePersonParams.Age = values[3]
			}
			if len(values) > 4 {
				updatePersonParams.Hobbies = strings.Split(values[4], ",")
			}

			updatedPerson, status := srv.db.Update(id, updatePersonParams)
			if !status {
				write(conn, "Person with ID %s not found", id)
			} else {
				write(conn, "1 row updated: \nID: %s, \nName: %s, \nAge: %s, \nHobbies: %v", updatedPerson.ID, updatedPerson.Name, updatedPerson.Age, updatedPerson.Hobbies)
			}
		case len(values) == 2 && values[0] == "delete":
			status := srv.db.Delete(values[1])
			if !status {
				write(conn, fmt.Sprintf("Person with ID %s not found", values[1]))
			} else {
				write(conn, "1 row deleted")
			}
		case len(values) == 1 && values[0] == "all":
			persons := srv.db.All()
			if len(*persons) == 0 {
				write(conn, "No persons found")
			} else {
				writeAll(conn, *persons)
			}
		case len(values) == 1 && values[0] == "exit":
			if err := conn.Close(); err != nil {
				fmt.Println("could not close connection", err.Error())
			}
		default:
			write(conn, fmt.Sprintf("UNKNOWN:  %s", l))
		}
	}
}

func (srv *Server) warnConnections(timeout time.Duration) {
	for _, conn := range srv.connections {
		write(conn, fmt.Sprintf("host want to shutdown the server in %s", timeout.String()))
	}
}

func (srv *Server) closeConnections() {
	fmt.Println("closing all connections")
	for id, conn := range srv.connections {
		err := conn.Close()
		if err != nil {
			fmt.Println("could not close connection with id:", id)
		}
	}
}

func (srv *Server) Stop() {
	fmt.Println("stopping the database server")
	close(srv.quit)
	<-srv.exited
	fmt.Println("database server successfully stoped")
}
