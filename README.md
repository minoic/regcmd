### regcmd 正则指令解析器

#### 用途

当一个服务需要在运行时对一些简要的后台指令做出响应时，regcmd 提供一种监听 stdin 的指令信息并执行代码的简要操作逻辑。您可以注册自己的指令和操作，regcmd 会自动生成对应的 help 指令，之后监听指定的数据流 （如 stdin）即可解析指令。

#### 使用

```go
package main

import (
    "context"
    "fmt"
    "github.com/MinoIC/regcmd"
    "math/rand"
    "os"
    "time"
)

func main() {
    regcmd.ShouldRegister("show", []string{"say hello world"}, func(ctx *regcmd.Context, args []string) {
        fmt.Println("hello world")
    })
    regcmd.ShouldRegister("show (.*)", []string{"user", "say hello to the given user name"}, func(ctx *regcmd.Context, args []string) {
        fmt.Println("hello", args[0])
    })
    err := regcmd.Register("sleep", []string{}, func(ctx *regcmd.Context, args []string) {
        id := rand.Int()
        fmt.Println(id, "starts to sleep")
        time.Sleep(5 * time.Second)
        fmt.Println(id, "sleeped 5 seconds")
    })
    if err != nil {
        fmt.Println(err)
    }
    regcmd.Listen(os.Stdin, regcmd.WithPoolSize(2),
        regcmd.WithLoggerFunc(func(s string) {
            fmt.Println("log: ", s)
        }),
        regcmd.WithContextGeneration(func() context.Context {
            return context.Background()
        }),
    )
}
```

运行以上代码，可以经由控制台完成一些简单的指令，结果如下图。

```bash
s
log:  Invalid command: s **Type <help> for commands help
help
log:  ---- all commands help ----
log:  sleep help // To get this help
log:  sleep 
log:  show <user> // say hello to the given user name
log:  show help   // To get this help
log:  show        // say hello world
log:  help // To get all commands help
sleep
5577006791947779410 starts to sleep
5577006791947779410 sleeped 5 seconds
sleep
8674665223082153551 starts to sleep
sleep
6129484611666145821 starts to sleep
sleep
sleep
8674665223082153551 sleeped 5 seconds
4037200794235010051 starts to sleep
6129484611666145821 sleeped 5 seconds
3916589616287113937 starts to sleep
4037200794235010051 sleeped 5 seconds
3916589616287113937 sleeped 5 seconds
```

