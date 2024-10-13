package command

type Command string

const (
	Start   Command = "/start"
	Profile Command = "/profile"
)

func FromString(str string) Command {
	return Command(str)
}
