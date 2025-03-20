package UsersManager

import (
	"Backend/Globals"
	Ext2 "Backend/Structs/ext2"
	"fmt"
)

func Logout() {
	if !Globals.ActiveUser.Status {
		fmt.Println("Error: No hay ningún usuario activo")
		return
	}

	// Resetear la sesión
	Globals.ActiveUser = Ext2.User{}

	fmt.Println("Sesión cerrada exitosamente.")
}
