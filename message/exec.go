package message

type Exec struct {
	Command string
	Args    []string
}

func NewExec(command string, args []string) *Exec {
	return &Exec{
		Command: command,
		Args:    args,
	}
}
