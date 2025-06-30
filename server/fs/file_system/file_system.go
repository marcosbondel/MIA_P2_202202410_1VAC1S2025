package file_system

import (
	"MIA_P2_202202410_1VAC1S2025/fs/structs"
	"MIA_P2_202202410_1VAC1S2025/fs/utils"
	"MIA_P2_202202410_1VAC1S2025/fs/utils_inodes"
	"MIA_P2_202202410_1VAC1S2025/models"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Mkfs(id string, type_ string, fs_ string, buffer_string *string) {
	fmt.Println("======Start MKFS======")
	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)
	fmt.Println("Fs:", fs_)

	*buffer_string += "======Start MKFS======\n"
	*buffer_string += "Id: " + id + "\n"
	*buffer_string += "Type: " + type_ + "\n"
	*buffer_string += "Fs: " + fs_ + "\n"

	driveletter := string(id[0])
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR structs.MRB
	if err := utils.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 && strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
			if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "1") {
				index = i
				break
			}
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		*buffer_string += "Partition not found\n"
		return
	}

	numerador := int32(TempMBR.Partitions[index].Size - int32(binary.Size(structs.Superblock{})))
	denrominador_base := int32(4 + int32(binary.Size(structs.Inode{})) + 3*int32(binary.Size(structs.Fileblock{})))
	temp := int32(0)
	if fs_ != "2fs" {
		temp = int32(binary.Size(structs.Journaling{}))
	}
	n := int32(numerador / (denrominador_base + temp))
	fmt.Println("N:", n)
	*buffer_string += fmt.Sprintf("N: %d\n", n)

	var newSuperblock structs.Superblock
	newSuperblock.S_inodes_count = n
	newSuperblock.S_blocks_count = 3 * n
	newSuperblock.S_free_blocks_count = 3*n - 2
	newSuperblock.S_free_inodes_count = n - 2
	copy(newSuperblock.S_mtime[:], "28/02/2024")
	copy(newSuperblock.S_umtime[:], "28/02/2024")
	newSuperblock.S_mnt_count = 0

	if fs_ == "2fs" {
		create_ext2(n, TempMBR.Partitions[index], newSuperblock, "16/06/2025", file, buffer_string)
	} else {
		create_ext3(n, TempMBR.Partitions[index], newSuperblock, "16/06/2025", file)
	}

	*buffer_string += "======End MKFS======\n"
	fmt.Println("======End MKFS======")
}

func create_ext2(n int32, partition structs.Partition, sb structs.Superblock, date string, file *os.File, buffer_string *string) {
	fmt.Println("======Start CREATE EXT2======")
	*buffer_string += "======Start CREATE EXT2======\n"

	sb.S_filesystem_type = 2
	sb.S_bm_inode_start = partition.Start + int32(binary.Size(structs.Superblock{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + n
	sb.S_inode_start = sb.S_bm_block_start + 3*n
	sb.S_block_start = sb.S_inode_start + n*int32(binary.Size(structs.Inode{}))

	// Asignar valores importantes del superbloque
	sb.S_magic = 0xEF53
	sb.S_inode_size = int32(binary.Size(structs.Inode{}))
	sb.S_block_size = int32(binary.Size(structs.Fileblock{}))
	sb.S_fist_ino = 2
	sb.S_first_blo = 2

	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, byte(0), int64(sb.S_bm_inode_start+i))
	}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, byte(0), int64(sb.S_bm_block_start+i))
	}
	for i := int32(0); i < n; i++ {
		var inode structs.Inode
		for j := range inode.I_block {
			inode.I_block[j] = -1
		}
		utils.WriteObject(file, inode, int64(sb.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
	}
	for i := int32(0); i < 3*n; i++ {
		var block structs.Fileblock
		utils.WriteObject(file, block, int64(sb.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
	}

	var inode0 structs.Inode
	inode0.I_uid = 1
	inode0.I_gid = 1
	inode0.I_size = 0
	copy(inode0.I_atime[:], date)
	copy(inode0.I_ctime[:], date)
	copy(inode0.I_mtime[:], date)
	copy(inode0.I_perm[:], "664")
	for i := range inode0.I_block {
		inode0.I_block[i] = -1
	}
	inode0.I_block[0] = 0
	inode0.I_type = [1]byte{'0'}

	var folderBlock0 structs.Folderblock
	folderBlock0.B_content[0] = structs.Content{B_inodo: 0}
	copy(folderBlock0.B_content[0].B_name[:], ".")
	folderBlock0.B_content[1] = structs.Content{B_inodo: 0}
	copy(folderBlock0.B_content[1].B_name[:], "..")
	folderBlock0.B_content[2] = structs.Content{B_inodo: 1}
	copy(folderBlock0.B_content[2].B_name[:], "users.txt")

	var data = "1,G,root\n1,U,root,root,123\n"
	var inode1 structs.Inode
	inode1.I_uid = 1
	inode1.I_gid = 1
	inode1.I_size = int32(len(data))
	copy(inode1.I_atime[:], date)
	copy(inode1.I_ctime[:], date)
	copy(inode1.I_mtime[:], date)
	copy(inode1.I_perm[:], "664")
	inode1.I_type = [1]byte{'1'}
	for i := range inode1.I_block {
		inode1.I_block[i] = -1
	}
	inode1.I_block[0] = 1

	var fileBlock1 structs.Fileblock
	copy(fileBlock1.B_content[:], []byte(data))

	utils.WriteObject(file, sb, int64(partition.Start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start+1))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start+1))
	utils.WriteObject(file, inode0, int64(sb.S_inode_start))
	utils.WriteObject(file, inode1, int64(sb.S_inode_start+int32(binary.Size(inode0))))
	utils.WriteObject(file, folderBlock0, int64(sb.S_block_start))
	utils.WriteObject(file, fileBlock1, int64(sb.S_block_start+int32(binary.Size(fileBlock1))))

	*buffer_string += "======End CREATE EXT2======\n"
	fmt.Println("======End CREATE EXT2======")
}

func create_ext3(n int32, partition structs.Partition, sb structs.Superblock, date string, file *os.File) {
	fmt.Println("======Start CREATE EXT3======")
	sb.S_filesystem_type = 3
	sb.S_bm_inode_start = partition.Start + int32(binary.Size(structs.Superblock{})) + n*int32(binary.Size(structs.Journaling{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + n
	sb.S_inode_start = sb.S_bm_block_start + 3*n
	sb.S_block_start = sb.S_inode_start + n*int32(binary.Size(structs.Inode{}))

	// Asignar valores importantes del superbloque
	sb.S_magic = 0xEF53
	sb.S_inode_size = int32(binary.Size(structs.Inode{}))
	sb.S_block_size = int32(binary.Size(structs.Fileblock{}))
	sb.S_fist_ino = 2
	sb.S_first_blo = 2

	// Inicializar journaling
	for i := int32(0); i < n; i++ {
		var emptyJ structs.Journaling
		utils.WriteObject(file, emptyJ, int64(partition.Start+int32(binary.Size(structs.Superblock{}))+i*int32(binary.Size(structs.Journaling{}))))
	}

	for i := int32(0); i < n; i++ {
		utils.WriteObject(file, byte(0), int64(sb.S_bm_inode_start+i))
	}
	for i := int32(0); i < 3*n; i++ {
		utils.WriteObject(file, byte(0), int64(sb.S_bm_block_start+i))
	}
	for i := int32(0); i < n; i++ {
		var inode structs.Inode
		for j := range inode.I_block {
			inode.I_block[j] = -1
		}
		utils.WriteObject(file, inode, int64(sb.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
	}
	for i := int32(0); i < 3*n; i++ {
		var block structs.Fileblock
		utils.WriteObject(file, block, int64(sb.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
	}

	var inode0 structs.Inode
	inode0.I_uid = 1
	inode0.I_gid = 1
	inode0.I_size = 0
	copy(inode0.I_atime[:], date)
	copy(inode0.I_ctime[:], date)
	copy(inode0.I_mtime[:], date)
	copy(inode0.I_perm[:], "664")
	inode0.I_type = [1]byte{'0'}
	for i := range inode0.I_block {
		inode0.I_block[i] = -1
	}
	inode0.I_block[0] = 0

	var folderBlock0 structs.Folderblock
	folderBlock0.B_content[0] = structs.Content{B_inodo: 0}
	copy(folderBlock0.B_content[0].B_name[:], ".")
	folderBlock0.B_content[1] = structs.Content{B_inodo: 0}
	copy(folderBlock0.B_content[1].B_name[:], "..")
	folderBlock0.B_content[2] = structs.Content{B_inodo: 1}
	copy(folderBlock0.B_content[2].B_name[:], "users.txt")

	var data = "1,G,root\n1,U,root,root,123\n"
	var inode1 structs.Inode
	inode1.I_uid = 1
	inode1.I_gid = 1
	inode1.I_size = int32(len(data))
	copy(inode1.I_atime[:], date)
	copy(inode1.I_ctime[:], date)
	copy(inode1.I_mtime[:], date)
	copy(inode1.I_perm[:], "664")
	inode1.I_type = [1]byte{'1'}
	for i := range inode1.I_block {
		inode1.I_block[i] = -1
	}
	inode1.I_block[0] = 1

	var fileBlock1 structs.Fileblock
	copy(fileBlock1.B_content[:], []byte(data))

	utils.WriteObject(file, sb, int64(partition.Start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_inode_start+1))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start))
	utils.WriteObject(file, byte(1), int64(sb.S_bm_block_start+1))
	utils.WriteObject(file, inode0, int64(sb.S_inode_start))
	utils.WriteObject(file, inode1, int64(sb.S_inode_start+int32(binary.Size(inode0))))
	utils.WriteObject(file, folderBlock0, int64(sb.S_block_start))
	utils.WriteObject(file, fileBlock1, int64(sb.S_block_start+int32(binary.Size(fileBlock1))))

	fmt.Println("======End CREATE EXT3======")
}

// ListDisks escanea la carpeta ./test y retorna una lista de discos disponibles
// controllers/disk_controller.go
func ListDisks() ([]models.DiskInfo, error) {
	files, err := ioutil.ReadDir("./fs/test")
	if err != nil {
		return nil, err
	}

	var disks []models.DiskInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".bin") {
			filePath := "./fs/test/" + file.Name()
			f, err := os.Open(filePath)
			if err != nil {
				continue
			}

			var mbr structs.MRB
			if err := utils.ReadObject(f, &mbr, 0); err != nil {
				f.Close()
				continue
			}
			f.Close()

			var partitions []models.PartitionInfo
			for _, p := range mbr.Partitions {
				if p.Size == 0 {
					continue
				}
				partitions = append(partitions, models.PartitionInfo{
					Status: strings.Trim(string(p.Status[:]), "\x00"),
					Type:   strings.ToUpper(strings.Trim(string(p.Type[:]), "\x00")),
					Fit:    strings.Trim(string(p.Fit[:]), "\x00"),
					Start:  p.Start,
					Size:   p.Size,
					Name:   strings.Trim(string(p.Name[:]), "\x00"),
				})
			}

			disk := models.DiskInfo{
				Name:       strings.TrimSuffix(file.Name(), ".bin"),
				Size:       mbr.MbrSize,
				Creation:   strings.Trim(string(mbr.CreationDate[:]), "\x00"),
				Fit:        strings.Trim(string(mbr.Fit[:]), "\x00"),
				Signature:  mbr.Signature,
				Partitions: partitions,
			}
			disks = append(disks, disk)
		}
	}

	return disks, nil
}

func GetDiskPartitions(driveLetter string) ([]models.PartitionInfo, error) {
	path := "./fs/test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo del disco: %v", err)
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		return nil, fmt.Errorf("no se pudo leer el MBR: %v", err)
	}

	var partitions []models.PartitionInfo

	for i, part := range mbr.Partitions {
		if part.Size == 0 {
			continue
		}

		partType := strings.ToLower(strings.Trim(string(part.Type[:]), "\x00"))
		partStatus := strings.Trim(string(part.Status[:]), "\x00")
		partFit := strings.Trim(string(part.Fit[:]), "\x00")
		partName := strings.Trim(string(part.Name[:]), "\x00")
		partId := strings.Trim(string(part.Id[:]), "\x00")
		mounted := partStatus == "1"

		if partType == "e" {
			// Añadir partición extendida
			extended := models.PartitionInfo{
				Status:         partStatus,
				Type:           "E",
				Fit:            partFit,
				Start:          part.Start,
				Size:           part.Size,
				Name:           partName,
				Id:             partId,
				Mounted:        mounted,
				PartitionIndex: i + 1,
			}
			partitions = append(partitions, extended)

			// Procesar EBRs correctamente
			next := part.Start
			visited := make(map[int32]bool)

			for {
				if visited[next] {
					fmt.Println("Ciclo detectado en EBR, abortando...")
					break
				}
				visited[next] = true

				var ebr structs.EBR
				if err := utils.ReadObject(file, &ebr, int64(next)); err != nil {
					break
				}

				if ebr.PartSize > 0 {
					p := models.PartitionInfo{
						Status:    fmt.Sprintf("%d", ebr.PartStatus),
						Type:      "L",
						Fit:       strings.Trim(string(ebr.PartFit[:]), "\x00"),
						Start:     ebr.PartStart,
						Size:      ebr.PartSize,
						Name:      strings.Trim(string(ebr.PartName[:]), "\x00"),
						Mounted:   ebr.PartStatus == 1,
						IsLogical: true,
						PartNext:  ebr.PartNext,
						EBRStart:  next,
					}
					partitions = append(partitions, p)
				}

				if ebr.PartNext <= 0 || ebr.PartNext == next {
					break
				}
				next = ebr.PartNext
			}
		} else {
			// Partición primaria
			p := models.PartitionInfo{
				Status:         partStatus,
				Type:           "P",
				Fit:            partFit,
				Start:          part.Start,
				Size:           part.Size,
				Name:           partName,
				Id:             partId,
				Mounted:        mounted,
				PartitionIndex: i + 1,
			}
			partitions = append(partitions, p)
		}
	}

	fmt.Printf("✅ Particiones encontradas (%s): %d\n", driveLetter, len(partitions))
	return partitions, nil
}

func GetFileSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	// Simular sesión activa
	id := "A110"
	driveLetter := string(id[0])
	binPath := "./fs/test/" + strings.ToUpper(driveLetter) + ".bin"

	file, err := os.Open(binPath)
	if err != nil {
		http.Error(w, "Error abriendo disco", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		http.Error(w, "Error leyendo MBR", http.StatusInternalServerError)
		return
	}

	index := int(id[1] - '1')
	part := mbr.Partitions[index]

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		http.Error(w, "Error leyendo superbloque", http.StatusInternalServerError)
		return
	}

	inodeIndex := utils_inodes.InitSearch(path, file, sb)
	if inodeIndex == -1 {
		http.Error(w, "Ruta no encontrada", http.StatusNotFound)
		return
	}

	var inode structs.Inode
	offset := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(inode))
	if err := utils.ReadObject(file, &inode, offset); err != nil {
		http.Error(w, "Error leyendo inodo", http.StatusInternalServerError)
		return
	}

	// Si es directorio
	if inode.I_type[0] == '0' {
		var children []map[string]interface{}

		for _, blockIndex := range inode.I_block {
			if blockIndex == -1 {
				continue
			}
			var folderBlock structs.Folderblock
			blockOffset := int64(sb.S_block_start) + int64(blockIndex)*int64(binary.Size(folderBlock))
			if err := utils.ReadObject(file, &folderBlock, blockOffset); err != nil {
				continue
			}

			for _, entry := range folderBlock.B_content {
				name := strings.Trim(string(entry.B_name[:]), "\x00")
				if name == "" || name == "." || name == ".." {
					continue
				}
				childOffset := int64(sb.S_inode_start) + int64(entry.B_inodo)*int64(binary.Size(inode))
				var childInode structs.Inode
				if err := utils.ReadObject(file, &childInode, childOffset); err != nil {
					continue
				}
				entryType := "file"
				if childInode.I_type[0] == '0' {
					entryType = "directory"
				}

				children = append(children, map[string]interface{}{
					"name":        name,
					"type":        entryType,
					"uid":         childInode.I_uid,
					"gid":         childInode.I_gid,
					"size":        childInode.I_size,
					"perm":        strings.Trim(string(childInode.I_perm[:]), "\x00"),
					"created_at":  strings.Trim(string(childInode.I_ctime[:]), "\x00"),
					"modified_at": strings.Trim(string(childInode.I_mtime[:]), "\x00"),
					"accessed_at": strings.Trim(string(childInode.I_atime[:]), "\x00"),
				})
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"type":        "directory",
			"path":        path,
			"uid":         inode.I_uid,
			"gid":         inode.I_gid,
			"perm":        strings.Trim(string(inode.I_perm[:]), "\x00"),
			"modified_at": strings.Trim(string(inode.I_mtime[:]), "\x00"),
			"children":    children,
		})

	} else {
		// Es archivo
		content := utils_inodes.GetInodeFileData(inode, file, sb)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"type":        "file",
			"path":        path,
			"uid":         inode.I_uid,
			"gid":         inode.I_gid,
			"perm":        strings.Trim(string(inode.I_perm[:]), "\x00"),
			"size":        inode.I_size,
			"created_at":  strings.Trim(string(inode.I_ctime[:]), "\x00"),
			"modified_at": strings.Trim(string(inode.I_mtime[:]), "\x00"),
			"accessed_at": strings.Trim(string(inode.I_atime[:]), "\x00"),
			"content":     content,
		})
	}
}
