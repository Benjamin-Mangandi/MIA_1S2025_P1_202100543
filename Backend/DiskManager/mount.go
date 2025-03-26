package DiskManager

import (
	"Backend/Globals"
	"Backend/Responsehandler"
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
		response := strings.Repeat("*", 30) + "\n" +
			"Error: Partición no encontrada o no es una partición primaria"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Generar el ID del disco y verificar si ya tiene particiones montadas
	diskID := generateDiskID(path)
	if FindPartition(path, name) {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: La particion ya esta montada: " + name + "\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}
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

	// Comentar la siguientes líneas para evitar escribir en el MBR
	// if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
	//     fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
	//     return
	// }
	response := strings.Repeat("-", 40) + "\n" +
		"Partición montada con ID:" + partitionID + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
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

func FindPartition(path, name string) bool {
	// Recorrer las particiones montadas en el mapa
	for _, partitions := range Globals.MountedPartitions {
		// Recorrer cada partición en la lista de particiones
		for _, partition := range partitions {
			// Verificar si tanto el Path como el Name coinciden
			if partition.Path == path && partition.Name == name {
				// Retornar la partición y true si se encontró
				return true
			}
		}
	}
	// Retornar nil y false si no se encontró la partición
	return false
}

func PrintMountedPartitions() {
	// Encabezado
	response := strings.Repeat("-", 40) + "\n" +
		"Particiones montadas:\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)

	// Verificar si hay particiones montadas
	if len(Globals.MountedPartitions) == 0 {
		response := strings.Repeat("-", 40) + "\n" +
			"No hay particiones montadas.\n"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Variable para acumular el string completo
	var fullResponse string

	// Recorrer discos montados
	for diskID, partitions := range Globals.MountedPartitions {
		fullResponse += strings.Repeat("*", 30) + "\n" +
			" * Disco: " + diskID + "\n"

		// Recorrer particiones montadas en el disco
		for _, partition := range partitions {
			partitionInfo := fmt.Sprintf(" - Partición Nombre: %s\nID: %s\nStatus: %c\n",
				partition.Name, partition.ID, partition.Status)
			fullResponse += partitionInfo
		}

	}

	// Al final, añadir toda la información acumulada al global response
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, fullResponse)
}
