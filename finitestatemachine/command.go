package finitestatemachine

import "fmt"

type CommandFunc func(*SMContext, Events) error

type Command struct {
	Name      string
	SkipError bool
	Fn        CommandFunc
}

type Commands []Command

func (cs Commands) exec(ctx *SMContext, es Events) {
	for _, c := range cs {
		err := c.Fn(ctx, es)
		if err != nil {
			fmt.Println(fmt.Sprintf("命令[%s]执行失败", c.Name), err)
			if !c.SkipError {
				return
			}
		}
	}
}
