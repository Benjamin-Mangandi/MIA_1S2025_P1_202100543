package DiskManager

import (
	"Backend/Globals"
	Disk "Backend/Structs/disk"
	"Backend/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Mkdisk crea un nuevo disco binario con el formato EXT2 simulado
func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("====== INICIO MKDISK ======")
	fmt.Printf("Size: %d\nFit: %s\nUnit: %s\nPath: %s\n", size, fit, unit, path)

	// Convertir fit y unit a minúsculas para evitar errores por mayúsculas
	fit = strings.ToLower(fit)
	unit = strings.ToLower(unit)

	// Validar fit (bf, wf, ff)
	validFits := map[string]bool{"bf": true, "wf": true, "ff": true}
	if !validFits[fit] {
		fmt.Println("Error: Fit debe ser 'bf', 'wf' o 'ff'.")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: El tamaño del disco debe ser mayor a 0.")
		return
	}

	// Validar unit (k o m)
	var multiplier int
	switch unit {
	case "k":
		multiplier = 1024
	case "m":
		multiplier = 1024 * 1024
	default:
		fmt.Println("Error: Las unidades válidas son 'k' (kilobytes) o 'm' (megabytes).")
		return
	}

	// Calcular el tamaño real en bytes
	sizeInBytes := size * multiplier

	// Crear archivo binario
	if err := Utilities.CreateFile(path); err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}

	// Abrir archivo binario
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Escribir bytes vacíos de manera eficiente
	emptyBlock := make([]byte, 1024) // Bloque de 1 KB lleno de ceros
	for i := 0; i < sizeInBytes; i += 1024 {
		_, err := file.Write(emptyBlock)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	}

	// Crear estructura MBR
	var newMBR Disk.MBR
	newMBR.Size = int32(sizeInBytes)
	newMBR.Signature = rand.Int31() // Número aleatorio único para el disco
	newMBR.Fit = fit[0]             // Almacenar solo el primer carácter del fit

	// Obtener la fecha en formato YYYY-MM-DD HH:MM y almacenarla en bytes
	currentTime := time.Now().Format("2006-01-02 15:04")
	copy(newMBR.CreationDate[:], currentTime)

	// Escribir el MBR en el archivo
	if err := Utilities.WriteObject(file, newMBR, 0); err != nil {
		fmt.Println("Error al escribir el MBR en el archivo:", err)
		return
	}

	// Leer el MBR para verificar su correcta escritura
	var tempMBR Disk.MBR
	if err := Utilities.ReadObject(file, &tempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Imprimir el MBR leído
	Disk.PrintMBR(tempMBR)

	fmt.Println("====== FIN MKDISK ======")
}

func RmDisk(path string) {

	// Verificar si el archivo existe antes de eliminarlo
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Error: No se encontró el disco en la ruta '%s'.\n", path)
		return
	}

	// Confirmar con el usuario antes de eliminar
	fmt.Printf("¿Está seguro de que desea eliminar el disco en '%s'? (s/n): ", path)
	var response string
	fmt.Scanln(&response)

	if response != "s" && response != "S" {
		fmt.Println("Operación cancelada.")
		return
	}

	// Intentar eliminar el archivo
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error: No se pudo eliminar el disco. %v\n", err)
		return
	}

	fmt.Println("Disco eliminado con éxito.")
}

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

// Montar una partición en un disco
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

	// Verificar si la partición ya está montada
	if partition.Status == '1' {
		fmt.Println("Error: La partición ya está montada")
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

	// Guardar cambios en el MBR
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	fmt.Printf("Partición montada con ID: %s\n", partitionID)
	fmt.Println("\nMBR actualizado:")
	Disk.PrintMBR(TempMBR)
}
