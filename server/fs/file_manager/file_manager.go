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
func Mkusr(user string, pass string, grp string) {
	fmt.Println("======Start MKUSR======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Grp:", grp)

	// Validaciones básicas
	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede crear nuevos usuarios.")
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"

	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		return
	}
	defer file.Close()

	// Leer MBR y SuperBloque
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		return
	}

	index := int(global.CurrentUser.ID[1] - '1') // Correlativo correcto
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice de partición fuera de rango")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el SuperBloque")
		return
	}

	// Obtener inodo de /users.txt
	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró /users.txt")
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("ERROR: No se pudo leer el inodo")
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
			return
		}
	}
	if !existsGroup {
		fmt.Println("ERROR: El grupo especificado no existe.")
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
		return
	}

	fmt.Println("Usuario creado exitosamente:", user)
	fmt.Println("======End MKUSR======")
}
func Mkgrp(grp string) {
	fmt.Println("======Start MKGRP======")
	fmt.Println("Grp:", grp)

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede crear grupos.")
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
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
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("Error al leer el superblock:", err)
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo /users.txt")
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("Error al leer el inodo:", err)
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
		return
	}

	fmt.Println("Grupo creado exitosamente:", grp)
	fmt.Println("======End MKGRP======")
}

func Rmgrp(grp string) {
	fmt.Println("======Start RMGRP======")

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede eliminar grupos.")
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
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
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[partIndex].Start)); err != nil {
		fmt.Println("Error al leer el Superblock:", err)
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo users.txt")
		return
	}

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, int64(sb.S_inode_start)+int64(inodeIndex)*int64(binary.Size(inode))); err != nil {
		fmt.Println("Error al leer el inodo:", err)
		return
	}

	data := utils_inodes.GetInodeFileData(inode, file, sb)
	if data == "" {
		fmt.Println("ERROR: Archivo users.txt vacío o corrupto")
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
			fields[0] = "0"
			found = true
		}
		updated = append(updated, strings.Join(fields, ","))
	}

	if !found {
		fmt.Println("ERROR: Grupo no encontrado.")
		return
	}

	newContent := strings.Join(updated, "\n") + "\n"
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("ERROR al actualizar el contenido:", err)
		return
	}

	fmt.Println("Grupo eliminado lógicamente:", grp)
	fmt.Println("======End RMGRP======")
}

func Rmusr(user string) {
	fmt.Println("======Start RMUSR======")

	if !global.CurrentUser.Status || global.CurrentUser.User != "root" {
		fmt.Println("ERROR: Solo el usuario root puede eliminar usuarios.")
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR al abrir el archivo binario:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR al leer el MBR:", err)
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR al leer el superbloque:", err)
		return
	}

	inodeIndex := utils_inodes.InitSearch("/users.txt", file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo /users.txt")
		return
	}

	var inode structs.Inode
	inodeOffset := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("ERROR al leer el inodo:", err)
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
		return
	}

	newContent := strings.Join(updated, "\n") + "\n"
	if err := utils_inodes.UpdateInodeFileData(inodeIndex, newContent, file, sb); err != nil {
		fmt.Println("ERROR al actualizar el archivo users.txt:", err)
		return
	}

	fmt.Println("Usuario eliminado (lógicamente):", user)
	fmt.Println("======End RMUSR======")
}

func Mkfile(path string, size int, r bool, cont string) {
	fmt.Println("======Start MKFILE======")
	fmt.Println("Path:", path)
	fmt.Println("Size:", size)
	fmt.Println("Recursive (r):", r)
	fmt.Println("Contenido:", cont)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión.")
		return
	}

	// Obtener la ruta del disco
	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		return
	}
	defer file.Close()

	// Leer MBR y Superblock
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superblock")
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
		return
	}

	fmt.Println("Archivo creado correctamente.")
	fmt.Println("======End MKFILE======")
}

func Mkdir(path string, p bool) {
	fmt.Println("======Start MKDIR======")
	fmt.Println("Path:", path)
	fmt.Println("Recursive (r):", p)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión.")
		return
	}

	driveletter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superblock")
		return
	}

	err = utils_inodes.CreateFolderRecursive(path, file, sb, p)
	if err != nil {
		fmt.Println("ERROR al crear carpeta:", err)
		return
	}

	fmt.Println("Carpeta creada correctamente.")
	fmt.Println("======End MKDIR======")
}

func Cat(path string) {
	fmt.Println("======Start CAT======")
	fmt.Println("Path:", path)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión primero.")
		return
	}

	// Obtener ruta del disco .bin
	driveLetter := string(global.CurrentUser.ID[0])
	filePath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filePath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo .bin")
		return
	}
	defer file.Close()

	// Leer el MBR y Superbloque
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR")
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	if index < 0 || index >= len(mbr.Partitions) {
		fmt.Println("ERROR: ID fuera de rango")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque")
		return
	}

	// Buscar el inodo del archivo
	inodeIndex := utils_inodes.InitSearch(path, file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo.")
		return
	}

	// Leer el inodo
	var inode structs.Inode
	offset := int64(sb.S_inode_start + inodeIndex*int32(binary.Size(structs.Inode{})))
	if err := utils.ReadObject(file, &inode, offset); err != nil {
		fmt.Println("ERROR: No se pudo leer el inodo")
		return
	}

	// Obtener y mostrar contenido
	content := utils_inodes.GetInodeFileData(inode, file, sb)
	fmt.Println("Contenido del archivo:\n" + content)

	fmt.Println("======End CAT======")
}

func Find(startPath string, name string) {
	fmt.Println("======Start FIND======")
	fmt.Println("Start path:", startPath)
	fmt.Println("Name to find:", name)

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		return
	}

	driveLetter := string(global.CurrentUser.ID[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo del disco.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(global.CurrentUser.ID[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("ERROR: Índice inválido.")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	startInodeIndex := utils_inodes.InitSearch(startPath, file, sb)
	if startInodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el path inicial.")
		return
	}

	matchAll := (name == "*")
	FindRecursive(file, sb, startInodeIndex, startPath, name, matchAll)

	fmt.Println("======End FIND======")
}

func FindRecursive(file *os.File, sb structs.Superblock, inodeIndex int32, currentPath string, target string, matchAll bool) {
	var inode structs.Inode
	inodeOffset := sb.S_inode_start + inodeIndex*int32(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, int64(inodeOffset)); err != nil {
		fmt.Println("ERROR al leer inodo:", err)
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
			}

			// Recursivamente explorar subcarpetas
			childInodeOffset := sb.S_inode_start + entry.B_inodo*int32(binary.Size(inode))
			var childInode structs.Inode
			if err := utils.ReadObject(file, &childInode, int64(childInodeOffset)); err != nil {
				continue
			}

			if string(childInode.I_type[:]) == "0" { // Es carpeta
				FindRecursive(file, sb, entry.B_inodo, childPath, target, matchAll)
			}
		}
	}
}

func Pause() {
	fmt.Println("====== PAUSE ======")
	fmt.Print("Presione ENTER para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
