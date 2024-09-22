//go mod init hello
//go mod tidy

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	//Ruta  mac /Users/johanospina/Downloads

	//Para Windows C:/Users/ospin/OneDrive/Imágenes && C:/Users/ospin/Downloads

	//Lista de los argumentos pasados incluyendo el archivo go
	args := os.Args

	fmt.Println("Archivo de ejecucion: " + args[0])

	//Verificacion de que se paso el directorio de las imagenes
	if len(args) < 2 {
		print("Falta parametros")
		return
	}

	dir := args[1]

	//Semilla para el valor aleatorio
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	//extension a filtrar
	var extImg = []string{
		// Formatos Rasterizados
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".heif", ".heic", ".ppm", ".pgm", ".pbm",
		// Formatos Vectoriales
		".svg", ".eps", ".ai",
		// Otros Formatos de Imágenes
		".ico", ".icns", ".psd", ".xcf", ".raw", ".cr2", ".nef", ".arw",
		// Formatos 3D o Especializados
		".dds", ".exr", ".hdr", ".avif",
	}

	var filterFiles []string

	//lee los archivos con os.ReadDir
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	//Itera sobre los archivos en la carpeta y guarda en el array
	var i = 0 //Iteracion para guardar al ritmo de que se encuentran imagenes
	for _, file := range files {
		if foundExt2(extImg, file.Name()) {
			fmt.Println(file.Name())
			filterFiles = append(filterFiles, file.Name())
			i++
		}
	}

	//rn es el numero aleatorio generado a partir del tamaño del array
	rn := r.Intn(len(filterFiles) - 1)
	rImg := filterFiles[rn]
	print("\n" + rImg)

	//Consultar nombre de host e imprimirlo
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error al querer obtener el nombre: ", err)
		return
	}
	fmt.Println("\n"+"Nombre del host: ", hostname)

	//PARA CODIFICAR A BASE64
	//Conseguir los bytes de la imagen
	pathImg := dir + "/" + rImg

	imgBytes, err := os.ReadFile(pathImg)
	if err != nil {
		fmt.Println("Error al leer la imagen: ", err)
		return
	}

	//Codificar a base 64
	b64String := base64.StdEncoding.EncodeToString(imgBytes)

	//Crea el archivo dado que si se imprime es demasiado largo
	file, err := os.Create("imgBase64.txt")
	if err != nil {
		fmt.Println("Error al crear el archivo: ", err)
		return
	}
	defer file.Close()

	//Escritura del archivo
	_, err = file.WriteString(b64String)
	if err != nil {
		fmt.Println("Error al escribir el archivo:", err)
		return
	}
	fmt.Println("Archivo creado")

	//Para subir el servicio web:
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for _, file := range filterFiles {
			fmt.Fprint(rw, file+"\n")
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}

func foundExt2(arr []string, name string) bool {
	aux := strings.Split(name, ".")
	for _, ind := range arr {
		if "."+aux[len(aux)-1] == ind {
			return true
		}
	}
	return false
}

// Se hizo asi solo por practica pero se sabe que es mas eficiente un bucle
//func foundExt(arr []string, name string, i int) bool {
//	if i >= len(arr) {
//		return false
//	}
//	aux := strings.Split(name, ".")
//	if "."+aux[len(aux)-1] == arr[i] {
//		return true
//	}
//	return foundExt(arr, name, i+1)
//}
