package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Variables globales
var (
	mu                  sync.Mutex
	todasLasImagenes    [][]string
	imagenesMostradas   map[string]bool
	ultimoIndiceCarpeta int
)

// Función principal
func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Uso: go run main.go <carpeta1> <carpeta2>")
		return
	}

	todasLasImagenes = make([][]string, len(args)-1)
	for i, carpeta := range args[1:] {
		cargarImagenesDeCarpeta(carpeta, i)
	}

	imagenesMostradas = make(map[string]bool)
	src := rand.NewSource(time.Now().UnixNano())
	ganadorAleatorio := rand.New(src)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if len(todasLasImagenes) == 0 {
			http.Error(rw, "No hay imágenes disponibles", http.StatusInternalServerError)
			return
		}

		indiceCarpetaActual := (ultimoIndiceCarpeta + 1) % len(todasLasImagenes)
		var imagenesDisponibles []string
		for _, img := range todasLasImagenes[indiceCarpetaActual] {
			if !imagenesMostradas[img] {
				imagenesDisponibles = append(imagenesDisponibles, img)
			}
		}

		// Si no hay imágenes disponibles, reinicia las mostradas
		if len(imagenesDisponibles) == 0 {
			fmt.Println("No hay imágenes disponibles en la carpeta actual. Cambiando de carpeta.")
			imagenesMostradas = make(map[string]bool)
			indiceCarpetaActual = (indiceCarpetaActual + 1) % len(todasLasImagenes)
			imagenesDisponibles = todasLasImagenes[indiceCarpetaActual]
		}

		// Seleccionar hasta 4 imágenes diferentes
		numImagenesAMostrar := 4
		if len(imagenesDisponibles) < numImagenesAMostrar {
			numImagenesAMostrar = len(imagenesDisponibles)
		}

		imagenesSeleccionadas := make([]string, 0, numImagenesAMostrar)
		imagenesMostradasTemp := make(map[string]bool)

		for len(imagenesSeleccionadas) < numImagenesAMostrar {
			numeroAleatorio := ganadorAleatorio.Intn(len(imagenesDisponibles))
			imagenSeleccionada := imagenesDisponibles[numeroAleatorio]

			if !imagenesMostradasTemp[imagenSeleccionada] {
				imagenesSeleccionadas = append(imagenesSeleccionadas, imagenSeleccionada)
				imagenesMostradasTemp[imagenSeleccionada] = true
				imagenesMostradas[imagenSeleccionada] = true
			}
		}

		// Leer y codificar las imágenes seleccionadas en Base64
		imagenesBase64 := make([]string, len(imagenesSeleccionadas))
		for i, img := range imagenesSeleccionadas {
			imagenBytes, err := os.ReadFile(img)
			if err != nil {
				fmt.Println("Error al leer la imagen: ", err)
				http.Error(rw, "Error al leer la imagen", http.StatusInternalServerError)
				return
			}
			imagenesBase64[i] = base64.StdEncoding.EncodeToString(imagenBytes)
		}

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "Desconocido"
		}

		// Respuesta HTML con Bootstrap
		fmt.Fprintf(rw, `<!DOCTYPE html>
			<html lang="es">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Servidor de imágenes</title>
				<!-- Bootstrap CSS -->
				<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
				<style>
					.footer {
						background-color: #00008B;
						color: #fff;
						padding: 10px;
						position: fixed;
						bottom: 0;
						width: 100%%;
						text-align: center;
					}
					.navbar {
						background-color: #00008B;
						color: #fff;
					}
					.container {
						margin-top: 20px;
					}
					.img-thumbnail {
						width: 100%%;
						height: 250px; /* Fijar altura para uniformidad */
						object-fit: cover; /* Para mantener la proporción de la imagen */
						margin-bottom: 10px;
					}
				</style>
			</head>
			<body>
				<!-- Navbar -->
				<nav class="navbar navbar-dark">
					<div class="container-fluid">
						<a class="navbar-brand" href="#">Servidor de imágenes</a>
					</div>
				</nav>

				<!-- Contenido -->
				<div class="container">
					<h1 class="text-center">Hostname: %s</h1>
					<h2 class="text-center">Tema: Fantasia Medieval</h2>
					<div class="row">
						<div class="col-md-6">
							<a href="/imagen?img=%s">
								<img src="data:image/jpg;base64,%s" alt="Imagen 1" class="img-thumbnail">
							</a>
							<p class="text-center">Imagen1.jpg</p>
						</div>
						<div class="col-md-6">
							<a href="/imagen?img=%s">
								<img src="data:image/jpg;base64,%s" alt="Imagen 2" class="img-thumbnail">
							</a>
							<p class="text-center">Imagen2.jpg</p>
						</div>
					</div>
					<div class="row">
						<div class="col-md-6">
							<a href="/imagen?img=%s">
								<img src="data:image/jpg;base64,%s" alt="Imagen 3" class="img-thumbnail">
							</a>
							<p class="text-center">Imagen3.jpg</p>
						</div>
						<div class="col-md-6">
							<a href="/imagen?img=%s">
								<img src="data:image/jpg;base64,%s" alt="Imagen 4" class="img-thumbnail">
							</a>
							<p class="text-center">Imagen4.jpg</p>
						</div>
					</div>
				</div>

				<!-- Footer -->
				<div class="footer">
					<p>Computación en la nube | Juan Pablo Alviz V - Johan Andres Ospina O | 2024-2</p>
				</div>

				<!-- Bootstrap JS -->
				<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js"></script>
			</body>
			</html>`, hostname, imagenesSeleccionadas[0], imagenesBase64[0], imagenesSeleccionadas[1], imagenesBase64[1], imagenesSeleccionadas[2], imagenesBase64[2], imagenesSeleccionadas[3], imagenesBase64[3])

		ultimoIndiceCarpeta = indiceCarpetaActual
	})

	// Servidor de detalles de imagen
	http.HandleFunc("/imagen", func(rw http.ResponseWriter, req *http.Request) {
		img := req.URL.Query().Get("img")
		if img == "" {
			http.Error(rw, "Imagen no encontrada", http.StatusBadRequest)
			return
		}

		// Leer la imagen
		imagenBytes, err := os.ReadFile(img)
		if err != nil {
			http.Error(rw, "Error al leer la imagen", http.StatusInternalServerError)
			return
		}

		// Obtener información de la imagen
		info, err := os.Stat(img)
		if err != nil {
			http.Error(rw, "Error al obtener detalles de la imagen", http.StatusInternalServerError)
			return
		}
		peso := info.Size()
		ext := filepath.Ext(img)
		tipo := http.DetectContentType(imagenBytes)

		// Respuesta con detalles de la imagen
		fmt.Fprintf(rw, `<!DOCTYPE html>
			<html lang="es">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Detalles de la Imagen</title>
				<!-- Bootstrap CSS -->
				<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
				<style>
					.footer {
						background-color: #00008B;
						color: #fff;
						padding: 10px;
						position: fixed;
						bottom: 0;
						width: 100%%;
						text-align: center;
					}
					.navbar {
						background-color: #00008B;
						color: #fff;
					}
					.container {
						margin-top: 20px;
					}
					.img-preview {
						width: 100%%;
						max-height: 500px;
						object-fit: cover;
						margin-bottom: 20px;
					}
					.info-box {
						background-color: #f8f9fa;
						border: 1px solid #ccc;
						padding: 15px;
						border-radius: 5px;
					}
				</style>
			</head>
			<body>
				<!-- Navbar -->
				<nav class="navbar navbar-dark">
					<div class="container-fluid">
						<a class="navbar-brand" href="#">Servidor de imágenes</a>
					</div>
				</nav>

				<!-- Contenido -->
				<div class="container">
					<h1 class="text-center">Detalles de la imagen</h1>
					<div class="row">
						<div class="col-md-12 text-center">
							<img src="data:image/jpg;base64,%s" alt="Imagen seleccionada" class="img-preview">
						</div>
					</div>
					<div class="row">
						<div class="col-md-12">
							<div class="info-box">
								<p><strong>Nombre de la imagen:</strong> %s</p>
								<p><strong>Peso:</strong> %d bytes</p>
								<p><strong>Tipo de imagen:</strong> %s</p>
								<p><strong>Extensión:</strong> %s</p>
							</div>
						</div>
					</div>
				</div>

				<!-- Footer -->
				<div class="footer">
					<p>Computación en la nube | Juan Pablo Alviz V - Johan Andres Ospina O | 2024-2</p>
				</div>

				<!-- Bootstrap JS -->
				<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js"></script>
			</body>
			</html>`, base64.StdEncoding.EncodeToString(imagenBytes), filepath.Base(img), peso, tipo, ext)
	})

	fmt.Println("Servidor iniciado en http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

// Función para cargar imágenes de una carpeta
func cargarImagenesDeCarpeta(carpeta string, indice int) {
	archivos, err := os.ReadDir(carpeta)
	if err != nil {
		log.Fatal(err)
	}

	var extensionesImagen = []string{".png", ".jpg", ".jpeg"}

	for _, archivo := range archivos {
		if tieneExtensionValida(extensionesImagen, archivo.Name()) {
			todasLasImagenes[indice] = append(todasLasImagenes[indice], carpeta+"/"+archivo.Name())
		}
	}
}

// Función para verificar si el archivo tiene una extensión de imagen válida
func tieneExtensionValida(extensiones []string, nombre string) bool {
	partes := strings.Split(nombre, ".")
	for _, ext := range extensiones {
		if "."+partes[len(partes)-1] == ext {
			return true
		}
	}
	return false
}
