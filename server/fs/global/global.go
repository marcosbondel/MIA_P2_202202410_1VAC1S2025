package global

type UserInfo struct {
	ID     string // ID de partición montada, como "A1"
	User   string // Nombre de usuario actual (por ejemplo, "root", "admin", etc.)
	Group  string // Grupo al que pertenece el usuario
	Status bool   // true si hay sesión activa
}

var CurrentUser UserInfo
