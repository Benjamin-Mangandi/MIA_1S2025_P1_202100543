package DiskManager

import (
	"Backend/Globals"
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"bytes"
	"fmt"
	"strings"
)

func Mount(path string, name string) {
	// Abrir el archivo del disco
	file, err := Utilities.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Disk.MBR
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Printf("Buscando partición con nombre: '%s'\n", name)

	// Buscar la partición por nombre
	nameBytes := [16]byte{}
	copy(nameBytes[:], []byte(name))
	partitionFound := false
	var partition Disk.Partition
	var partitionIndex int

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Type == 'p' && bytes.Equal(TempMBR.Partitions[i].Name[:], nameBytes[:]) {
			partition = TempMBR.Partitions[i]
			partitionIndex = i
			partitionFound = true
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: Partición no encontrada o no es una partición primaria")
		return
	}

	// Generar el ID del disco y verificar si ya tiene particiones montadas
	diskID := generateDiskID(path)
	mountedPartitionsInDisk := Globals.MountedPartitions[diskID]

	// Determinar la letra para el disco
	var letter byte
	if len(mountedPartitionsInDisk) == 0 {
		lastDiskID := getLastDiskID()
		if lastDiskID == "" {
			letter = 'a'
		} else {
			lastLetter := Globals.MountedPartitions[lastDiskID][0].ID[len(Globals.MountedPartitions[lastDiskID][0].ID)-1]
			letter = lastLetter + 1
		}
	} else {
		letter = mountedPartitionsInDisk[0].ID[len(mountedPartitionsInDisk[0].ID)-1]
	}

	// Crear el ID de la partición
	carnet := "202100543"
	lastTwoDigits := carnet[len(carnet)-2:]
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)

	// Montar la partición
	partition.Status = '1'
	copy(partition.Id[:], partitionID)
	TempMBR.Partitions[partitionIndex] = partition

	Globals.MountedPartitions[diskID] = append(Globals.MountedPartitions[diskID], Disk.MountedPartition{
		Path:   path,
		Name:   name,
		ID:     partitionID,
		Status: '1',
		Start:  partition.Start,
	})

	// Comentar o eliminar la siguiente línea para evitar escribir en el MBR
	// if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
	//     fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
	//     return
	// }

	fmt.Printf("Partición montada con ID: %s\n", partitionID)
}

// Obtener todas las particiones montadas
func GetMountedPartitions() map[string][]Disk.MountedPartition {
	return Globals.MountedPartitions
}

func GetMountedPartitionByID(id string) Disk.MountedPartition {
	for _, partitions := range Globals.MountedPartitions { // Itera sobre los valores del mapa
		for _, partition := range partitions { // Itera sobre la lista de particiones
			if partition.ID == id {
				return partition
			}
		}
	}
	return Disk.MountedPartition{} // Devuelve un objeto vacío si no se encuentra
}
func getLastDiskID() string {
	for diskID := range Globals.MountedPartitions {
		return diskID // Devuelve el primer encontrado
	}
	return "" // Si no hay discos montados
}

func generateDiskID(path string) string {
	return strings.ToLower(path)
}

func PrintMountedPartitions() {
	fmt.Println("Particiones montadas:")

	if len(Globals.MountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	for diskID, partitions := range Globals.MountedPartitions {
		fmt.Printf("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			fmt.Printf(" - Partición Name: %s, ID: %s, Path: %s, Status: %c\n",
				partition.Name, partition.ID, partition.Path, partition.Status)
		}
	}
	fmt.Println("")
}
