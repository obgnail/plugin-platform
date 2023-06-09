package common_type

type HttpRequest struct {
	Method   string
	QueryMap map[string]string
	Url      string
	Path     string
	Headers  map[string][]string
	Body     []byte
	Root     bool
}

type HttpResponse struct {
	Err        PluginError
	Headers    map[string][]string
	Body       []byte
	StatusCode int
}

type HttpContext struct {
	Params   map[string]string
	Request  *HttpRequest
	Response *HttpResponse
}

type EventPublisher interface {
	Subscribe(condition []string) PluginError                                       // 不支持通配符, ex: project.task.create
	SubscribeWithFilter(condition []string, filter map[string][]string) PluginError // 支持过滤条件, ex: filter: project_uuid_in:[""]
	Unsubscribe(condition []string) PluginError                                     // 支持通配符
}

type AsyncInvokeCallbackParams func(PluginError)

type Workspace interface {
	CreateFile(string) PluginError
	MakeDir(string) PluginError
	Rename(string, string) PluginError
	Remove(string) PluginError
	IsExist(string) (bool, PluginError)
	IsDir(string) (bool, PluginError)
	Copy(string, string) PluginError
	Read(string) ([]byte, PluginError)
	ReadLines(string, int32, int32) ([]byte, PluginError)
	WriteBytes(string, []byte) PluginError
	AppendBytes(string, []byte) PluginError
	WriteStrings(string, []string) PluginError
	AppendStrings(string, []string) PluginError
	Zip(string, []string) PluginError
	UnZip(string, string) PluginError
	Gz(string) PluginError
	UnGz(string, string) PluginError
	Hash(string) ([]byte, PluginError)
	List(string) ([]string, PluginError)

	AsyncCopy(string, string, AsyncInvokeCallbackParams)
	AsyncZip(string, []string, AsyncInvokeCallbackParams)
	AsyncUnZip(string, string, AsyncInvokeCallbackParams)
	AsyncGz(string, AsyncInvokeCallbackParams)
	AsyncUnGz(string, string, AsyncInvokeCallbackParams)
}

type NetworkCallBack func(*HttpResponse, PluginError)

type Network interface {
	Fetch(*HttpRequest) *HttpResponse

	AsyncFetch(*HttpRequest, NetworkCallBack)
}

type APICore interface {
	Fetch(*HttpRequest) *HttpResponse
}

type DBCallBack func([]*RawData, []*ColumnDesc, PluginError)

type RawData struct {
	Cell [][]byte
}

type ColumnDesc struct {
	Index int64
	Name  string
	Type  string
}

type SysDB interface {
	Select(db, sql string) ([]*RawData, []*ColumnDesc, PluginError)
	AsyncSelect(db, sql string, callback DBCallBack)
	Exec(db, sql string) PluginError
	// Unmarshal eg:
	//	type User struct {
	//		UUID string `orm:"uuid"`
	//		Name string `orm:"name"`
	//	}
	//	users := make([]*User, 0)
	//	err = Unmarshal(rawData, colDesc, &users)
	Unmarshal(rawData []*RawData, columnDesc []*ColumnDesc, v interface{}) PluginError
}

type LocalDB interface {
	Select(sql string) ([]*RawData, []*ColumnDesc, PluginError)
	AsyncSelect(sql string, callback DBCallBack)
	// Unmarshal eg:
	//	type User struct {
	//		UUID string `orm:"uuid"`
	//		Name string `orm:"name"`
	//	}
	//	users := make([]*User, 0)
	//	err = Unmarshal(rawData, colDesc, &users)
	Unmarshal(rawData []*RawData, columnDesc []*ColumnDesc, v interface{}) PluginError
	Exec(sql string) PluginError
	ImportSQL(sqlFilePath string) PluginError
}

type Ability interface {
	Call(kind string, req []byte) ([]byte, PluginError)
}
