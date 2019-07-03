package output

import "fmt"

type ConsoleOutput struct {
	OutputBase
}

func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (c *ConsoleOutput) Write() error {
	pairs := c.OutputBase.pairs

	fmt.Println("===== SAME FILES =====")
	idx := 0
	for _, p := range pairs {

		//Don't output different files
		if !p.Compare() {
			continue
		}

		//Files are different. Output
		idx += 1
		fmt.Printf("[%d] %s\n", idx, p.FileOne.GetSubPath())

	}

	fmt.Println()
	fmt.Println("===== DIFFERENT FILES =====")
	idx = 0
	for _, p := range pairs {

		//Don't output same files
		if p.Compare() {
			continue
		}

		//Files are different. Output
		idx += 1
		fmt.Printf("[%d] %s\n", idx, p.FileOne.GetSubPath())

	}

	return nil
}
