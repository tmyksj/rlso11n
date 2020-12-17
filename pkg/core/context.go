package core

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"os"
	"sync"
)

type Context struct {
	dir           string
	finalizerList []func()
	finalizerMu   sync.Mutex
	managerNode   *object.Node
	myNode        *object.Node
	nodeList      []*object.Node
	ready         bool
}

var context = &Context{
	ready: false,
}

func GetContext() *Context {
	return context
}

func (context *Context) AddFinalizer(fun func()) {
	context.finalizerMu.Lock()
	defer context.finalizerMu.Unlock()

	context.finalizerList = append(context.finalizerList, fun)
}

func (context *Context) CheckReady() error {
	if context.ready {
		return nil
	} else {
		return errors.By(nil, "context is not ready")
	}
}

func (context *Context) Dir() string {
	return context.dir
}

func (context *Context) DockerSock() string {
	return context.Dir() + "/runtime/docker.sock"
}

func (context *Context) Env() []string {
	return append(os.Environ(),
		"XDG_CACHE_HOME="+context.Dir()+"/cache",
		"XDG_CONFIG_HOME="+context.Dir()+"/config",
		"XDG_DATA_HOME="+context.Dir()+"/data",
		"XDG_RUNTIME_DIR="+context.Dir()+"/runtime")
}

func (context *Context) ManagerNode() *object.Node {
	return context.managerNode
}

func (context *Context) MyNode() *object.Node {
	return context.myNode
}

func (context *Context) NodeList() []*object.Node {
	return context.nodeList
}

func (context *Context) RpcPort() int {
	return 50128
}
