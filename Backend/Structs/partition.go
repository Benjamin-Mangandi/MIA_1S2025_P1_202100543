package Structs

import "fmt"

type Partition struct {
	Status      byte     // '0' o '1' según si está montada
	Type        byte     // 'P' (Primaria) o 'E' (Extendida)
	Fit         byte     // 'B', 'F' o 'W' (Best, First, Worst)
	Start       int32    // Byte donde inicia la partición
	Size        int32    // Tamaño total de la partición en bytes
	Name        [16]byte // Nombre de la partición
	Correlative int32    // -1 por defecto, incrementa al montar
	Id          [4]byte  // ID de la partición montada
}

func PrintPartition(data Partition) {
	fmt.Println(fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s", string(data.Name[:]), string(data.Type), data.Start, data.Size, string(data.Status), string(data.Id[:])))
}
