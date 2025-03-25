package Disk

import (
	"Backend/Responsehandler"
	"bytes"
	"fmt"
)

type Partition struct {
	Status      byte     // '0' o '1' según si está montada
	Type        byte     // 'P' (Primaria), 'E' (Extendida) o 'L' (Logica)
	Fit         byte     // 'B', 'F' o 'W' (Best, First, Worst)
	Start       int32    // Byte donde inicia la partición
	Size        int32    // Tamaño total de la partición en bytes
	Name        [16]byte // Nombre de la partición
	Correlative int32    // -1 por defecto, incrementa al montar
	Id          [4]byte  // ID de la partición montada
}

func PrintPartition(data Partition) {
	nameStr := string(bytes.Trim(data.Name[:], "\x00"))
	idStr := string(bytes.Trim(data.Id[:], "\x00"))
	answer := "---------------------\n" +
		"Partición creada correctamente\n" +
		"Nombre: " + nameStr + "\n" +
		"Tamaño: " + fmt.Sprintf("%d", data.Size) + " bytes\n" +
		"Tipo: " + fmt.Sprintf("%c", data.Type) + "\n" +
		"Estado: " + fmt.Sprintf("%c", data.Status) + "\n" +
		"Fit: " + fmt.Sprintf("%c", data.Fit) + "\n" +
		"Correlativo: " + fmt.Sprintf("%d", data.Correlative) + "\n" +
		"ID: " + idStr + "\n" +
		"---------------------"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, answer)
}
