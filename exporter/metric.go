package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getNebulaMetrics(ipAddress string, port int32) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/stats", ipAddress, port))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	metrics := strings.Split(strings.TrimSpace(string(bytes)), "\n")

	return metrics, nil
}

func isNebulaComponentRunning(ipAddress string, port int32) bool {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/status", ipAddress, port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	type nebulaStatus struct {
		GitInfoSha string `json:"git_info_sha"`
		Status     string `json:"status"`
	}

	var status nebulaStatus
	if err := json.Unmarshal(bytes, &status); err != nil {
		return false
	}

	return status.Status == "running"
}
