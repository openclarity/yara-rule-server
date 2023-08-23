package rules

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func generateIndex(indexGenPATH string, logger *logrus.Entry) error {
	indexGen := exec.Command(indexGenPATH)

	resultJsonB, err := indexGen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run index_gen.sh: %v, %v", err, string(resultJsonB))
	}
	logger.Debugln(string(resultJsonB))

	return nil
}

func compile(yaracPATH, input, output string, logger *logrus.Entry) error {
	yarac := exec.Command(yaracPATH, input, output)

	resultJsonB, err := yarac.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run yarac: %v, %v", err, string(resultJsonB))
	}
	logger.Debugln(string(resultJsonB))

	return nil
}
