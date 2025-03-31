package FileSystem

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkdir(path string, p string) {
	// Verificar sesión activa (tu código existente)
	if !Globals.ActiveUser.Status {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No hay un usuario activo."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Obtener partición montada (tu código existente)
	mountedPartition := DiskManager.GetMountedPartitionByID(Globals.ActiveUser.PartitionID)
	if mountedPartition.ID == "" {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No se encontró la partición montada."
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Abrir archivo del disco (tu código existente)
	file, err := Utilities.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo del disco.")
		return
	}
	defer file.Close()

	// Leer superbloque (tu código existente)
	var sb Ext2.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque.")
		return
	}
	var createParents bool
	if p == "true" {
		createParents = true
	} else {
		createParents = false
	}
	createDirectory(path, createParents, &sb, file, mountedPartition)
	response := strings.Repeat("-", 40) + "\n" +
		fmt.Sprintf("Directorio creado exitosamente: %s\n", path)
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}

func createDirectory(dirPath string, createParents bool, sb *Ext2.Superblock, file *os.File, mountedPartition Disk.MountedPartition) {
	// Obtener directorios padres y el destino del directorio
	parentDirs, destDir := Utilities.GetParentDirectories(dirPath)

	// Si el parámetro -p está habilitado (createParents == true), crear los directorios intermedios
	if createParents {
		for _, parentDir := range parentDirs {
			fmt.Println(parentDir)
			CreateFolder(file, parentDirs[1:], parentDir, sb, mountedPartition)
		}
	}

	// Crear el directorio final
	CreateFolder(file, parentDirs, destDir, sb, mountedPartition)

	// Escribe el superbloque actualizado al archivo
	if err := Utilities.WriteObject(file, sb, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al escribir el Superblock: ", err)
		return
	}
}

func CreateFolder(file *os.File, parentsDir []string, destDir string, sb *Ext2.Superblock, mountedPartition Disk.MountedPartition) {
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		createFolderInInode(file, 0, parentsDir, destDir, sb, mountedPartition)
		return
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ { //Desde el inodo 0
		createFolderInInode(file, i, parentsDir, destDir, sb, mountedPartition)
	}
}

func createFolderInInode(file *os.File, inodeIndex int32, parentsDir []string, destDir string, sb *Ext2.Superblock, mountedPartition Disk.MountedPartition) {
	inode := Ext2.Inode{}
	err := Utilities.ReadObject(file, &inode, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] != '0' {

		return
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}
		// Crear un nuevo bloque de carpeta
		block := Ext2.Folderblock{}

		err := Utilities.ReadObject(file, &block, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
		if err != nil {
			return
		}

		// Iterar sobre cada contenido del bloque, desde el índice 2 (evitamos . y ..)
		for indexContent := 0; indexContent < len(block.B_content); indexContent++ {
			content := block.B_content[indexContent]

			// Si hay más carpetas padres en la ruta
			if len(parentsDir) != 0 {
				// Si el contenido está vacío, salir
				if content.B_inodo == -1 {
					fmt.Printf("No se encontró carpeta padre en inodo %d en la posición %d, terminando.\n", inodeIndex, indexContent) // Depuración
					break
				}

				// Obtener la carpeta padre más cercana
				parentDir, err := Utilities.First(parentsDir)
				if err != nil {
					return
				}

				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				parentDirName := strings.Trim(parentDir, "\x00 ")

				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					fmt.Printf("Carpeta padre '%s' encontrada en inodo %d. Recursion para crear el siguiente directorio.\n", parentDirName, content.B_inodo) // Depuración
					// Llamada recursiva para seguir creando carpetas
					createFolderInInode(file, content.B_inodo, Utilities.RemoveElement(parentsDir, 0), destDir, sb, mountedPartition)
					return
				}
			} else {
				if content.B_inodo != -1 {
					continue
				}

				fmt.Printf("Asignando el nombre del directorio '%s' al bloque en la posición %d\n", destDir, indexContent) // Depuración
				// Actualizar el contenido del bloque con el nuevo directorio
				copy(content.B_name[:], destDir)
				content.B_inodo = sb.S_fist_ino

				// Actualizar el bloque con el nuevo contenido
				block.B_content[indexContent] = content

				// Escribir el bloque actualizado de vuelta en el archivo
				err = Utilities.WriteObject(file, block, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					fmt.Println("Error al escribir el bloque de carpeta:", err)
					return
				}
				newBlockIndex := sb.S_first_blo // Obtener el primer bloque libre
				// Crear el inodo de la nueva carpeta
				folderInode := Ext2.Inode{
					I_uid:   1,
					I_gid:   1,
					I_size:  0,
					I_block: [15]int32{newBlockIndex, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
					I_type:  [1]byte{'0'}, // Tipo carpeta
					I_perm:  [3]byte{'6', '6', '4'},
				}
				currentTime := time.Now().Format("2006-01-02 15:04")
				var now [16]byte
				copy(now[:], []byte(currentTime)) // Convertimos el string a bytes correctamente.

				copy(folderInode.I_atime[:], now[:])
				copy(folderInode.I_ctime[:], now[:])
				copy(folderInode.I_mtime[:], now[:])

				err = Utilities.WriteObject(file, folderInode, int64(sb.S_inode_start)+int64(sb.S_fist_ino)*int64(sb.S_inode_size))
				if err != nil {
					return
				}

				//Crear el bloque para la nueva carpeta
				folderBlock := Ext2.Folderblock{
					B_content: [4]Ext2.Content{
						{B_name: [12]byte{'.'}, B_inodo: content.B_inodo},
						{B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
						{B_name: [12]byte{'-'}, B_inodo: -1},
						{B_name: [12]byte{'-'}, B_inodo: -1},
					},
				}

				err = Utilities.WriteObject(file, folderBlock, int64(sb.S_block_start)+int64(sb.S_first_blo)*int64(sb.S_block_size))
				if err != nil {
					return
				}
				for j, block := range inode.I_block {
					if block == -1 {
						inode.I_block[j] = sb.S_fist_ino
						break
					}
				}
				inodePosition := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(Ext2.Inode{}))
				err = Utilities.WriteObject(file, &inode, inodePosition)
				if err != nil {
					return
				}

				// Actualizar el bitmap de bloques
				UpdateBitmap(file, *sb, int(sb.S_fist_ino), int(sb.S_first_blo))
				UpdateSuperblockAfterBlockAllocation(sb)
				UpdateSuperblockAfterInodeAllocation(sb)
				err = Utilities.WriteObject(file, sb, int64(mountedPartition.Start))
				if err != nil {
					return
				}
				return
			}
		}
	}

	fmt.Printf("No se encontraron bloques disponibles para crear la carpeta '%s' en inodo %d\n", destDir, inodeIndex) // Depuración
}

func UpdateBitmap(file *os.File, sb Ext2.Superblock, inodeposition int, blockposition int) {
	inodeBitmap := make([]byte, (sb.S_inodes_count+7)/8) // Redondeamos al byte más cercano
	blockBitmap := make([]byte, (sb.S_blocks_count+7)/8) // Redondeamos al byte más cercano

	// Leer los bitmaps desde el archivo
	if _, err := file.ReadAt(inodeBitmap, int64(sb.S_bm_inode_start)); err != nil {
		fmt.Println("Error al leer bitmap de inodos:", err)
		return
	}
	if _, err := file.ReadAt(blockBitmap, int64(sb.S_bm_block_start)); err != nil {
		fmt.Println("Error al leer bitmap de bloques:", err)
		return
	}

	// Función auxiliar para marcar un bit en una posición específica
	markBit := func(bitmap []byte, pos int) {
		bytePos := pos / 8               // Encuentra el byte correspondiente
		bitPos := pos % 8                // Encuentra la posición del bit dentro del byte
		bitmap[bytePos] |= (1 << bitPos) // Marca el bit (ponlo en 1)
	}

	// Marcar los inodos usados
	inodePositions := []int{
		inodeposition, // Inodo raíz "/"
	}
	for _, pos := range inodePositions {
		markBit(inodeBitmap, pos)
	}

	// Marcar los bloques usados
	blockPositions := []int{
		blockposition, // Bloque de la carpeta raíz
	}
	for _, pos := range blockPositions {
		markBit(blockBitmap, pos)
	}

	// Escribir los bitmaps modificados de vuelta al archivo
	if _, err := file.WriteAt(inodeBitmap, int64(sb.S_bm_inode_start)); err != nil {
		fmt.Println("Error al escribir bitmap de inodos:", err)
		return
	}
	if _, err := file.WriteAt(blockBitmap, int64(sb.S_bm_block_start)); err != nil {
		fmt.Println("Error al escribir bitmap de bloques:", err)
		return
	}
}

func UpdateSuperblockAfterBlockAllocation(sb *Ext2.Superblock) {
	sb.S_free_blocks_count--
	sb.S_first_blo++
}

func UpdateSuperblockAfterInodeAllocation(sb *Ext2.Superblock) {
	sb.S_free_inodes_count--
	sb.S_fist_ino++
}
