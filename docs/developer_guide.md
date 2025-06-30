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

  - Generación de reportes gráficos (rep) usando Graphviz sobre: MBR (mbr), Disco (disk), SuperBloque (sb), Bitmaps (bm\_inode, bm\_block), Tabla de Inodos (inode), Bloques Usados (block), Árbol de Directorios/Archivos (tree), Contenido de Archivo (file), Listado tipo ls -l (ls).

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

-----

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

-----

## 🧱 Estructuras de Datos en Go

A continuación, se documentan detalladamente las estructuras de datos clave utilizadas para simular el sistema de archivos, incluyendo su propósito y los campos que las componen.

### MBR

La estructura **`MRB`** (Master Boot Record) representa el primer sector de un disco virtual. Contiene información vital para el arranque y la gestión de particiones.

```go
type MRB struct {
    MbrSize      int32        // Tamaño total del disco virtual en bytes.
    CreationDate [10]byte     // Fecha y hora de creación del disco.
    Signature    int32        // Número de firma único para el MBR.
    Fit          [1]byte      // Indica la estrategia de ajuste (ajuste en 'B'est, 'F'irst, o 'W'orst).
    Partitions   [4]Partition // Arreglo de 4 particiones, que representan las particiones primarias en el disco.
}
```

### EBR

La estructura **`EBR`** (Extended Boot Record) es utilizada para gestionar particiones lógicas dentro de una partición extendida. Cada partición lógica tiene su propio EBR.

```go
type EBR struct {
    PartStatus int32     // Estado de la partición (ej. 1 para activa, 0 para inactiva).
    PartFit    [1]byte   // Estrategia de ajuste para la partición lógica.
    PartStart  int32     // Posición inicial de la partición en el disco (en bytes).
    PartSize   int32     // Tamaño de la partición en bytes.
    PartNext   int32     // Puntero al siguiente EBR en la lista enlazada de particiones lógicas. Si es 0, no hay más.
    PartName   [16]byte  // Nombre de la partición.
}
```

### Partition

La estructura **`Partition`** define una partición individual dentro del MBR.

```go
type Partition struct {
    Status      [1]byte  // Estado de la partición ('1' para activa, '0' para inactiva).
    Type        [1]byte  // Tipo de partición ('P'rimaria, 'E'xtendida o 'L'ógica).
    Fit         [1]byte  // Estrategia de ajuste para la partición.
    Start       int32    // Posición de inicio de la partición en el disco (en bytes).
    Size        int32    // Tamaño de la partición en bytes.
    Name        [16]byte // Nombre de la partición.
    Correlative int32    // Número correlativo de la partición.
    Id          [4]byte  // Identificador único de la partición.
}
```

### Superblock

El **`Superblock`** contiene los metadatos más importantes del sistema de archivos, como el tamaño de la partición, el conteo de inodos y bloques, y los punteros a las tablas de inodos y bloques.

```go
type Superblock struct {
    S_filesystem_type   int32      // Tipo de sistema de archivos (ej. 2 para EXT2).
    S_inodes_count      int32      // Número total de inodos en el sistema de archivos.
    S_blocks_count      int32      // Número total de bloques de datos en el sistema de archivos.
    S_free_blocks_count int32      // Número de bloques de datos libres.
    S_free_inodes_count int32      // Número de inodos libres.
    S_mtime             [17]byte   // Última fecha y hora de montaje.
    S_umtime            [17]byte   // Última fecha y hora de desmontaje.
    S_mnt_count         int32      // Conteo de montajes.
    S_magic             int32      // Número mágico para identificar el sistema de archivos (0xEF53 para EXT2).
    S_inode_size        int32      // Tamaño de un inodo en bytes.
    S_block_size        int32      // Tamaño de un bloque en bytes.
    S_fist_ino          int32      // Primer inodo libre.
    S_first_blo         int32      // Primer bloque libre.
    S_bm_inode_start    int32      // Posición de inicio del bitmap de inodos.
    S_bm_block_start    int32      // Posición de inicio del bitmap de bloques.
    S_inode_start       int32      // Posición de inicio de la tabla de inodos.
    S_block_start       int32      // Posición de inicio de la tabla de bloques.
}
```

### Inode

La estructura **`Inode`** representa un nodo de índice, que almacena los metadatos de un archivo o directorio, como permisos, tamaño, propietario y los punteros a los bloques de datos.

```go
type Inode struct {
    I_uid   int32       // ID del usuario propietario.
    I_gid   int32       // ID del grupo propietario.
    I_size  int32       // Tamaño del archivo en bytes.
    I_atime [17]byte    // Última fecha y hora de acceso.
    I_ctime [17]byte    // Fecha y hora de creación.
    I_mtime [17]byte    // Última fecha y hora de modificación.
    I_block [15]int32   // Arreglo de 15 punteros a bloques.
                        // Los primeros 12 son directos, el 13 es de indirección simple, el 14 doble, y el 15 triple (no implementado).
    I_type  [1]byte     // Tipo de inodo ('1' para archivo, '2' para directorio).
    I_perm  [3]byte     // Permisos de usuario, grupo y otros.
}
```

### Bloques de datos

El sistema de archivos utiliza diferentes tipos de bloques para almacenar la información.

#### Fileblock

Un **`Fileblock`** es un bloque que almacena el contenido de un archivo.

```go
type Fileblock struct {
    B_content [64]byte // Contenido del archivo, hasta 64 bytes.
}
```

#### Folderblock

Un **`Folderblock`** almacena las entradas de un directorio, vinculando nombres de archivos/directorios a sus inodos correspondientes.

```go
type Content struct {
    B_name  [12]byte // Nombre del archivo o directorio.
    B_inodo int32    // Número de inodo asociado.
}

type Folderblock struct {
    B_content [4]Content // Arreglo de 4 estructuras Content, cada una representando una entrada.
}
```

#### Pointerblock

Un **`Pointerblock`** es un bloque de indirección que contiene punteros a otros bloques de datos o de punteros, permitiendo almacenar archivos grandes.

```go
type Pointerblock struct {
    B_pointers [16]int32 // Arreglo de 16 punteros a otros bloques.
}
```

### Journaling

El sistema de **`Journaling`** se utiliza para registrar operaciones antes de que se escriban en el disco, permitiendo la recuperación del sistema de archivos en caso de fallos.

#### Content\_J

La estructura **`Content_J`** representa una entrada individual en el registro de transacciones (journal).

```go
type Content_J struct {
    Operation [10]byte  // Tipo de operación realizada (ej. "MKDIR", "MKFILE").
    Path      [100]byte // Ruta del archivo o directorio afectado.
    Content   [100]byte // Contenido o información adicional de la operación.
    Date      [17]byte  // Fecha y hora de la operación.
}
```

#### Journaling

La estructura **`Journaling`** representa el registro completo de transacciones.

```go
type Journaling struct {
    Size      int32        // Tamaño total del registro.
    Ultimo    int32        // Índice de la última entrada de contenido.
    Contenido [50]Content_J // Arreglo de 50 entradas de registro.
}
```

-----

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

-----

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

-----

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

-----

## 📂 Comandos soportados

  - `mkdisk`, `rmdisk`
  - `fdisk`, `mount`, `unmount`
  - `mkfs`, `login`, `logout`
  - `mkgrp`, `rmgrp`, `mkusr`, `rmusr`
  - `mkfile`, `cat`, `mkdir`, `find`, `rep`

Todos los comandos también pueden ser representados mediante peticiones HTTP al API, y su salida puede ser transformada a visualizaciones usando Graphviz u otros formatos.

-----

## ☁️ Consideraciones de Despliegue

  - El backend fue desarrollado en Go 1.21+ y se despliega en EC2 con Ubuntu 22.04.
  - Los discos `.bin` se almacenan en la ruta `/home/ubuntu/app/fs/test/`.
  - El frontend se construye con `npm run build` y se sube a un bucket S3 configurado como sitio web estático.
  - La comunicación se habilita con CORS abiertos para permitir el acceso desde el bucket S3.

-----

## 📎 Observaciones Finales

  - Este sistema simula de forma realista un sistema EXT3 con control de bloques, bitmaps, árboles de carpetas, usuarios, grupos y permisos.
  - El uso de estructuras como el MBR, Superbloque, inodos y bloques permite practicar la gestión de archivos como en un sistema operativo real.
  - El modelo cliente-servidor facilita la integración con interfaces gráficas modernas como React.