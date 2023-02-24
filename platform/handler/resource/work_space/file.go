package work_space

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	. "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/file_utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SpaceOperation struct {
	AppID      string
	InstanceID string
}

func (o *SpaceOperation) CreateFile(name string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(CreateFileFailure, CreateFileFailureError.Error(), err.Error())
	}
	newFile, err := os.Create(path)
	if err != nil {
		return NewPluginError(CreateFileFailure, CreateFileFailureError.Error(), err.Error())
	}
	defer newFile.Close()
	return nil
}

func (o *SpaceOperation) MakeDir(name string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(MakeDirFailure, MakeDirFailureError.Error(), err.Error())
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return NewPluginError(MakeDirFailure, MakeDirFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) Rename(originalPath string, newPath string) PluginError {
	oldPath, err := GetFilePath(o.AppID, o.InstanceID, originalPath)
	if err != nil {
		return NewPluginError(ReNameFileFailure, ReNameFileFailureError.Error(), err.Error())
	}
	path, err := GetFilePath(o.AppID, o.InstanceID, newPath)
	if err != nil {
		return NewPluginError(ReNameFileFailure, ReNameFileFailureError.Error(), err.Error())
	}

	err = os.Rename(oldPath, path)
	if err != nil {
		return NewPluginError(ReNameFileFailure, ReNameFileFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) Remove(name string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(RemoveFileFailure, RemoveFileFailureError.Error(), err.Error())
	}
	err = os.Remove(path)
	if err != nil {
		return NewPluginError(RemoveFileFailure, RemoveFileFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) IsExist(name string) (bool, PluginError) {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return false, NewPluginError(IsExistFileFailure, IsExistFileFailureError.Error(), err.Error())
	}
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, NewPluginError(IsExistFileFailure, IsExistFileFailureError.Error(), err.Error())
	}
	return true, nil
}

func (o *SpaceOperation) IsDir(name string) (bool, PluginError) {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return false, NewPluginError(IsDirFailure, IsDirFailureError.Error(), err.Error())
	}

	s, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return s.IsDir(), nil
}

func (o *SpaceOperation) Copy(originalPath string, newPath string) PluginError {
	_originalPath, err := GetFilePath(o.AppID, o.InstanceID, originalPath)
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}
	_newPath, err := GetFilePath(o.AppID, o.InstanceID, newPath)
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}

	originalFile, err := os.Open(_originalPath)
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}
	defer originalFile.Close()

	newFile, err := os.Create(_newPath)
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}

	err = newFile.Sync()
	if err != nil {
		return NewPluginError(CopyFileFailure, CopyFileFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) Read(name string) ([]byte, PluginError) {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return []byte{}, NewPluginError(ReadFailure, ReadFailureError.Error(), err.Error())
	}

	file, err := os.Open(path)
	if err != nil {
		return []byte{}, NewPluginError(ReadFailure, ReadFailureError.Error(), err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, NewPluginError(ReadFailure, ReadFailureError.Error(), err.Error())
	}

	return data, nil
}

func (o *SpaceOperation) ReadLines(name string, lineBegin, lineEnd int32) ([]byte, PluginError) {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return nil, NewPluginError(ReadLinesFailure, ReadLinesFailureError.Error(), err.Error())
	}

	file, err := os.Open(path)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return nil, NewPluginError(ReadLinesFailure, ReadLinesFailureError.Error(), err.Error())
	}
	fileScanner := bufio.NewScanner(file)
	var lineCount int32 = 1
	var fileByte []byte
	for fileScanner.Scan() {
		if lineCount >= lineBegin && lineCount <= lineEnd {
			bt := fileScanner.Bytes()
			fileByte = append(fileByte, bt...)
		}
		if lineCount > lineEnd {
			break
		}
		lineCount++
	}

	return fileByte, nil
}

func (o *SpaceOperation) WriteBytes(name string, byteSlice []byte) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(WriteBytesFailure, WriteBytesFailureError.Error(), err.Error())
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return NewPluginError(WriteBytesFailure, WriteBytesFailureError.Error(), err.Error())
	}
	defer file.Close()

	_, err = file.Write(byteSlice)
	if err != nil {
		return NewPluginError(WriteBytesFailure, WriteBytesFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) AppendBytes(filePath string, byteSlice []byte) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, filePath)
	if err != nil {
		return NewPluginError(AppendBytesFailure, AppendBytesFailureError.Error(), err.Error())
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return NewPluginError(AppendBytesFailure, AppendBytesFailureError.Error(), err.Error())
	}
	defer file.Close()
	_, err = file.Write(byteSlice)
	if err != nil {
		return NewPluginError(AppendBytesFailure, AppendBytesFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) WriteStrings(name string, content []string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(WriteStringsFailure, WriteStringsFailureError.Error(), err.Error())
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return NewPluginError(WriteStringsFailure, WriteStringsFailureError.Error(), err.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return NewPluginError(WriteStringsFailure, WriteStringsFailureError.Error(), err.Error())
		}
	}
	write.Flush()
	return nil
}

func (o *SpaceOperation) AppendStrings(filePath string, content []string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, filePath)
	if err != nil {
		return NewPluginError(AppendStringsFailure, AppendStringsFailureError.Error(), err.Error())
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return NewPluginError(AppendStringsFailure, AppendStringsFailureError.Error(), err.Error())
	}

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return NewPluginError(AppendStringsFailure, AppendStringsFailureError.Error(), err.Error())
		}
	}
	write.Flush()
	return nil
}

func (o *SpaceOperation) Zip(outFileName string, targetFiles []string) PluginError {
	var TargetFiles []string
	for _, v := range targetFiles {
		path, err := GetFilePath(o.AppID, o.InstanceID, v)
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}
		TargetFiles = append(TargetFiles, path)
	}
	OutFileName, err := GetFilePath(o.AppID, o.InstanceID, outFileName)
	if err != nil {
		return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
	}

	type targetFile struct {
		Name string
		Body []byte
	}
	var filesToArchive []targetFile

	for _, targetFilePath := range TargetFiles {
		tFile, err := os.Open(targetFilePath)
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}
		defer tFile.Close()
		info, err := tFile.Stat()
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}
		fileName := info.Name()
		content, err := o.Read(targetFilePath)
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}

		filesToArchive = append(filesToArchive, targetFile{
			fileName,
			content,
		})
	}

	outFile, err := os.Create(OutFileName)
	if err != nil {
		return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	for _, file := range filesToArchive {
		fileWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}
		_, err = fileWriter.Write(file.Body)
		if err != nil {
			return NewPluginError(ZipFailure, ZipFailureError.Error(), err.Error())
		}
	}
	zipWriter.Close()
	return nil
}

func (o *SpaceOperation) UnZip(name string, targetDir string) PluginError {
	zipPath, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
	}
	targetDirPath, err := GetFilePath(o.AppID, o.InstanceID, targetDir)
	if err != nil {
		return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
	}

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
	}
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			targetDirPath,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return NewPluginError(UnZipFailure, UnZipFailureError.Error(), err.Error())
			}
		}
	}
	return nil
}

func (o *SpaceOperation) Gz(name string) PluginError {
	gzFileName := name + ".gz"
	gzPath, err := GetFilePath(o.AppID, o.InstanceID, gzFileName)
	if err != nil {
		return NewPluginError(GzFailure, GzFailureError.Error(), err.Error())
	}
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(GzFailure, GzFailureError.Error(), err.Error())
	}

	outputFile, err := os.Create(gzPath)
	if err != nil {
		return NewPluginError(GzFailure, GzFailureError.Error(), err.Error())
	}

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	content, err := o.Read(path)
	if err != nil {
		return NewPluginError(GzFailure, GzFailureError.Error(), err.Error())
	}
	_, err = gzipWriter.Write(content)
	if err != nil {
		return NewPluginError(GzFailure, GzFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) UnGz(name string, targetFile string) PluginError {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}
	targetPath, err := GetFilePath(o.AppID, o.InstanceID, targetFile)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}

	gzipFile, err := os.Open(path)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}
	defer gzipReader.Close()

	outfileWriter, err := os.Create(targetPath)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}
	defer outfileWriter.Close()

	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return NewPluginError(UnGzFailure, UnGzFailureError.Error(), err.Error())
	}
	return nil
}

func (o *SpaceOperation) Hash(name string) ([]byte, PluginError) {
	path, err := GetFilePath(o.AppID, o.InstanceID, name)
	if err != nil {
		return nil, NewPluginError(HashFailure, HashFailureError.Error(), err.Error())
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, NewPluginError(HashFailure, HashFailureError.Error(), err.Error())
	}
	hash := md5.Sum(data)
	fileByte := hash[:]
	return fileByte, nil
}

func (o *SpaceOperation) List(dirPath string) ([]string, PluginError) {
	// 获取到真实路径
	realPath, err := GetFilePath(o.AppID, o.InstanceID, dirPath)
	if err != nil {
		return nil, NewPluginError(MakeDirFailure, MakeDirFailureError.Error(), err.Error())
	}

	var dirList []string
	err = filepath.Walk(realPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}

			// 对用户隐藏路径前缀
			p := strings.Split(path, realPath+"/")
			if len(p) > 1 && len(p[1]) > 0 {
				dirList = append(dirList, p[1])
			}
			return nil
		})
	if err != nil {
		return dirList, NewPluginError(ListFileFailure, ListFileFailureError.Error(), err.Error())
	}
	return dirList, nil
}

func GetFilePath(appID, instanceID, name string) (path string, err error) {
	pathPrefix := []string{config.StringOrPanic("platform.plugin_runtime_dir"), appID, instanceID}

	var filename string
	if strings.Contains(name, string(os.PathSeparator)) {
		fileSlice := strings.Split(name, string(os.PathSeparator))
		filename = fileSlice[len(fileSlice)-1]
		fileDirName := fileSlice[:len(fileSlice)-1]
		pathPrefix = append(pathPrefix, fileDirName...)
		dirPrefix := filepath.Join(pathPrefix...)
		exist, err := file_utils.PathExists(dirPrefix)
		if err != nil {
			return "", errors.Trace(err)
		}
		if !exist {
			if err := os.MkdirAll(dirPrefix, os.ModePerm); err != nil {
				return "", errors.Trace(err)
			}
		}
	} else {
		exist, err := file_utils.PathExists(filepath.Join(pathPrefix...))
		if err != nil {
			return "", errors.Trace(err)
		}
		if !exist {
			if err := os.MkdirAll(filepath.Join(pathPrefix...), os.ModePerm); err != nil {
				return "", errors.Trace(err)
			}
		}

		filename = name
	}

	pathPrefix = append(pathPrefix, filename)
	path = filepath.Join(pathPrefix...)
	return path, nil
}
