package core

type GlobalOptions struct {
	ConfigFile string
}

func NewGlobalOptions(configFile string) *GlobalOptions {
	return &GlobalOptions{
		ConfigFile: configFile,
	}
}
