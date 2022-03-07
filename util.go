package weblogic_buildpack

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsJvmApplicationPackage(appPath string) (bool, error) {
	if appXMLPresent, err := FileExists(filepath.Join(appPath, "META-INF", "application.xml")); err != nil {
		return false, fmt.Errorf("unable to check application.xml\n%w", err)
	} else if appXMLPresent {
		return true, nil
	}

	if webInfPresent, err := DirExists(filepath.Join(appPath, "WEB-INF")); err != nil {
		return false, fmt.Errorf("unable to check WEB-INF\n%w", err)
	} else if webInfPresent {
		return true, nil
	}

	return false, nil
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func DirExists(path string) (bool, error) {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func RecursiveZip(sourcePath, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(strings.TrimPrefix(filePath, filepath.Dir(sourcePath)), "/")
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func CleanupDir(sourcePath string) error {
	files, err := filepath.Glob(filepath.Join(sourcePath, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyContent is the same as fs.Copy with some small differences
// 1. the destination directory is not deleted (so can use /workspace with no errors)
// 2. does not fail if the destination directory exists already (so can use /workspace with no errors)
func CopyContent(source, destination string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = copyDirectory(source, destination)
		if err != nil {
			return err
		}
	} else {
		err = copyFile(source, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	info, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	err = os.Chmod(destination, info.Mode())
	if err != nil {
		return err
	}

	return nil
}

func copyDirectory(source, destination string) error {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		switch {
		case info.IsDir():
			err = os.Mkdir(filepath.Join(destination, path), os.ModePerm)
			if os.IsExist(err) {
				return nil
			}
			if err != nil {
				return err
			}

		case (info.Mode() & os.ModeSymlink) != 0:
			err = copyLink(source, destination, path)
			if err != nil {
				return err
			}

		default:
			err = copyFile(filepath.Join(source, path), filepath.Join(destination, path))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func copyLink(source, destination, path string) error {
	link, err := os.Readlink(filepath.Join(source, path))
	if err != nil {
		return err
	}

	err = os.Symlink(link, filepath.Join(destination, path))
	if err != nil {
		return err
	}

	return nil
}
