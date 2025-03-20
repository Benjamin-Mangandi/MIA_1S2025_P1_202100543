package Responsehandler // Renombrado para mayor claridad

// Estructura que representa la respuesta de la aplicación
type AppResponse struct {
	Content string
}

// Función para concatenar contenido en la respuesta
func AppendContent(prevResponse *AppResponse, newContent string) {
	prevResponse.Content += newContent
}
func Clear() {
	GlobalResponse.Content = ""
}

// Variable global para almacenar la respuesta acumulada
var GlobalResponse AppResponse
