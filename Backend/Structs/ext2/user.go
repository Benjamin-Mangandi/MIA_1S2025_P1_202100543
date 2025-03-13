package Ext2

import "fmt"

type User struct {
	ID       string
	Type     string
	Group    string
	Name     string
	Password string
	Status   bool
}

// PrintUser imprime los datos del usuario en un formato legible
func PrintUser(user User) {
	if !user.Status {
		fmt.Println("Usuario eliminado o inactivo")
		return
	}

	fmt.Println("====== Información del Usuario ======")
	fmt.Printf("ID:       %s\n", user.ID)
	fmt.Printf("Tipo:     %s\n", user.Type)
	fmt.Printf("Grupo:    %s\n", user.Group)
	fmt.Printf("Nombre:   %s\n", user.Name)
	fmt.Printf("Password: %s\n", user.Password)
	fmt.Println("=====================================")
}

// deleteUser marca el usuario como eliminado
func deleteUser(user *User) {
	user.ID = "0"
	user.Status = false
}
