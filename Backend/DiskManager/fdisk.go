package DiskManager

import (
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Fdisk(size int, path string, name string, unit string, type_ string, fit string) {
	fmt.Println("====== Start FDISK ======")

	// Convertir el tamaño a bytes
	size = convertSize(size, unit)

	// Abrir el archivo de disco
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", path)
		return
	}

	// Leer el MBR
	var mbr Disk.MBR
	if err := Utilities.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Contar particiones y verificar espacio
	primaryCount, extendedCount, usedSpace, totalPartitions := countPartitions(mbr)

	// Validar si la partición ya existe
	if partitionExists(mbr, file, name) {
		fmt.Printf("Error: Ya existe una partición con el nombre '%s'\n", name)
		return
	}

	// Validar restricciones de particionamiento
	if err := validatePartitionStructure(type_, size, mbr, primaryCount, extendedCount, usedSpace, totalPartitions); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Crear la nueva partición
	createPartition(file, &mbr, size, name, type_, fit, totalPartitions)

	// Guardar el MBR actualizado
	if err := Utilities.WriteObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al escribir el MBR actualizado:", err)
		return
	}
	defer file.Close()
	fmt.Println("====== FIN FDISK ======")
}

// convertSize - Convierte el tamaño a bytes
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

// countPartitions - Cuenta el número de particiones y el espacio usado
func countPartitions(mbr Disk.MBR) (primary, extended int, usedSpace int32, total int) {
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

func partitionExists(mbr Disk.MBR, file *os.File, name string) bool {
	for i := 0; i < 4; i++ {
		// Eliminar caracteres nulos y espacios en blanco del nombre de la partición
		existingName := string(bytes.TrimRight(mbr.Partitions[i].Name[:], "\x00"))
		existingName = strings.TrimSpace(existingName)

		fmt.Printf("Comparando partición existente: '%s' con '%s'\n", existingName, name)

		if existingName == name {
			fmt.Printf("Partición con nombre '%s' ya existe.\n", name)
			return true
		}

		// Si es una partición extendida, revisar dentro de la EBR
		if mbr.Partitions[i].Type == 'e' {
			if logicalPartitionExists(file, mbr.Partitions[i], name) {
				return true
			}
		}
	}
	return false
}

func logicalPartitionExists(file *os.File, extendedPartition Disk.Partition, name string) bool {
	var ebr Disk.EBR
	pos := extendedPartition.Start

	for pos != -1 {
		// Leer el EBR
		if err := Utilities.ReadObject(file, &ebr, int64(pos)); err != nil {
			fmt.Println("Error al leer el EBR:", err)
			return false
		}

		// Eliminar caracteres nulos y espacios en blanco del nombre de la partición lógica
		existingName := string(bytes.TrimRight(ebr.PartName[:], "\x00"))
		existingName = strings.TrimSpace(existingName)

		fmt.Printf("Comparando partición lógica existente: '%s' con '%s'\n", existingName, name)

		// Comparar nombres
		if existingName == name {
			fmt.Printf("Partición lógica con nombre '%s' ya existe.\n", name)
			return true
		}

		// Mover al siguiente EBR
		pos = ebr.PartNext
	}

	return false
}

// Valida la estructura de particiones antes de crear una nueva
func validatePartitionStructure(type_ string, size int, mbr Disk.MBR, primary, extended int, usedSpace int32, total int) error {
	if total >= 4 && type_ != "l" {
		return fmt.Errorf("no se pueden crear más de 4 particiones primarias o extendidas")
	}
	if type_ == "e" && extended > 0 {
		return fmt.Errorf("solo se permite una partición extendida por disco")
	}
	if type_ == "l" && extended == 0 {
		return fmt.Errorf("no se pueden crear particiones lógicas sin una extendida")
	}
	if usedSpace+int32(size) > mbr.Size {
		return fmt.Errorf("no hay suficiente espacio en el disco para crear esta partición")
	}
	return nil
}

func createPartition(file *os.File, mbr *Disk.MBR, size int, name, type_, fit string, totalPartitions int) {
	startPos := int32(binary.Size(*mbr))
	if totalPartitions > 0 {
		startPos = mbr.Partitions[totalPartitions-1].Start + mbr.Partitions[totalPartitions-1].Size
	}

	// Buscar un espacio disponible en el MBR para la nueva partición primaria o extendida
	if type_ != "l" {
		for i := 0; i < 4; i++ {
			if mbr.Partitions[i].Size == 0 {
				mbr.Partitions[i] = Disk.Partition{
					Size:   int32(size),
					Start:  startPos,
					Type:   type_[0],
					Fit:    fit[0],
					Status: 0,
				}
				// Copiar el nombre correctamente
				copy(mbr.Partitions[i].Name[:], []byte(name))
				fmt.Printf("Nombre copiado: '%s'\n", string(mbr.Partitions[i].Name[:]))

				// Si es extendida, inicializar el primer EBR
				if type_ == "e" {
					initEBR(file, startPos)
				}
				Disk.PrintPartition(mbr.Partitions[i])
				return
			}
		}
		fmt.Println("Error: No hay espacio disponible para una nueva partición.")
		return
	}

	// Si es una partición lógica, buscar dentro de la extendida
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Type == 'e' {
			createLogicalPartition(file, mbr.Partitions[i], size, name, fit)
			return
		}
	}

	fmt.Println("Error: No se encontró una partición extendida donde agregar la partición lógica.")
}

// Crea una nueva partición lógica dentro de una extendida
func createLogicalPartition(file *os.File, extendedPartition Disk.Partition, size int, name, fit string) {
	var ebr Disk.EBR
	pos := extendedPartition.Start
	prevPos := int32(-1)

	// Buscar el último EBR en la lista
	for pos != -1 {
		if err := Utilities.ReadObject(file, &ebr, int64(pos)); err != nil {
			fmt.Println("Error al leer el EBR:", err)
			return
		}
		if ebr.PartNext == -1 {
			break
		}
		prevPos = pos
		pos = ebr.PartNext
	}

	// Crear nueva partición lógica
	newEBR := Disk.EBR{
		PartFit:   fit[0],
		PartStart: pos + int32(binary.Size(ebr)),
		PartSize:  int32(size),
		PartNext:  -1,
	}
	copy(newEBR.PartName[:], name)

	// Escribir el nuevo EBR en el disco
	Utilities.WriteObject(file, newEBR, int64(newEBR.PartStart))

	// Enlazar con el EBR anterior si existía
	if prevPos != -1 {
		ebr.PartNext = newEBR.PartStart
		Utilities.WriteObject(file, ebr, int64(prevPos))
	}

	fmt.Println("Partición lógica creada correctamente:", name)
}

// Inicializa un EBR en una partición extendida
func initEBR(file *os.File, start int32) {
	ebr := Disk.EBR{
		PartFit:   'b',
		PartStart: start,
		PartSize:  0,
		PartNext:  -1,
	}
	Utilities.WriteObject(file, ebr, int64(start))
}
