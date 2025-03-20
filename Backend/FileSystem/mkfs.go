package FileSystem

import (
	"Backend/DiskManager"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("====== INICIO MKFS ======")
	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)
	fmt.Println("Fs:", fs_)

	// Buscar la partición montada por ID
	mountedPartition, found := findMountedPartition(id)
	if !found {
		fmt.Println("Error: Partición no encontrada o no montada.")
		return
	}

	// Abrir el archivo binario
	file, err := Utilities.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close() // Asegura que el archivo se cierre al finalizar

	// Leer el MBR
	var TempMBR Disk.MBR
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Buscar la partición dentro del MBR
	index := findPartitionName(TempMBR, mountedPartition.Name)
	if index == -1 {
		fmt.Println("Error: Partición no encontrada en el MBR.")
		return
	}

	Disk.PrintPartition(TempMBR.Partitions[index])

	// Validar si la partición ya está marcada como montada
	if TempMBR.Partitions[index].Status == '1' {
		fmt.Println("Advertencia: La partición ya está marcada como montada (Status = 1). No se realizará el formateo.")
		return // Detener el proceso de formateo
	}

	// Validar tipo de sistema de archivos
	if fs_ != "2fs" {
		fmt.Println("Error: Solo está disponible EXT2 (2FS).")
		return
	}

	// Cálculo de inodos
	n := calculateInodes(TempMBR.Partitions[index])
	fmt.Println("INODOS:", n)

	// Crear y configurar el superbloque
	newSuperblock := createSuperblock(n, TempMBR.Partitions[index])

	// Crear el sistema de archivos EXT2
	create_ext2(n, TempMBR.Partitions[index], newSuperblock, "07/03/2025", file)

	// Actualizar el ID y el estado de la partición en el MBR
	copy(TempMBR.Partitions[index].Id[:], []byte(mountedPartition.ID))
	TempMBR.Partitions[index].Status = '1' // '1' indica que está montada

	// Escribir el MBR actualizado en el archivo
	if err := Utilities.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	fmt.Println("Partición actualizada en el MBR:")
	Disk.PrintPartition(TempMBR.Partitions[index])

	fmt.Println("====== FIN MKFS ======")
}

// Busca una partición montada por ID
func findMountedPartition(id string) (Disk.MountedPartition, bool) {
	for _, partitions := range DiskManager.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.ID == id {
				if partition.Status == '1' { // Verifica si está montada
					return partition, true
				}
				fmt.Println("Error: La partición aún no está montada.")
				return Disk.MountedPartition{}, false
			}
		}
	}
	return Disk.MountedPartition{}, false
}

// Encuentra el índice de la partición en el MBR
func findPartitionIndex(mbr Disk.MBR, id string) int {
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Size != 0 && strings.Contains(string(mbr.Partitions[i].Id[:]), id) {
			return i
		}
	}
	return -1
}

func findPartitionName(mbr Disk.MBR, name string) int {
	for i := 0; i < 4; i++ {
		if mbr.Partitions[i].Size != 0 && strings.Contains(string(mbr.Partitions[i].Name[:]), name) {
			return i
		}
	}
	return -1
}

// Calcula la cantidad de inodos basada en el tamaño de la partición
func calculateInodes(partition Disk.Partition) int32 {
	numerador := partition.Size - int32(binary.Size(Ext2.Superblock{}))
	denominador := int32(4 + binary.Size(Ext2.Inode{}) + 3*binary.Size(Ext2.Fileblock{}))
	return numerador / denominador
}

// Crea y configura el superbloque para EXT2
func createSuperblock(n int32, partition Disk.Partition) Ext2.Superblock {
	return Ext2.Superblock{
		S_filesystem_type:   2, // EXT2
		S_inodes_count:      n,
		S_blocks_count:      3 * n,
		S_free_blocks_count: 3*n - 2,
		S_free_inodes_count: n - 2,
		S_mnt_count:         1,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(Ext2.Inode{})),
		S_block_size:        int32(binary.Size(Ext2.Fileblock{})),
		S_bm_inode_start:    partition.Start + int32(binary.Size(Ext2.Superblock{})),
		S_bm_block_start:    partition.Start + int32(binary.Size(Ext2.Superblock{})) + n,
		S_inode_start:       partition.Start + int32(binary.Size(Ext2.Superblock{})) + n + 3*n,
		S_block_start:       partition.Start + int32(binary.Size(Ext2.Superblock{})) + n + 3*n + n*int32(binary.Size(Ext2.Inode{})),
	}
}

func create_ext2(n int32, partition Disk.Partition, newSuperblock Ext2.Superblock, date string, file *os.File) {
	fmt.Println("====== Start CREATE EXT2 ======")
	fmt.Println("INODOS:", n)

	// Imprimir Superblock inicial
	Ext2.PrintSuperblock(newSuperblock)
	fmt.Println("Date:", date)

	// Obtener tamaños de estructuras para evitar cálculos repetidos
	inodeSize := int32(binary.Size(Ext2.Inode{}))
	blockSize := int32(binary.Size(Ext2.Fileblock{}))
	folderBlockSize := int32(binary.Size(Ext2.Folderblock{}))

	// Escribir bitmaps de inodos y bloques con una sola llamada a WriteObject por bloque
	inodeBitmap := make([]byte, n)
	blockBitmap := make([]byte, 3*n)
	if err := Utilities.WriteObject(file, inodeBitmap, int64(newSuperblock.S_bm_inode_start)); err != nil {
		fmt.Println("Error al escribir bitmap de inodos: ", err)
		return
	}
	if err := Utilities.WriteObject(file, blockBitmap, int64(newSuperblock.S_bm_block_start)); err != nil {
		fmt.Println("Error al escribir bitmap de bloques: ", err)
		return
	}

	// Inicializa inodos y bloques con valores predeterminados
	if err := initInodesAndBlocks(n, newSuperblock, file); err != nil {
		fmt.Println("Error al inicializar inodos y bloques: ", err)
		return
	}

	// Crea la carpeta raíz y el archivo users.txt
	if err := createRootAndUsersFile(newSuperblock, date, file); err != nil {
		fmt.Println("Error al crear la carpeta raíz y users.txt: ", err)
		return
	}

	// Escribe el superbloque actualizado al archivo
	if err := Utilities.WriteObject(file, newSuperblock, int64(partition.Start)); err != nil {
		fmt.Println("Error al escribir el Superblock: ", err)
		return
	}

	// Marca los primeros inodos y bloques como usados
	if err := markUsedInodesAndBlocks(newSuperblock, file); err != nil {
		fmt.Println("Error al marcar inodos y bloques usados: ", err)
		return
	}

	// Leer e imprimir los inodos después de formatear
	fmt.Println("====== Imprimiendo Inodos ======")
	for i := int32(0); i < n; i++ {
		var inode Ext2.Inode
		offset := int64(newSuperblock.S_inode_start + i*inodeSize)
		if err := Utilities.ReadObject(file, &inode, offset); err != nil {
			fmt.Println("Error al leer inodo:", err)
			return
		}
		Ext2.PrintInode(inode)
	}

	// Leer e imprimir los Folderblocks y Fileblocks después de formatear
	fmt.Println("====== Imprimiendo Folderblocks y Fileblocks ======")

	// Imprimir Folderblock en la primera posición
	var folderblock Ext2.Folderblock
	folderOffset := int64(newSuperblock.S_block_start)
	if err := Utilities.ReadObject(file, &folderblock, folderOffset); err != nil {
		fmt.Println("Error al leer Folderblock:", err)
		return
	}
	Ext2.PrintFolderblock(folderblock)

	// Imprimir Fileblock en la siguiente posición
	var fileblock Ext2.Fileblock
	fileOffset := int64(newSuperblock.S_block_start + folderBlockSize)
	if err := Utilities.ReadObject(file, &fileblock, fileOffset); err != nil {
		fmt.Println("Error al leer Fileblock:", err)
		return
	}
	Ext2.PrintFileblock(fileblock)

	// Imprimir los siguientes Fileblocks (si hay más de uno)
	for i := int32(1); i < 3*n; i++ {
		var fileblock Ext2.Fileblock
		fileOffset := int64(newSuperblock.S_block_start + folderBlockSize + i*blockSize)
		if err := Utilities.ReadObject(file, &fileblock, fileOffset); err != nil {
			fmt.Println("Error al leer Fileblock:", err)
			return
		}
		Ext2.PrintFileblock(fileblock)
	}

	// Imprimir el Superblock final
	Ext2.PrintSuperblock(newSuperblock)

	fmt.Println("====== End CREATE EXT2 ======")
}

// Función auxiliar para inicializar inodos y bloques
func initInodesAndBlocks(n int32, newSuperblock Ext2.Superblock, file *os.File) error {
	// Inicializar un inodo por defecto
	newInode := Ext2.Inode{I_block: [15]int32{-1}}

	// Escribir todos los inodos en un solo bucle
	inodeSize := int32(binary.Size(Ext2.Inode{}))
	for i := int32(0); i < n; i++ {
		offset := int64(newSuperblock.S_inode_start + i*inodeSize)
		if err := Utilities.WriteObject(file, newInode, offset); err != nil {
			return err
		}
	}

	// Inicializar un bloque de archivo vacío
	newFileblock := Ext2.Fileblock{}

	// Escribir todos los bloques en un solo bucle
	blockSize := int32(binary.Size(Ext2.Fileblock{}))
	for i := int32(0); i < 3*n; i++ {
		offset := int64(newSuperblock.S_block_start + i*blockSize)
		if err := Utilities.WriteObject(file, newFileblock, offset); err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock Ext2.Superblock, date string, file *os.File) error {
	// Inicializar los inodos
	inode0 := Ext2.Inode{I_block: [15]int32{-1}}
	inode1 := Ext2.Inode{I_block: [15]int32{-1}}

	initInode(&inode0, date)
	initInode(&inode1, date)

	inode0.I_block[0] = 0
	inode1.I_block[0] = 1
	inode1.I_size = int32(len("1,G,root\n1,U,root,root,123\n")) // Asignar tamaño real

	// Crear bloque de carpeta
	var folderBlock Ext2.Folderblock
	folderBlock.B_content[0].B_inodo = 0
	copy(folderBlock.B_content[0].B_name[:], ".")
	folderBlock.B_content[1].B_inodo = 0
	copy(folderBlock.B_content[1].B_name[:], "..")
	folderBlock.B_content[2].B_inodo = 1
	copy(folderBlock.B_content[2].B_name[:], "users.txt")

	// Crear bloque de archivo con datos
	var fileBlock Ext2.Fileblock
	copy(fileBlock.B_content[:], "1,G,root\n1,U,root,root,123\n")

	// Escribir en el archivo
	objects := []struct {
		data   interface{}
		offset int64
	}{
		{inode0, int64(newSuperblock.S_inode_start)},
		{inode1, int64(newSuperblock.S_inode_start + int32(binary.Size(Ext2.Inode{})))},
		{folderBlock, int64(newSuperblock.S_block_start)},
		{fileBlock, int64(newSuperblock.S_block_start + int32(binary.Size(Ext2.Folderblock{})))},
	}

	for _, obj := range objects {
		if err := Utilities.WriteObject(file, obj.data, obj.offset); err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para inicializar un inodo
func initInode(inode *Ext2.Inode, date string) {
	*inode = Ext2.Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_block: [15]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	}
	copy(inode.I_atime[:], date)
	copy(inode.I_ctime[:], date)
	copy(inode.I_mtime[:], date)
	copy(inode.I_perm[:], "664")
}

// Función auxiliar para marcar los inodos y bloques usados
func markUsedInodesAndBlocks(newSuperblock Ext2.Superblock, file *os.File) error {
	positions := []int64{
		int64(newSuperblock.S_bm_inode_start),
		int64(newSuperblock.S_bm_inode_start + 1),
		int64(newSuperblock.S_bm_block_start),
		int64(newSuperblock.S_bm_block_start + 1),
	}

	for _, pos := range positions {
		if err := Utilities.WriteObject(file, byte(1), pos); err != nil {
			return err
		}
	}
	return nil
}
