package main

// import (
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
// )
// func sayHelloHandler(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("*****************Redirect*********************")
//	http.Redirect(w, r, "http://www.baidu.com", http.StatusFound)
//	content := []byte("hello world")
//	err := ioutil.WriteFile("test.txt", content, 0644)
//	if err != nil {
//		panic(err)
//	}
// }
// func main() {
//	http.HandleFunc("/", sayHelloHandler)
//	log.Fatal(http.ListenAndServe(":50280", nil))
// }

import (
	"fmt"
	"os"

	servicename "github.com/NpoolPlatform/third-login-gateway/pkg/service-name"

	"github.com/NpoolPlatform/go-service-framework/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	rabbitmqconst "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const"
	redisconst "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"

	cli "github.com/urfave/cli/v2"
)

func main() {
	commands := cli.Commands{
		runCmd,
	}

	description := fmt.Sprintf("my %v service cli\nFor help on any individual command run <%v COMMAND -h>\n",
		servicename.ServiceName, servicename.ServiceName)
	err := app.Init(
		servicename.ServiceName,
		description,
		"",
		"",
		"./",
		nil,
		commands,
		mysqlconst.MysqlServiceName,
		rabbitmqconst.RabbitMQServiceName,
		redisconst.RedisServiceName,
	)
	if err != nil {
		logger.Sugar().Errorf("fail to create %v: %v", servicename.ServiceName, err)
		return
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Sugar().Errorf("fail to run %v: %v", servicename.ServiceName, err)
	}
}
