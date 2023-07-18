package imagetag

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"harborgetag/tools"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
)

type ImageTag struct {
	token       struct{ Token string }
	Tags        []string
	authService [3]string
	Scheme      string
	Username    string
	Password    string
	VerifySSL   bool
	Image       string
	Registry    string
	Order       string
	Filter      string
}

func (i *ImageTag) getToken() error {
	err := i.getAuthService()
	if err != nil {
		return fmt.Errorf("get auth service error: %v", err)
	}
	realm, service := i.authService[1], i.authService[2]
	url := realm + "?" + "service=" + service + "&scope=" + "repository:" + i.Image + ":pull"

	log.Printf("get token url is: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(i.Username, i.Password)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !(i.VerifySSL),
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("get token failed:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("get token failed:", err)
		return err
	}
	json.Unmarshal(body, &(i.token))
	if i.token.Token == "" {
		log.Println("get token failed:", string(body))
		return fmt.Errorf("%v", string(body))
	}
	return nil
}

func (i *ImageTag) getAuthService() error {
	url := i.Scheme + i.Registry + "/v2/"
	rtn := [3]string{}
	rtn[0] = "" // type
	rtn[1] = "" // realm
	rtn[2] = "" // service

	log.Printf("get auth service url is: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	headerValue := resp.Header.Get("Www-Authenticate")
	typePattern, _ := regexp.Compile(`^(\S+)`)
	typeMatcher := ""
	typeFind := typePattern.FindStringSubmatch(headerValue)
	if len(typeFind) != 0 {
		typeMatcher = typeFind[1]
	}
	if typeMatcher == "Bearer" {
		macherPattern, _ := regexp.Compile(`Bearer realm="(\S+)",service="([\S ]+)"`)
		matchFind := macherPattern.FindStringSubmatch(headerValue)
		if len(matchFind) != 0 {
			rtn[0] = "Bearer"
			rtn[1] = matchFind[1]
			rtn[2] = matchFind[2]
			log.Println("authService: type=Bearer, realm=" + rtn[1] + ", service=" + rtn[2])
		} else {
			log.Println("no authService available from ", url)
		}
	}
	i.authService = rtn
	return nil
}

func (i *ImageTag) GetImageTagsFromRegistry() error {
	err := i.getToken()
	if err != nil {
		return fmt.Errorf("get token error: %v", err)
	}
	url := i.Scheme + i.Registry + "/v2/" + i.Image + "/tags/list"

	log.Printf("get image tags url is: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer"+" "+i.token.Token)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !(i.VerifySSL),
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	json.Unmarshal(body, i)
	if len(i.Tags) == 0 {
		return fmt.Errorf(string(body))
	}
	i.Tags = tools.Filter(i.Filter, i.Tags)
	switch i.Order {
	case "Reverse-Natural-Ordering":
		func(s []string) {
			for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
				s[i], s[j] = s[j], s[i]
			}
		}(i.Tags)
	case "Descending-Versions":
		sort.Sort(tools.OrderStringsSlice(i.Tags))
	case "Ascending-Versions":
		sort.Sort(sort.Reverse(tools.OrderStringsSlice(i.Tags)))
	}
	return nil
}

// NewGetImageTag is the constructor that fills the ImageTag struct with default values .
func NewGetImageTag(username, password, register, image, order, filter string, verifySSL bool) *ImageTag {
	return &ImageTag{
		Username:  username,
		Password:  password,
		Registry:  register,
		Image:     image,
		Order:     order,
		Filter:    filter,
		VerifySSL: verifySSL,
	}
}
