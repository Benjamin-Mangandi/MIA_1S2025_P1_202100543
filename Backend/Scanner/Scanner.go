package Scanner

import (
	"Backend/Scanner/CommandsDisk"
	"Backend/Scanner/CommandsFileSystem"
	"Backend/Scanner/CommandsReports"
	"Backend/Scanner/CommandsUsers"
	"fmt"
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
func Scan(input string) {
	lines := strings.Split(input, "\n") // Separar por saltos de línea
	for _, line := range lines {
		command, params := getCommandAndParams(line)
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
	case strings.EqualFold(command, "cat"):
		CommandsFileSystem.Cat(params)
	case strings.EqualFold(command, "mkgrp"):
		CommandsUsers.Mkgrp(params)
	case strings.EqualFold(command, "rmgrp"):
		CommandsUsers.Rmgrp(params)
	case strings.EqualFold(command, "mkusr"):
		CommandsUsers.Mkusr(params)
	case strings.EqualFold(command, "rmusr"):
		CommandsUsers.Rmusr(params)
	case strings.EqualFold(command, "rep"):
		CommandsReports.Report(params)
	default:
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}
