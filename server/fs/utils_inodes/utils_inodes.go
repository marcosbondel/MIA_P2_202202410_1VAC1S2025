package utils_inodes

import (
	"MIA_P2_202202410_1VAC1S2025/fs/structs"
	"MIA_P2_202202410_1VAC1S2025/fs/utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// login -user=root -pass=123 -id=A119
func InitSearch(path string, file *os.File, tempSuperblock structs.Superblock) int32 {
	fmt.Println("======Start INITSEARCH======")
	fmt.Println("path:", path)

	// Filtrar componentes vacíos y espacios
	var StepsPath []string
	for _, part := range strings.Split(path, "/") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			StepsPath = append(StepsPath, trimmed)
		}
	}

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 structs.Inode
	if err := utils.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		fmt.Println("ERROR: Cannot read root inode:", err)
		return -1
	}

	fmt.Println("======End INITSEARCH======")
	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(steps []string, inode structs.Inode, file *os.File, sb structs.Superblock) int32 {
	currentInode := inode
	currentInodeIndex := int32(0)

	for _, step := range steps {
		found := false
		var nextInode structs.Inode
		var nextInodeIndex int32 = -1

		for _, block := range currentInode.I_block {
			if block == -1 {
				continue
			}

			var folder structs.Folderblock
			offset := sb.S_block_start + block*int32(binary.Size(folder))
			if err := utils.ReadObject(file, &folder, int64(offset)); err != nil {
				continue
			}

			for _, content := range folder.B_content {
				name := strings.Trim(string(content.B_name[:]), "\x00 ")

				if name == step && content.B_inodo != -1 {
					nextInodeIndex = content.B_inodo
					inodeOffset := sb.S_inode_start + nextInodeIndex*int32(binary.Size(currentInode))
					if err := utils.ReadObject(file, &nextInode, int64(inodeOffset)); err != nil {
						return -1
					}

					currentInode = nextInode
					currentInodeIndex = nextInodeIndex
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if !found {
			fmt.Println("❌ Carpeta no encontrada:", step)
			return -1
		}
	}

	return currentInodeIndex
}

func GetInodeIndex(inode structs.Inode, file *os.File, sb structs.Superblock) int32 {
	var temp structs.Inode
	for i := int32(0); i < sb.S_inodes_count; i++ {
		offset := sb.S_inode_start + i*int32(binary.Size(temp))
		if err := utils.ReadObject(file, &temp, int64(offset)); err != nil {
			continue
		}
		if temp == inode {
			return i
		}
	}
	return -1
}

// func GetInodeFileData(inode structs.Inode, file *os.File, sb structs.Superblock) string {
// 	fmt.Println("======Start GETINODEFILEDATA======")
// 	var content string
// 	blockSize := int32(binary.Size(structs.Fileblock{}))

// 	for i := 0; i < 12; i++ { // Directos
// 		block := inode.I_block[i]
// 		if block == -1 {
// 			break
// 		}

// 		var fileBlock structs.Fileblock
// 		offset := sb.S_block_start + int32(block)*blockSize
// 		if err := binaryRead(file, offset, &fileBlock); err != nil {
// 			fmt.Println("Error leyendo bloque:", err)
// 			break
// 		}

// 		content += string(fileBlock.B_content[:])
// 	}
// 	fmt.Println("======End GETINODEFILEDATA======")
// 	return content
// }

func GetInodeFileData(inode structs.Inode, file *os.File, sb structs.Superblock) string {
	var content []byte
	bytesRemaining := int(inode.I_size)

	for _, ptr := range inode.I_block {
		if ptr == -1 || bytesRemaining <= 0 {
			break
		}

		offset := int64(sb.S_block_start) + int64(ptr)*int64(binary.Size(structs.Fileblock{}))
		var block structs.Fileblock
		if err := utils.ReadObject(file, &block, offset); err != nil {
			break
		}

		for _, b := range block.B_content {
			if bytesRemaining <= 0 {
				break
			}
			content = append(content, b)
			bytesRemaining--
		}
	}

	return string(content)
}

func binaryRead(file *os.File, pos int32, data interface{}) error {
	_, err := file.Seek(int64(pos), 0)
	if err != nil {
		return err
	}
	return binary.Read(file, binary.LittleEndian, data)
}

func UpdateInodeFileData(inodeIndex int32, content string, file *os.File, sb structs.Superblock) error {
	fmt.Println("======Start UPDATEINODEFILEDATA======")

	blockSize := int(binary.Size(structs.Fileblock{}))
	neededBlocks := (len(content) + blockSize - 1) / blockSize

	// === 1. Leer inodo actual ===
	var inode structs.Inode
	inodePos := sb.S_inode_start + inodeIndex*int32(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, int64(inodePos)); err != nil {
		return fmt.Errorf("error leyendo inodo: %v", err)
	}

	// === 2. Liberar bloques anteriores ===
	for _, oldPtr := range inode.I_block {
		if oldPtr == -1 {
			continue
		}
		bitmapPos := sb.S_bm_block_start + oldPtr
		_, _ = file.WriteAt([]byte{0}, int64(bitmapPos)) // liberar bitmap
		// opcional: limpiar bloque físicamente
		var empty structs.Fileblock
		blockOffset := sb.S_block_start + oldPtr*int32(blockSize)
		utils.WriteObject(file, empty, int64(blockOffset))
	}

	// === 3. Buscar nuevos bloques libres ===
	usedBlocks := 0
	var assignedBlocks []int32
	for i := 0; i < int(sb.S_blocks_count) && usedBlocks < neededBlocks; i++ {
		bitmapPos := sb.S_bm_block_start + int32(i)
		var value byte
		if err := utils.ReadObject(file, &value, int64(bitmapPos)); err != nil {
			continue
		}
		if value == 0 {
			// marcar como usado
			file.WriteAt([]byte{1}, int64(bitmapPos))
			assignedBlocks = append(assignedBlocks, int32(i))
			usedBlocks++
		}
	}

	if usedBlocks < neededBlocks {
		return fmt.Errorf("no hay bloques libres suficientes")
	}

	// === 4. Escribir nuevo contenido ===
	for i, block := range assignedBlocks {
		start := i * blockSize
		end := start + blockSize
		if end > len(content) {
			end = len(content)
		}
		var fb structs.Fileblock
		copy(fb.B_content[:], []byte(content[start:end]))
		blockOffset := sb.S_block_start + block*int32(blockSize)
		utils.WriteObject(file, fb, int64(blockOffset))
	}

	// === 5. Actualizar inodo ===
	for i := 0; i < 15; i++ {
		if i < len(assignedBlocks) {
			inode.I_block[i] = assignedBlocks[i]
		} else {
			inode.I_block[i] = -1
		}
	}
	inode.I_size = int32(len(content))

	if err := utils.WriteObject(file, inode, int64(inodePos)); err != nil {
		return fmt.Errorf("error escribiendo inodo: %v", err)
	}

	fmt.Println("======End UPDATEINODEFILEDATA======")
	return nil
}

func WriteObject(file *os.File, data interface{}, pos int64) error {
	_, err := file.Seek(pos, 0)
	if err != nil {
		return err
	}
	return binary.Write(file, binary.LittleEndian, data)
}

func CreateFileWithPath(path string, content []byte, file *os.File, sb structs.Superblock, r bool) error {
	fmt.Println("== CreateFileWithPath ==")

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path inválido: debe comenzar con /")
	}

	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")

	// Extraer parentPath y fileName
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 || lastSlash == len(path)-1 {
		return fmt.Errorf("ruta inválida: debe contener al menos una carpeta y nombre de archivo")
	}

	parentPath := path[:lastSlash]
	fileName := path[lastSlash+1:]

	fmt.Println("parentPath:", parentPath)

	// Buscar carpeta contenedora
	parentInodeIndex := InitSearch(parentPath, file, sb)
	fmt.Println("parentInodeIndex:", parentInodeIndex)

	if parentInodeIndex == -1 {
		if r {
			fmt.Println("La carpeta padre no existe, pero -r está activado. Intentando crear directorios...")
			if err := CreateFolderRecursive(parentPath, file, sb, true); err != nil {
				return fmt.Errorf("no se pudieron crear carpetas intermedias: %v", err)
			}
			parentInodeIndex = InitSearch(parentPath, file, sb)
			if parentInodeIndex == -1 {
				return fmt.Errorf("la carpeta contenedora no se encontró tras crearla")
			}
		} else {
			return fmt.Errorf("la carpeta padre no existe y no se especificó -r")
		}
	}

	// Leer inodo del directorio contenedor
	var parentInode structs.Inode
	offset := int64(sb.S_inode_start + parentInodeIndex*int32(binary.Size(structs.Inode{})))
	if err := utils.ReadObject(file, &parentInode, offset); err != nil {
		return fmt.Errorf("error al leer carpeta contenedora")
	}

	// Verificar si ya existe el archivo
	if ExistsInFolder(fileName, parentInode, file, sb) {
		return fmt.Errorf("el archivo %s ya existe", fileName)
	}

	// Crear inodo para el archivo
	newInode := structs.Inode{
		I_uid:  1,
		I_gid:  1,
		I_size: int32(len(content)),
		I_type: [1]byte{'1'},
		I_perm: [3]byte{'6', '6', '4'},
	}
	for i := range newInode.I_block {
		newInode.I_block[i] = -1
	}

	// Asignar bloque
	newBlockIndex := GetFreeBlockIndex(file, sb)
	if newBlockIndex == -1 {
		return fmt.Errorf("no hay bloques libres disponibles")
	}
	var fileBlock structs.Fileblock
	copy(fileBlock.B_content[:], content)

	// Escribir bloque
	blockOffset := int64(sb.S_block_start + newBlockIndex*int32(binary.Size(fileBlock)))
	if err := utils.WriteObject(file, fileBlock, blockOffset); err != nil {
		return fmt.Errorf("error al escribir bloque")
	}

	// Marcar bitmap
	utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start+newBlockIndex))
	sb.S_free_blocks_count--

	// Actualizar inodo
	newInode.I_block[0] = newBlockIndex
	newInodeIndex := GetFreeInodeIndex(file, sb)
	if newInodeIndex == -1 {
		return fmt.Errorf("no hay inodos libres disponibles")
	}
	inodeOffset := int64(sb.S_inode_start + newInodeIndex*int32(binary.Size(newInode)))
	if err := utils.WriteObject(file, newInode, inodeOffset); err != nil {
		return fmt.Errorf("error al escribir inodo")
	}
	utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start+newInodeIndex))
	sb.S_free_inodes_count--

	// Añadir a la carpeta
	if err := AddToFolder(fileName, parentInodeIndex, newInodeIndex, file, sb); err != nil {
		return fmt.Errorf("error al agregar archivo en carpeta: %v", err)
	}

	// Guardar superbloque actualizado
	if err := utils.WriteObject(file, sb, int64(sb.S_block_start)-int64(binary.Size(sb))); err != nil {
		return fmt.Errorf("error al actualizar el superbloque")
	}

	fmt.Println("Archivo creado exitosamente:", path)
	return nil
}

func GetFreeBlockIndex(file *os.File, sb structs.Superblock) int32 {
	for i := int32(0); i < sb.S_blocks_count; i++ {
		var b byte
		offset := int64(sb.S_bm_block_start + i)
		if err := binaryReadAt(file, &b, offset); err != nil {
			continue
		}
		if b == 0 {
			return i
		}
	}
	return -1
}

func GetFreeInodeIndex(file *os.File, sb structs.Superblock) int32 {
	for i := int32(0); i < sb.S_inodes_count; i++ {
		var b byte
		pos := int64(sb.S_bm_inode_start + i)
		if err := binaryReadAt(file, &b, pos); err != nil {
			continue
		}
		if b == 0 {
			return i
		}
	}
	return -1
}

func binaryReadAt(file *os.File, out interface{}, offset int64) error {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return err
	}
	return binary.Read(file, binary.LittleEndian, out)
}

func GetFreeBlock(file *os.File, sb structs.Superblock) int32 {
	for i := int32(0); i < sb.S_blocks_count; i++ {
		var b byte
		utils.ReadObject(file, &b, int64(sb.S_bm_block_start+i))
		if b == 0 {
			return i
		}
	}
	return -1
}

func GetFreeInode(file *os.File, sb structs.Superblock) int32 {
	for i := int32(0); i < sb.S_inodes_count; i++ {
		var b byte
		utils.ReadObject(file, &b, int64(sb.S_bm_inode_start+i))
		if b == 0 {
			return i
		}
	}
	return -1
}

func ExistsInFolder(name string, inode structs.Inode, file *os.File, sb structs.Superblock) bool {
	for _, block := range inode.I_block {
		if block != -1 {
			var folderBlock structs.Folderblock
			offset := int64(sb.S_block_start + block*int32(binary.Size(structs.Folderblock{})))
			if err := utils.ReadObject(file, &folderBlock, offset); err != nil {
				continue
			}
			for _, entry := range folderBlock.B_content {
				entryName := strings.TrimRight(string(entry.B_name[:]), "\x00")
				if entryName == name {
					return true
				}
			}
		}
	}
	return false
}

func AddToFolder(name string, parentInodeIndex int32, childInodeIndex int32, file *os.File, sb structs.Superblock) error {
	inodeOffset := int64(sb.S_inode_start + parentInodeIndex*int32(binary.Size(structs.Inode{})))

	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
		return err
	}

	for i := 0; i < len(inode.I_block); i++ {
		block := inode.I_block[i]
		if block == -1 {
			newBlockIndex := GetFreeBlockIndex(file, sb)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay bloques disponibles")
			}

			var folderBlock structs.Folderblock
			for j := 0; j < 4; j++ {
				folderBlock.B_content[j].B_inodo = -1
			}
			copy(folderBlock.B_content[0].B_name[:], []byte(name))
			folderBlock.B_content[0].B_inodo = childInodeIndex

			blockOffset := int64(sb.S_block_start + newBlockIndex*int32(binary.Size(folderBlock)))
			utils.WriteObject(file, folderBlock, blockOffset)
			utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start+newBlockIndex))
			sb.S_free_blocks_count--

			inode.I_block[i] = newBlockIndex
			utils.WriteObject(file, inode, inodeOffset)
			return nil
		} else {
			var folderBlock structs.Folderblock
			blockOffset := int64(sb.S_block_start + block*int32(binary.Size(folderBlock)))
			if err := utils.ReadObject(file, &folderBlock, blockOffset); err != nil {
				continue
			}

			for j := 0; j < 4; j++ {
				if strings.Trim(string(folderBlock.B_content[j].B_name[:]), "\x00") == "" {
					copy(folderBlock.B_content[j].B_name[:], []byte(name))
					folderBlock.B_content[j].B_inodo = childInodeIndex
					utils.WriteObject(file, folderBlock, blockOffset)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no hay espacio para agregar más entradas")
}

// Crea carpetas intermedias si no existen (como mkdir -p)
func CreateFolderRecursive(path string, file *os.File, sb structs.Superblock, recursive bool) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("el path debe comenzar con /")
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return fmt.Errorf("ruta inválida")
	}

	var rootInode structs.Inode
	if err := utils.ReadObject(file, &rootInode, int64(sb.S_inode_start)); err != nil {
		return fmt.Errorf("error al leer el inodo raíz")
	}

	currentInode := rootInode
	currentInodeIndex := int32(0)

	for _, part := range parts {
		found := false
		var nextInodeIndex int32 = -1
		var nextInode structs.Inode

		for _, block := range currentInode.I_block {
			if block == -1 {
				continue
			}

			var folderBlock structs.Folderblock
			blockOffset := int64(sb.S_block_start + block*int32(binary.Size(structs.Folderblock{})))
			if err := utils.ReadObject(file, &folderBlock, blockOffset); err != nil {
				continue
			}

			for _, entry := range folderBlock.B_content {
				entryName := strings.TrimRight(string(entry.B_name[:]), "\x00")
				if entryName == part && entry.B_inodo != -1 {
					nextInodeIndex = entry.B_inodo
					inodeOffset := int64(sb.S_inode_start + nextInodeIndex*int32(binary.Size(structs.Inode{})))
					if err := utils.ReadObject(file, &nextInode, inodeOffset); err != nil {
						return err
					}
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if found {
			currentInode = nextInode
			currentInodeIndex = nextInodeIndex
		} else if recursive {
			newInodeIndex := GetFreeInodeIndex(file, sb)
			if newInodeIndex == -1 {
				return fmt.Errorf("no hay inodos disponibles")
			}

			newBlockIndex := GetFreeBlockIndex(file, sb)
			if newBlockIndex == -1 {
				return fmt.Errorf("no hay bloques disponibles")
			}

			// Crear nuevo inodo
			newInode := structs.Inode{
				I_uid:  1,
				I_gid:  1,
				I_size: 0,
				I_type: [1]byte{'0'},
				I_perm: [3]byte{'7', '7', '7'},
			}
			for i := 0; i < len(newInode.I_block); i++ {
				newInode.I_block[i] = -1
			}
			newInode.I_block[0] = newBlockIndex

			// Crear folder block
			var folderBlock structs.Folderblock
			for i := 0; i < 4; i++ {
				folderBlock.B_content[i].B_inodo = -1
			}
			copy(folderBlock.B_content[0].B_name[:], []byte("."))
			folderBlock.B_content[0].B_inodo = newInodeIndex
			copy(folderBlock.B_content[1].B_name[:], []byte(".."))
			folderBlock.B_content[1].B_inodo = currentInodeIndex

			// Escribir
			utils.WriteObject(file, newInode, int64(sb.S_inode_start+newInodeIndex*int32(binary.Size(newInode))))
			utils.WriteObject(file, folderBlock, int64(sb.S_block_start+newBlockIndex*int32(binary.Size(folderBlock))))
			utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start+newInodeIndex))
			utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start+newBlockIndex))
			sb.S_free_inodes_count--
			sb.S_free_blocks_count--

			// Vincular con el padre
			if err := AddToFolder(part, currentInodeIndex, newInodeIndex, file, sb); err != nil {
				return err
			}

			currentInode = newInode
			currentInodeIndex = newInodeIndex
		} else {
			return fmt.Errorf("la carpeta '%s' no existe y -r no fue especificado", part)
		}
	}

	return nil
}
