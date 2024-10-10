package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Variables globales
var (
	mu                  sync.Mutex      // Mutex para proteger el acceso concurrente a recursos compartidos
	todasLasImagenes    [][]string      // Lista de imágenes de ambas carpetas (lista de listas)
	imagenesMostradas   map[string]bool // Mapa para almacenar imágenes que ya se han mostrado
	ultimoIndiceCarpeta int             // Índice de la última carpeta utilizada
)

// Función principal
func main() {
	// Lista de argumentos pasados por línea de comandos
	args := os.Args

	// Verificación de que se pasaron suficientes argumentos
	if len(args) < 3 {
		fmt.Println("Uso: go run main.go <carpeta1> <carpeta2>")
		return
	}

	// Cargar imágenes de ambas carpetas
	todasLasImagenes = make([][]string, len(args)-1) // Inicializar la lista de imágenes
	for i, carpeta := range args[1:] {
		cargarImagenesDeCarpeta(carpeta, i)
	}

	// Inicializar el mapa de imágenes mostradas
	imagenesMostradas = make(map[string]bool)

	// Semilla para el generador de números aleatorios
	src := rand.NewSource(time.Now().UnixNano())
	ganadorAleatorio := rand.New(src)

	// Función manejadora para el servidor HTTP
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		mu.Lock() // Bloquear acceso concurrente
		if len(todasLasImagenes) == 0 {
			http.Error(rw, "No hay imágenes disponibles", http.StatusInternalServerError)
			mu.Unlock()
			return
		}

		// Asegurarse de que la carpeta seleccionada no sea la misma que la última
		indiceCarpetaActual := (ultimoIndiceCarpeta + 1) % len(todasLasImagenes)

		// Filtrar las imágenes de la carpeta actual que no han sido mostradas
		var imagenesDisponibles []string
		for _, img := range todasLasImagenes[indiceCarpetaActual] {
			if !imagenesMostradas[img] {
				imagenesDisponibles = append(imagenesDisponibles, img)
			}
		}

		// Si no hay imágenes disponibles en la carpeta actual, cambiar a la otra carpeta
		if len(imagenesDisponibles) == 0 {
			fmt.Println("No hay imágenes disponibles en la carpeta actual. Cambiando de carpeta.")
			imagenesMostradas = make(map[string]bool)                               // Reiniciar las imágenes mostradas
			indiceCarpetaActual = (indiceCarpetaActual + 1) % len(todasLasImagenes) // Cambiar de carpeta
			imagenesDisponibles = todasLasImagenes[indiceCarpetaActual]             // Usar todas las imágenes de la nueva carpeta
		}

		// Seleccionar una imagen aleatoria de las disponibles
		numeroAleatorio := ganadorAleatorio.Intn(len(imagenesDisponibles))
		imagenSeleccionada := imagenesDisponibles[numeroAleatorio]

		// Marcar la imagen como mostrada
		imagenesMostradas[imagenSeleccionada] = true

		// Leer los bytes de la imagen seleccionada
		imagenBytes, err := os.ReadFile(imagenSeleccionada)
		if err != nil {
			fmt.Println("Error al leer la imagen: ", err)
			http.Error(rw, "Error al leer la imagen", http.StatusInternalServerError)
			mu.Unlock()
			return
		}

		// Codificar la imagen en Base64
		imagenBase64 := base64.StdEncoding.EncodeToString(imagenBytes)

		// Obtener el nombre del host
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "Desconocido"
		}

		// Generar la respuesta HTML que muestre la imagen en Base64
		fmt.Fprintf(rw, `<!DOCTYPE html>
			<html lang="es">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Servidor de imágenes</title>
			</head>
			<body>
				<h1>Servidor de imágenes</h1>
				<h2>Hostname: %s</h2>
				<img src="data:image/jpg;base64,%s" alt="Imagen aleatoria" style="max-width: 100%%; height: auto;">
				<h3>Computación en la nube</h3>
				<p>Estudiantes:</p>
				<ul>
					<li>Juan Pablo Alviz Velasquez</li>
					<li>Johan Andres Ospina Ospina</li>
				</ul>
			</body>
			</html>`, hostname, imagenBase64)

		// Actualizar el índice de la última carpeta utilizada
		ultimoIndiceCarpeta = indiceCarpetaActual
		mu.Unlock() // Desbloquear acceso
	})

	// Iniciar el servidor en el puerto 3000
	fmt.Println("Servidor iniciado en http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

// Función para cargar imágenes de una carpeta
func cargarImagenesDeCarpeta(carpeta string, indice int) {
	archivos, err := os.ReadDir(carpeta)
	if err != nil {
		log.Fatal(err)
	}

	// Extensiones de imágenes a filtrar
	var extensionesImagen = []string{".png", ".jpg", ".jpeg"}

	for _, archivo := range archivos {
		if tieneExtensionValida(extensionesImagen, archivo.Name()) {
			todasLasImagenes[indice] = append(todasLasImagenes[indice], carpeta+"/"+archivo.Name())
		}
	}
}

// Función para verificar si el archivo tiene una extensión de imagen válida
func tieneExtensionValida(extensiones []string, nombre string) bool {
	// Separar el nombre por el punto para obtener la extensión
	partes := strings.Split(nombre, ".")
	for _, ext := range extensiones {
		if "."+partes[len(partes)-1] == ext {
			return true
		}
	}
	return false
}
