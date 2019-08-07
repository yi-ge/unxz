package unxz

import (
	"archive/tar"
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ulikunitz/xz"
)

var (
	linkFiles     = []linkFile{}
	errHeaderType = errors.New("unxz: invalid tar header type")
)

// Unxz - Unxz struct.
type Unxz struct {
	Src  string
	Dest string
}

type linkFile struct {
	filePath   string
	targetPath string
}

// New - Create a new Unxz.
func New(src string, dest string) Unxz {
	return Unxz{src, dest}
}

// Extract - Extract *.tar.xz package.
func (uz Unxz) Extract() error {
	destPath := uz.Dest

	xzFile, err := os.Open(uz.Src)
	if err != nil {
		return err
	}
	defer xzFile.Close()

	r := bufio.NewReader(xzFile)
	xr, err := xz.NewReader(r)

	if err != nil {
		return err
	}

	tr := tar.NewReader(xr)

	os.Mkdir(destPath, 0755)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		err = untarFile(tr, hdr, destPath, true)
		if err != nil {
			return err
		}
	}

	for _, item := range linkFiles {
		err = writeLink(item.filePath, item.targetPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeFile(filePath string, in io.Reader, fm os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = out.Chmod(fm)
	if err != nil && runtime.GOOS != "windows" {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func writeSymbolicLink(filePath string, targetPath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	err = os.Symlink(targetPath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func writeLink(filePath string, targetPath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	err = os.Link(targetPath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func untarFile(tr *tar.Reader, header *tar.Header, destPath string, stripInnerFolder bool) error {
	fileName := header.Name
	filePath := filepath.Join(destPath, fileName)

	if stripInnerFolder {
		slashIndex := strings.Index(fileName, "/")
		if slashIndex != -1 {
			fileName = fileName[slashIndex+1:]
		}
	}

	switch header.Typeflag {
	case tar.TypeDir:
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			return err
		}

		return nil
	case tar.TypeReg, tar.TypeRegA:
		return writeFile(filePath, tr, header.FileInfo().Mode())
	case tar.TypeLink:
		linkFiles = append(linkFiles, linkFile{
			filePath,
			filepath.Join(destPath, header.Linkname),
		})
		return nil
	case tar.TypeSymlink:
		return writeSymbolicLink(filePath, header.Linkname)
	case tar.TypeXGlobalHeader:
		return nil
	default:
		return errHeaderType
	}
}
