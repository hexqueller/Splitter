# `splitter`
Пакет `splitter` предоставляет функции для разбиения файлов на части, объединения частей в исходный файл и удаления файлов или их частей.

---

## Функции

### `SplitFileByParts`
```go
func SplitFileByParts(filePath string, numParts int64, size int64, fileName string) error
```
Разбивает файл на заданное количество частей.

**Входные данные:**
- `filePath` (string): Путь к исходному файлу.
- `numParts` (int64): Количество частей, на которое нужно разбить файл.
- `size` (int64): Размер исходного файла в байтах.
- `fileName` (string): Имя файла для частей.

**Выходные данные:**
- (error): Ошибка, если произошёл сбой при разбиении файла.

**Описание работы:**
1. Открывает исходный файл для чтения.
2. Рассчитывает размер каждой части и остаток.
3. Создаёт части, записывая соответствующий объём данных в каждый файл.
4. Именует части с использованием пакета `base64names`.

---

### `MergeFileParts`
```go
func MergeFileParts(sortedFilePaths []string, outputFilePath string) error
```
Объединяет части файлов в один исходный файл.

**Входные данные:**
- `sortedFilePaths` ([]string): Список путей к частям, отсортированный по порядку частей.
- `outputFilePath` (string): Путь для сохранения объединённого файла.

**Выходные данные:**
- (error): Ошибка, если произошёл сбой при объединении частей.

**Описание работы:**
1. Создаёт выходной файл для записи.
2. Читает содержимое каждой части файла и записывает его в выходной файл.
3. Закрывает файлы после обработки.

---

### `DeleteFile`
```go
func DeleteFile(filePath string)
```
Удаляет файл по указанному пути.

**Входные данные:**
- `filePath` (string): Путь к файлу для удаления.

**Выходные данные:**
- Нет. Паникует в случае ошибки.

---

### `DeleteFileArray`
```go
func DeleteFileArray(parts []string)
```
Удаляет список файлов по указанным путям.

**Входные данные:**
- `parts` ([]string): Список путей к файлам для удаления.

**Выходные данные:**
- Нет. Паникует в случае ошибки.

**Описание работы:**
1. Проходит по списку путей.
2. Вызывает функцию `DeleteFile` для каждого пути.
