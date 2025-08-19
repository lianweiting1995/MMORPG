package contract

type IConfig interface {
	GetServer() IServer
	GetRegistry() IRegistry
}

type IServer interface {
	GetName() string
	GetVersion() string
}

type IRegistry interface {
	GetConsul() IConsul
}

type IConsul interface {
	GetAddrs() []string
}
