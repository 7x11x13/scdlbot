package download

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"

	cookierelayclient "github.com/7x11x13/scdlbot/cookie-relay-client"
)

func SoundCloud(downloadDir string, url string) (*string, error) {
	userID := os.Getenv("SOUNDCLOUD_USER_ID")
	token, err := cookierelayclient.GetCookieValueWithName("soundcloud", userID, "oauth_token")
	if err != nil {
		return nil, err
	}
	scdlCmd := exec.Command("scdl", "-l", url, "--auth-token", *token, "--max-size", "25mb", "--original-art", "--path", downloadDir, "--hide-progress", "--no-original")
	log.Println("Calling SCDL: ", scdlCmd)
	_, err = scdlCmd.Output()
	var eerr *exec.ExitError
	if errors.As(err, &eerr) {
		log.Printf("SCDL output: %s", string(eerr.Stderr))
	}
	if err != nil {
		return nil, err
	}

	matches, err := fs.Glob(os.DirFS(downloadDir), "*")
	if err != nil {
		return nil, err
	}
	if matches == nil || len(matches) == 0 {
		return nil, fmt.Errorf("Download failed")
	}
	downloaded := path.Join(downloadDir, matches[0])
	return &downloaded, nil
}
