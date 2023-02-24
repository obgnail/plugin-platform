package local

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/file_utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var _ common_type.Workspace = (*Space)(nil)

type Space struct {
	plugin common_type.IPlugin
}

func NewSpace(plugin common_type.IPlugin) common_type.Workspace {
	return &Space{plugin: plugin}
}

func (s *Space) CreateFile(name string) common_type.PluginError {
	newFile, err := os.Create(s.getSpacePath(name))
	if err != nil {
		return common_type.NewPluginError(common_type.CreateFileFailure, err.Error(), common_type.CreateFileFailureError.Error())
	}
	defer newFile.Close()
	return nil
}

func (s *Space) MakeDir(name string) common_type.PluginError {
	_name := s.getSpacePath(name)
	err := os.MkdirAll(_name, os.ModePerm)
	if err != nil {
		return common_type.NewPluginError(common_type.MakeDirFailure, err.Error(), common_type.MakeDirFailureError.Error())
	}
	return nil
}

func (s *Space) Rename(originalPath string, newPath string) common_type.PluginError {
	_originalPath := s.getSpacePath(originalPath)
	_newPath := s.getSpacePath(newPath)

	err := os.Rename(_originalPath, _newPath)
	if err != nil {
		return common_type.NewPluginError(common_type.ReNameFileFailure, err.Error(), common_type.ReNameFileFailureError.Error())
	}
	return nil
}

func (s *Space) Remove(name string) common_type.PluginError {
	err := os.Remove(s.getSpacePath(name))
	if err != nil {
		return common_type.NewPluginError(common_type.RemoveFileFailure, err.Error(), common_type.RemoveFileFailureError.Error())
	}
	return nil
}

func (s *Space) IsExist(name string) (bool, common_type.PluginError) {
	_, err := os.Stat(s.getSpacePath(name))
	if err != nil && os.IsNotExist(err) {
		return false, common_type.NewPluginError(common_type.IsExistFileFailure, err.Error(), common_type.IsExistFileFailureError.Error())
	}
	return true, nil
}

func (s *Space) IsDir(name string) (bool, common_type.PluginError) {
	path := s.getSpacePath(name)
	stat, err := os.Stat(path)
	if err != nil {
		return false, common_type.NewPluginError(common_type.IsDirFailure, err.Error(), common_type.IsDirFailureError.Error())
	}
	return stat.IsDir(), nil
}

func (s *Space) Copy(originalPath string, newPath string) common_type.PluginError {
	_originalPath := s.getSpacePath(originalPath)
	_newPath := s.getSpacePath(newPath)

	originalFile, err := os.Open(_originalPath)
	if err != nil {
		return common_type.NewPluginError(common_type.CopyFileFailure, err.Error(), common_type.CopyFileFailureError.Error())
	}
	defer originalFile.Close()

	newFile, err := os.Create(_newPath)
	if err != nil {
		return common_type.NewPluginError(common_type.CopyFileFailure, err.Error(), common_type.CopyFileFailureError.Error())
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return common_type.NewPluginError(common_type.CopyFileFailure, err.Error(), common_type.CopyFileFailureError.Error())
	}

	err = newFile.Sync()
	if err != nil {
		return common_type.NewPluginError(common_type.CopyFileFailure, err.Error(), common_type.CopyFileFailureError.Error())
	}
	return nil
}

func (s *Space) Read(name string) ([]byte, common_type.PluginError) {
	file, err := os.Open(s.getSpacePath(name))
	if err != nil {
		return []byte{}, common_type.NewPluginError(common_type.ReadFailure, err.Error(), common_type.ReadFailureError.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, common_type.NewPluginError(common_type.ReadFailure, err.Error(), common_type.ReadFailureError.Error())
	}

	return data, nil
}

func (s *Space) ReadLines(name string, lineBegin, lineEnd int32) ([]byte, common_type.PluginError) {
	path := s.getSpacePath(name)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, common_type.NewPluginError(common_type.ReadLinesFailure, err.Error(), common_type.ReadLinesFailureError.Error())
	}
	fileScanner := bufio.NewScanner(file)
	var lineCount int32 = 1
	var fileByte []byte
	for fileScanner.Scan() {
		if lineCount >= lineBegin && lineCount <= lineEnd {
			bs := fileScanner.Bytes()
			for _, v := range bs {
				fileByte = append(fileByte, v)
			}
		}
		if lineCount > lineEnd {
			break
		}
		lineCount++
	}

	return fileByte, nil
}

func (s *Space) WriteBytes(name string, byteSlice []byte) common_type.PluginError {
	file, err := os.OpenFile(s.getSpacePath(name), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return common_type.NewPluginError(common_type.WriteBytesFailure, err.Error(), common_type.WriteBytesFailureError.Error())
	}
	defer file.Close()

	_, err = file.Write(byteSlice)
	if err != nil {
		return common_type.NewPluginError(common_type.WriteBytesFailure, err.Error(), common_type.WriteBytesFailureError.Error())
	}
	return nil
}

func (s *Space) AppendBytes(filePath string, byteSlice []byte) common_type.PluginError {
	_filePath := s.getSpacePath(filePath)

	file, err := os.OpenFile(_filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return common_type.NewPluginError(common_type.AppendBytesFailure, err.Error(), common_type.AppendBytesFailureError.Error())
	}
	defer file.Close()
	_, err = file.Write(byteSlice)
	if err != nil {
		return common_type.NewPluginError(common_type.AppendBytesFailure, err.Error(), common_type.AppendBytesFailureError.Error())
	}
	return nil
}

func (s *Space) WriteStrings(name string, content []string) common_type.PluginError {
	file, err := os.OpenFile(s.getSpacePath(name), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return common_type.NewPluginError(common_type.WriteStringsFailure, err.Error(), common_type.WriteStringsFailureError.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return common_type.NewPluginError(common_type.WriteStringsFailure, err.Error(), common_type.WriteStringsFailureError.Error())
		}
	}
	write.Flush()
	return nil
}

func (s *Space) AppendStrings(filePath string, content []string) common_type.PluginError {
	_filePath := s.getSpacePath(filePath)

	file, err := os.OpenFile(_filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return common_type.NewPluginError(common_type.AppendStringsFailure, err.Error(), common_type.AppendStringsFailureError.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return common_type.NewPluginError(common_type.AppendStringsFailure, err.Error(), common_type.AppendStringsFailureError.Error())
		}
	}
	write.Flush()
	return nil
}

func (s *Space) Zip(outFileName string, targetFiles []string) common_type.PluginError {
	type targetFile struct {
		Name string
		Body []byte
	}
	var filesToArchive []targetFile

	for _, targetFilePath := range targetFiles {
		tFile, err := os.Open(s.getSpacePath(targetFilePath))
		if err != nil {
			return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
		}
		defer tFile.Close()
		info, err := tFile.Stat()
		if err != nil {
			return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
		}
		fileName := info.Name()
		content, err := s.Read(targetFilePath)
		if err != nil {
			return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
		}

		filesToArchive = append(filesToArchive, targetFile{
			fileName,
			content,
		})
	}

	_outFileName := s.getSpacePath(outFileName)
	outFile, err := os.Create(_outFileName)
	if err != nil {
		return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	for _, file := range filesToArchive {
		fileWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
		}
		_, err = fileWriter.Write(file.Body)
		if err != nil {
			return common_type.NewPluginError(common_type.ZipFailure, err.Error(), common_type.ZipFailureError.Error())
		}
	}
	zipWriter.Close()
	return nil
}

func (s *Space) UnZip(name string, targetDir string) common_type.PluginError {
	_name := s.getSpacePath(name)
	_targetDir := s.getSpacePath(targetDir)

	zipReader, err := zip.OpenReader(_name)
	if err != nil {
		return common_type.NewPluginError(common_type.UnZipFailure, err.Error(), common_type.UnZipFailureError.Error())
	}
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return common_type.NewPluginError(common_type.UnZipFailure, err.Error(), common_type.UnZipFailureError.Error())
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			_targetDir,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return common_type.NewPluginError(common_type.UnZipFailure, err.Error(), common_type.UnZipFailureError.Error())
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return common_type.NewPluginError(common_type.UnZipFailure, err.Error(), common_type.UnZipFailureError.Error())
			}
		}
	}
	return nil
}

func (s *Space) Gz(name string) common_type.PluginError {
	gzFileName := name + ".gz"
	_gzFileName := s.getSpacePath(gzFileName)

	outputFile, err := os.Create(_gzFileName)
	if err != nil {
		return common_type.NewPluginError(common_type.GzFailure, err.Error(), common_type.GzFailureError.Error())
	}

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	content, err := s.Read(name)
	if err != nil {
		return common_type.NewPluginError(common_type.GzFailure, err.Error(), common_type.GzFailureError.Error())
	}
	_, err = gzipWriter.Write(content)
	if err != nil {
		return common_type.NewPluginError(common_type.GzFailure, err.Error(), common_type.GzFailureError.Error())
	}
	return nil
}

func (s *Space) UnGz(name string, targetFile string) common_type.PluginError {
	_name := s.getSpacePath(name)

	gzipFile, err := os.Open(_name)
	if err != nil {
		return common_type.NewPluginError(common_type.UnGzFailure, err.Error(), common_type.UnGzFailureError.Error())
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return common_type.NewPluginError(common_type.UnGzFailure, err.Error(), common_type.UnGzFailureError.Error())
	}
	defer gzipReader.Close()

	outfileWriter, err := os.Create(s.getSpacePath(targetFile))
	if err != nil {
		return common_type.NewPluginError(common_type.UnGzFailure, err.Error(), common_type.UnGzFailureError.Error())
	}
	defer outfileWriter.Close()

	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return common_type.NewPluginError(common_type.UnGzFailure, err.Error(), common_type.UnGzFailureError.Error())
	}
	return nil
}

func (s *Space) Hash(name string) ([]byte, common_type.PluginError) {
	data, err := ioutil.ReadFile(s.getSpacePath(name))

	if err != nil {
		return nil, common_type.NewPluginError(common_type.HashFailure, err.Error(), common_type.HashFailureError.Error())
	}
	hash := md5.Sum(data)
	fileByte := hash[:]
	return fileByte, nil
}

func (s *Space) List(dirPath string) ([]string, common_type.PluginError) {
	var dirList []string
	_dirPath := s.getSpacePath(dirPath)

	err := filepath.Walk(_dirPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			dirList = append(dirList, path)
			return nil
		})
	if err != nil {
		return dirList, common_type.NewPluginError(common_type.ListFileFailure, err.Error(), common_type.ListFileFailureError.Error())
	}
	return dirList, nil
}

func (s *Space) AsyncCopy(originalPath string, newPath string, callback common_type.AsyncInvokeCallbackParams) {
	errChan := make(chan common_type.PluginError, 1)
	go func() {
		errChan <- s.Copy(originalPath, newPath)
	}()

	s.callback(errChan, callback)
}

func (s *Space) AsyncZip(outFileName string, targetFiles []string, callback common_type.AsyncInvokeCallbackParams) {
	errChan := make(chan common_type.PluginError, 1)
	go func() {
		errChan <- s.Zip(outFileName, targetFiles)
	}()

	s.callback(errChan, callback)
}

func (s *Space) AsyncUnZip(name string, targetDir string, callback common_type.AsyncInvokeCallbackParams) {
	errChan := make(chan common_type.PluginError, 1)
	go func() {
		errChan <- s.UnZip(name, targetDir)
	}()

	s.callback(errChan, callback)
}

func (s *Space) AsyncGz(name string, callback common_type.AsyncInvokeCallbackParams) {
	errChan := make(chan common_type.PluginError, 1)
	go func() {
		errChan <- s.Gz(name)
	}()

	s.callback(errChan, callback)
}

func (s *Space) AsyncUnGz(name string, targetFile string, callback common_type.AsyncInvokeCallbackParams) {
	errChan := make(chan common_type.PluginError, 1)
	go func() {
		errChan <- s.UnGz(name, targetFile)
	}()

	s.callback(errChan, callback)
}

func (s *Space) getSpacePath(name string) string {
	appUUID := s.plugin.GetPluginDescription().PluginDescription().ApplicationID()
	instanceUUID := s.plugin.GetPluginDescription().InstanceID()
	_dir := filepath.Join(config.StringOrPanic("host.plugin_workspace_dir"), appUUID, instanceUUID)

	exist, err := file_utils.PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
	}

	if !exist {
		// 创建文件夹
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}

	var _filePath string
	if name == "" {
		_filePath = _dir
	} else {
		_filePath = _dir + "/" + name
	}

	return _filePath
}

func (s *Space) callback(errChan chan common_type.PluginError, callback common_type.AsyncInvokeCallbackParams) {
	var err common_type.PluginError
	select {
	case <-time.After(3 * time.Second):
		err = common_type.NewPluginError(common_type.MsgTimeOut, common_type.MsgTimeOutError.Error(), "timeout")
	case err = <-errChan:
	}
	callback(err)
}
