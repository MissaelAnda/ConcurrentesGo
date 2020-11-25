package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

func main() {
	go server()

	var input string
	fmt.Scanln(&input)
}

type Calificacion struct {
	Materia string
	Alumno  string
	Calif   float32
}

type Server struct {
	materias map[string]map[string]float32
}

func (this *Server) CalificaAlumno(datos Calificacion, reply *string) error {
	if this.materias == nil {
		this.materias = make(map[string]map[string]float32)
	}
	if datos.Calif < 0.0 || datos.Calif > 100.0 {
		return errors.New("Calificacion invalida.")
	}

	mat, existe := this.materias[datos.Materia]
	if !existe {
		mat = make(map[string]float32)
		this.materias[datos.Materia] = mat
	}

	cal, existe := mat[datos.Alumno]
	if existe {
		return errors.New("El alumno " + datos.Alumno + " ya tiene una calificacion en " + datos.Materia + ": " + strconv.FormatFloat(float64(cal), 'f', 2, 32))
	}
	mat[datos.Alumno] = datos.Calif
	*reply = "Calificacion asignada."
	return nil
}

func (this *Server) PromedioAlumno(alumno string, reply *float32) error {
	var suma float32 = 0
	materias := 0
	for _, v := range this.materias {
		cal, ok := v[alumno]
		if ok {
			materias++
			suma += cal
		}
	}

	if materias == 0 {
		return errors.New(alumno + " no se encuentra en ninguna clase.")
	}

	*reply = suma / float32(materias)
	return nil
}

func (this *Server) PromedioGeneral(_ int32, reply *float32) error {
	materias := 0
	var total float32 = 0
	for _, v := range this.materias {
		alumnos := 0
		var local float32 = 0
		for _, cal := range v {
			alumnos++
			local += cal
		}
		materias++
		if alumnos > 0 {
			total += local / float32(alumnos)
		}
	}
	if materias > 0 {
		*reply = total / float32(materias)
		return nil
	} else {
		return errors.New("No hay nada para promediar.")
	}
}

func (this *Server) PromedioMateria(materia string, reply *float32) error {
	mat, existe := this.materias[materia]
	if !existe {
		return errors.New("La materia " + materia + " no existe.")
	}
	alumnos := 0
	var total float32 = 0
	for _, cal := range mat {
		alumnos++
		total += cal
	}
	if alumnos > 0 {
		*reply = total / float32(alumnos)
		return nil
	} else {
		return errors.New("No hay alumnos para promediar.")
	}
}

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}
