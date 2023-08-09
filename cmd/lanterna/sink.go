package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// postJSON issues a HTTP POST to url, with msg as body, JSON-encoded.
//
// Unfortunately the GChat API for incoming webhooks wants the secret as URL
// parameter, thus we must redact the URL when logging.
func postJSON(url string, msg map[string]string) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("JSON encoding: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("http new request: %w", RedactErrorURL(err))
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %s", RedactErrorURL(err))
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %s; URL: %s; body: %s",
			resp.Status, RedactURL(req.URL), strings.TrimSpace(string(respBody)))
	}

	return nil
}

// RedactURL returns a _best effort_ redacted copy of theURL.
// Use this workaround only when you are forced to use an API that encodes
// secrets in the URL instead of setting them in the request header.
// If you have control of the API, please never encode secrets in the URL.
// Redaction is applied as follows:
// - removal of all query parameters
// - removal of "username:password@" HTTP Basic Authentication
// Warning: it is still possible that the redacted URL contains secrets, for
// example if the secret is encoded in the path. Don't do this.
func RedactURL(theURL *url.URL) *url.URL {
	urlCopy := *theURL

	// remove all query parameters
	if urlCopy.RawQuery != "" {
		urlCopy.RawQuery = "REDACTED"
	}
	// remove password in user:password@host
	if _, ok := urlCopy.User.Password(); ok {
		urlCopy.User = url.UserPassword("REDACTED", "REDACTED")
	}

	return &urlCopy
}

// RedactErrorURL returns a _best effort_ redacted copy of err. See
// RedactURL for caveats and limitations.
// In case err is not of type url.Error, then it returns the error untouched.
func RedactErrorURL(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		urlErr.URL = RedactURLString(urlErr.URL)
		return urlErr
	}
	return err
}

// RedactURLString returns a _best effort_ redacted copy of theURL. See
// RedactURL for caveats and limitations.
// In case theURL cannot be parsed, then return the parse error string.
func RedactURLString(theURL string) string {
	urlo, err := url.Parse(theURL)
	if err != nil {
		return err.Error()
	}
	return RedactURL(urlo).String()
}
