### regcmd 正则指令解析器

#### 用途

当一个服务需要在运行时对一些简要的后台指令做出响应时，regcmd 提供一种监听 stdin 的指令信息并执行代码的简要操作逻辑。您可以注册自己的指令和操作，regcmd 会自动生成对应的 help 指令，之后监听指定的数据流 （如 stdin）即可解析指令。

#### 使用

```go
package main

import (
    "fmt"
    "github.com/MinoIC/regcmd"
    "os"
)

func main(){
    _ = regcmd.Register("show",[]string{"say hello world"}, func(args []string) {
        fmt.Println("hello world")
    })
    _ = regcmd.Register("show (.*)",[]string{"user","say hello to the given user name"}, func(args []string) {
        fmt.Println("hello",args[0])
    })
    regcmd.Listen(os.Stdin)
}
```

运行以上代码，可以经由控制台完成一些简单的指令，结果如下图。其中 `? ` `help` `show` `show minoic` 为用户的输入，regcmd 对每一条指令给出了应答。

```
?
Invalid command: ? **Type <help> for commands help
help
---- all commands help ----
help // To get all commands help
show <user> // say hello to the given user name
show help   // To get this help
show        // say hello world
show
hello world
show minoic
hello minoic
```

