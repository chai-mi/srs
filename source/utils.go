package source

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getInput(src string) ([]byte, error) {
	URL, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("request error! ")
	}
	return io.ReadAll(resp.Body)
}

func removeComment(line string) string {
	before, _, ok := strings.Cut(line, "#")
	if !ok {
		return line
	}
	return strings.TrimSpace(before)
}
