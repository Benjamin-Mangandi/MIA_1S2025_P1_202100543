package Scanner

import (
	"Backend/DiskManager"
	"Backend/FileSystem"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Expresión regular para capturar parámetros en el formato -key=value
var regex = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// Extrae el comando y sus parámetros de la entrada del usuario
func getCommandAndParams(input string) (string, string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}

	parts := strings.Fields(input)
	command := strings.ToLower(parts[0]) // Normalizar a minúsculas
	params := strings.Join(parts[1:], " ")

	return command, params
}

// Escanea y procesa los comandos ingresados por el usuario
func Scan() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("======================")
		fmt.Print("Ingrese comando: ")

		// Leer la línea completa de entrada del usuario
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}

		// Procesar el comando
		command, params := getCommandAndParams(input)

		if command == "" {
			fmt.Println("Error: No se ingresó un comando válido.")
			continue
		}

		fmt.Printf("Comando: %s\nParámetros: %s\n", command, params)

		// Llamar a la función para analizar el comando
		AnalyzeCommand(command, params)

	}
}

func AnalyzeCommand(command string, params string) {
	switch {
	case strings.EqualFold(command, "mkdisk"):
		Commandmkdisk(params)
	case strings.EqualFold(command, "fdisk"):
		Commandfdisk(params)
	case strings.EqualFold(command, "rmdisk"):
		CommandRmDisk(params)
	case strings.EqualFold(command, "mounted"):
		Commandmounted(params)
	case strings.EqualFold(command, "mount"):
		Commandmount(params)
	case strings.EqualFold(command, "mkfs"):
		Commandmkfs(params)
	case strings.EqualFold(command, "login"):
		CommandLogin(params)
	case strings.EqualFold(command, "Logout"):
		CommandLogOut(params)
	default:
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}

// Procesa el comando mkdisk
func Commandmkdisk(params string) {
	// Definir flags con valores predeterminados
	fs := flag.NewFlagSet("mkdisk", flag.ContinueOnError)
	fit := fs.String("fit", "ff", "Tipo de ajuste (bf, ff, wf)")
	unit := fs.String("unit", "m", "Unidad de tamaño (k, m)")
	path := fs.String("path", "", "Ruta donde se creará el disco")
	size := fs.Int("size", 0, "Tamaño del disco")

	fs.Parse(os.Args[1:])
	matches := regex.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas
		fs.Set(flagName, flagValue)
	}
	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor a 0")
		return
	}

	if *path == "" {
		fmt.Println("Error: Se requiere especificar un path")
		return
	}

	// Validar unidad
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unidad inválida, debe ser 'k' o 'm'")
		return
	}

	// Llamar a la función para crear el disco
	DiskManager.Mkdisk(*size, *fit, *unit, *path)
	//mkdisk -Size=3000 -unit=K -path=/home/benjamin/discos/Disco1.mia
	//mkdisk -path=/home/benjamin/discos/Disco2.mia -Unit=K -size=3000
	//mkdisk -size=5 -unit=M -path="/home/benjamin/discos/Disco3.mia"​
	//mkdisk -size=10 -Path="/home/benjamin/discos/Disco4.mia"
	//mkdisk -size=1024 -unit=M -fit=BF -path="/home/benjamin/discos/DiscoDefinitivo.mia"
}

func Commandfdisk(input string) {
	// Definir flags con valores predeterminados
	fs := flag.NewFlagSet("fdisk", flag.ContinueOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "k", "Unidad (b/k/m)")
	type_ := fs.String("type", "p", "Tipo (p/e/l)")
	fit := fs.String("fit", "wf", "Ajuste (bf/ff/wf)")

	// Expresión regular para capturar los flags en el input
	matches := regex.FindAllStringSubmatch(input, -1)

	// Mapas para validaciones rápidas
	validUnits := map[string]bool{"b": true, "k": true, "m": true}
	validTypes := map[string]bool{"p": true, "e": true, "l": true}
	validFits := map[string]bool{"bf": true, "ff": true, "wf": true}

	// Procesar los flags extraídos
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre a minúsculas
		flagValue := strings.ToLower(match[2])    // Convertir valor a minúsculas si aplica
		flagValue = strings.Trim(flagValue, "\"") // Eliminar comillas

		// Intentar asignar el valor al flag correspondiente
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Advertencia: Flag desconocido '%s'\n", flagName)
		}
	}

	// Validaciones de parámetros
	if *size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor a 0")
		return
	}
	if *path == "" {
		fmt.Println("Error: Se requiere la ruta del archivo")
		return
	}
	if !validUnits[*unit] {
		fmt.Println("Error: Unidad inválida, debe ser 'b', 'k' o 'm'")
		return
	}
	if !validTypes[*type_] {
		fmt.Println("Error: Tipo inválido, debe ser 'p' (primaria), 'e' (extendida) o 'l' (logica)")
		return
	}
	if !validFits[*fit] {
		fmt.Println("Errora: Ajuste inválido, debe ser 'bf', 'ff' o 'wf'")
		return
	}

	// Llamar a la función Fdisk con los valores procesados
	DiskManager.Fdisk(*size, *path, *name, *unit, *type_, *fit)
	//fdisk -size=600 -path="/home/benjamin/discos/Disco1.mia" -name=Particion1
	//fdisk -Size=2000 -path=/home/benjamin/discos/Disco1.mia -name=Particion2
	//fdisk -type=E -path=/home/benjamin/discos/Disco1.mia -Unit=K -name=Particion3 -size=300
	//fdisk -size=1 -type=L -unit=M -fit=BF -path=/home/benjamin/discos/Disco1.mia -name="Particion4"
	//fdisk -path=/home/benjamin/discos/Disco2.mia -name=Part3 -Unit=K -size=200
}

func CommandRmDisk(input string) {
	// Definir flags
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco a eliminar")

	// Buscar y extraer los flags del input
	args := strings.Fields(input)
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-path=") {
			*path = strings.TrimPrefix(args[i], "-path=")
			*path = strings.Trim(*path, "\"") // Eliminar comillas si las hay
		}
	}

	// Validar si se proporcionó el path
	if *path == "" {
		fmt.Println("Error: El parámetro -path es obligatorio.")
		return
	}

	// Llamar a la función en DiskManager para eliminar el disco
	DiskManager.RmDisk(*path)
	//rmdisk -path="/home/benjamin/discos/disco1.mia"
}

func Commandmount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco")
	name := fs.String("name", "", "Nombre de la partición")

	// Extraer los parámetros usando regex
	matches := regex.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])     // Convertir nombre del flag a minúsculas
		flagValue := strings.Trim(match[2], "\"") // Eliminar comillas

		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
	}

	// Validación de campos obligatorios
	if *path == "" || *name == "" {
		fmt.Println("Error: Los parámetros '-path' y '-name' son obligatorios")
		return
	}

	// Llamar a la función Mount con el nombre en minúsculas
	DiskManager.Mount(*path, strings.ToLower(*name))

}

//mount -path="/home/benjamin/discos/disco1.mia" -name=particion1
//mount -path="/home/benjamin/discos/disco1.mia" -name=particion2
//mount -path="/home/benjamin/discos/disco2.mia" -name=part3
//mount -path="/home/benjamin/discos/disco1.mia" -name=part3

func Commandmounted(params string) {
	// Verificar que no se pasen parámetros
	if strings.TrimSpace(params) != "" {
		fmt.Println("Error: El comando 'mounted' no acepta parámetros")
		return
	}
	// Llamar a la función para imprimir las particiones montadas
	DiskManager.PrintMountedPartitions()
	//mounted
}

func Commandmkfs(params string) {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "full", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Convertir parámetros a minúsculas para hacer la búsqueda insensible a mayúsculas
	params = strings.ToLower(params)

	// Expresión regular para extraer parámetros
	matches := regex.FindAllStringSubmatch(params, -1)
	validFlags := map[string]*string{"id": id, "type": type_, "fs": fs_}

	for _, match := range matches {
		flagName, flagValue := strings.ToLower(match[1]), strings.Trim(match[2], "\"")

		if ptr, exists := validFlags[flagName]; exists {
			*ptr = flagValue
		} else {
			fmt.Printf("Error: Flag '%s' no encontrada\n", flagName)
			return
		}
	}

	// Validar que 'id' no esté vacío
	if *id == "" {
		fmt.Println("Error: El parámetro 'id' es obligatorio.")
		return
	}

	FileSystem.Mkfs(strings.ToLower(*id), *type_, *fs_)
}

//mkfs -id=431a
//mkfs -type=full -id=431a

func CommandLogin(params string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "ID Particion")

	matches := regex.FindAllStringSubmatch(params, -1)

	// Mapa para almacenar los valores ingresados por el usuario
	parsedFlags := make(map[string]string)

	for _, match := range matches {
		flagName := match[1]                      // Flag tal cual fue escrito
		flagValue := strings.Trim(match[2], "\"") // Quita comillas si las tiene

		// Asigna el flag en la estructura fs
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: No se pudo establecer el flag '%s'\n", flagName)
			return
		}
		parsedFlags[flagName] = flagValue // Guardar para depuración
	}

	// Imprimir parámetros detectados para depuración
	fmt.Println("====== Parámetros Escaneados ======")
	for key, value := range parsedFlags {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("===================================")

	// Validación de campos obligatorios
	if *user == "" || *pass == "" || *id == "" {
		fmt.Println("Error: Los parámetros '-user', '-pass' y '-id' son obligatorios")
		return
	}
	FileSystem.Login(*user, *pass, *id)
}

//login -user=root -pass=123 -id=431a

func CommandLogOut(params string) {
	// Verificar que no se pasen parámetros
	if strings.TrimSpace(params) != "" {
		fmt.Println("Error: El comando 'mounted' no acepta parámetros")
		return
	}
	// Llamar a la función para imprimir las particiones montadas
	FileSystem.Logout()
}
