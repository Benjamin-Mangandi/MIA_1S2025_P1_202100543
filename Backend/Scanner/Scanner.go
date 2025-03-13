package Scanner

import (
	"Backend/Scanner/CommandsDisk"
	"Backend/Scanner/CommandsFileSystem"
	"Backend/Scanner/CommandsUsers"
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		CommandsDisk.MkDisk(params)
	case strings.EqualFold(command, "fdisk"):
		CommandsDisk.FDisk(params)
	case strings.EqualFold(command, "rmdisk"):
		CommandsDisk.RmDisk(params)
	case strings.EqualFold(command, "mounted"):
		CommandsDisk.Mounted(params)
	case strings.EqualFold(command, "mount"):
		CommandsDisk.Mount(params)
	case strings.EqualFold(command, "mkfs"):
		CommandsFileSystem.Mkfs(params)
	case strings.EqualFold(command, "login"):
		CommandsUsers.Login(params)
	case strings.EqualFold(command, "Logout"):
		CommandsUsers.LogOut(params)
	default:
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}
