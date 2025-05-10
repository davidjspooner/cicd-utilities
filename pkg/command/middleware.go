package command

type Middeware func(Command)

func LogicalGroup(cmd Command) {
	get, ok := cmd.(getCommonImpl)
	if !ok {
		panic("command is not of type getCommon")
	}
	get.common().logicalGroup = true
}

func PostRun[T any](postRun ExecuteFunc[T]) Middeware {
	return func(cmd Command) {
		impl, ok := cmd.(*commandImpl[T])
		if !ok {
			panic("command is not of type commandImpl")
		}
		impl.postRun = postRun
	}
}

func Aliases[T any](aliases ...string) Middeware {
	return func(cmd Command) {
		impl, ok := cmd.(*commandImpl[T])
		if !ok {
			panic("command is not of type commandImpl")
		}
		impl.aliases = aliases
	}
}
