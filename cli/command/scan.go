package command

type (
	Scan struct {
		Directory string `short:"d" long:"directory"`
	}
)

func (s *Scan) Execute(args []string) error {
	return nil
}
