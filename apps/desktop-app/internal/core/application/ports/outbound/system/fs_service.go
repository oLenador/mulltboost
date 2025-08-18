package outbound

import "context"

type FileSystemService interface {
	CreateFile(ctx context.Context, path string, content []byte) error
	ReadFile(ctx context.Context, path string) ([]byte, error)
	DeleteFile(ctx context.Context, path string) error
	CopyFile(ctx context.Context, src, dst string) error
	MoveFile(ctx context.Context, src, dst string) error
	
	// Operações de diretório
	CreateDirectory(ctx context.Context, path string) error
	DeleteDirectory(ctx context.Context, path string) error
	ListDirectory(ctx context.Context, path string) ([]string, error)
	
	// Permissões
	SetFilePermissions(ctx context.Context, path string, permissions uint32) error
	GetFilePermissions(ctx context.Context, path string) (uint32, error)
	
	// Informações
	GetFileSize(ctx context.Context, path string) (int64, error)
	FileExists(ctx context.Context, path string) (bool, error)
}