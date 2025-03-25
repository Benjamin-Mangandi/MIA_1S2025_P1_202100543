package Disk

import (
	"Backend/Responsehandler"
	"bytes"
	"fmt"
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

	nameStr := string(bytes.Trim(data.PartName[:], "\x00"))
	response := strings.Repeat("-", 40) + "\n" +
		"Partición lógica creada correctamente\n" +
		"Nombre: " + nameStr + "\n" +
		"Tamaño: " + fmt.Sprintf("%d", data.PartSize) + " bytes\n" +
		"Tipo: l\n" +
		"Estado: " + fmt.Sprintf("%c", data.PartMount) + "\n" +
		"Fit: " + fmt.Sprintf("%c", data.PartFit) + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
