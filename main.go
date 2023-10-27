package main

import (
	//"fmt"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

// Funcion para la conexion de la base de datos
func ConexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	Name := "sistema"

	conexion, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/Crear", Crear)
	http.HandleFunc("/Insertar", Insertar)
	http.HandleFunc("/Eliminar", Eliminar)
	http.HandleFunc("/Editar", Editar)
	http.HandleFunc("/Actualizar", Actualizar)

	log.Println("Ejecutando servidor...")
	http.ListenAndServe(":3030", nil)
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	ConexionExitosa := ConexionBD()
	Registros, err := ConexionExitosa.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arrayEmpleado := []Empleado{}

	for Registros.Next() {
		var id int
		var nombre, correo string
		err = Registros.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arrayEmpleado = append(arrayEmpleado, empleado)
	}

	plantillas.ExecuteTemplate(w, "Inicio", arrayEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "Crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		ConexionExitosa := ConexionBD()
		RegistroInsertado, err := ConexionExitosa.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		RegistroInsertado.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	ConexionExitosa := ConexionBD()
	RegistroEditado, err := ConexionExitosa.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for RegistroEditado.Next() {
		var id int
		var nombre, correo string
		err = RegistroEditado.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "Editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		ConexionExitosa := ConexionBD()
		RegistroActualizado, err := ConexionExitosa.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		RegistroActualizado.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}

func Eliminar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	ConexionExitosa := ConexionBD()
	RegistroEliminado, err := ConexionExitosa.Prepare("DELETE FROM empleados WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	RegistroEliminado.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}
