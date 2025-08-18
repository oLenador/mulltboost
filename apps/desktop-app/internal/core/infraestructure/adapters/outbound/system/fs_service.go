//go:build windows
// +build windows
package system

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	outbound "github.com/oLenador/mulltbost/internal/core/application/ports/outbound/system"
	"golang.org/x/sys/windows"
)

type WindowsFileSystemService struct {
	fileOps FileOperations
}

type FileOperations interface {
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	ReadDir(name string) ([]os.DirEntry, error)
	Stat(name string) (os.FileInfo, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
	Copy(src, dst string) error
}

type DefaultFileOperations struct{}

func (d DefaultFileOperations) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (d DefaultFileOperations) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (d DefaultFileOperations) Remove(name string) error {
	return os.Remove(name)
}

func (d DefaultFileOperations) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (d DefaultFileOperations) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (d DefaultFileOperations) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (d DefaultFileOperations) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (d DefaultFileOperations) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (d DefaultFileOperations) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (d DefaultFileOperations) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (d DefaultFileOperations) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (d DefaultFileOperations) Copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// NewWindowsFileSystemService cria uma nova instância do serviço
func NewWindowsFileSystemService() outbound.FileSystemService {
	return &WindowsFileSystemService{
		fileOps: DefaultFileOperations{},
	}
}

// NewWindowsFileSystemServiceWithOperations permite injeção de dependências para testes
func NewWindowsFileSystemServiceWithOperations(fileOps FileOperations) outbound.FileSystemService {
	return &WindowsFileSystemService{
		fileOps: fileOps,
	}
}

// CreateFile implementa a criação de arquivos
func (w *WindowsFileSystemService) CreateFile(ctx context.Context, path string, content []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Normaliza o caminho para Windows
	normalizedPath := filepath.Clean(path)
	
	// Cria diretórios pai se necessário
	dir := filepath.Dir(normalizedPath)
	if err := w.fileOps.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Cria e escreve o arquivo
	if err := w.fileOps.WriteFile(normalizedPath, content, 0644); err != nil {
		return fmt.Errorf("failed to create file %s: %w", normalizedPath, err)
	}

	return nil
}

// ReadFile implementa a leitura de arquivos
func (w *WindowsFileSystemService) ReadFile(ctx context.Context, path string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	content, err := w.fileOps.ReadFile(normalizedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", normalizedPath, err)
	}

	return content, nil
}

// DeleteFile implementa a exclusão de arquivos
func (w *WindowsFileSystemService) DeleteFile(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	if err := w.fileOps.Remove(normalizedPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", normalizedPath, err)
	}

	return nil
}

// CopyFile implementa a cópia de arquivos
func (w *WindowsFileSystemService) CopyFile(ctx context.Context, src, dst string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedSrc := filepath.Clean(src)
	normalizedDst := filepath.Clean(dst)

	// Verifica se o arquivo de origem existe
	if exists, err := w.FileExists(ctx, normalizedSrc); err != nil {
		return fmt.Errorf("failed to check source file existence: %w", err)
	} else if !exists {
		return fmt.Errorf("source file %s does not exist", normalizedSrc)
	}

	// Cria diretórios pai do destino se necessário
	dstDir := filepath.Dir(normalizedDst)
	if err := w.fileOps.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	if err := w.fileOps.Copy(normalizedSrc, normalizedDst); err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %w", normalizedSrc, normalizedDst, err)
	}

	return nil
}

// MoveFile implementa a movimentação de arquivos
func (w *WindowsFileSystemService) MoveFile(ctx context.Context, src, dst string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedSrc := filepath.Clean(src)
	normalizedDst := filepath.Clean(dst)

	// Verifica se o arquivo de origem existe
	if exists, err := w.FileExists(ctx, normalizedSrc); err != nil {
		return fmt.Errorf("failed to check source file existence: %w", err)
	} else if !exists {
		return fmt.Errorf("source file %s does not exist", normalizedSrc)
	}

	// Cria diretórios pai do destino se necessário
	dstDir := filepath.Dir(normalizedDst)
	if err := w.fileOps.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	if err := w.fileOps.Rename(normalizedSrc, normalizedDst); err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", normalizedSrc, normalizedDst, err)
	}

	return nil
}

// CreateDirectory implementa a criação de diretórios
func (w *WindowsFileSystemService) CreateDirectory(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	if err := w.fileOps.MkdirAll(normalizedPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", normalizedPath, err)
	}

	return nil
}

// DeleteDirectory implementa a exclusão de diretórios
func (w *WindowsFileSystemService) DeleteDirectory(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	if err := w.fileOps.RemoveAll(normalizedPath); err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", normalizedPath, err)
	}

	return nil
}

// ListDirectory implementa a listagem de conteúdo de diretórios
func (w *WindowsFileSystemService) ListDirectory(ctx context.Context, path string) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	entries, err := w.fileOps.ReadDir(normalizedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory %s: %w", normalizedPath, err)
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}

// SetFilePermissions implementa a configuração de permissões específicas do Windows
func (w *WindowsFileSystemService) SetFilePermissions(ctx context.Context, path string, permissions uint32) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	// Converte path para UTF-16 para APIs do Windows
	pathPtr, err := windows.UTF16PtrFromString(normalizedPath)
	if err != nil {
		return fmt.Errorf("failed to convert path to UTF-16: %w", err)
	}

	// Define as permissões usando a API SetFileAttributes do Windows
	err = setFileAttributesW(pathPtr, permissions)
	if err != nil {
		return fmt.Errorf("failed to set permissions for %s: %w", normalizedPath, err)
	}

	return nil
}

// GetFilePermissions implementa a obtenção de permissões específicas do Windows
func (w *WindowsFileSystemService) GetFilePermissions(ctx context.Context, path string) (uint32, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	// Converte path para UTF-16 para APIs do Windows
	pathPtr, err := windows.UTF16PtrFromString(normalizedPath)
	if err != nil {
		return 0, fmt.Errorf("failed to convert path to UTF-16: %w", err)
	}

	// Obtém as permissões usando a API GetFileAttributes do Windows
	attributes, err := getFileAttributesW(pathPtr)
	if err != nil {
		return 0, fmt.Errorf("failed to get permissions for %s: %w", normalizedPath, err)
	}

	return attributes, nil
}

// GetFileSize implementa a obtenção do tamanho do arquivo
func (w *WindowsFileSystemService) GetFileSize(ctx context.Context, path string) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	info, err := w.fileOps.Stat(normalizedPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info for %s: %w", normalizedPath, err)
	}

	return info.Size(), nil
}

// FileExists implementa a verificação de existência de arquivo
func (w *WindowsFileSystemService) FileExists(ctx context.Context, path string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	normalizedPath := filepath.Clean(path)
	
	_, err := w.fileOps.Stat(normalizedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence for %s: %w", normalizedPath, err)
	}

	return true, nil
}

// APIs específicas do Windows para manipulação de atributos de arquivo
var (
	kernel32           = windows.NewLazyDLL("kernel32.dll")
	procSetFileAttribW = kernel32.NewProc("SetFileAttributesW")
	procGetFileAttribW = kernel32.NewProc("GetFileAttributesW")
)



func setFileAttributesW(pathPtr *uint16, attributes uint32) error {
    err := windows.SetFileAttributes(pathPtr, attributes)
    if err != nil {
        return err
    }
    return nil
}

func getFileAttributesW(pathPtr *uint16) (uint32, error) {
    r1, err := windows.GetFileAttributes(pathPtr)
    if err != nil {
        return 0, err
    }
    return r1, nil
}

const (
	FILE_ATTRIBUTE_READONLY            = 0x1
	FILE_ATTRIBUTE_HIDDEN              = 0x2
	FILE_ATTRIBUTE_SYSTEM              = 0x4
	FILE_ATTRIBUTE_DIRECTORY           = 0x10
	FILE_ATTRIBUTE_ARCHIVE             = 0x20
	FILE_ATTRIBUTE_DEVICE              = 0x40
	FILE_ATTRIBUTE_NORMAL              = 0x80
	FILE_ATTRIBUTE_TEMPORARY           = 0x100
	FILE_ATTRIBUTE_SPARSE_FILE         = 0x200
	FILE_ATTRIBUTE_REPARSE_POINT       = 0x400
	FILE_ATTRIBUTE_COMPRESSED          = 0x800
	FILE_ATTRIBUTE_OFFLINE             = 0x1000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED = 0x2000
	FILE_ATTRIBUTE_ENCRYPTED           = 0x4000
)