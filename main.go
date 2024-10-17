package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

// Proyecto representa un proyecto de residencia
type Proyecto struct {
	ID            int    `json:"id"`
	Titulo        string `json:"titulo"`
	Descripcion   string `json:"descripcion"`
	Estudiante    int    `json:"estudiante"`
	FechaRegistro string `json:"fecha_registro"`
	Estatus       string `json:"estatus"`
}

// Usuario representa un usuario
type Usuario struct {
	ID         int    `json:"id"`
	Usuario    string `json:"usuario"`
	Nombre     string `json:"nombre"`
	Apellidos  string `json:"apellidos"`
	Contrasena string `json:"contrasena"`
	Carrera    string `json:"carrera"`
	Semestre   int    `json:"semestre"`
}

func main() {
	// Configuración de la conexión a la base de datos
	connString := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creando la conexión:", err)
		return
	}
	defer db.Close()
	// Verifica si puedes hacer ping a la base de datos
	err = db.Ping()
	if err != nil {
		log.Fatal("Error conectándose a la base de datos: ", err.Error())
	}

	log.Println("Conexión a la base de datos exitosa")

	r := gin.Default()

	// Endpoints para proyectos

	r.GET("/proyectos/:id_usuario", func(c *gin.Context) {
		idUsuario, err := strconv.Atoi(c.Param("id_usuario"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
			return
		}
		proyectos, err := getProyectosByUsuario(db, idUsuario)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proyectos)
	})

	r.GET("/proyectos", func(c *gin.Context) {
		proyectos, err := getProyectos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, proyectos)
	})

	r.POST("/proyectos", func(c *gin.Context) {
		var nuevoProyecto Proyecto
		if err := c.ShouldBindJSON(&nuevoProyecto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := createProyecto(db, &nuevoProyecto); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, nuevoProyecto)
	})

	r.PUT("/proyectos/:id", func(c *gin.Context) {
		var proyecto Proyecto
		if err := c.ShouldBindJSON(&proyecto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}
		proyecto.ID = id

		if err := updateProyecto(db, &proyecto); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, proyecto)
	})

	r.DELETE("/proyectos/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		if err := deleteProyecto(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Proyecto eliminado"})
	})

	// Endpoints para usuarios

	r.POST("/usuarios", func(c *gin.Context) {
		var nuevoUsuario Usuario
		if err := c.ShouldBindJSON(&nuevoUsuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := createUsuario(db, &nuevoUsuario); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, nuevoUsuario)
	})

	// Endpoint para autenticar usuarios
	r.POST("/auth", func(c *gin.Context) {
		var credenciales struct {
			Usuario    string `json:"usuario"`
			Contrasena string `json:"contrasena"`
		}
		if err := c.ShouldBindJSON(&credenciales); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		usuario, err := authenticateUsuario(db, credenciales.Usuario, credenciales.Contrasena)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}

		c.JSON(http.StatusOK, usuario)
	})

	// Ejecutar el servidor
	r.Run(":8080") // Cambia el puerto si es necesario
}

// Función para obtener todos los proyectos
func getProyectos(db *sql.DB) ([]Proyecto, error) {
	rows, err := db.Query("SELECT id_proyecto, titulo, descripcion, IDestudiante, fecha_registro, estatus FROM proyectos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proyectos []Proyecto
	for rows.Next() {
		var p Proyecto
		if err := rows.Scan(&p.ID, &p.Titulo, &p.Descripcion, &p.Estudiante, &p.FechaRegistro, &p.Estatus); err != nil {
			return nil, err
		}
		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

// Función para obtener proyectos por ID de usuario
func getProyectosByUsuario(db *sql.DB, idUsuario int) ([]Proyecto, error) {
	query := "SELECT id_proyecto, titulo, descripcion, IDestudiante, fecha_registro, estatus FROM proyectos WHERE IDestudiante = @IDestudiante"
	rows, err := db.Query(query, sql.Named("IDestudiante", idUsuario))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proyectos []Proyecto
	for rows.Next() {
		var p Proyecto
		if err := rows.Scan(&p.ID, &p.Titulo, &p.Descripcion, &p.Estudiante, &p.FechaRegistro, &p.Estatus); err != nil {
			return nil, err
		}
		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

// Función para crear un nuevo proyecto
func createProyecto(db *sql.DB, proyecto *Proyecto) error {
	query := "INSERT INTO proyectos (titulo, descripcion, IDestudiante, fecha_registro, estatus) OUTPUT INSERTED.id_proyecto VALUES (@titulo, @descripcion, @IDestudiante, @fecha_registro, @estatus)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		sql.Named("titulo", proyecto.Titulo),
		sql.Named("descripcion", proyecto.Descripcion),
		sql.Named("IDestudiante", proyecto.Estudiante),
		sql.Named("fecha_registro", proyecto.FechaRegistro),
		sql.Named("estatus", "Abierto"),
	).Scan(&proyecto.ID)
	if err != nil {
		return err
	}
	return nil
}

// Función para actualizar un proyecto
func updateProyecto(db *sql.DB, proyecto *Proyecto) error {
	query := "UPDATE proyectos SET titulo = @titulo, descripcion = @descripcion, IDestudiante = @IDestudiante, fecha_registro = @fecha_registro, estatus = @estatus WHERE id_proyecto = @id"
	_, err := db.Exec(query,
		sql.Named("titulo", proyecto.Titulo),
		sql.Named("descripcion", proyecto.Descripcion),
		sql.Named("IDestudiante", proyecto.Estudiante),
		sql.Named("fecha_registro", proyecto.FechaRegistro),
		sql.Named("estatus", proyecto.Estatus),
		sql.Named("id", proyecto.ID),
	)
	if err != nil {
		return err
	}
	return nil
}

// Función para eliminar un proyecto
func deleteProyecto(db *sql.DB, id int) error {
	query := "DELETE FROM proyectos WHERE id_proyecto = @id"
	_, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

// Función para crear un nuevo usuario
func createUsuario(db *sql.DB, usuario *Usuario) error {
	query := "INSERT INTO usuarios (usuario, nombre, apellidos, contrasena, carrera, semestre) OUTPUT INSERTED.IDestudiante VALUES (@usuario, @nombre, @apellidos, @contrasena, @carrera, @semestre)"
	err := db.QueryRow(query,
		sql.Named("usuario", usuario.Usuario),
		sql.Named("nombre", usuario.Nombre),
		sql.Named("apellidos", usuario.Apellidos),
		sql.Named("contrasena", usuario.Contrasena),
		sql.Named("carrera", usuario.Carrera),
		sql.Named("semestre", usuario.Semestre),
	).Scan(&usuario.ID)
	if err != nil {
		return err
	}
	return nil
}

// Función para autenticar un usuario
func authenticateUsuario(db *sql.DB, usuario string, contrasena string) (*Usuario, error) {
	var u Usuario
	query := "SELECT IDestudiante, usuario, nombre, apellidos, contrasena, carrera, semestre FROM usuarios WHERE usuario = @usuario AND contrasena = @contrasena"
	err := db.QueryRow(query,
		sql.Named("usuario", usuario),
		sql.Named("contrasena", contrasena),
	).Scan(&u.ID, &u.Usuario, &u.Nombre, &u.Apellidos, &u.Contrasena, &u.Carrera, &u.Semestre)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
