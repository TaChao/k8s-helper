package main

import (
	"os"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"strconv"
)

func main() {
	router := iris.New()
	router.Get("/nginx/{name:string}/c/{count:int}", nginx)
	router.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func nginx(ctx iris.Context) {

	var count, err = ctx.Params().GetInt("count")
	result := iris.Map{
		"statu":   500,
		"massage": err,
	}
	if err != nil {
		goto END
	}
	for i := 0; i < count; i++ {
		nginxProcess(ctx.Params().Get("name") + "-" + strconv.Itoa(i))
	}
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(result)
		}
	}()
	result = iris.Map{
		"statu":   200,
		"massage": "success",
	}
END:
	ctx.JSON(result)
}
func nginxProcess(name string) {
	os.Args = []string{os.Args[0], "exec", "-i", name, "--", "nginx", "-s", "reload"}
	cmd := cmd.NewKubectlCommand(cmdutil.NewFactory(nil), os.Stdin, os.Stdout, os.Stderr)
	cmd.Execute()
}
