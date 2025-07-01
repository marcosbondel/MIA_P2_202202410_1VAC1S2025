package analyzer

import (
	"MIA_P2_202202410_1VAC1S2025/fs/disk_management"
	"MIA_P2_202202410_1VAC1S2025/fs/file_manager"
	"MIA_P2_202202410_1VAC1S2025/fs/file_system"
	"MIA_P2_202202410_1VAC1S2025/fs/reports_generator"
	"MIA_P2_202202410_1VAC1S2025/fs/user"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func Analyze() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("=== MIA Console - Marcos Bonifasi (202202410) ===")
	fmt.Println("Write commands (or 'exit' for closing the app):")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())

		// Salir del bucle
		if line == "exit" {
			fmt.Println("Saliendo del sistema...")
			break
		}

		// Ignorar comentarios o líneas vacías
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		command, params := getCommandAndParams(line)
		fmt.Printf("Command: %s, Params: %s\n", command, params)
		// AnalyzeCommnad(command, params, )
	}
}

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

// func analyzeHTTPInput(input string) (string, string) {
// 	parts := strings.Fields(input)
// 	if len(parts) > 0 {
// 		command := strings.ToLower(parts[0])
// 		params := strings.Join(parts[1:], " ")
// 		return command, params
// 	}
// 	return "", input
// }

func AnalyzeHTTPInput(input string) string {
	buffer_response := "\n"
	command, params := getCommandAndParams(input)

	AnalyzeCommnad(command, params, &buffer_response)

	return buffer_response
}

func AnalyzeCommnad(command string, params string, buffer_response *string) {
	// fmt.Printf("Command: %s, Params: %s\n", command, params)
	*buffer_response += "\n"
	*buffer_response += "\n"
	*buffer_response += "Executing command...\n"
	*buffer_response += fmt.Sprintf("Command: %s, Params: %s\n", command, params)

	switch command {
	case "mkdisk":
		fn_mkdisk(params, buffer_response)
	case "fdisk":
		fn_fdisk(params, buffer_response)
	case "rmdisk":
		fn_rmdisk(params, buffer_response)
	case "mount":
		fn_mount(params, buffer_response)
	case "unmount":
		fn_unmount(params, buffer_response)
	case "mkfs":
		fn_mkfs(params, buffer_response)
	case "login":
		fn_login(params, buffer_response)
	case "logout":
		fn_logout(buffer_response)
	case "mkgrp":
		fn_mkgrp(params, buffer_response)
	case "rmgrp":
		fn_rmgrp(params, buffer_response)
	case "mkusr":
		fn_mkusr(params, buffer_response)
	case "rmusr":
		fn_rmusr(params, buffer_response)
	case "mkfile":
		fn_mkfile(params, buffer_response)
	case "mkdir":
		fn_mkdir(params, buffer_response)
	case "find":
		fn_find(params, buffer_response)
	case "cat":
		fn_cat(params, buffer_response)
	case "pause":
		file_manager.Pause()
	case "rep":
		fn_rep(params, buffer_response)
	case "execute":
		fn_execute(params, buffer_response)
	case "exit":
		fmt.Println("Exiting the program.")
		os.Exit(0)
	default:
		fmt.Println("Error: Command not recognized.")
	}
}

func fn_mount(params string, buffer_response *string) {
	// Define flags
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	name := fs.String("driveletter", "", "Letter of the drive")
	path := fs.String("name", "", "Name of the partition")

	// get flags values
	managementFlags(fs, params)

	// Call the function
	disk_management.Mount(*name, *path, buffer_response)
}

func fn_unmount(params string, buffer_response *string) {
	// Define flags
	fs := flag.NewFlagSet("unmount", flag.ExitOnError)
	id_partition := fs.String("id", "", "Id of the partition")

	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_response += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	disk_management.Unmount(*id_partition, buffer_response)
}

func fn_fdisk(params string, buffer_response *string) {
	// Define flags
	// fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	// size := fs.Int("size", 0, "Size")
	// driveLetter := fs.String("driveletter", "", "Letra de unidad")
	// name := fs.String("name", "", "Nombre de la partición")
	// type_ := fs.String("type", "p", "Tipo de partición (p/e)")
	// fit := fs.String("fit", "", "Fit")
	// unit := fs.String("unit", "m", "Unit")
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Size")
	add_ := fs.Int("add", 0, "Aumentar o reducir tamaño de partición")
	delete := fs.String("delete", "", "Eliminar partición")
	driveLetter := fs.String("driveletter", "", "Letra del disco")
	name := fs.String("name", "", "Nombre de la partición")
	type_ := fs.String("type", "p", "Tipo de partición")
	fit := fs.String("fit", "", "Fit")
	unit := fs.String("unit", "k", "Unidad (k/m)")

	// get flags values
	managementFlags(fs, params)

	// Call the function
	disk_management.Fdisk(*size, *driveLetter, *name, *type_, *fit, *unit, *add_, *delete, buffer_response)
}

func fn_mkdisk(params string, buffer_response *string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Size")
	fit := fs.String("fit", "", "Fit")
	unit := fs.String("unit", "m", "Unit")

	// get flags values
	managementFlags(fs, params)

	// Call the function
	disk_management.Mkdisk(*size, *fit, *unit, buffer_response)

}

func managementFlags(fs *flag.FlagSet, params string) {
	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Obtener los nombres de todas las flags
	var flagNames []string
	fs.VisitAll(func(f *flag.Flag) {
		flagNames = append(flagNames, f.Name)
	})

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		if contains(flagNames, flagName) {
			fs.Set(flagName, flagValue)
		} else {
			fmt.Println("Error: Flag not found:", flagName)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func fn_rmdisk(params string, buffer_response *string) {

	// Define flags
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	driveLetter := fs.String("driveletter", "d", "Disk letter")
	// Parse the flags
	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_response += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	disk_management.Rmdisk(*driveLetter, buffer_response)
}

func fn_mkfs(input string, buffer_response *string) {
	// Define flags
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_response += "Error: Flag not found: " + flagName + "\n"
		}
	}

	// Call the function
	file_system.Mkfs(*id, *type_, *fs_, buffer_response)

}

func fn_login(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user_ := fs.String("user", "", "User")
	pass := fs.String("pass", "", "Password")
	id := fs.String("id", "", "Id")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	user.Login(*user_, *pass, *id, buffer_string)

}

func fn_logout(buffer_string *string) {
	user.Logout(buffer_string)
}

func fn_mkusr(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	grp := fs.String("grp", "", "grupo")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "grp":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Mkusr(*user, *pass, *grp, buffer_string)

}

func fn_mkgrp(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("mkgrp", flag.ExitOnError)
	grp := fs.String("name", "", "Group name")
	// Parse the flags
	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Mkgrp(*grp, buffer_string)
}

func fn_rep(params string, buffer_string *string) {
	flagSet := flag.NewFlagSet("rep", flag.ContinueOnError)
	name := flagSet.String("name", "", "Nombre del reporte")
	path := flagSet.String("path", "", "Ruta de salida del reporte")
	id := flagSet.String("id", "", "ID de la partición")
	ruta := flagSet.String("ruta", "", "Ruta de salida del reporte")

	args := strings.Fields(params)
	if err := flagSet.Parse(args); err != nil {
		fmt.Println("Error:", err)
		*buffer_string += "Error: " + err.Error() + "\n"
		return
	}

	if *name == "" || *path == "" || *id == "" {
		fmt.Println("ERROR: Faltan parámetros en rep.")
		*buffer_string += "ERROR: Faltan parámetros en rep.\n"
		return
	}

	reports_generator.GenerarReporte(*name, *path, *id, *ruta, buffer_string)
}

func fn_execute(params string, buffer_string *string) {
	flagSet := flag.NewFlagSet("execute", flag.ContinueOnError)
	path := flagSet.String("path", "", "Ruta del archivo .sdaa")

	args := strings.Fields(params)
	if err := flagSet.Parse(args); err != nil {
		fmt.Println("Error:", err)
		// buffer_string += "Error: " + err.Error() + "\n"
		return
	}

	if *path == "" {
		fmt.Println("ERROR: Falta parámetro -path")
		// buffer_string += "ERROR: Falta parámetro -path\n"
		return
	}

	file, err := os.ReadFile(*path)
	if err != nil {
		fmt.Println("ERROR: No se pudo leer el archivo:", *path)
		// buffer_string += "ERROR: No se pudo leer el archivo: " + *path + "\n"
		return
	}

	lines := strings.Split(string(file), "\n")
	fmt.Println("===Start===")
	// buffer_string += "===Start===\n"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fmt.Println("Executing Line:", line)
		// buffer_string += "Executing Line: " + line + "\n"
		command, params := getCommandAndParams(line)
		AnalyzeCommnad(command, params, buffer_string) // Ejecuta el comando
	}
	fmt.Println("===End===")
}

func fn_rmgrp(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("rmgrp", flag.ExitOnError)
	grp := fs.String("name", "", "Group	name")
	// Parse the flags
	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Rmgrp(*grp, buffer_string)
}

func fn_rmusr(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("rmusr", flag.ExitOnError)
	user_ := fs.String("user", "", "User")

	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Rmusr(*user_, buffer_string)

}

func fn_mkfile(input string, buffer_string *string) {
	// Inicializa los flags con valores por defecto
	fs := flag.NewFlagSet("mkfile", flag.ExitOnError)
	path := fs.String("path", "", "Path of the file")
	size := fs.Int("size", 0, "Size of the file")
	r := fs.Bool("r", false, "Recursive creation of folders")
	cont := fs.String("cont", "", "File content path")

	// Parsea los argumentos (por si vienen por consola)
	fs.Parse(os.Args[1:])

	// Revisa manualmente si el flag -r aparece en el input (como palabra suelta)
	if strings.Contains(strings.ToLower(input), "-r") && !strings.Contains(strings.ToLower(input), "-r=") {
		*r = true
	}

	// Extrae y asigna valores manualmente usando regex
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|[^ ]+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "path", "size", "cont":
			fs.Set(flagName, flagValue)
		// No hagas nada con -r aquí, ya lo evaluaste antes
		default:
			fmt.Printf("Warning: Flag -%s no reconocido\n", flagName)
			*buffer_string += fmt.Sprintf("Warning: Flag -%s no reconocido\n", flagName)
		}
	}

	// Llama a la función
	file_manager.Mkfile(*path, *size, *r, *cont, buffer_string)
}

func fn_cat(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("cat", flag.ExitOnError)
	path := fs.String("file1", "", "file1 of the file")
	// Parse the flags
	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "file1":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Cat(*path, buffer_string)
}

func fn_mkdir(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("mkdir", flag.ExitOnError)
	path := fs.String("path", "", "Path of the directory")
	r := fs.Bool("r", false, "Create parent directories if they do not exist")
	// Parse the flags
	fs.Parse(os.Args[1:])

	// Revisa manualmente si el flag -r aparece en el input (como palabra suelta)
	if strings.Contains(strings.ToLower(input), "-r") && !strings.Contains(strings.ToLower(input), "-r=") {
		*r = true
	}

	// Extrae y asigna valores manualmente usando regex
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|[^ ]+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path", "r":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}

	// Call the function
	file_manager.Mkdir(*path, *r, buffer_string)
}

func fn_find(input string, buffer_string *string) {
	// Define flags
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	path := fs.String("path", "", "Path of the file")
	name := fs.String("name", "", "Name of the file")
	// Parse the flags
	fs.Parse(os.Args[1:])
	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)
	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
			*buffer_string += fmt.Sprintf("Error: Flag not found: %s\n", flagName)
		}
	}
	// Call the function
	file_manager.Find(*path, *name, buffer_string)

}
