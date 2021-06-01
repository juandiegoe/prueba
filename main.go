package main

import (
	//"fmt"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD()(conexion *sql.DB){

	conexion, err := sql.Open("mysql", ConnectionString)

	if err != nil{
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	http.HandleFunc("/",Inicio)
	http.HandleFunc("/editar",Editar)
	http.HandleFunc("/crear",Crear)
	http.HandleFunc("/insertar",Insertar)
	http.HandleFunc("/borrar",Borrar)
	http.HandleFunc("/actualizar",Actualizar)
	
	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8000",nil)
}

type Empleado struct{
	Id int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request){

	conexionEstablecida := conexionBD()
	registros,err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil{
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next(){

		var id int
		var nombre, correo string

		err = registros.Scan(&id,&nombre,&correo)

		if err != nil{
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	plantillas.ExecuteTemplate(w,"inicio",arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request){
	plantillas.ExecuteTemplate(w,"crear",nil)
}

func Insertar(w http.ResponseWriter, r *http.Request){

	if r.Method == "POST"{

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros,err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) Values(?,?)")

		if err != nil{
			panic(err.Error())
		}

		insertarRegistros.Exec(nombre,correo)

		http.Redirect(w,r,"/",301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request){

	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	borrarRegistro,err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil{
		panic(err.Error())
	}

	borrarRegistro.Exec(idEmpleado)

	http.Redirect(w,r,"/",301)
}

func Editar(w http.ResponseWriter, r *http.Request){

	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	registro,err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?",idEmpleado)

	empleado := Empleado{}

	for registro.Next(){

		var id int
		var nombre, correo string

		err = registro.Scan(&id,&nombre,&correo)

		if err != nil{
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	plantillas.ExecuteTemplate(w,"editar",empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request){

	if r.Method == "POST"{

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		modificarRegistros,err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?,correo=? WHERE id=?")

		if err != nil{
			panic(err.Error())
		}

		modificarRegistros.Exec(nombre,correo,id)

		http.Redirect(w,r,"/",301)
	}
}