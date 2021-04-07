package pkg

// TestSetup 测试全局准备
func TestSetup(confFilepath string) {
	Conf = NewConfig(confFilepath)
	Log = NewLog(Conf.LogPath, DebugLevel)
	DB = NewDB(Conf.ArangoDB.Url, Conf.ArangoDB.Username, Conf.ArangoDB.Passwd,
		Conf.ArangoDB.Database, Conf.ArangoDB.ModelColNameMap)
}
