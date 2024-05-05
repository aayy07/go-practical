package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Person struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

var db *sql.DB

func main() {

	var err error
	// PostgreSQL connection parameters
	const (
		host     = "satao.db.elephantsql.com"
		port     = 5432
		user     = "cgswxlgr"
		password = "2ORg0JFnH7fKSxtavG-egu8sQVUK7xm4"
		dbname   = "cgswxlgr"
	)

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the PostgreSQL database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")

	router := gin.Default()

	router.GET("/person/:person_id/info", getPersonInfoHandler)

	router.GET("/person/create", createPerson)

	router.Run(":8080")
}

//GET API: TASK 1

func getPersonInfoHandler(c *gin.Context) {

	personID := c.Param("person_id")
	person, err := getPersonByID(personID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}
		panic(err)
	}
	c.JSON(http.StatusOK, person)
}

func getPersonByID(personID string) (Person, error) {
	person := Person{}

	query := `select
	p.name as name,
	ph."number" as phone_number,
	ad.city as city,
	ad.state as state,
	ad.street1 as street1,
	ad.street2 as street2,
	ad.zip_code as zip_code
	from
	address_join aj
	join person p on
	aj.person_id = p.id
	join phone ph on
	p.id = ph.person_id
	join address ad on
	ad.id = aj.address_id
	where p.id = $1;`

	row := db.QueryRow(query, personID)
	err := row.Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)
	return person, err
}

//POST API: TASK 2

func createPerson(c *gin.Context) {
	var query string
	var person Person
	var err error
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query = `INSERT INTO person(id, name)
	VALUES($1, $2)`
	_, err = db.Exec(query, rand.Intn(101), person.Name)
	if err != nil {
		log.Println("*****1", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}

	query = `INSERT INTO phone("id", "number")
	VALUES($1, $2)`
	_, err = db.Exec(query, rand.Intn(101), person.PhoneNumber)
	if err != nil {
		log.Println("*****2", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}

	query = `INSERT INTO address(id, city, state, street1, street2, zip_code)
	VALUES($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, rand.Intn(101), person.City, person.State, person.Street1, person.Street2, person.ZipCode)
	if err != nil {
		log.Println("*****3", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}
