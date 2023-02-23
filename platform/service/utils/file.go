package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"github.com/obgnail/plugin-platform/common/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func ReadFile(filePath string) ([]byte, error) {
	f, _ := os.Open(filePath)
	defer f.Close()
	v, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return v, nil
}

func GetYamlFromFile(fileHeader *multipart.FileHeader, file multipart.File) (content string, err error) {
	fileBuf := make([]byte, fileHeader.Size)
	_, err = file.Read(fileBuf)
	if err != nil {
		return "", errors.Trace(err)
	}

	fType := magicNumber(fileBuf, 0)
	if fType == "" {
		return "", errors.PluginUploadError(errors.FileMalformed)
	}

	reader := bytes.NewReader(fileBuf)
	switch fType {
	case FileTypeZip:
		zr, err := zip.NewReader(reader, fileHeader.Size)
		if err != nil {
			return "", errors.Trace(err)
		}

		for _, f := range zr.File {
			if IsPluginYamlPath(f.Name) {
				rc, err := f.Open()
				if err != nil {
					return "", errors.Trace(err)
				}
				defer rc.Close()
				v, err := ioutil.ReadAll(rc)
				if err != nil {
					return "", errors.Trace(err)
				}
				result := string(v)
				return result, nil
			}
		}
	case FileTypeGzip:
		gr, err := gzip.NewReader(reader)
		if err != nil {
			return "", errors.Trace(err)
		}
		defer gr.Close()

		tr := tar.NewReader(gr)
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", errors.Trace(err)
			}

			if IsPluginYamlPath(h.Name) {
				result, err := ioutil.ReadAll(tr)
				if err != nil {
					return "", errors.Trace(err)
				}
				return string(result), nil
			}
		}
	}

	return "", nil
}

const (
	FileTypeTar  = "tar"
	FileTypeZip  = "zip"
	FileTypeGzip = "gzip"
	FileTypeBzip = "bzip"
	FileTypeXz   = "xz"
)

var (
	magicZIP  = []byte{0x50, 0x4b, 0x03, 0x04}
	magicGZ   = []byte{0x1f, 0x8b}
	magicBZIP = []byte{0x42, 0x5a}
	magicTAR  = []byte{0x75, 0x73, 0x74, 0x61, 0x72} // at offset 257
	magicXZ   = []byte{0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00}
)

func magicNumber(headerBytes []byte, offset int) string {
	magic := headerBytes[offset : offset+6]
	if bytes.Equal(magicTAR, magic[0:5]) {
		return FileTypeTar
	}
	if bytes.Equal(magicZIP, magic[0:4]) {
		return FileTypeZip
	}
	if bytes.Equal(magicGZ, magic[0:2]) {
		return FileTypeGzip
	}
	if bytes.Equal(magicBZIP, magic[0:2]) {
		return FileTypeBzip
	}
	if bytes.Equal(magicXZ, magic) {
		return FileTypeXz
	}
	return ""
}

func SaveDecompressedFiles(fileHeader *multipart.FileHeader, path string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return errors.Trace(err)
	}
	defer src.Close()
	reader, err := zip.NewReader(src, fileHeader.Size)
	if err != nil {
		return errors.Trace(err)
	}
	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		dst := filepath.Join(path, f.Name)
		if err := saveFile(f, dst); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func saveFile(file *zip.File, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	if filepath.Ext(dst) == "" {
		if err := outFile.Chmod(os.ModePerm); err != nil {
			return err
		}
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, rc)
	return err
}
