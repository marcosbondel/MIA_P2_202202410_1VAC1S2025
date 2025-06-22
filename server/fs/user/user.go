package user

import (
	//  "os"
	"MIA_P2_202202410_1VAC1S2025/fs/global"
	"MIA_P2_202202410_1VAC1S2025/fs/structs"
	"MIA_P2_202202410_1VAC1S2025/fs/utils"
	"MIA_P2_202202410_1VAC1S2025/fs/utils_inodes"
	"encoding/binary"
	"fmt"
	"strings"
)

// // login -user=root -pass=123 -id=A119
func Login(user string, pass string, id string, buffer_string *string) bool {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	*buffer_string += "======Start LOGIN======\n"
	*buffer_string += "User: " + user + "\n"
	*buffer_string += "Pass: " + pass + "\n"
	*buffer_string += "Id: " + id + "\n"

	if global.CurrentUser.Status {
		fmt.Println("User already logged in")
		*buffer_string += "User already logged in\n"
		return false
	}

	var login bool = false
	driveletter := string(id[0])

	// Open bin file
	filepath := "./fs/test/" + strings.ToUpper(driveletter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		return false
	}

	var TempMBR structs.MRB
	// Read object from bin file
	if err := utils.ReadObject(file, &TempMBR, 0); err != nil {
		return false
	}

	// Print object
	structs.PrintMBR(TempMBR)

	fmt.Println("-------------")
	*buffer_string += "-------------\n"

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				*buffer_string += "Partition found\n"
				if strings.Contains(string(TempMBR.Partitions[i].Status[:]), "1") {
					fmt.Println("Partition is mounted")
					*buffer_string += "Partition is mounted\n"
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					*buffer_string += "Partition is not mounted\n"
					return false
				}
				break
			}
		}
	}

	if index != -1 {
		structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		*buffer_string += "Partition not found\n"
		return false
	}

	var tempSuperblock structs.Superblock
	// Read object from bin file
	if err := utils.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return false
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1
	indexInode := utils_inodes.InitSearch("/users.txt", file, tempSuperblock)

	// indexInode := int32(1)

	var crrInode structs.Inode
	// Read object from bin file
	if err := utils.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		return false
	}

	// read file data
	data := utils_inodes.GetInodeFileData(crrInode, file, tempSuperblock)

	fmt.Println("Fileblock------------")
	*buffer_string += "Fileblock------------\n"
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// login -user=root -pass=123 -id=A119

	// Iterar a través de las líneas
	for _, line := range lines {
		// Imprimir cada línea
		// fmt.Println(line)
		words := strings.Split(line, ",")
		fmt.Println("Words:", words)
		*buffer_string += fmt.Sprintf("Words: %v\n", words)
		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true

				break
			}
		}
	}

	// Print object
	fmt.Println("Inode", crrInode.I_block)
	*buffer_string += "Inode: \n"
	*buffer_string += fmt.Sprintf("%v\n", crrInode.I_block)

	// Close bin file
	defer file.Close()

	if login {
		fmt.Println("User logged in")
		*buffer_string += "User logged in\n"
		global.CurrentUser.ID = id
		global.CurrentUser.Status = true
		global.CurrentUser.User = user
	} else {
		fmt.Println("User not found or invalid credentials")
		*buffer_string += "User not found or invalid credentials\n"
		return false
	}

	*buffer_string += "======End LOGIN======\n"
	fmt.Println("======End LOGIN======")
	return true
}

func Logout(buffer_string *string) bool {
	fmt.Println("======Start LOGOUT======")
	*buffer_string += "======Start LOGOUT======\n"

	if global.CurrentUser.Status {
		global.CurrentUser.ID = ""
		global.CurrentUser.Status = false
		fmt.Println("User logged out")
		*buffer_string += "User logged out\n"
		*buffer_string += "======End LOGOUT======\n"
		fmt.Println("======End LOGOUT======")
		return true
	} else {
		fmt.Println("No user logged in")
		*buffer_string += "No user logged in\n"
		*buffer_string += "======End LOGOUT======\n"
		fmt.Println("======End LOGOUT======")
		return false
	}
	// *buffer_string += "======End LOGOUT======\n"
	// fmt.Println("======End LOGOUT======")
}
