package Structs

import (
	"fmt"
	"strings"
)

// Estructura del EBR según la imagen
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
