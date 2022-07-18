package gdrive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sut-storage-go/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type DriveHandler struct {
	Config *config.Config
}

func (d DriveHandler) NewDriveService() (*drive.Service, error) {
	conf, _ := config.LoadConfig()

	ctx := context.Background()
	b, err := ioutil.ReadFile(conf.ServiceAccount)
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, err
	}
	client := getClient(config, d.Config.TokenPath)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (d *DriveHandler) RegenerateToken() (*oauth2.Token, error) {
	var client = &http.Client{}
	payload := map[string]string{
		"client_id":     d.Config.ClientId,
		"client_secret": d.Config.ClientSecret,
		"refresh_token": d.Config.GdriveRefreshToken,
		"grant_type":    "refresh_token",
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://accounts.google.com/o/oauth2/token?access_type=offline", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	tok := &oauth2.Token{}
	err = json.NewDecoder(resp.Body).Decode(tok)

	if err != nil {
		return nil, err
	}

	return tok, nil
}
