package DiskManager

import (
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Mkdisk crea un nuevo disco binario
func Mkdisk(size int, fit string, unit string, path string) {
	// Convertir fit y unit a minúsculas
	fit = strings.ToLower(fit)
	unit = strings.ToLower(unit)

	// Validar fit (bf, wf, ff)
	validFits := map[string]bool{"bf": true, "wf": true, "ff": true}
	if !validFits[fit] {
		fmt.Println("Error: Fit debe ser 'bf', 'wf' o 'ff'.")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: El tamaño del disco debe ser mayor a 0.")
		return
	}

	// Validar unit (k o m)
	var multiplier int
	switch unit {
	case "k":
		multiplier = 1024
	case "m":
		multiplier = 1024 * 1024
	default:
		fmt.Println("Error: Las unidades válidas son 'k' (kilobytes) o 'm' (megabytes).")
		return
	}

	// Calcular el tamaño real en bytes
	sizeInBytes := size * multiplier

	// Crear archivo binario
	if err := Utilities.CreateFile(path); err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}

	// Abrir archivo binario
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Escribir bytes vacíos
	emptyBlock := make([]byte, 1024) // Bloque de 1 KB lleno de ceros
	for i := 0; i < sizeInBytes; i += 1024 {
		_, err := file.Write(emptyBlock)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	}

	// Crear estructura MBR
	var newMBR Disk.MBR
	newMBR.Size = int32(sizeInBytes)
	newMBR.Signature = rand.Int31() // Número aleatorio único para el disco
	newMBR.Fit = fit[0]             // Almacenar solo el primer carácter del fit

	// Obtener la fecha en formato YYYY-MM-DD HH:MM y almacenarla en bytes
	currentTime := time.Now().Format("2006-01-02 15:04")
	copy(newMBR.CreationDate[:], currentTime)

	// Escribir el MBR en el archivo
	if err := Utilities.WriteObject(file, newMBR, 0); err != nil {
		fmt.Println("Error al escribir el MBR en el archivo:", err)
		return
	}

	// Leer el MBR
	var tempMBR Disk.MBR
	if err := Utilities.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Imprimir el MBR leído
	Disk.PrintMBR(tempMBR)
}
