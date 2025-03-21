package Disk

import (
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Estructura del EBR
type EBR struct {
	PartMount byte     // Indica si la partición está montada
	PartFit   byte     // Ajuste de la partición (B, F o W)
	PartStart int32    // Byte donde inicia la partición
	PartSize  int32    // Tamaño total de la partición en bytes
	PartNext  int32    // Byte donde está el próximo EBR (-1 si no hay siguiente)
	PartName  [16]byte // Nombre de la partición
}

// Función para imprimir la estructura EBR
func PrintEBR(data EBR) {
	// Convertir PartName a string y eliminar caracteres nulos
	partName := strings.TrimRight(string(data.PartName[:]), "\x00")

	fmt.Printf("Name: %s, Fit: %c, Start: %d, Size: %d, Next: %d, Mount: %c\n",
		partName,
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount)
}

func ReadEBR(start int32, file *os.File) (*EBR, error) {
	fmt.Printf("Leyendo EBR desde el archivo en la posición: %d\n", start)

	// Leer bytes crudos para depuración
	rawData := make([]byte, binary.Size(EBR{}))
	_, err := file.Seek(int64(start), 0)
	if err != nil {
		return nil, fmt.Errorf("fallo al posicionarse en %d: %v", start, err)
	}

	_, err = file.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("fallo al leer datos en %d: %v", start, err)
	}

	fmt.Printf("Bytes crudos leídos en %d: %v\n", start, rawData)

	// Ahora decodifica normalmente
	return Decode(file, int64(start))
}

// Decode deserializa la estructura EBR desde un archivo en la posición especificada
func Decode(file *os.File, position int64) (*EBR, error) {
	ebr := &EBR{}

	// Verificar que la posición no sea negativa y esté dentro del rango del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %v", err)
	}
	if position < 0 || position >= fileInfo.Size() {
		return nil, fmt.Errorf("posición inválida para EBR: %d", position)
	}

	err = Utilities.ReadFromFile(file, position, ebr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("EBR decoded from position %d with success.\n", position)
	return ebr, nil
}
