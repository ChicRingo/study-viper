package main

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	//建立默认值
	viper.SetDefault("fileDir", "./")

	//读取配置文件
	//viper.SetConfigFile("./conf/config.yaml") // 指定配置文件，要写完整目录，和下面三行不要同时使用，会覆盖配置
	viper.AddConfigPath("./conf/") // 还可以在工作目录中查找配置
	viper.SetConfigName("config")  // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")    // 如果配置文件的名称中没有扩展名，则需要配置此项

	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	/*
		viper.WriteConfig() // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径
		viper.SafeWriteConfig()
		viper.WriteConfigAs("/path/to/my/.config")
		viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
		viper.SafeWriteConfigAs("/path/to/my/.other_config")
	*/
	/*
		监控并重新读取配置文件
		Viper支持在运行时实时读取配置文件的功能。

		需要重新启动服务器以使配置生效的日子已经一去不复返了，viper驱动的应用程序可以在运行时读取配置文件的更新，而不会错过任何消息。

		只需告诉viper实例watchConfig。可选地，你可以为Viper提供一个回调函数，以便在每次发生更改时运行。

		确保在调用WatchConfig()之前添加了所有的配置路径。
	*/
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})

	err = r.Run()
	if err != nil {
		return
	}

}

//Viper在后台使用github.com/mitchellh/mapstructure来解析值，其默认情况下使用mapstructuretag。
type config struct {
	Port    int
	Name    string
	PathMap string `mapstructure:"path_map"`
}
