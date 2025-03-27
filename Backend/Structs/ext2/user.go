package Ext2

import (
	"Backend/Responsehandler"
	"fmt"
	"strings"
)

type User struct {
	ID          string
	Type        string
	Group       string
	Name        string
	Password    string
	Status      bool
	PartitionID string
}

// PrintUser imprime los datos del usuario en un formato legible y los agrega al Responsehandler
func PrintUser(user User) {
	if !user.Status {
		response := "Usuario eliminado o inactivo"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		fmt.Println(response)
		return
	}

	response := strings.Repeat("-", 40) + "\n" +
		"Inicio de sesión exitoso.\n" +
		strings.Repeat("-", 40) + "\n" +
		fmt.Sprintf("ID:       %s\n", user.ID) +
		fmt.Sprintf("Tipo:     %s\n", user.Type) +
		fmt.Sprintf("Grupo:    %s\n", user.Group) +
		fmt.Sprintf("Nombre:   %s\n", user.Name) +
		fmt.Sprintf("Password: %s\n", user.Password)

	// Agregar la información al response
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
