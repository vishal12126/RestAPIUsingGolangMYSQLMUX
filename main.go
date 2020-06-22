package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//Employee : struct
type Employee struct {
	ID             int    `json:"id"`
	Employeename   string `json:"employee_name"`
	Employeesalary int    `json:"employee_salary"`
	Employeeage    int    `json:"employee_age"`
}

var db *sql.DB
var err error

func main() {
	fmt.Printf("hello world\n")
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/mygodb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

//Get All
func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employees []Employee
	result, err := db.Query("SELECT id, employee_name, employee_salary,employee_age from employee")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var Employee Employee
		err := result.Scan(&Employee.ID, &Employee.Employeename, &Employee.Employeesalary, &Employee.Employeeage)
		if err != nil {
			panic(err.Error())
		}
		employees = append(employees, Employee)
	}
	json.NewEncoder(w).Encode(employees)
}

//Get
func getEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, employee_name, employee_salary,employee_age from employee WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var Employee Employee
	for result.Next() {
		err := result.Scan(&Employee.ID, &Employee.Employeename, &Employee.Employeesalary, &Employee.Employeeage)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(Employee)
}

//Create
func createEmployee(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO employee(employee_name,employee_salary,employee_age) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	keyVal1 := make(map[string]int)
	json.Unmarshal(body, &keyVal)
	json.Unmarshal(body, &keyVal1)
	//id := keyVal["id"]
	name := keyVal["employee_name"]
	salary := keyVal1["employee_salary"]
	age := keyVal1["employee_age"]
	_, err = stmt.Exec(name, salary, age)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New employee was created")
}

//Update
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE employee SET employee_name = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println("in between")
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["employee_name"]
	_, err = stmt.Exec(newName, params["id"])
	log.Println("INSERT: newName: " + newName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Employee with ID = %s was updated", params["id"])
}

//Delete

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM employee WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "employee with ID = %s was deleted", params["id"])
}
