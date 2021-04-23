package file

type Provider interface {
	CreateDirectory(path string) error
	CreateFileWithContent(path, content string) error
	DeleteDirectory(path string) error
}
