package rules

import (
	"fmt"
	"os/exec"
)

func Compile(output string) error {
	if err := generateIndex(); err != nil {
		return fmt.Errorf("failed to generate index.yar: %v", err)
	}
	if err := compile("index.yar", output); err != nil {
		return fmt.Errorf("failed to copile index.yar: %v", err)
	}

	return nil
}

func generateIndex() error {
	indexGen := exec.Command("index_gen.sh")

	resultJsonB, err := indexGen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run index_gen.sh: %v, %v", err, string(resultJsonB))
	}

	return nil
}

func compile(input, output string) error {
	yarac := exec.Command("yarac", input, output)

	resultJsonB, err := yarac.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run yarac: %v, %v", err, string(resultJsonB))
	}

	return nil
}
