package docker_compose_yaml

type NodeYaml struct {
	Peer     PeerYaml
	Order    OrderYaml
	Cli      CliYaml
	CouchDB  CouchDBYaml
	Explorer ExplorerYaml
}

type Ports struct {
	Order7050    string
	Order7053    string
	Peer7051     string
	Peer7052     string
	Explorer8080 string
	Couchdb5984  string
}

type File struct {
	Name string //文件名称
	Data []byte //文件内容
}

func NewNodeYaml() *NodeYaml {
	return &NodeYaml{}
}
