package Utilities

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
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
