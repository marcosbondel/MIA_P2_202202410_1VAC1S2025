package file_manager

import (
	"MIA_P2_202202410_1VAC1S2025/fs/global"
	"MIA_P2_202202410_1VAC1S2025/fs/structs"
	"MIA_P2_202202410_1VAC1S2025/fs/utils"
	"MIA_P2_202202410_1VAC1S2025/fs/utils_inodes"
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// login -user=root -pass=123 -id=A119
// mkusr -user=user1 -pass=CurrentUser -grp=CurrentUsers
func Mkusr(user string, pass string, grp string, buffer_string *string) {
	fmt.Println("======Start MKUSR======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Grp:", grp)

	*buffer_string += "======Start MKUSR======\n"
	*buffer_string += fmt.Sprintf("User: %s\n", user)
	*buffer_string += fmt.Sprintf("Pass: %s\n", pass)
	*buffer_string += fmt.Sprintf("Grp: %s\n", grp)

	// Validaciones básicas
	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede crear nuevos usuarios.")
		*buffer_string += "ERROR: Solo el usuario root puede crear nuevos usuarios.\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"

	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		*buffer_string += "ERROR: No se pudo abrir el archivo .bin\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}
	defer file.Close()

	// Leer MBR y SuperBloque
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		*buffer_string += "ERROR: No se pudo leer el MBR\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1') // Correlativo correcto
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice de partición fuera de rango")
		*buffer_string += "ERROR: Índice de partición fuera de rango\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el SuperBloque")
		*buffer_string += "ERROR: No se pudo leer el SuperBloque\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	// Obtener inodo de /users.txt
	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró /users.txt")
		*buffer_string += "ERROR: No se encontró /users.txt\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("ERROR: No se pudo leer el inodo")
		*buffer_string += "ERROR: No se pudo leer el inodo\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	// Leer el contenido actual de /users.txt
	data := utils_inodes.GetInodeFileData(inode, file, sb)
	lines := strings.Split(data, "\n")

	// Verificar si grupo existe y si el usuario ya existe
	existsGroup := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) >= 3 && parts[1] == "G" && parts[0] != "0" && parts[2] == grp {
			existsGroup = true
		}
		if len(parts) >= 4 && parts[1] == "U" && parts[0] != "0" && parts[3] == user {
			fmt.Println("ERROR: El usuario ya existe.")
			*buffer_string += "ERROR: El usuario ya existe.\n"
			*buffer_string += "======End MKUSR======\n"
			return
		}
	}
	if !existsGroup {
		fmt.Println("ERROR: El grupo especificado no existe.")
		*buffer_string += "ERROR: El grupo especificado no existe.\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	// Generar nuevo UID
	uid := 1
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) >= 2 && parts[1] == "U" && parts[0] != "0" {
			uid++
		}
	}

	newLine := fmt.Sprintf("%d,U,%s,%s,%s\n", uid, grp, user, pass)
	newContent := data + newLine

	// Guardar cambios en /users.txt
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("ERROR: No se pudo actualizar el archivo /users.txt")
		*buffer_string += "ERROR: No se pudo actualizar el archivo /users.txt\n"
		*buffer_string += "======End MKUSR======\n"
		return
	}

	*buffer_string += fmt.Sprintf("Usuario creado exitosamente: %s\n", user)
	*buffer_string += "======End MKUSR======\n"

	fmt.Println("Usuario creado exitosamente:", user)
	fmt.Println("======End MKUSR======")
}
func Mkgrp(grp string, buffer_string *string) {
	fmt.Println("======Start MKGRP======")
	fmt.Println("Grp:", grp)

	*buffer_string += "======Start MKGRP======\n"
	*buffer_string += fmt.Sprintf("Grp: %s\n", grp)

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede crear grupos.")
		*buffer_string += "ERROR: Solo el usuario root puede crear grupos.\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		*buffer_string += "Error al abrir el archivo: " + err.Error() + "\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		*buffer_string += "Error al leer el MBR: " + err.Error() + "\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	// Buscar la partición montada usando ID exacto
	index := -1
	for i := 0; i < 4; i++ {
		if strings.TrimSpace(string(mbr.Partitions[i].Id[:])) == global.CurrentUser.ID {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("ERROR: No se encontró la partición montada.")
		*buffer_string += "ERROR: No se encontró la partición montada.\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("Error al leer el superblock:", err)
		*buffer_string += "Error al leer el superblock: " + err.Error() + "\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo /users.txt")
		*buffer_string += "ERROR: No se encontró el archivo /users.txt\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("Error al leer el inodo:", err)
		*buffer_string += "Error al leer el inodo: " + err.Error() + "\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	data := utils_inodes.GetInodeFileData(inode, file, sb)
	lines := strings.Split(data, "\n")

	// Verificar si ya existe un grupo activo con ese nombre
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) >= 3 && fields[0] != "0" && fields[1] == "G" && fields[2] == grp {
			fmt.Println("ERROR: El grupo ya existe.")
			*buffer_string += "ERROR: El grupo ya existe.\n"
			*buffer_string += "======End MKGRP======\n"
			return
		}
	}

	// Calcular nuevo ID (el mayor ID existente + 1)
	newID := 1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) >= 3 && fields[1] == "G" {
			if idParsed := strings.TrimSpace(fields[0]); idParsed != "0" {
				newID++
			}
		}
	}

	newLine := fmt.Sprintf("%d,G,%s\n", newID, grp)
	newContent := data + newLine
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("Error al escribir el archivo:", err)
		*buffer_string += "Error al escribir el archivo: " + err.Error() + "\n"
		*buffer_string += "======End MKGRP======\n"
		return
	}

	*buffer_string += fmt.Sprintf("Grupo creado exitosamente: %s\n", grp)
	*buffer_string += "======End MKGRP======\n"
	fmt.Println("Grupo creado exitosamente:", grp)
	fmt.Println("======End MKGRP======")
}

func Rmgrp(grp string, buffer_string *string) {
	fmt.Println("======Start RMGRP======")
	*buffer_string += "======Start RMGRP======\n"
	*buffer_string += fmt.Sprintf("Grp: %s\n", grp)

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede eliminar grupos.")
		*buffer_string += "ERROR: Solo el usuario root puede eliminar grupos.\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		*buffer_string += "Error al abrir el disco: " + err.Error() + "\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		*buffer_string += "Error al leer el MBR: " + err.Error() + "\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	// Buscar la partición por ID exacto
	partIndex := -1
	for i := 0; i < 4; i++ {
		if strings.TrimSpace(string(mbr.Partitions[i].Id[:])) == global.CurrentUser.ID {
			partIndex = i
			break
		}
	}
	if partIndex == -1 {
		fmt.Println("ERROR: No se encontró la partición montada.")
		*buffer_string += "ERROR: No se encontró la partición montada.\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[partIndex].Start)); err != nil {
		fmt.Println("Error al leer el Superblock:", err)
		*buffer_string += "Error al leer el Superblock: " + err.Error() + "\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo users.txt")
		*buffer_string += "ERROR: No se encontró el archivo users.txt\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("Error al leer el inodo:", err)
		*buffer_string += "Error al leer el inodo: " + err.Error() + "\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	data := utils_inodes.GetInodeFileData(inode, file, sb)
	if data == "" {
		fmt.Println("ERROR: Archivo users.txt vacío o corrupto")
		*buffer_string += "ERROR: Archivo users.txt vacío o corrupto\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	lines := strings.Split(data, "\n")
	var updated []string
	found := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")

		// Ignora registros lógicamente eliminados
		if fields[0] == "0" {
			updated = append(updated, line)
			continue
		}

		if len(fields) >= 3 && fields[1] == "G" && fields[2] == grp {
			fmt.Println("Grupo encontrado, marcando como eliminado:", grp)
			*buffer_string += fmt.Sprintf("Grupo encontrado, marcando como eliminado: %s\n", grp)
			fields[0] = "0"
			found = true
		}
		updated = append(updated, strings.Join(fields, ","))
	}

	if !found {
		fmt.Println("ERROR: Grupo no encontrado.")
		*buffer_string += "ERROR: Grupo no encontrado.\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}

	newContent := strings.Join(updated, "\n") + "\n"
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("ERROR al actualizar el contenido:", err)
		*buffer_string += "ERROR al actualizar el contenido: " + err.Error() + "\n"
		*buffer_string += "======End RMGRP======\n"
		return
	}
	*buffer_string += fmt.Sprintf("Grupo eliminado lógicamente: %s\n", grp)
	*buffer_string += "======End RMGRP======\n"
	fmt.Println("Grupo eliminado lógicamente:", grp)
	fmt.Println("======End RMGRP======")
}

func Rmusr(user string, buffer_string *string) {
	fmt.Println("======Start RMUSR======")
	*buffer_string += "======Start RMUSR======\n"

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede eliminar usuarios.")
		*buffer_string += "ERROR: Solo el usuario root puede eliminar usuarios.\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR al abrir el archivo binario:", err)
		*buffer_string += "ERROR al abrir el archivo binario: " + err.Error() + "\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR al leer el MBR:", err)
		*buffer_string += "ERROR al leer el MBR: " + err.Error() + "\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR al leer el superbloque:", err)
		*buffer_string += "ERROR al leer el superbloque: " + err.Error() + "\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo /users.txt")
		*buffer_string += "ERROR: No se encontró el archivo /users.txt\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	var inode structs.Inode
	inodeOffset := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("ERROR al leer el inodo:", err)
		*buffer_string += "ERROR al leer el inodo: " + err.Error() + "\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	data := utils_inodes.GetInodeFileData(inode, file, sb)
	lines := strings.Split(data, "\n")

	var updated []string
	found := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) >= 4 && fields[1] == "U" && strings.TrimSpace(fields[3]) == user {
			fields[0] = "0"
			found = true
		}
		updated = append(updated, strings.Join(fields, ","))
	}

	if !found {
		fmt.Println("ERROR: Usuario no encontrado.")
		*buffer_string += "ERROR: Usuario no encontrado.\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	newContent := strings.Join(updated, "\n") + "\n"
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("ERROR al actualizar el archivo users.txt:", err)
		*buffer_string += "ERROR al actualizar el archivo users.txt: " + err.Error() + "\n"
		*buffer_string += "======End RMUSR======\n"
		return
	}

	*buffer_string += fmt.Sprintf("Usuario eliminado lógicamente: %s\n", user)
	*buffer_string += "======End RMUSR======\n"
	fmt.Println("Usuario eliminado (lógicamente):", user)
	fmt.Println("======End RMUSR======")
}

func Mkfile(path string, size int, r bool, cont string, buffer_string *string) {
	fmt.Println("======Start MKFILE======")
	fmt.Println("Path:", path)
	fmt.Println("Size:", size)
	fmt.Println("Recursive (r):", r)
	fmt.Println("Contenido:", cont)

	*buffer_string += "======Start MKFILE======\n"
	*buffer_string += fmt.Sprintf("Path: %s\n", path)
	*buffer_string += fmt.Sprintf("Size: %d\n", size)
	*buffer_string += fmt.Sprintf("Recursive (r): %t\n", r)
	*buffer_string += fmt.Sprintf("Contenido: %s\n", cont)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión.")
		*buffer_string += "ERROR: Debes iniciar sesión.\n"
		*buffer_string += "======End MKFILE======\n"
		return
	}

	// Obtener la ruta del disco
	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		*buffer_string += "ERROR: No se pudo abrir el archivo .bin\n"
		*buffer_string += "======End MKFILE======\n"
		return
	}
	defer file.Close()

	// Leer MBR y Superblock
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		*buffer_string += "ERROR: No se pudo leer el MBR\n"
		*buffer_string += "======End MKFILE======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superblock")
		*buffer_string += "ERROR: No se pudo leer el Superblock\n"
		*buffer_string += "======End MKFILE======\n"
		return
	}

	// Determinar el contenido
	var content []byte
	if cont != "" {
		content = []byte(cont)
	} else {
		content = make([]byte, size)
		for i := range content {
			content[i] = byte('0')
		}
	}

	// Crear archivo
	err = utils_inodes.CreateFileWithPath(path, content, file, sb, r)
	if err != nil {
		fmt.Println("ERROR:", err)
		*buffer_string += "ERROR: " + err.Error() + "\n"
		*buffer_string += "======End MKFILE======\n"
		return
	}

	*buffer_string += "Archivo creado correctamente.\n"
	*buffer_string += "======End MKFILE======\n"
	fmt.Println("Archivo creado correctamente.")
	fmt.Println("======End MKFILE======")
}

func Mkdir(path string, p bool, buffer_string *string) {
	fmt.Println("======Start MKDIR======")
	fmt.Println("Path:", path)
	fmt.Println("Recursive (r):", p)

	*buffer_string += "======Start MKDIR======\n"
	*buffer_string += fmt.Sprintf("Path: %s\n", path)
	*buffer_string += fmt.Sprintf("Recursive (r): %t\n", p)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión.")
		*buffer_string += "ERROR: Debes iniciar sesión.\n"
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		*buffer_string += "ERROR: No se pudo abrir el archivo .bin\n"
		*buffer_string += "======End MKDIR======\n"
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		*buffer_string += "ERROR: No se pudo leer el MBR\n"
		*buffer_string += "======End MKDIR======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superblock")
		*buffer_string += "ERROR: No se pudo leer el Superblock\n"
		*buffer_string += "======End MKDIR======\n"
		return
	}

	err = utils_inodes.CreateFolderRecursive(path, file, sb, p)
	if err != nil {
		fmt.Println("ERROR al crear carpeta:", err)
		*buffer_string += "ERROR al crear carpeta: " + err.Error() + "\n"
		*buffer_string += "======End MKDIR======\n"
		return
	}

	*buffer_string += "Carpeta creada correctamente.\n"
	*buffer_string += "======End MKDIR======\n"

	fmt.Println("Carpeta creada correctamente.")
	fmt.Println("======End MKDIR======")
}

func Cat(path string, buffer_string *string) {
	fmt.Println("======Start CAT======")
	fmt.Println("Path:", path)

	*buffer_string += "======Start CAT======\n"
	*buffer_string += fmt.Sprintf("Path: %s\n", path)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión primero.")
		*buffer_string += "ERROR: Debes iniciar sesión primero.\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	// Obtener ruta del disco .bin
	driveLetter := string(global.CurrentUser.ID[0])
	filePath := "./fs/test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filePath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		*buffer_string += "ERROR: No se pudo abrir el archivo .bin\n"
		*buffer_string += "======End CAT======\n"
		return
	}
	defer file.Close()

	// Leer el MBR y Superbloque
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		*buffer_string += "ERROR: No se pudo leer el MBR\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	if index < 0 || index >= len(mbr.Partitions) {
		fmt.Println("ERROR: ID fuera de rango")
		*buffer_string += "ERROR: ID fuera de rango\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque")
		*buffer_string += "ERROR: No se pudo leer el Superbloque\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	// Buscar el inodo del archivo
	inodeIndex := utils_inodes.InitSearch(path, file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo.")
		*buffer_string += "ERROR: No se encontró el archivo.\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	// Leer el inodo
	var inode structs.Inode
	offset := int64(sb.S_inode_start + inodeIndex*int32(binary.Size(structs.Inode{})))
	if err := utils.ReadObject(file, &inode, offset); err != nil {
		fmt.Println("ERROR: No se pudo leer el inodo")
		*buffer_string += "ERROR: No se pudo leer el inodo\n"
		*buffer_string += "======End CAT======\n"
		return
	}

	// Obtener y mostrar contenido
	content := utils_inodes.GetInodeFileData(inode, file, sb)
	fmt.Println("Contenido del archivo:\n" + content)
	*buffer_string += "Contenido del archivo:\n" + content + "\n"
	*buffer_string += "======End CAT======\n"

	fmt.Println("======End CAT======")
}

func Find(startPath string, name string, buffer_string *string) {
	fmt.Println("======Start FIND======")
	fmt.Println("Start path:", startPath)
	fmt.Println("Name to find:", name)

	*buffer_string += "======Start FIND======\n"
	*buffer_string += fmt.Sprintf("Start path: %s\n", startPath)
	*buffer_string += fmt.Sprintf("Name to find: %s\n", name)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		*buffer_string += "ERROR: No hay sesión activa.\n"
		*buffer_string += "======End FIND======\n"
		return
	}

	driveLetter := string(global.CurrentUser.ID[0])
	filepath := "./fs/test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo del disco.")
		*buffer_string += "ERROR: No se pudo abrir el archivo del disco.\n"
		*buffer_string += "======End FIND======\n"
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		*buffer_string += "ERROR: No se pudo leer el MBR.\n"
		*buffer_string += "======End FIND======\n"
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("ERROR: Índice inválido.")
		*buffer_string += "ERROR: Índice inválido.\n"
		*buffer_string += "======End FIND======\n"
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		*buffer_string += "ERROR: No se pudo leer el Superbloque.\n"
		*buffer_string += "======End FIND======\n"
		return
	}

	startInodeIndex := utils_inodes.InitSearch(startPath, file, sb)
	if startInodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el path inicial.")
		*buffer_string += "ERROR: No se encontró el path inicial.\n"
		*buffer_string += "======End FIND======\n"
		return
	}

	matchAll := (name == "*")
	FindRecursive(file, sb, startInodeIndex, startPath, name, matchAll, buffer_string)

	*buffer_string += "======End FIND======\n"
	fmt.Println("======End FIND======")
}

func FindRecursive(file *os.File, sb structs.Superblock, inodeIndex int32, currentPath string, target string, matchAll bool, buffer_string *string) {
	var inode structs.Inode
	inodeOffset := sb.S_inode_start + inodeIndex*int32(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, int64(inodeOffset)); err != nil {
		fmt.Println("ERROR al leer inodo:", err)
		*buffer_string += "ERROR al leer inodo: " + err.Error() + "\n"
		return
	}

	for _, ptr := range inode.I_block {
		if ptr == -1 {
			continue
		}

		blockOffset := sb.S_block_start + ptr*int32(binary.Size(structs.Folderblock{}))
		var folder structs.Folderblock
		if err := utils.ReadObject(file, &folder, int64(blockOffset)); err != nil {
			continue
		}

		for _, entry := range folder.B_content {
			name := strings.Trim(string(entry.B_name[:]), "\x00")
			if name == "" || name == "." || name == ".." {
				continue
			}

			childPath := currentPath + "/" + name

			if matchAll || name == target {
				fmt.Println("FOUND:", childPath)
				*buffer_string += fmt.Sprintf("FOUND: %s\n", childPath)
			}

			// Recursivamente explorar subcarpetas
			childInodeOffset := sb.S_inode_start + entry.B_inodo*int32(binary.Size(inode))
			var childInode structs.Inode
			if err := utils.ReadObject(file, &childInode, int64(childInodeOffset)); err != nil {
				continue
			}

			if string(childInode.I_type[:]) == "0" { // Es carpeta
				FindRecursive(file, sb, entry.B_inodo, childPath, target, matchAll, buffer_string)
			}
		}
	}
}

func Pause() {
	fmt.Println("====== PAUSE ======")
	fmt.Print("Presione ENTER para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
