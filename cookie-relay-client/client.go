package cookierelayclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
)

type CookiePartitionKey struct {
	TopLevelSite string `json:"topLevelSite"`
}

type Cookie struct {
	Domain           string             `json:"domain"`
	ExpirationDate   int64              `json:"expirationDate"`
	FirstPartyDomain string             `json:"firstPartyDomain,omitempty"`
	HostOnly         bool               `json:"hostOnly"`
	HttpOnly         bool               `json:"httpOnly"`
	Name             string             `json:"name"`
	PartitionKey     CookiePartitionKey `json:"partitionKey,omitempty"`
	Path             string             `json:"path"`
	Secure           bool               `json:"secure"`
	Session          bool               `json:"session"`
	SameSite         string             `json:"sameSite"`
	StoreId          string             `json:"storeId"`
	Value            string             `json:"value"`
}

func GetCookies(website string, userID string) ([]Cookie, error) {
	apiKey := os.Getenv("COOKIE_RELAY_API_KEY")
	url, err := url.Parse(os.Getenv("COOKIE_RELAY_URL"))
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, "cookies", website, userID)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	req.Header.Set("Cookie-Relay-API-Key", apiKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Got response status from cookie relay: %s", res.Status)
	}

	defer res.Body.Close()

	var cookies []Cookie
	err = json.NewDecoder(res.Body).Decode(&cookies)
	if err != nil {
		return nil, err
	}
	return cookies, nil
}

func GetCookieValueWithName(website string, userID string, name string) (*string, error) {
	cookies, err := GetCookies(website, userID)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		if cookie.Name == name {
			return &cookie.Value, nil
		}
	}
	return nil, fmt.Errorf("Could not find cookie with name '%s'", name)
}
