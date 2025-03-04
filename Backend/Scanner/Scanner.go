package Scanner

import (
	"Backend/DiskManager"
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
		//mkdisk -size=1024 -unit=M -fit=BF -path="/home/benjamin/discos/disk1.mia"
	}
}

func AnalyzeCommand(command string, params string) {
	switch {
	case strings.EqualFold(command, "mkdisk"):
		Commandmkdisk(params)
	case strings.EqualFold(command, "fdisk"):
		Commandfdisk(params)
	case strings.EqualFold(command, "mount"):
		Commandmount(params)
	default:
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}

// Procesa el comando mkdisk
func Commandmkdisk(params string) {
	// Definir flags con valores predeterminados
	fs := flag.NewFlagSet("mkdisk", flag.ContinueOnError)
	size := fs.Int("size", 0, "Tamaño del disco")
	fit := fs.String("fit", "ff", "Tipo de ajuste (bf, ff, wf)")
	unit := fs.String("unit", "m", "Unidad de tamaño (k, m)")
	path := fs.String("path", "", "Ruta donde se creará el disco")

	_ = fs.Parse([]string{})

	// Extraer parámetros de la entrada del usuario
	matches := regex.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]                     // Nombre del flag
		flagValue := strings.Trim(match[2], `"`) // Valor del flag sin comillas

		// Normalizar los valores en minúsculas cuando corresponda
		if flagName != "path" {
			flagValue = strings.ToLower(flagValue)
		}

		// Intentar asignar el valor al flag correspondiente
		if err := fs.Set(flagName, flagValue); err != nil {
			fmt.Printf("Error: Flag desconocida '%s'\n", flagName)
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor a 0")
		return
	}

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf', 'ff' o 'wf'")
		return
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: La unidad debe ser 'k' o 'm'")
		return
	}

	if *path == "" {
		fmt.Println("Error: Se requiere una ruta válida")
		return
	}

	// Llamar a la función para crear el disco
	DiskManager.Mkdisk(*size, *fit, *unit, *path)
}

func Commandfdisk(input string) {
	// Definir flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	//name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "", "Ajuste") // Dejar fit vacío por defecto

	// Parsear los flags
	fs.Parse(os.Args[1:])

	// Encontrar los flags en el input
	matches := regex.FindAllStringSubmatch(input, -1)

	// Procesar el input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	if *path == "" {
		fmt.Println("Error: Path is required")
		return
	}

	// Si no se proporcionó un fit, usar el valor predeterminado "w"
	if *fit == "" {
		*fit = "w"
	}

	// Validar fit (b/w/f)
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: Fit must be 'b', 'f', or 'w'")
		return
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be 'k' or 'm'")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: Type must be 'p', 'e', or 'l'")
		return
	}

	// Llamar a la función
	//DiskManagement.Fdisk(*size, *path, *name, *unit, *type_, *fit)
}

func Commandmount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := regex.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2]) // Convertir todo a minúsculas
		flagValue = strings.Trim(flagValue, "\"")
		fs.Set(flagName, flagValue)
	}

	if *path == "" || *name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}

	// Convertir el nombre a minúsculas antes de pasarlo al Mount
	//lowercaseName := strings.ToLower(*name)
	//DiskManagement.Mount(*path, lowercaseName)
}
