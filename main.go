package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//Employee : struct
type Employee struct {
	ID             int    `json:"id"`
	Employeename   string `json:"employee_name"`
	Employeesalary int64  `json:"employee_salary"`
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
	//	router.HandleFunc("/employees", createEmployee).Methods("POST")
	//	router.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	//	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	//	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employees []Employee
	result, err := db.Query("SELECT id, employee_name, employee_salary,Employee_age from employee")
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
