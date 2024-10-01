package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Estructura para pasar datos a la plantilla HTML
type PageData struct {
	HostName    string
	RandomImage string
	Files       []string
}

func main() {

	args := os.Args

	fmt.Println("Archivo de ejecucion: " + args[0])

	// Verificación de que se pasó el directorio de las imágenes
	if len(args) < 2 {
		print("Faltan parámetros")
		return
	}

	dir := args[1]

	// Semilla para el valor aleatorio
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Extensiones a filtrar
	var extImg = []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".heif", ".heic", ".ppm", ".pgm", ".pbm",
		".svg", ".eps", ".ai", ".ico", ".icns", ".psd", ".xcf", ".raw", ".cr2", ".nef", ".arw",
		".dds", ".exr", ".hdr", ".avif",
	}

	var filterFiles []string

	// Lee los archivos del directorio
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Itera sobre los archivos y guarda en el array
	for _, file := range files {
		if foundExt2(extImg, file.Name()) {
			filterFiles = append(filterFiles, file.Name())
		}
	}

	// Genera un número aleatorio para seleccionar una imagen
	rn := r.Intn(len(filterFiles))
	rImg := filterFiles[rn]

	// Consultar nombre de host
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error al obtener el nombre del host: ", err)
		return
	}

	// Para codificar una imagen a base64 (opcional)
	pathImg := dir + "/" + rImg
	imgBytes, err := os.ReadFile(pathImg)
	if err != nil {
		fmt.Println("Error al leer la imagen: ", err)
		return
	}
	b64String := base64.StdEncoding.EncodeToString(imgBytes)

	// Crear el archivo codificado en base64
	file, err := os.Create("imgBase64.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo: ", err)
		return
	}
	defer file.Close()
	file.WriteString(b64String)

	// Servir archivos estáticos (como el CSS de Tailwind)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Ruta principal para mostrar la página
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			HostName:    hostname,
			RandomImage: rImg,
			Files:       filterFiles,
		}
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, data)
	})

	// Inicia el servidor en el puerto 3000
	fmt.Println("Servidor en http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

// Función para verificar si el archivo tiene una de las extensiones filtradas
func foundExt2(arr []string, name string) bool {
	aux := strings.Split(name, ".")
	for _, ext := range arr {
		if "."+aux[len(aux)-1] == ext {
			return true
		}
	}
	return false
}
