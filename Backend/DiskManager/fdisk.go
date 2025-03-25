package DiskManager

import (
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
)

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Ajustar el tamaño en bytes
	size = convertSize(size, unit)

	// Abrir el archivo binario en la ruta proporcionada
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error: Could not open file at path:", path)
		return
	}
	defer file.Close() // Cerrar el archivo al finalizar

	// Leer el objeto MBR
	var TempMBR Disk.MBR
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file")
		return
	}
	fmt.Println("-------------")

	// Validar particiones y espacio usado
	totalPartitions, _, extendedCount, usedSpace := countPartitions(&TempMBR)
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}

	if type_ == "e" && extendedCount > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}

	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}

	if usedSpace+int32(size) > TempMBR.Size {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	// Determinar la posición de inicio (gap) para la nueva partición
	gap := calculateGap(&TempMBR, totalPartitions)

	// Crear partición primaria o extendida
	if type_ == "p" || type_ == "e" {
		emptyIndex := findEmptyPartition(&TempMBR)
		if emptyIndex < 0 {
			fmt.Println("Error: No se encontró espacio en la tabla de particiones.")
			return
		}
		TempMBR.Partitions[emptyIndex].Size = int32(size)
		TempMBR.Partitions[emptyIndex].Start = gap
		copy(TempMBR.Partitions[emptyIndex].Name[:], name)
		TempMBR.Partitions[emptyIndex].Fit = fit[0]
		TempMBR.Partitions[emptyIndex].Status = '0'
		TempMBR.Partitions[emptyIndex].Type = type_[0]
		TempMBR.Partitions[emptyIndex].Correlative = int32(totalPartitions + 1)
		Disk.PrintPartition(TempMBR.Partitions[emptyIndex])
		if type_ == "e" {
			// Inicializar el primer EBR en la partición extendida
			ebr := Disk.EBR{
				PartFit:   fit[0],
				PartStart: gap, // El primer EBR se coloca al inicio de la partición extendida
				PartSize:  0,
				PartNext:  -1,
			}
			copy(ebr.PartName[:], "")
			Utilities.WriteObject(file, ebr, int64(gap))
		}
	} else if type_ == "l" {
		// Crear partición lógica dentro de la partición extendida
		if err := createLogicalPartition(file, &TempMBR, size, name, fit); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Sobrescribir el MBR actualizado
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	// Verificar y mostrar el MBR actualizado
	var TempMBR2 Disk.MBR
	if err := Utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return
	}

	fmt.Println("======FIN FDISK======")
	fmt.Println("")
}

// convertSize ajusta el tamaño en bytes según la unidad (k o m)
func convertSize(size int, unit string) int {
	switch unit {
	case "k":
		return size * 1024
	case "m":
		return size * 1024 * 1024
	default:
		return size
	}
}

// countPartitions recorre las particiones del MBR y retorna el total, la cantidad de primarias, extendidas y el espacio usado
func countPartitions(mbr *Disk.MBR) (total int, primary int, extended int, usedSpace int32) {
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Size != 0 {
			total++
			usedSpace += mbr.Partitions[i].Size
			if mbr.Partitions[i].Type == 'p' {
				primary++
			} else if mbr.Partitions[i].Type == 'e' {
				extended++
			}
		}
	}
	return
}

// calculateGap determina la posición de inicio para la nueva partición
func calculateGap(mbr *Disk.MBR, totalPartitions int) int32 {
	if totalPartitions == 0 {
		return int32(binary.Size(*mbr))
	}
	// Se asume que la última partición ocupada es la última en el arreglo
	last := mbr.Partitions[totalPartitions-1]
	return last.Start + last.Size
}

// findEmptyPartition retorna el índice de la primera partición vacía, o -1 si no se encuentra
func findEmptyPartition(mbr *Disk.MBR) int {
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Size == 0 {
			return i
		}
	}
	return -1
}

// createLogicalPartition maneja la creación de particiones lógicas dentro de una partición extendida
func createLogicalPartition(file *os.File, mbr *Disk.MBR, size int, name string, fit string) error {
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Type == 'e' {
			ebrPos := mbr.Partitions[i].Start
			var ebr Disk.EBR
			// Recorrer la cadena de EBR hasta encontrar el último
			for {
				if err := Utilities.ReadObject(file, &ebr, int64(ebrPos)); err != nil {
					return fmt.Errorf("error al leer EBR: %v", err)
				}
				if ebr.PartNext == -1 {
					break
				}
				ebrPos = ebr.PartNext
			}

			// Calcular la posición para el nuevo EBR y la partición lógica
			newEBRPos := ebr.PartStart + ebr.PartSize
			logicalPartitionStart := newEBRPos + int32(binary.Size(ebr))

			// Actualizar el EBR anterior para enlazar al nuevo EBR
			ebr.PartNext = newEBRPos
			Utilities.WriteObject(file, ebr, int64(ebrPos))

			// Crear y escribir el nuevo EBR
			newEBR := Disk.EBR{
				PartFit:   fit[0],
				PartStart: logicalPartitionStart,
				PartSize:  int32(size),
				PartNext:  -1,
				PartMount: '0',
			}
			copy(newEBR.PartName[:], name)
			Utilities.WriteObject(file, newEBR, int64(newEBRPos))

			// Imprimir el nuevo EBR creado
			Disk.PrintEBR(newEBR)
			break
		}
	}
	return nil
}
