package UsersManager

import (
	"Backend/Globals"
	"Backend/Responsehandler"
	Ext2 "Backend/Structs/ext2"
	"strings"
)

func Logout() {
	if !Globals.ActiveUser.Status {
		response := strings.Repeat("*", 30) + "\n" +
			"Error: No hay ningún usuario activo"
		Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
		return
	}

	// Resetear la sesión
	Globals.ActiveUser = Ext2.User{}
	response := strings.Repeat("-", 40) + "\n" +
		"Sesión cerrada exitosamente." + "\n"
	Responsehandler.AppendContent(&Responsehandler.GlobalResponse, response)
}
