package Utilities

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Crea un archivo Binario
func CreateFile(name string) error {
	// Se verifica si el archivo existe
	if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
		return err
	}

	// Crear el archivo
	if _, err := os.Create(name); err != nil {
		return err
	}

	return nil
}

// Abre un archivo binario en modo lectura y escritura
func OpenFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR, 0644)
}

// Escribe un objeto en un archivo binario en una posición específica
func WriteObject(file *os.File, data interface{}, position int64) error {
	// Moverse a la posición especificada
	if _, err := file.Seek(position, io.SeekStart); err != nil {
		return err
	}

	// Escribir los datos en el archivo
	return binary.Write(file, binary.LittleEndian, data)
}

// Lee un objeto desde un archivo binario en una posición específica
func ReadObject(file *os.File, data interface{}, position int64) error {
	// Moverse a la posición especificada
	if _, err := file.Seek(position, io.SeekStart); err != nil {
		return err
	}

	// Leer los datos del archivo
	return binary.Read(file, binary.LittleEndian, data)
}

// readFromFile lee datos desde un archivo binario en la posición especificada
func ReadFromFile(file *os.File, offset int64, data interface{}) error {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to offset %d: %w", offset, err)
	}

	err = binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		return fmt.Errorf("failed to read data from file: %w", err)
	}

	return nil
}

func CreateParentDirs(path string) error {
	dir := filepath.Dir(path)
	// os.MkdirAll no sobrescribe las carpetas existentes, solo crea las que no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear las carpetas padre: %v", err)
	}
	return nil
}

// GetParentDirectories obtiene las carpetas padres y el directorio de destino
func GetParentDirectories(path string) ([]string, string) {
	// Normalizar el path
	path = filepath.Clean(path)

	// Dividir el path en sus componentes
	components := strings.Split(path, string(filepath.Separator))

	// Lista para almacenar las rutas de las carpetas padres
	var parentDirs []string

	// Construir las rutas de las carpetas padres, excluyendo la última carpeta
	for i := 1; i < len(components)-1; i++ {
		parentDirs = append(parentDirs, components[i])
	}

	// La última carpeta es la carpeta de destino
	destDir := components[len(components)-1]

	return parentDirs, destDir
}

// First devuelve el primer elemento de un slice
func First[T any](slice []T) (T, error) {
	if len(slice) == 0 {
		var zero T
		return zero, errors.New("el slice está vacío")
	}
	return slice[0], nil
}

// RemoveElement elimina un elemento de un slice en el índice dado
func RemoveElement[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice // Índice fuera de rango, devolver el slice original
	}
	return append(slice[:index], slice[index+1:]...)
}

// splitStringIntoChunks divide una cadena en partes de tamaño chunkSize y las almacena en una lista
func SplitStringIntoChunks(s string) []string {
	var chunks []string
	for i := 0; i < len(s); i += 64 {
		end := i + 64
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
