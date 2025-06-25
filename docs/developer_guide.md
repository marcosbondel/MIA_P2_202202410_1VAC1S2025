# 📘 Manual Técnico - Sistema de Archivos Web

## 📌  Objetivos
### Objetivo General
- Describir detalladamente la implementación, arquitectura, estructuras de datos y lógica de funcionamiento del simulador de sistema de archivos EXT2, facilitando su comprensión, uso y mantenimiento futuro en el contexto del curso de Manejo e Implementación de Archivos.
### Objetivos Específicos
- Explicar en detalle las estructuras de datos clave que simulan los componentes de un disco y un sistema de archivos EXT2 (MBR, Particiones, EBR, SuperBloque, Inodos, Bloques de Datos/Carpetas/Punteros, Bitmaps).
- Presentar la arquitectura Cliente-Servidor utilizada (Frontend React.js y Backend Go).
- Detallar el proceso de serialización y deserialización de estas estructuras hacia/desde un archivo binario que representa el disco virtual.

## 🌐 Alcances del Sistema
Este sistema simula la creación y gestión de discos virtuales y sistemas de archivos basados en la estructura EXT2. El alcance funcional incluye:

### Gestión de Discos
- Creación de discos virtuales (mkdisk) con tamaño y ajuste especificados (BF, FF, WF).
- Eliminación de discos virtuales (rmdisk).

### Gestión de Particiones
- Creación de particiones Primarias (fdisk -type=P).
- Creación de partición Extendida (fdisk -type=E) (máximo una por disco).
- Creación de particiones Lógicas (fdisk -type=L) dentro de la Extendida, utilizando EBRs encadenados.

### Montaje
- Montaje de particiones Primarias y Lógicas (mount) asignando un ID único.
- Listado de particiones montadas (mounted).

### Formateo
- Creación de un sistema de archivos EXT2 (mkfs) en una partición montada, incluyendo:
  - Cálculo de número de inodos y bloques (n).
  - Escritura del SuperBloque.
  - Inicialización de Bitmaps de Inodos y Bloques.
  - Creación del Inodo raíz (/).
  - Creación del archivo /users.txt con el usuario root inicial.


### Gestión de Usuarios y Grupos 
- Inicio de sesión (login) validando contra /users.txt.
- Cierre de sesión (logout).
- Creación de grupos (mkgrp).
- Eliminación de grupos (rmgrp) (excepto root).
- Creación de usuarios (mkusr) asignados a un grupo existente.
- Eliminación de usuarios (rmusr) (excepto root).
- Cambio de grupo para un usuario (chgrp).

### Gestión de Archivos y Directorios
- Creación de directorios (mkdir), incluyendo creación recursiva de padres (-p).
- Creación de archivos (mkfile), con contenido opcional desde tamaño (-size) o archivo local (-cont), y creación recursiva de padres (-r). Soporta indirección simple y doble (triple pendiente).
- Visualización de contenido de archivos (cat).

### Generación de Reportes
- Generación de reportes gráficos (rep) usando Graphviz sobre: MBR (mbr), Disco (disk), SuperBloque (sb), Bitmaps (bm_inode, bm_block), Tabla de Inodos (inode), Bloques Usados (block), Árbol de Directorios/Archivos (tree), Contenido de Archivo (file), Listado tipo ls -l (ls).


## 🔄 Especificaciones técnicas
### Requisitos de Hardware
- **Memoria RAM:** 2GB (Recomendado 4GB+ para ejecución fluida, especialmente con discos grandes).
- **Espacio en Disco:** 1GB libre (para el código fuente, Go, Node.js, y los discos virtuales generados).

- **Procesador:** 1GHz x64 o superior.
- Pantalla
- Teclado
- Mouse (opcional)
### Requisitos de Software 
- **Sistema Operativo:** Compatible con Go y Node.js (Linux [Mint/Ubuntu recomendado], macOS, Windows).
- **Go:** Versión 1.18 o superior
- **Node.js:** Versión LTS recomendada (Verificar con node -v). Incluye npm.
- **Vue CLI:** (Si se usa para gestionar el frontend Vue)
- **Graphviz:** Necesario para generar los reportes gráficos. Debe estar instalado y el comando dot accesible desde el PATH del sistema (Verificar con dot -V).
- **IDE/Editor:** Visual Studio Code (recomendado) con extensiones para Go y React.js, u otro editor/IDE de preferencia.
- **Terminal/Consola:** Para compilar y ejecutar el backend/frontend.
---

## 📐 Arquitectura del Sistema

### Estructura General

El sistema está compuesto por dos grandes módulos:

- **Frontend (React)**: Interfaz de usuario que simula un explorador de archivos.
- **Backend (Go)**: Simula un sistema de archivos EXT3, controlando discos, particiones, archivos, carpetas, permisos, usuarios y reportes.

### Comunicación

La comunicación entre frontend y backend se realiza mediante HTTP utilizando `fetch` o `axios` desde React hacia un API REST creada con `mux` en Go.

### Despliegue en AWS

- **Frontend**: Deploy estático en un bucket de S3 configurado como sitio web público.
- **Backend**: Instancia EC2 con Go instalado y el backend ejecutándose como servicio.


---

## 🧱 Estructuras de Datos en Go

### MBR
```go
type MRB struct {
    MbrSize     int32
    CreationDate [16]byte
    Signature     int32
    Fit           [2]byte
    Partitions    [4]Partition
}
```

### Partition
```go
type Partition struct {
    Status [1]byte
    Type   [1]byte
    Fit    [2]byte
    Start  int32
    Size   int32
    Name   [16]byte
    Id     [16]byte
}
```

### Superblock
```go
type Superblock struct {
    S_inodes_count int32
    S_blocks_count int32
    S_free_blocks_count int32
    S_free_inodes_count int32
    S_mtime [16]byte
    S_umtime [16]byte
    S_mnt_count int32
    S_magic int32
    S_inode_size int32
    S_block_size int32
    S_fist_ino int32
    S_first_blo int32
    S_bm_inode_start int32
    S_bm_block_start int32
    S_inode_start int32
    S_block_start int32
    S_filesystem_type int32
}
```

### Inodo y Bloques
```go
type Inode struct {
    I_uid int32
    I_gid int32
    I_size int32
    I_atime [16]byte
    I_ctime [16]byte
    I_mtime [16]byte
    I_block [15]int32
    I_type [1]byte
    I_perm [4]byte
}

type Fileblock struct {
    B_content [64]byte
}

type Folderblock struct {
    B_content [4]Content
}

type Content struct {
    B_name [12]byte
    B_inodo int32
}
```


---

## 🔁 Endpoints REST API (Go + Mux)

### Autenticación
```go
POST /api/auth/login
{
  "user": "root",
  "pass": "123",
  "id": "A110"
}
```

### Discos
```go
GET /api/disks
GET /api/disks/{driveLetter}/partitions
```

### Exploración de archivos
```go
GET /api/fs?path=/home/user
GET /api/fs/content?path=/home/user/docs/file.txt
```

### Subida de archivos SDAA
```go
POST /api/files/upload
Content-Type: multipart/form-data
```


---

## 🧰 Funciones importantes en Go

### Login
```go
func Login(user, pass, id string) {
  // Abre el archivo binario, busca el inodo de /users.txt
  // y verifica credenciales para habilitar sesión
}
```

### Mkfile
```go
func Mkfile(path string, size int, r bool, cont string) {
  // Crea un archivo en una ruta especificada simulando bloques EXT3
}
```

### Rmusr
```go
func Rmusr(user string) {
  // Lógica para marcar lógicamente a un usuario con '0' en users.txt
}
```

### UpdateInodeFileData
```go
func UpdateInodeFileData(...) error {
  // Sobrescribe bloques e inodo con contenido nuevo
}
```

### Generación de reportes
```go
func GenerarReporteMBR(path string, id string) {}
func GenerarReporteBMBlock(path string, id string) {}
func GenerarReporteFile(id, path, ruta string) {}
```


---

## 💻 Código en React

### Login
```jsx
const handleLogin = async () => {
  const res = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user, pass, id })
  });
  const data = await res.json();
  console.log(data);
};
```

### Explorador de archivos
```jsx
useEffect(() => {
  fetch(`/api/fs?path=${currentPath}`)
    .then(res => res.json())
    .then(setFiles);
}, [currentPath]);
```

### Subir archivos
```jsx
const handleUpload = async (file) => {
  const formData = new FormData();
  formData.append('file', file);
  await fetch('/api/files/upload', {
    method: 'POST',
    body: formData
  });
};
```


---

## 📂 Comandos soportados

- `mkdisk`, `rmdisk`
- `fdisk`, `mount`, `unmount`
- `mkfs`, `login`, `logout`
- `mkgrp`, `rmgrp`, `mkusr`, `rmusr`
- `mkfile`, `cat`, `mkdir`, `find`, `rep`

Todos los comandos también pueden ser representados mediante peticiones HTTP al API, y su salida puede ser transformada a visualizaciones usando Graphviz u otros formatos.


---

## ☁️ Consideraciones de Despliegue

- El backend fue desarrollado en Go 1.21+ y se despliega en EC2 con Ubuntu 22.04.
- Los discos `.bin` se almacenan en la ruta `/home/ubuntu/app/fs/test/`.
- El frontend se construye con `npm run build` y se sube a un bucket S3 configurado como sitio web estático.
- La comunicación se habilita con CORS abiertos para permitir el acceso desde el bucket S3.


---

## 📎 Observaciones Finales

- Este sistema simula de forma realista un sistema EXT3 con control de bloques, bitmaps, árboles de carpetas, usuarios, grupos y permisos.
- El uso de estructuras como el MBR, Superbloque, inodos y bloques permite practicar la gestión de archivos como en un sistema operativo real.
- El modelo cliente-servidor facilita la integración con interfaces gráficas modernas como React.
