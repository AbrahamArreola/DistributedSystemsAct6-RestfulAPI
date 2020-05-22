package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var students = make(map[string]map[string]float32)

//Student is...
type Student struct {
	Student string
	Subject string
	Score   float32
}

func addScore(student Student) []byte {
	var responseString string

	if _, exists := students[student.Student]; !exists {
		students[student.Student] = make(map[string]float32)
	}

	if _, exists := students[student.Student][student.Subject]; !exists {
		students[student.Student][student.Subject] = float32(student.Score)
		responseString = `{"code": "Calificación agregada"}`
	} else {
		responseString = `{"code": "Error: Calificación existente"}`
	}

	return []byte(responseString)
}

func getStudentscores() ([]byte, error) {
	jsonData, err := json.MarshalIndent(students, "", "    ")
	if err != nil {
		return jsonData, nil
	}
	return jsonData, err
}

func getStudentScore(studentName string) ([]byte, error) {
	jsonData := []byte(`{}`)

	if student, exists := students[studentName]; exists {
		jsonData, err := json.MarshalIndent(student, "", "    ")
		if err != nil {
			return jsonData, nil
		}
		return jsonData, err
	}
	return jsonData, nil
}

func deleteStudent(student string) []byte {
	if _, exists := students[student]; exists {
		delete(students, student)
		return []byte(`{"code": "Estudiante eliminado"}`)
	}
	return []byte(`{"code": "Error: Estudiante inexistente"}`)
}

func updateScore(student Student) []byte {
	if _, exists := students[student.Student][student.Subject]; exists {
		students[student.Student][student.Subject] = float32(student.Score)
		return []byte(`{"code": "Calificación modificada"}`)
	}
	return []byte(`{"code": "Error: alumno o materia inexistente"}`)
}

func student(response http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Method)
	switch request.Method {
	case "POST":
		var student Student
		err := json.NewDecoder(request.Body).Decode(&student)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON := addScore(student)
		response.Header().Set(
			"Content-Type",
			"application/json",
		)
		response.Write(responseJSON)

	case "GET":
		responseJSON, err := getStudentscores()
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Header().Set(
			"Content-Type",
			"application/json",
		)
		response.Write(responseJSON)

	case "PUT":
		var student Student
		err := json.NewDecoder(request.Body).Decode(&student)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON := updateScore(student)
		response.Header().Set(
			"Content-Type",
			"application/json",
		)
		response.Write(responseJSON)
	}
}

func studentID(response http.ResponseWriter, request *http.Request) {
	student := strings.TrimPrefix(request.URL.Path, "/student/")
	fmt.Println(request.Method, ":", student)
	switch request.Method {
	case "GET":
		responseJSON, err := getStudentScore(student)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
		response.Header().Set(
			"Content-Type",
			"application/json",
		)
		response.Write(responseJSON)

	case "DELETE":
		responseJSON := deleteStudent(student)
		response.Header().Set(
			"Content-Type",
			"application/json",
		)
		response.Write(responseJSON)
	}
}

func main() {
	http.HandleFunc("/student", student)
	http.HandleFunc("/student/", studentID)
	fmt.Println("Corriendo RESTful API...")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
