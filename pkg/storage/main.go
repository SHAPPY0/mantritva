package storage

type FileMeta struct {
	OriginalName string
	SavedName    string
	FullPath     string
	AbsolutePath string
	Extenstion   string
	DownloadURL  string
}

type Storage interface {
	Save(file *FileMeta) error
	Delete(filePath string) error
	Get(fileName string) (*FileMeta, error)
	List() ([]FileMeta, error)
}
