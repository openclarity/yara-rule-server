package rules

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func download(rulePath string, urls []string, logger *logrus.Entry) []string {
	archives := make([]string, 0)
	for i := range urls {
		u, err := url.Parse(urls[i])
		if err != nil {
			logger.Errorf("Failed to parse url=%s: %v", urls[i], err)
			continue
		}

		dirName := fmt.Sprintf("%s%s", path.Dir(rulePath), path.Dir(u.Path))
		if err := os.MkdirAll(dirName, 0755); err != nil {
			logger.Errorf("failed to create directory=%s: %v", dirName, err)
		}
		fileName := fmt.Sprintf("%s/%s", dirName, path.Base(u.Path))
		logger.Infof("Downloading %s into %s", urls[i], fileName)
		if err := downloadFile(fileName, urls[i]); err != nil {
			logger.Errorf("Failed to download file from url=%s: %v", urls[i], err)
		}
		archives = append(archives, fileName)
	}
	return archives
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create ouptut file=%s: %v", filepath, err)
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get url=%s: %s", url, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file=%s: %v", filepath, err)
	}

	return nil
}
