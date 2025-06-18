package reports_generator

import (
	"MIA_P1_202202410_1VAC1S2025/global"
	"MIA_P1_202202410_1VAC1S2025/structs"
	"MIA_P1_202202410_1VAC1S2025/utils"
	"MIA_P1_202202410_1VAC1S2025/utils_inodes"
	"encoding/binary"
	"fmt"
	"html"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func GenerarReporte(name string, path string, id string, ruta string) {
	switch name {
	case "mbr":
		GenerarReporteMBR(path, id)
	// case "tree":
	// 	GenerarReporteTrees(path, id)
	case "disk":
		GenerarReporteDISK(path, id)
	case "inode":
		GenerarReporteINODE(path, id)
	case "bm_inode":
		GenerarReporteBMInode(path, id)
	case "bm_bloc":
		GenerarReporteBMBlock(path, id)
	case "tree":
		GenerarReporteTREE(path, id)
	case "block":
		GenerarReporteBlock(path, id)
	case "sb":
		GenerarReporteSB(path, id)
	case "file":
		GenerarReporteFile(id, path, ruta)
	case "ls":
		GenerarReporteLS(id, path, ruta)
	default:
		fmt.Println("ERROR: Reporte no reconocido:", name)
	}
}

func GenerarReporteTree1(path string, id string) {
	dotContent := `digraph G {
		node [shape=record]
		tree [label="Reporte Tree\n(users.txt)"]
	}`

	// Guardar archivo .dot temporal
	dotPath := "tree.dot"
	if err := os.WriteFile(dotPath, []byte(dotContent), 0644); err != nil {
		fmt.Println("Error creando archivo .dot:", err)
		return
	}

	// Ejecutar Graphviz para crear .png
	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte generado:", path)
}

// =================== REPORTE MBR ===================

func GenerarReporteMBR(path string, id string) {
	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
		return
	}

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=none fontname=\"Helvetica\"];\n")
	dot.WriteString("mbr [label=<\n")
	dot.WriteString("<table border='1' cellborder='1' cellspacing='0' cellpadding='4'>\n")
	dot.WriteString("<tr><td bgcolor='purple' colspan='2'><font color='white'><b>REPORTE DE MBR</b></font></td></tr>\n")

	dot.WriteString(fmt.Sprintf("<tr><td bgcolor='#f2e6ff'><b>mbr_tamano</b></td><td>%d</td></tr>\n", mbr.MbrSize))
	dot.WriteString(fmt.Sprintf("<tr><td bgcolor='#f2e6ff'><b>mbr_fecha_creacion</b></td><td>%s</td></tr>\n", strings.Trim(string(mbr.CreationDate[:]), "\x00")))
	dot.WriteString(fmt.Sprintf("<tr><td bgcolor='#f2e6ff'><b>mbr_disk_signature</b></td><td>0x%X</td></tr>\n", mbr.Signature))

	for i := 0; i < 4; i++ {
		part := mbr.Partitions[i]
		if part.Size == 0 {
			continue
		}

		partName := html.EscapeString(strings.TrimRight(string(part.Name[:]), "\x00"))
		partType := strings.Trim(string(part.Type[:]), "\x00")

		dot.WriteString("<tr><td bgcolor='indigo' colspan='2'><font color='white'><b>Partition</b></font></td></tr>\n")
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_status</b></td><td>%s</td></tr>\n", strings.Trim(string(part.Status[:]), "\x00")))
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_type</b></td><td>%s</td></tr>\n", partType))
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_fit</b></td><td>%s</td></tr>\n", strings.Trim(string(part.Fit[:]), "\x00")))
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_start</b></td><td>%d</td></tr>\n", part.Start))
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_size</b></td><td>%d</td></tr>\n", part.Size))
		dot.WriteString(fmt.Sprintf("<tr><td><b>part_name</b></td><td>%s</td></tr>\n", partName))

		if strings.ToLower(partType) == "e" {
			var next int32 = part.Start
			visited := make(map[int32]bool)

			for next != -1 {
				if visited[next] {
					fmt.Println("⚠️  Detenido: ciclo detectado en lista de EBRs")
					break
				}
				visited[next] = true

				var ebr structs.EBR
				if err := utils.ReadObject(file, &ebr, int64(next)); err != nil {
					fmt.Println("Error al leer EBR:", err)
					break
				}

				if ebr.PartSize <= 0 {
					break
				}

				logicalName := html.EscapeString(strings.TrimRight(string(ebr.PartName[:]), "\x00"))
				logicalFit := html.EscapeString(strings.Trim(string(ebr.PartFit[:]), "\x00"))

				dot.WriteString("<tr><td bgcolor='#ff9999' colspan='2'><b>Particion Logica</b></td></tr>\n")
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_status</b></td><td>%d</td></tr>\n", ebr.PartStatus))
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_next</b></td><td>%d</td></tr>\n", ebr.PartNext))
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_fit</b></td><td>%s</td></tr>\n", logicalFit))
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_start</b></td><td>%d</td></tr>\n", ebr.PartStart))
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_size</b></td><td>%d</td></tr>\n", ebr.PartSize))
				dot.WriteString(fmt.Sprintf("<tr><td><b>part_name</b></td><td>%s</td></tr>\n", logicalName))

				if ebr.PartNext <= 0 || ebr.PartNext == next {
					break
				}
				next = ebr.PartNext
			}
		}
	}

	dot.WriteString("</table>>];\n")
	dot.WriteString("}\n")

	dotPath := "mbr.dot"
	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
		fmt.Println("Error al escribir archivo DOT:", err)
		return
	}

	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error ejecutando Graphviz:", err)
		return
	}

	fmt.Println("✅ Reporte MBR generado exitosamente:", path)
}

// func GenerarReporteDISK(path string, id string) {
// 	driveLetter := string(id[0])
// 	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

// 	file, err := utils.OpenFile(filepath)
// 	if err != nil {
// 		fmt.Println("Error al abrir disco:", err)
// 		return
// 	}
// 	defer file.Close()

// 	var mbr structs.MRB
// 	if err := utils.ReadObject(file, &mbr, 0); err != nil {
// 		fmt.Println("Error al leer MBR:", err)
// 		return
// 	}

// 	diskSize := mbr.MbrSize
// 	var dot strings.Builder
// 	dot.WriteString("digraph G {\n")
// 	dot.WriteString("node [shape=plaintext fontname=\"Helvetica\"];\n")
// 	dot.WriteString("struct [label=<\n")
// 	dot.WriteString("<table border='1' cellborder='1' cellspacing='0'>\n")
// 	dot.WriteString("<tr><td colspan='100'><b>Disco</b></td></tr>\n")
// 	dot.WriteString("<tr>")

// 	// MBR
// 	dot.WriteString("<td><b>MBR</b></td>\n")
// 	actualPos := int32(binary.Size(mbr))

// 	// Ordenar particiones por Start
// 	particiones := make([]structs.Partition, 0)
// 	for _, p := range mbr.Partitions {
// 		if p.Size > 0 {
// 			particiones = append(particiones, p)
// 		}
// 	}
// 	sort.Slice(particiones, func(i, j int) bool {
// 		return particiones[i].Start < particiones[j].Start
// 	})

// 	for _, part := range particiones {
// 		if actualPos < part.Start {
// 			// Espacio libre antes de la partición
// 			freeSize := part.Start - actualPos
// 			porcentaje := float64(freeSize) / float64(diskSize) * 100
// 			dot.WriteString(fmt.Sprintf("<td>Libre<br/>%.2f%%</td>\n", porcentaje))
// 			actualPos += freeSize
// 		}

// 		partType := strings.Trim(string(part.Type[:]), "\x00")
// 		partName := strings.TrimRight(string(part.Name[:]), "\x00")
// 		porcentaje := float64(part.Size) / float64(diskSize) * 100

// 		if partType == "e" {
// 			// Extendida
// 			dot.WriteString(fmt.Sprintf("<td><table border='0'><tr><td colspan='10'><b>Extendida</b></td></tr><tr>"))

// 			var nextEBR int32 = part.Start
// 			for {
// 				var ebr structs.EBR
// 				if err := utils.ReadObject(file, &ebr, int64(nextEBR)); err != nil {
// 					break
// 				}
// 				if ebr.PartSize == 0 {
// 					break
// 				}
// 				ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
// 				ebrSize := float64(ebr.PartSize) / float64(diskSize) * 100

// 				dot.WriteString(fmt.Sprintf("<td><b>Lógica</b><br/>%s<br/>%.2f%%</td>\n", ebrName, ebrSize))

// 				if ebr.PartNext == -1 {
// 					break
// 				}
// 				if ebr.PartStart+ebr.PartSize < ebr.PartNext {
// 					// Espacio libre entre EBRs
// 					free := ebr.PartNext - (ebr.PartStart + ebr.PartSize)
// 					porcFree := float64(free) / float64(diskSize) * 100
// 					dot.WriteString(fmt.Sprintf("<td>Libre<br/>%.2f%%</td>\n", porcFree))
// 				}
// 				nextEBR = ebr.PartNext
// 			}

// 			dot.WriteString("</tr></table></td>\n")
// 		} else {
// 			// Primaria
// 			dot.WriteString(fmt.Sprintf("<td><b>Primaria</b><br/>%s<br/>%.2f%%</td>\n", partName, porcentaje))
// 		}

// 		actualPos += part.Size
// 	}

// 	// Espacio libre final
// 	if actualPos < diskSize {
// 		freeSize := diskSize - actualPos
// 		porcentaje := float64(freeSize) / float64(diskSize) * 100
// 		dot.WriteString(fmt.Sprintf("<td>Libre<br/>%.2f%%</td>\n", porcentaje))
// 	}

// 	dot.WriteString("</tr></table>>];\n}")
// 	dot.WriteString("}\n")

// 	// Escribir archivo .dot
// 	dotPath := "disk.dot"
// 	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
// 		fmt.Println("Error al escribir archivo DOT:", err)
// 		return
// 	}

// 	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)

// 	if err := cmd.Run(); err != nil {
// 		fmt.Println("Error ejecutando Graphviz:", err)
// 		return
// 	}

// 	fmt.Println("Reporte DISK generado exitosamente:", path)
// }

func GenerarReporteDISK(path string, id string) {
	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
		return
	}

	diskSize := mbr.MbrSize
	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext fontname=\"Helvetica\"];\n")
	dot.WriteString("struct [label=<\n")
	dot.WriteString("<table border='1' cellborder='1' cellspacing='0'>\n")
	dot.WriteString("<tr><td colspan='100'><b>Disco</b></td></tr>\n")
	dot.WriteString("<tr>")

	// MBR
	dot.WriteString("<td><b>MBR</b></td>\n")
	actualPos := int32(binary.Size(mbr))

	// Ordenar particiones por Start
	particiones := make([]structs.Partition, 0)
	for _, p := range mbr.Partitions {
		if p.Size > 0 {
			particiones = append(particiones, p)
		}
	}
	sort.Slice(particiones, func(i, j int) bool {
		return particiones[i].Start < particiones[j].Start
	})

	for _, part := range particiones {
		if actualPos < part.Start {
			// Espacio libre antes de la partición
			freeSize := part.Start - actualPos
			porcentaje := float64(freeSize) / float64(diskSize) * 100
			dot.WriteString(fmt.Sprintf("<td>Libre<br/>%.2f%%</td>\n", porcentaje))
			actualPos += freeSize
		}

		partType := strings.Trim(string(part.Type[:]), "\x00")
		partName := strings.TrimRight(string(part.Name[:]), "\x00")
		porcentaje := float64(part.Size) / float64(diskSize) * 100

		if partType == "e" {
			// Extendida con contenido directo (simplificado, sin tabla anidada)
			dot.WriteString(fmt.Sprintf("<td><b>Extendida</b><br/>%s<br/>%.2f%%", partName, porcentaje))

			var nextEBR int32 = part.Start
			for {
				var ebr structs.EBR
				if err := utils.ReadObject(file, &ebr, int64(nextEBR)); err != nil {
					break
				}
				if ebr.PartSize == 0 {
					break
				}
				ebrName := strings.TrimRight(string(ebr.PartName[:]), "\x00")
				ebrSize := float64(ebr.PartSize) / float64(diskSize) * 100

				dot.WriteString(fmt.Sprintf("<br/><b>Lógica</b>: %s %.2f%%", ebrName, ebrSize))

				if ebr.PartNext == -1 {
					break
				}
				if ebr.PartStart+ebr.PartSize < ebr.PartNext {
					// Espacio libre entre EBRs
					free := ebr.PartNext - (ebr.PartStart + ebr.PartSize)
					porcFree := float64(free) / float64(diskSize) * 100
					dot.WriteString(fmt.Sprintf("<br/>Libre: %.2f%%", porcFree))
				}
				nextEBR = ebr.PartNext
			}
			dot.WriteString("</td>\n")
		} else {
			// Primaria
			dot.WriteString(fmt.Sprintf("<td><b>Primaria</b><br/>%s<br/>%.2f%%</td>\n", partName, porcentaje))
		}

		actualPos += part.Size
	}

	// Espacio libre final
	if actualPos < diskSize {
		freeSize := diskSize - actualPos
		porcentaje := float64(freeSize) / float64(diskSize) * 100
		dot.WriteString(fmt.Sprintf("<td>Libre<br/>%.2f%%</td>\n", porcentaje))
	}

	dot.WriteString("</tr></table>>];\n")
	dot.WriteString("}\n")

	// Escribir archivo .dot
	dotPath := "disk.dot"
	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
		fmt.Println("Error al escribir archivo DOT:", err)
		return
	}

	// Ejecutar Graphviz
	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte DISK generado exitosamente:", path)
}

func GenerarReporteINODE(path string, id string) {
	fmt.Println("======Start REP INODE======")

	if !global.CurrentUser.Status {
		fmt.Println("ERROR: Debes iniciar sesión.")
		return
	}

	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el disco.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("ERROR: Índice fuera de rango.")
		return
	}

	part := mbr.Partitions[index]
	sb := structs.Superblock{}
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	// Crear DOT
	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext fontname=\"Courier\"];\n")
	dot.WriteString("inodetable [label=<\n")
	dot.WriteString("<table border='1' cellborder='1' cellspacing='0'>\n")
	dot.WriteString("<tr><td colspan='20'><b>Reporte de Inodos</b></td></tr>\n")
	dot.WriteString("<tr><td><b>#</b></td><td><b>UID</b></td><td><b>GID</b></td><td><b>Tamaño</b></td><td><b>Tipo</b></td><td><b>Permisos</b></td><td colspan='15'><b>Bloques</b></td></tr>\n")

	for i := int32(0); i < sb.S_inodes_count; i++ {
		var inode structs.Inode
		offset := int64(sb.S_inode_start) + int64(i)*int64(binary.Size(inode))
		if err := utils.ReadObject(file, &inode, offset); err != nil {
			break
		}

		// Si es inodo vacío, omitir
		if inode.I_size == 0 {
			continue
		}

		dot.WriteString("<tr>")
		dot.WriteString(fmt.Sprintf("<td>%d</td><td>%d</td><td>%d</td><td>%d</td>", i, inode.I_uid, inode.I_gid, inode.I_size))

		tipo := "Archivo"
		if inode.I_type[0] == '0' {
			tipo = "Carpeta"
		}
		dot.WriteString(fmt.Sprintf("<td>%s</td>", tipo))
		dot.WriteString(fmt.Sprintf("<td>%s</td>", string(inode.I_perm[:])))

		for j := 0; j < 15; j++ {
			dot.WriteString(fmt.Sprintf("<td>%d</td>", inode.I_block[j]))
		}

		dot.WriteString("</tr>\n")
	}

	// dot.WriteString("</table>>];\n}")
	// dot.WriteString("}\n")

	dot.WriteString("</table>>];\n")
	dot.WriteString("}\n")

	// Guardar el .dot
	dotPath := "inode.dot"
	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
		fmt.Println("ERROR: No se pudo guardar el archivo DOT.")
		return
	}

	// Generar imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte INODE generado exitosamente:", path)
}

func GenerarReporteBMInode(path string, id string) {
	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		return
	}

	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo binario.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("ERROR: Índice fuera de rango.")
		return
	}

	sb := structs.Superblock{}
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	n := sb.S_inodes_count
	bitmap := make([]byte, n)
	if _, err := file.ReadAt(bitmap, int64(sb.S_bm_inode_start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el bitmap de inodos.")
		return
	}

	var builder strings.Builder
	for i, b := range bitmap {
		builder.WriteString(fmt.Sprintf("%d", b))
		if (i+1)%20 == 0 {
			builder.WriteString("\n")
		}
	}

	if err := os.WriteFile(path, []byte(builder.String()), 0644); err != nil {
		fmt.Println("ERROR: No se pudo guardar el archivo del reporte.")
		return
	}

	fmt.Println("Reporte bm_inode generado exitosamente:", path)
}

func GenerarReporteBMBlock(path string, id string) {
	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		return
	}

	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo binario.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("ERROR: Índice fuera de rango.")
		return
	}

	sb := structs.Superblock{}
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	n := sb.S_blocks_count
	bitmap := make([]byte, n)
	if _, err := file.ReadAt(bitmap, int64(sb.S_bm_block_start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el bitmap de bloques.")
		return
	}

	var builder strings.Builder
	for i, b := range bitmap {
		if b == 1 {
			builder.WriteByte('1')
		} else {
			builder.WriteByte('0')
		}
		if (i+1)%20 == 0 {
			builder.WriteByte('\n')
		}
	}

	if err := os.WriteFile(path, []byte(builder.String()), 0644); err != nil {
		fmt.Println("ERROR: No se pudo guardar el archivo del reporte.")
		return
	}

	fmt.Println("Reporte bm_block generado exitosamente:", path)
}

func GenerarReporteTREE(path string, id string) {
	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		return
	}

	driveLetter := string(id[0])
	filePath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filePath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo binario.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice fuera de rango.")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	var dot strings.Builder
	dot.WriteString("digraph Tree {\n")
	dot.WriteString("rankdir=LR;\n")
	dot.WriteString("node [shape=plaintext fontname=\"Helvetica\"];\n")

	for i := int32(0); i < sb.S_inodes_count; i++ {
		var inode structs.Inode
		offset := int64(sb.S_inode_start) + int64(i)*int64(binary.Size(inode))
		if err := utils.ReadObject(file, &inode, offset); err != nil {
			continue
		}
		if inode.I_uid == 0 && inode.I_gid == 0 {
			continue
		}

		nodo := fmt.Sprintf("inode%d", i)
		dot.WriteString(fmt.Sprintf("%s [label=<\n<table border='1' cellborder='1' cellspacing='0'>\n", nodo))
		dot.WriteString(fmt.Sprintf("<tr><td colspan='2'><b>Inodo %d</b></td></tr>\n", i))
		dot.WriteString(fmt.Sprintf("<tr><td><b>UID</b></td><td>%d</td></tr>\n", inode.I_uid))
		dot.WriteString(fmt.Sprintf("<tr><td><b>GID</b></td><td>%d</td></tr>\n", inode.I_gid))
		dot.WriteString(fmt.Sprintf("<tr><td><b>SIZE</b></td><td>%d</td></tr>\n", inode.I_size))
		dot.WriteString(fmt.Sprintf("<tr><td><b>TYPE</b></td><td>%s</td></tr>\n", string(inode.I_type[:])))
		for b := 0; b < 15; b++ {
			dot.WriteString(fmt.Sprintf("<tr><td><b>I_BLOCK[%d]</b></td><td>%d</td></tr>\n", b, inode.I_block[b]))
		}
		dot.WriteString("</table>>];\n")

		for j := 0; j < 15; j++ {
			ptr := inode.I_block[j]
			if ptr == -1 {
				continue
			}
			blockName := fmt.Sprintf("block%d", ptr)
			blockOffset := int64(sb.S_block_start) + int64(ptr)*int64(binary.Size(structs.Folderblock{}))

			// Intentar Folderblock
			var folder structs.Folderblock
			if err := utils.ReadObject(file, &folder, blockOffset); err == nil {
				dot.WriteString(fmt.Sprintf("%s [label=<\n<table border='1' cellborder='1' cellspacing='0'>\n", blockName))
				dot.WriteString(fmt.Sprintf("<tr><td colspan='3'><b>FolderBlock %d</b></td></tr>\n", ptr))
				dot.WriteString("<tr><td><b>Index</b></td><td><b>Name</b></td><td><b>Inode Ptr</b></td></tr>\n")
				for k, content := range folder.B_content {
					name := strings.Trim(string(content.B_name[:]), "\x00 ")
					if name != "" {
						dot.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%d</td></tr>\n", k, name, content.B_inodo))
					}
				}
				dot.WriteString("</table>>];\n")
				dot.WriteString(fmt.Sprintf("inode%d -> %s;\n", i, blockName))
				continue
			}

			// Intentar Fileblock
			var fileblock structs.Fileblock
			if err := utils.ReadObject(file, &fileblock, blockOffset); err == nil {
				content := strings.Trim(string(fileblock.B_content[:]), "\x00")
				dot.WriteString(fmt.Sprintf("%s [label=<\n<table border='1' cellborder='1' cellspacing='0'>\n", blockName))
				dot.WriteString(fmt.Sprintf("<tr><td><b>FileBlock %d</b></td></tr><tr><td>%s</td></tr>\n", ptr, content))
				dot.WriteString("</table>>];\n")
				dot.WriteString(fmt.Sprintf("inode%d -> %s;\n", i, blockName))
			}
		}
	}
	dot.WriteString("}\n")

	dotPath := "tree.dot"
	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
		fmt.Println("ERROR: No se pudo guardar el archivo .dot.")
		return
	}

	cmd := exec.Command("dot", "-Tjpg", dotPath, "-o", path)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte TREE generado exitosamente:", path)
}

func GenerarReporteBlock(pathSalida string, id string) {
	if !global.CurrentUser.Status {
		fmt.Println("ERROR: No hay sesión activa.")
		return
	}

	driveLetter := string(id[0])
	filePath := "./test/" + strings.ToUpper(driveLetter) + ".bin"
	file, err := utils.OpenFile(filePath)
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo binario.")
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("ERROR: No se pudo leer el MBR.")
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice fuera de rango.")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(mbr.Partitions[index].Start)); err != nil {
		fmt.Println("ERROR: No se pudo leer el Superbloque.")
		return
	}

	var dot strings.Builder
	dot.WriteString("digraph BlockReport {\n")
	dot.WriteString("node [shape=plaintext fontname=\"Courier\"];\n")

	inodo := structs.Inode{}
	for i := int32(0); i < sb.S_inodes_count; i++ {
		offset := int64(sb.S_inode_start) + int64(i)*int64(binary.Size(inodo))
		if err := utils.ReadObject(file, &inodo, offset); err != nil || (inodo.I_uid == 0 && inodo.I_gid == 0) {
			continue
		}

		for _, blockPtr := range inodo.I_block {
			if blockPtr == -1 {
				continue
			}
			bloqueOffset := int64(sb.S_block_start) + int64(blockPtr)*int64(binary.Size(structs.Folderblock{}))

			// Intentar leer como carpeta
			var folder structs.Folderblock
			if err := utils.ReadObject(file, &folder, bloqueOffset); err == nil {
				dot.WriteString(fmt.Sprintf("blk_folder%d [label=<\n<table border='1' cellborder='1'>\n<tr><td colspan='2'>Bloque Carpeta %d</td></tr>\n<tr><td><b>b_name</b></td><td><b>b_inodo</b></td></tr>\n", blockPtr, blockPtr))
				for _, entry := range folder.B_content {
					name := strings.Trim(string(entry.B_name[:]), "\x00 ")
					if name != "" {
						dot.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>\n", name, entry.B_inodo))
					}
				}
				dot.WriteString("</table>>];\n")
				continue
			}

			// Intentar leer como archivo
			var fblock structs.Fileblock
			if err := utils.ReadObject(file, &fblock, bloqueOffset); err == nil {
				cont := strings.Trim(string(fblock.B_content[:]), "\x00")
				dot.WriteString(fmt.Sprintf("blk_file%d [label=<\n<table border='1' cellborder='1'>\n<tr><td>Bloque Archivo %d</td></tr><tr><td>%s</td></tr></table>>];\n", blockPtr, blockPtr, cont))
				continue
			}
		}
	}
	dot.WriteString("}\n")

	if err := os.WriteFile("block.dot", []byte(dot.String()), 0644); err != nil {
		fmt.Println("ERROR: No se pudo guardar el archivo .dot.")
		return
	}

	cmd := exec.Command("dot", "-Tjpg", "block.dot", "-o", pathSalida)
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte Block generado exitosamente:", pathSalida)
}

func GenerarReporteSB(path string, id string) {
	driveLetter := string(id[0])
	filepath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

	file, err := utils.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index >= 4 {
		fmt.Println("Índice de partición inválido.")
		return
	}

	start := mbr.Partitions[index].Start
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(start)); err != nil {
		fmt.Println("Error al leer Superbloque:", err)
		return
	}

	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("node [shape=plaintext fontname=\"Helvetica\"];\n")
	dot.WriteString("sb [label=<\n")
	dot.WriteString("<table border='1' cellborder='1' cellspacing='0'>\n")
	dot.WriteString("<tr><td colspan='2'><b>Superbloque</b></td></tr>\n")

	// Añadir filas
	dot.WriteString(fmt.Sprintf("<tr><td>S_filesystem_type</td><td>%d</td></tr>\n", sb.S_filesystem_type))
	dot.WriteString(fmt.Sprintf("<tr><td>S_inodes_count</td><td>%d</td></tr>\n", sb.S_inodes_count))
	dot.WriteString(fmt.Sprintf("<tr><td>S_blocks_count</td><td>%d</td></tr>\n", sb.S_blocks_count))
	dot.WriteString(fmt.Sprintf("<tr><td>S_free_blocks_count</td><td>%d</td></tr>\n", sb.S_free_blocks_count))
	dot.WriteString(fmt.Sprintf("<tr><td>S_free_inodes_count</td><td>%d</td></tr>\n", sb.S_free_inodes_count))
	dot.WriteString(fmt.Sprintf("<tr><td>S_mtime</td><td>%s</td></tr>\n", strings.Trim(string(sb.S_mtime[:]), "\x00")))
	dot.WriteString(fmt.Sprintf("<tr><td>S_umtime</td><td>%s</td></tr>\n", strings.Trim(string(sb.S_umtime[:]), "\x00")))
	dot.WriteString(fmt.Sprintf("<tr><td>S_mnt_count</td><td>%d</td></tr>\n", sb.S_mnt_count))
	dot.WriteString(fmt.Sprintf("<tr><td>S_magic</td><td>%d</td></tr>\n", sb.S_magic))
	dot.WriteString(fmt.Sprintf("<tr><td>S_inode_size</td><td>%d</td></tr>\n", sb.S_inode_size))
	dot.WriteString(fmt.Sprintf("<tr><td>S_block_size</td><td>%d</td></tr>\n", sb.S_block_size))
	dot.WriteString(fmt.Sprintf("<tr><td>S_fist_ino</td><td>%d</td></tr>\n", sb.S_fist_ino))
	dot.WriteString(fmt.Sprintf("<tr><td>S_first_blo</td><td>%d</td></tr>\n", sb.S_first_blo))
	dot.WriteString(fmt.Sprintf("<tr><td>S_bm_inode_start</td><td>%d</td></tr>\n", sb.S_bm_inode_start))
	dot.WriteString(fmt.Sprintf("<tr><td>S_bm_block_start</td><td>%d</td></tr>\n", sb.S_bm_block_start))
	dot.WriteString(fmt.Sprintf("<tr><td>S_inode_start</td><td>%d</td></tr>\n", sb.S_inode_start))
	dot.WriteString(fmt.Sprintf("<tr><td>S_block_start</td><td>%d</td></tr>\n", sb.S_block_start))

	dot.WriteString("</table>>];\n")
	dot.WriteString("}\n")

	dotPath := "sb.dot"
	if err := os.WriteFile(dotPath, []byte(dot.String()), 0644); err != nil {
		fmt.Println("Error al escribir archivo DOT:", err)
		return
	}

	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", path)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error ejecutando Graphviz:", err)
		return
	}

	fmt.Println("Reporte SB generado exitosamente:", path)
}

func GenerarReporteFile(id string, path string, ruta string) {
	fmt.Println("======Start REP FILE======")
	fmt.Println("ID:", id)
	fmt.Println("Path:", path)
	fmt.Println("Ruta:", ruta)

	driveLetter := string(id[0])
	binPath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

	// Abrir archivo binario
	file, err := utils.OpenFile(binPath)
	if err != nil {
		fmt.Println("Error al abrir disco:", err)
		return
	}
	defer file.Close()

	// Leer MBR
	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error leyendo el MBR:", err)
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice de partición inválido.")
		return
	}

	part := mbr.Partitions[index]
	if part.Size == 0 {
		fmt.Println("ERROR: La partición no existe.")
		return
	}

	// Leer Superbloque
	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		fmt.Println("Error leyendo el superbloque:", err)
		return
	}

	// Buscar el inodo del archivo
	inodeIndex := utils_inodes.InitSearch(ruta, file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró el archivo:", ruta)
		return
	}

	inodeOffset := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(structs.Inode{}))
	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("Error leyendo el inodo:", err)
		return
	}

	// Obtener contenido
	contenido := utils_inodes.GetInodeFileData(inode, file, sb)
	if contenido == "" {
		fmt.Println("Advertencia: archivo vacío o error de lectura.")
	}
	fmt.Println("Contenido del archivo:", contenido)
	// Guardar el contenido en un archivo de salida
	if err := os.WriteFile(path, []byte(contenido), 0644); err != nil {
		fmt.Println("Error escribiendo el archivo de salida:", err)
		return
	}

	fmt.Println("Reporte FILE generado correctamente en:", path)
}

func GenerarReporteLS(id string, path string, ruta string) {
	fmt.Println("======Start REP LS======")

	driveLetter := string(id[0])
	binPath := "./test/" + strings.ToUpper(driveLetter) + ".bin"

	file, err := utils.OpenFile(binPath)
	if err != nil {
		fmt.Println("Error al abrir disco:", err)
		return
	}
	defer file.Close()

	var mbr structs.MRB
	if err := utils.ReadObject(file, &mbr, 0); err != nil {
		fmt.Println("Error leyendo MBR:", err)
		return
	}

	index := int(id[1] - '1')
	if index < 0 || index > 3 {
		fmt.Println("ERROR: Índice de partición inválido.")
		return
	}

	part := mbr.Partitions[index]
	if part.Size == 0 {
		fmt.Println("ERROR: La partición no existe.")
		return
	}

	var sb structs.Superblock
	if err := utils.ReadObject(file, &sb, int64(part.Start)); err != nil {
		fmt.Println("Error leyendo superbloque:", err)
		return
	}

	inodeIndex := utils_inodes.InitSearch(ruta, file, sb)
	if inodeIndex == -1 {
		fmt.Println("ERROR: No se encontró la ruta:", ruta)
		return
	}

	inodeOffset := int64(sb.S_inode_start) + int64(inodeIndex)*int64(binary.Size(structs.Inode{}))
	var inode structs.Inode
	if err := utils.ReadObject(file, &inode, inodeOffset); err != nil {
		fmt.Println("Error leyendo inodo:", err)
		return
	}

	var builder strings.Builder
	builder.WriteString("Permisos\tOwner\tGrupo\tSize(Bytes)\tFecha\t\tHora\t\tTipo\t\tName\n")
	builder.WriteString("------------------------------------------------------------------------------------------\n")

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

			childOffset := int64(sb.S_inode_start) + int64(entry.B_inodo)*int64(binary.Size(structs.Inode{}))
			var childInode structs.Inode
			if err := utils.ReadObject(file, &childInode, childOffset); err != nil {
				continue
			}

			tipo := "Archivo"
			if childInode.I_type[0] == '0' {
				tipo = "Carpeta"
			}

			permisos := string(childInode.I_perm[:])
			fecha := strings.Trim(string(childInode.I_mtime[:10]), "\x00")
			hora := strings.Trim(string(childInode.I_mtime[11:]), "\x00")

			builder.WriteString(fmt.Sprintf("%s\tUser%d\tGrupo%d\t%d\t\t%s\t%s\t%s\t%s\n",
				permisos,
				childInode.I_uid,
				childInode.I_gid,
				childInode.I_size,
				fecha,
				hora,
				tipo,
				name))
		}
	}

	// Guardar en archivo de texto
	if err := os.WriteFile(path, []byte(builder.String()), 0644); err != nil {
		fmt.Println("Error escribiendo el archivo de salida:", err)
		return
	}

	fmt.Println("✅ Reporte LS generado como archivo de texto en:", path)
}
