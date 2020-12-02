package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	go server()

	var input string
	fmt.Scanln(&input)
}

var materias map[string]map[string]float32

func server() {
	materias = make(map[string]map[string]float32)

	http.HandleFunc("/", root)
	http.HandleFunc("/calificar", CalificaAlumno)
	http.HandleFunc("/promedio/general", PromedioGeneral)
	http.HandleFunc("/promedio/alumno", PromedioAlumno)
	http.HandleFunc("/promedio/materia", PromedioMateria)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}

func root(res http.ResponseWriter, req *http.Request) {
	HTML("index.html", &res)
}

func CalificaAlumno(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		var errors []string
		materia, ok := req.PostForm["materia"]
		if !ok || materia[0] == "" {
			errors = append(errors, "Materia es requerida")
		}

		nombre, ok := req.PostForm["nombre"]
		if !ok || nombre[0] == "" {
			errors = append(errors, "Materia es requerida")
		}

		calif, ok := req.PostForm["calif"]
		if !ok || calif[0] == "" {
			errors = append(errors, "Materia es requerida")
		}

		califfloat, err := strconv.ParseFloat(calif[0], 32)
		if err != nil {
			errors = append(errors, "La calificacion no es un numero valido")
		} else if califfloat < 0 || califfloat > 100 {
			errors = append(errors, "La calificacion esta fuera de rango")
		}

		mat, ok := materias[materia[0]]
		if !ok {
			mat = make(map[string]float32)
			materias[materia[0]] = mat
		}

		cal, existe := mat[nombre[0]]
		if existe {
			errors = append(errors, "El alumno "+nombre[0]+" ya tiene una calificacion en "+materia[0]+": "+strconv.FormatFloat(float64(cal), 'f', 2, 32))
		} else {
			mat[nombre[0]] = float32(califfloat)
		}

		if len(errors) > 0 {

			alert := "<div class=\"alert alert-danger\" role=\"alert\">\n"
			for _, err := range errors {
				alert += "<p>" + err + "</p>\n"
			}
			alert += "</div>"

			HTML("form.html", &res, alert)
		} else {
			HTML("calificado.html", &res)
		}
	} else {
		HTML("form.html", &res, "")
	}
}

func PromedioGeneral(res http.ResponseWriter, req *http.Request) {
	matcount := 0
	var total float32 = 0
	for _, v := range materias {
		alumnos := 0
		var local float32 = 0
		for _, cal := range v {
			alumnos++
			local += cal
		}
		if alumnos > 0 {
			matcount++
			total += local / float32(alumnos)
		}
	}
	var msg string
	if matcount > 0 {
		total = total / float32(matcount)
		msg = "El promedio general es: " + strconv.FormatFloat(float64(total), 'f', 2, 32)
	} else {
		msg = "No hay materias para promediar <a href=\"/calificar\">califique estudiantes</a> para poder promediar."
	}

	HTML("promedio.html", &res, msg)
}

func PromedioAlumno(res http.ResponseWriter, req *http.Request) {
	if req.FormValue("nombre") == "" {
		HTML("promedio.alumno.html", &res, "", "", "")
		return
	}

	var suma float32 = 0
	matcount := 0
	alumno := req.FormValue("nombre")
	for _, v := range materias {
		cal, ok := v[alumno]
		if ok {
			matcount++
			suma += cal
		}
	}

	if matcount == 0 {
		HTML("promedio.alumno.html", &res, "is-invalid",
			"<div class=\"invalid-feedback\">\n"+alumno+" no se encuentra en ninguna clase.\n</div>", "")
	} else {
		total := suma / float32(matcount)
		HTML("promedio.alumno.html", &res, "", "", "<h2>El promedio del alumno es: "+strconv.FormatFloat(float64(total), 'f', 2, 32)+"</h2>")
	}
}

func PromedioMateria(res http.ResponseWriter, req *http.Request) {
	if req.FormValue("materia") == "" {
		HTML("promedio.materia.html", &res, "", "", "")
		return
	}

	materia := req.FormValue("materia")
	mat, existe := materias[materia]
	if !existe {
		HTML("promedio.materia.html", &res, "is-invalid",
			"<div class=\"invalid-feedback\">\n"+materia+" no existe.\n</div>", "")
		return
	}

	alumnos := 0
	var total float32 = 0
	for _, cal := range mat {
		alumnos++
		total += cal
	}

	if alumnos == 0 {
		HTML("promedio.alumno.html", &res, "is-invalid",
			"<div class=\"invalid-feedback\">\n"+materia+" no tiene alumnos para promediar.\n</div>", "")
	} else {
		total := total / float32(alumnos)
		HTML("promedio.materia.html", &res, "", "", "<h2>El promedio de la materia es: "+strconv.FormatFloat(float64(total), 'f', 2, 32)+"</h2>")
	}
}

func HTML(path string, res *http.ResponseWriter, msgs ...interface{}) {
	(*res).Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := ioutil.ReadFile(path)

	fmt.Fprintf(*res, string(html), msgs...)
}
