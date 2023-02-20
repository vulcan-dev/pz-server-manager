package pz

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	IsDir     bool      `json:"is_dir"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Size      int64     `json:"size,omitempty"`
	Files     []File    `json:"files,omitempty"` // Note: Files can also be a directory
}

type Filesystem struct {
	Root  string
	Cache map[string]File
}

func NewFilesystem(root string) *Filesystem {
	return &Filesystem{
		Root:  root,
		Cache: make(map[string]File),
	}
}

func directorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func (f *Filesystem) List(path string, ignoreCache bool) (File, error) {
	if path == "" {
		path = f.Root
	} else {
		log.Debug("Joining ", f.Root, " and ", path)
		path = filepath.Join(f.Root, path)
	}

	log.Info("Listing files in ", path)

	if !ignoreCache {
		if file, ok := f.Cache[path]; ok {
			file.Path = path
			return file, nil
		}
	}

	file, err := f.list(path)
	if err != nil {
		return File{}, err
	}

	if !ignoreCache {
		for i := range file.Files {
			iterFile := file.Files[i]
			_, ok := f.Cache[iterFile.Path]
			if !ok {
				log.Info("Caching ", iterFile.Path)

				if !iterFile.IsDir {
					f.Cache[iterFile.Path] = iterFile
				} else {
					log.Debug("Caching directory ", iterFile.Path)
					size, err := directorySize(iterFile.Path)
					if err != nil {
						log.Error(err)
					}

					log.Debug("Directory size is ", size)
					iterFile.Size = size
					f.Cache[iterFile.Path] = iterFile
				}
			}
		}
	}

	return file, nil
}

func (f *Filesystem) Get(path string) (File, error) {
	if path == "" {
		path = f.Root
	} else {
		log.Debug("Joining ", f.Root, " and ", path)
		path = filepath.Join(f.Root, path)
	}

	if file, ok := f.Cache[path]; ok {
		file.Path = path
		return file, nil
	}

	return File{}, os.ErrNotExist
}

func (f *File) Content() (string, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	bytes := make([]byte, stat.Size())
	_, err = file.Read(bytes)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (f *Filesystem) list(path string) (File, error) {
	// Check if the path exists
	_, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	file := File{
		Name:  path,
		Path:  path,
		Files: []File{},
	}

	// Recursively list files in the filesystem
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == file.Path {
			return nil
		}

		// If type is directory, repeat 👀
		if info.IsDir() {
			subFile, err := f.list(path)
			if err != nil {
				return err
			}

			file.Files = append(file.Files, subFile)
			file.IsDir = true
			size, err := directorySize(path)
			if err != nil {
				log.Error(err)
			}

			file.Size = size
			return filepath.SkipDir
		} else {
			file.Files = append(file.Files, File{
				Name:      info.Name(),
				Path:      path,
				IsDir:     false,
				CreatedAt: info.ModTime(),
				UpdatedAt: info.ModTime(),
				Size:      info.Size(),
			})
		}

		return nil
	})

	return file, err
}
