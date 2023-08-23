package rules

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/yara-rule-server/pkg/config"
)

func DownloadAndCompile(cfg *config.Config, logger *logrus.Entry) error {
	archives := download(cfg.RulePath, cfg.RuleURLs, logger)
	unarchive(archives, logger)

	if err := generateIndex(cfg.IndexGenPath, logger); err != nil {
		return fmt.Errorf("failed to generate index.yar: %v", err)
	}
	if err := compile(cfg.YaracPath, "index.yar", cfg.RulePath, logger); err != nil {
		return fmt.Errorf("failed to copile index.yar: %v", err)
	}

	return nil
}

func ScheduledDownload(cfg *config.Config, logger *logrus.Entry) error {
	s := gocron.NewScheduler(time.UTC)

	// cron expressions supported
	if _, err := s.Cron(cfg.RuleUpdateSchedule).Do(DownloadAndCompile, cfg, logger); err != nil {
		return fmt.Errorf("failed to create cron job: %v", err)
	}

	// starts the scheduler asynchronously
	s.StartAsync()

	return nil
}
