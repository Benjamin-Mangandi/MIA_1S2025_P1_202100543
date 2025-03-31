package FileSystem

import (
	"Backend/DiskManager"
	"Backend/Globals"
	"Backend/Responsehandler"
	Disk "Backend/Structs/disk"
	Ext2 "Backend/Structs/ext2"
	"Backend/Utilities"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Mkfile(path string, r string, size string, cont string) {
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

	dirPath, _ := GetDirectoryAndFile(path)
	// Leer superbloque (tu código existente)
	var sb Ext2.Superblock
	if err := Utilities.ReadObject(file, &sb, int64(mountedPartition.Start)); err != nil {
		fmt.Println("Error al leer el superbloque.")
		return
	}
	exists, _ := directoryExists(&sb, file, 0, dirPath) // Usamos el inodo raíz (0) para empezar la búsqueda
	if !exists {
		return
	}

	// Si -r está habilitado y el directorio no existe, creamos los directorios intermedios
	if r == "true" && !exists {

	}

	// string to int
	newsize, err := strconv.Atoi(size)
	if err != nil {
		panic(err)
	}
	if cont == "" {
		cont = generateContent(newsize)
	}
	// Crear el archivo usando el archivo de partición abierto
	createFile(path, newsize, cont, &sb, file, mountedPartition)
}

// generateContent genera una cadena de números del 0 al 9 hasta cumplir el tamaño ingresado
func generateContent(size int) string {
	content := ""
	for len(content) < size {
		content += "0123456789"
	}
	return content[:size] // Recorta la cadena al tamaño exacto
}

func GetDirectoryAndFile(path string) (string, string) {
	// Obtener la carpeta donde se creará el archivo
	dir := filepath.Dir(path)
	// Obtener el nombre del archivo
	file := filepath.Base(path)
	return dir, file
}

func directoryExists(sb *Ext2.Superblock, file *os.File, inodeIndex int32, dirName string) (bool, int32) {

	// Deserializar el inodo correspondiente
	inode := Ext2.Inode{}
	err := Utilities.ReadObject(file, &inode, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return false, -1
	}

	// Verificar si el inodo es de tipo carpeta (I_type == '0') para continuar
	if inode.I_type[0] != '0' {
		return false, -1
	}

	// Iterar sobre los bloques del inodo para buscar el directorio o archivo
	for _, blockIndex := range inode.I_block {
		if blockIndex == -1 {
			break // Si no hay más bloques asignados, terminamos la búsqueda
		}

		// Deserializar el bloque de directorio
		block := Ext2.Folderblock{}
		err := Utilities.ReadObject(file, &block, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
		if err != nil {
			return false, -1
		}

		// Iterar sobre los contenidos del bloque para verificar si el nombre coincide
		for _, content := range block.B_content {
			contentName := strings.Trim(string(content.B_name[:]), "\x00 ") // Convertir el nombre y eliminar los caracteres nulos
			if strings.EqualFold(contentName, dirName) && content.B_inodo != -1 {
				fmt.Printf("Directorio o archivo '%s' encontrado en inodo %d\n", dirName, content.B_inodo) // Depuración
				return true, content.B_inodo                                                               // Devolver true si el directorio/archivo fue encontrado
			}
		}
	}

	fmt.Printf("Directorio o archivo '%s' no encontrado en inodo %d\n", dirName, inodeIndex) // Depuración
	return false, -1                                                                         // No se encontró el directorio/archivo
}

// createFile ahora usa el archivo de partición ya abierto
func createFile(filePath string, size int, content string, sb *Ext2.Superblock, file *os.File, mountedPartition Disk.MountedPartition) {

	// Obtener los directorios padres y el destino
	parentDirs, destDir := Utilities.GetParentDirectories(filePath)
	// Obtener contenido por chunks
	chunks := Utilities.SplitStringIntoChunks(content)

	// Crear el archivo en el sistema de archivos
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentDirs) == 0 {
		createFileInInode(file, 0, parentDirs, destDir, size, chunks, sb, mountedPartition)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		createFileInInode(file, i, parentDirs, destDir, size, chunks, sb, mountedPartition)
	}

	err := Utilities.WriteObject(file, sb, int64(mountedPartition.Start))
	if err != nil {
		return
	}
}

func createFileInInode(file *os.File, inodeIndex int32, parentsDir []string, destFile string, fileSize int, fileContent []string, sb *Ext2.Superblock, mountedPartition Disk.MountedPartition) {
	inode := Ext2.Inode{}
	err := Utilities.ReadObject(file, &inode, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] == '1' {
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
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
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
					fmt.Printf("Encontrada carpeta padre '%s' en inodo %d\n", parentDirName, content.B_inodo) // Depuración
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					createFileInInode(file, content.B_inodo, Utilities.RemoveElement(parentsDir, 0), destFile, fileSize, fileContent, sb, mountedPartition)
					return
				} else {
					if content.B_inodo != -1 {
						continue
					}

					// Actualizar el contenido del bloque con el nuevo directorio
					copy(content.B_name[:], []byte(destFile))
					content.B_inodo = sb.S_fist_ino

					// Actualizar el bloque con el nuevo contenido
					block.B_content[indexContent] = content

					// Escribir el bloque actualizado de vuelta en el archivo
					err = Utilities.WriteObject(file, block, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
					if err != nil {
						fmt.Println("Error al escribir el bloque de carpeta:", err)
						return
					}
					// Crear el inodo de la nueva carpeta
					fileInode := Ext2.Inode{
						I_uid:   1,
						I_gid:   1,
						I_size:  0,
						I_block: [15]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
						I_type:  [1]byte{'0'}, // Tipo carpeta
						I_perm:  [3]byte{'6', '6', '4'},
					}
					currentTime := time.Now().Format("2006-01-02 15:04")
					var now [16]byte
					copy(now[:], []byte(currentTime)) // Convertimos el string a bytes correctamente.

					copy(fileInode.I_atime[:], now[:])
					copy(fileInode.I_ctime[:], now[:])
					copy(fileInode.I_mtime[:], now[:])
					// Crear los bloques del archivo
					for i := 0; i < len(fileContent); i++ {
						fileInode.I_block[i] = sb.S_blocks_count

						// Crear el bloque del archivo
						fileBlock := Ext2.Fileblock{
							B_content: [64]byte{},
						}
						copy(fileBlock.B_content[:], fileContent[i])

						err = Utilities.WriteObject(file, fileBlock, int64(sb.S_block_start)+int64(sb.S_first_blo)*int64(sb.S_block_size))
						if err != nil {
							return
						}

						fmt.Printf("Bloque de archivo '%s' serializado correctamente.\n", destFile) // Depuración

						UpdateSuperblockAfterBlockAllocation(sb)
					}
					err = Utilities.WriteObject(file, fileInode, int64(sb.S_inode_start)+int64(sb.S_fist_ino)*int64(sb.S_inode_size))
					if err != nil {
						return
					}

					// Actualizar el bitmap de inodos
					UpdateBitmap(file, *sb, int(sb.S_fist_ino), int(sb.S_first_blo))

					// Actualizar el superbloque
					UpdateSuperblockAfterInodeAllocation(sb)
					err = Utilities.WriteObject(file, sb, int64(mountedPartition.Start))
					if err != nil {
						return
					}
					return
				}
			}
		}
	}
}
