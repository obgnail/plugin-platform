package local

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/host/config"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var _ common.Workspace = (*Space)(nil)

type Space struct {
	plugin common.IPlugin
}

func NewSpace(plugin common.IPlugin) common.Workspace {
	return &Space{plugin: plugin}
}

func (s *Space) CreateFile(name string) common.PluginError {
	newFile, err := os.Create(s.GetSpacePath(name))
	if err != nil {
		return common.NewPluginError(common.CreateFileFailure, err.Error(), common.CreateFileFailureError.Error())
	}
	defer newFile.Close()
	return nil
}

func (s *Space) MakeDir(name string) common.PluginError {
	_name := s.GetSpacePath(name)
	err := os.MkdirAll(_name, os.ModePerm)
	if err != nil {
		return common.NewPluginError(common.MakeDirFailure, err.Error(), common.MakeDirFailureError.Error())
	}
	return nil
}

func (s *Space) Rename(originalPath string, newPath string) common.PluginError {
	_originalPath := s.GetSpacePath(originalPath)
	_newPath := s.GetSpacePath(newPath)

	err := os.Rename(_originalPath, _newPath)
	if err != nil {
		return common.NewPluginError(common.ReNameFileFailure, err.Error(), common.ReNameFileFailureError.Error())
	}
	return nil
}

func (s *Space) Remove(name string) common.PluginError {
	err := os.Remove(s.GetSpacePath(name))
	if err != nil {
		return common.NewPluginError(common.RemoveFileFailure, err.Error(), common.RemoveFileFailureError.Error())
	}
	return nil
}

func (s *Space) IsExist(name string) (bool, common.PluginError) {
	_, err := os.Stat(s.GetSpacePath(name))
	if err != nil && os.IsNotExist(err) {
		return false, common.NewPluginError(common.IsExistFileFailure, err.Error(), common.IsExistFileFailureError.Error())
	}
	return true, nil
}

func (s *Space) IsDir(name string) (bool, common.PluginError) {
	path := s.GetSpacePath(name)
	stat, err := os.Stat(path)
	if err != nil {
		return false, common.NewPluginError(common.IsDirFailure, err.Error(), common.IsDirFailureError.Error())
	}
	return stat.IsDir(), nil
}

func (s *Space) Copy(originalPath string, newPath string) common.PluginError {
	_originalPath := s.GetSpacePath(originalPath)
	_newPath := s.GetSpacePath(newPath)

	originalFile, err := os.Open(_originalPath)
	if err != nil {
		return common.NewPluginError(common.CopyFileFailure, err.Error(), common.CopyFileFailureError.Error())
	}
	defer originalFile.Close()

	newFile, err := os.Create(_newPath)
	if err != nil {
		return common.NewPluginError(common.CopyFileFailure, err.Error(), common.CopyFileFailureError.Error())
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return common.NewPluginError(common.CopyFileFailure, err.Error(), common.CopyFileFailureError.Error())
	}

	err = newFile.Sync()
	if err != nil {
		return common.NewPluginError(common.CopyFileFailure, err.Error(), common.CopyFileFailureError.Error())
	}
	return nil
}

func (s *Space) Read(name string) ([]byte, common.PluginError) {
	file, err := os.Open(s.GetSpacePath(name))
	if err != nil {
		return []byte{}, common.NewPluginError(common.ReadFailure, err.Error(), common.ReadFailureError.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, common.NewPluginError(common.ReadFailure, err.Error(), common.ReadFailureError.Error())
	}

	return data, nil
}

func (s *Space) ReadLines(name string, lineBegin, lineEnd int32) ([]byte, common.PluginError) {
	path := s.GetSpacePath(name)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, common.NewPluginError(common.ReadLinesFailure, err.Error(), common.ReadLinesFailureError.Error())
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

func (s *Space) WriteBytes(name string, byteSlice []byte) common.PluginError {
	file, err := os.OpenFile(s.GetSpacePath(name), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return common.NewPluginError(common.WriteBytesFailure, err.Error(), common.WriteBytesFailureError.Error())
	}
	defer file.Close()

	_, err = file.Write(byteSlice)
	if err != nil {
		return common.NewPluginError(common.WriteBytesFailure, err.Error(), common.WriteBytesFailureError.Error())
	}
	return nil
}

func (s *Space) AppendBytes(filePath string, byteSlice []byte) common.PluginError {
	_filePath := s.GetSpacePath(filePath)

	file, err := os.OpenFile(_filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return common.NewPluginError(common.AppendBytesFailure, err.Error(), common.AppendBytesFailureError.Error())
	}
	defer file.Close()
	_, err = file.Write(byteSlice)
	if err != nil {
		return common.NewPluginError(common.AppendBytesFailure, err.Error(), common.AppendBytesFailureError.Error())
	}
	return nil
}

func (s *Space) WriteStrings(name string, content []string) common.PluginError {
	file, err := os.OpenFile(s.GetSpacePath(name), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return common.NewPluginError(common.WriteStringsFailure, err.Error(), common.WriteStringsFailureError.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return common.NewPluginError(common.WriteStringsFailure, err.Error(), common.WriteStringsFailureError.Error())
		}
	}
	write.Flush()
	return nil
}

func (s *Space) AppendStrings(filePath string, content []string) common.PluginError {
	_filePath := s.GetSpacePath(filePath)

	file, err := os.OpenFile(_filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return common.NewPluginError(common.AppendStringsFailure, err.Error(), common.AppendStringsFailureError.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, val := range content {
		_, err = write.WriteString(val)
		if err != nil {
			return common.NewPluginError(common.AppendStringsFailure, err.Error(), common.AppendStringsFailureError.Error())
		}
	}
	write.Flush()
	return nil
}

func (s *Space) Zip(outFileName string, targetFiles []string) common.PluginError {
	type targetFile struct {
		Name string
		Body []byte
	}
	var filesToArchive []targetFile

	for _, targetFilePath := range targetFiles {
		tFile, err := os.Open(s.GetSpacePath(targetFilePath))
		if err != nil {
			return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
		}
		defer tFile.Close()
		info, err := tFile.Stat()
		if err != nil {
			return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
		}
		fileName := info.Name()
		content, err := s.Read(targetFilePath)
		if err != nil {
			return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
		}

		filesToArchive = append(filesToArchive, targetFile{
			fileName,
			content,
		})
	}

	_outFileName := s.GetSpacePath(outFileName)
	outFile, err := os.Create(_outFileName)
	if err != nil {
		return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	for _, file := range filesToArchive {
		fileWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
		}
		_, err = fileWriter.Write(file.Body)
		if err != nil {
			return common.NewPluginError(common.ZipFailure, err.Error(), common.ZipFailureError.Error())
		}
	}
	zipWriter.Close()
	return nil
}

func (s *Space) UnZip(name string, targetDir string) common.PluginError {
	_name := s.GetSpacePath(name)
	_targetDir := s.GetSpacePath(targetDir)

	zipReader, err := zip.OpenReader(_name)
	if err != nil {
		return common.NewPluginError(common.UnZipFailure, err.Error(), common.UnZipFailureError.Error())
	}
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return common.NewPluginError(common.UnZipFailure, err.Error(), common.UnZipFailureError.Error())
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
				return common.NewPluginError(common.UnZipFailure, err.Error(), common.UnZipFailureError.Error())
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return common.NewPluginError(common.UnZipFailure, err.Error(), common.UnZipFailureError.Error())
			}
		}
	}
	return nil
}

func (s *Space) Gz(name string) common.PluginError {
	gzFileName := name + ".gz"
	_gzFileName := s.GetSpacePath(gzFileName)

	outputFile, err := os.Create(_gzFileName)
	if err != nil {
		return common.NewPluginError(common.GzFailure, err.Error(), common.GzFailureError.Error())
	}

	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	content, err := s.Read(name)
	if err != nil {
		return common.NewPluginError(common.GzFailure, err.Error(), common.GzFailureError.Error())
	}
	_, err = gzipWriter.Write(content)
	if err != nil {
		return common.NewPluginError(common.GzFailure, err.Error(), common.GzFailureError.Error())
	}
	return nil
}

func (s *Space) UnGz(name string, targetFile string) common.PluginError {
	_name := s.GetSpacePath(name)

	gzipFile, err := os.Open(_name)
	if err != nil {
		return common.NewPluginError(common.UnGzFailure, err.Error(), common.UnGzFailureError.Error())
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return common.NewPluginError(common.UnGzFailure, err.Error(), common.UnGzFailureError.Error())
	}
	defer gzipReader.Close()

	outfileWriter, err := os.Create(s.GetSpacePath(targetFile))
	if err != nil {
		return common.NewPluginError(common.UnGzFailure, err.Error(), common.UnGzFailureError.Error())
	}
	defer outfileWriter.Close()

	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return common.NewPluginError(common.UnGzFailure, err.Error(), common.UnGzFailureError.Error())
	}
	return nil
}

func (s *Space) Hash(name string) ([]byte, common.PluginError) {
	data, err := ioutil.ReadFile(s.GetSpacePath(name))

	if err != nil {
		return nil, common.NewPluginError(common.HashFailure, err.Error(), common.HashFailureError.Error())
	}
	hash := md5.Sum(data)
	fileByte := hash[:]
	return fileByte, nil
}

func (s *Space) List(dirPath string) ([]string, common.PluginError) {
	var dirList []string
	_dirPath := s.GetSpacePath(dirPath)

	err := filepath.Walk(_dirPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			dirList = append(dirList, path)
			return nil
		})
	if err != nil {
		return dirList, common.NewPluginError(common.ListFileFailure, err.Error(), common.ListFileFailureError.Error())
	}
	return dirList, nil
}

func (s *Space) AsyncCopy(
	originalPath string,
	newPath string,
	object interface{},
	callback common.AsyncInvokeCallbackParams,
	timeoutCallback common.AsyncInvokeTimeoutCallback,
) {
	returnData := make(chan common.PluginError, 1)
	go func() {
		err := s.Copy(originalPath, newPath)
		returnData <- err
	}()

	select {
	case <-time.After(time.Duration(1000 * 30000000)):
		err := common.NewPluginError(common.MsgTimeOut, common.MsgTimeOutError.Error(), "timeout")
		timeoutCallback(err, object)
	case err := <-returnData:
		callback(err, object)
	}
}

func (s *Space) AsyncZip(
	outFileName string,
	targetFiles []string,
	object interface{},
	callback common.AsyncInvokeCallbackParams,
	timeoutCallback common.AsyncInvokeTimeoutCallback,
) {
	returnData := make(chan common.PluginError, 1)
	go func() {
		err := s.Zip(outFileName, targetFiles)
		returnData <- err
	}()

	select {
	case <-time.After(time.Duration(1000 * 30000000)):
		err := common.NewPluginError(common.MsgTimeOut, common.MsgTimeOutError.Error(), "timeout")
		timeoutCallback(err, object)
	case err := <-returnData:
		callback(err, object)
	}
}

func (s *Space) AsyncUnZip(
	name string,
	targetDir string,
	object interface{},
	callback common.AsyncInvokeCallbackParams,
	timeoutCallback common.AsyncInvokeTimeoutCallback,
) {
	returnData := make(chan common.PluginError, 1)
	go func() {
		err := s.UnZip(name, targetDir)
		returnData <- err
	}()

	select {
	case <-time.After(time.Duration(1000 * 30000000)):
		err := common.NewPluginError(common.MsgTimeOut, common.MsgTimeOutError.Error(), "timeout")
		timeoutCallback(err, object)
	case err := <-returnData:
		callback(err, object)
	}
}

func (s *Space) AsyncGz(
	name string,
	object interface{},
	callback common.AsyncInvokeCallbackParams,
	timeoutCallback common.AsyncInvokeTimeoutCallback,
) {
	returnData := make(chan common.PluginError, 1)
	go func() {
		err := s.Gz(name)
		returnData <- err
	}()

	select {
	case <-time.After(time.Duration(1000 * 30000000)):
		err := common.NewPluginError(common.MsgTimeOut, common.MsgTimeOutError.Error(), "timeout")
		timeoutCallback(err, object)
	case err := <-returnData:
		callback(err, object)
	}
}

func (s *Space) AsyncUnGz(
	name string,
	targetFile string,
	object interface{},
	callback common.AsyncInvokeCallbackParams,
	timeoutCallback common.AsyncInvokeTimeoutCallback,
) {
	returnData := make(chan common.PluginError, 1)
	go func() {
		err := s.UnGz(name, targetFile)
		returnData <- err
	}()

	select {
	case <-time.After(time.Duration(1000 * 30000000)):
		err := common.NewPluginError(common.MsgTimeOut, common.MsgTimeOutError.Error(), "timeout")
		timeoutCallback(err, object)
	case err := <-returnData:
		callback(err, object)
	}
}

func (s *Space) GetSpacePath(name string) string {
	fmt.Println("=============", s.plugin.GetPluginDescription())
	appUUID := s.plugin.GetPluginDescription().PluginDescription().ApplicationID()
	fmt.Println("xxx", appUUID)
	_dir := filepath.Join(config.StringOrPanic("runtime_work_space"), appUUID)

	exist, err := PathExists(_dir)
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

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
