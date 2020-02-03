package vrcapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	VrcApiBase         string = "https://api.vrchat.cloud/api/1/"
	CurrentUserDetails string = VrcApiBase + "auth/user"
	APIKey             string = "JlE5Jldo5Jibnk5O5hTx6XVqsJu4WJ26"
	FriendRequest      string = VrcApiBase + "user/%s/friendRequest"
)

type UserDetails struct {
	AcceptedTOSVersion             int           `json:"acceptedTOSVersion"`
	ActiveFriends                  []string      `json:"activeFriends"`
	AllowAvatarCopying             bool          `json:"allowAvatarCopying"`
	Bio                            string        `json:"bio"`
	BioLinks                       []interface{} `json:"bioLinks"`
	CurrentAvatar                  string        `json:"currentAvatar"`
	CurrentAvatarAssetURL          string        `json:"currentAvatarAssetUrl"`
	CurrentAvatarImageURL          string        `json:"currentAvatarImageUrl"`
	CurrentAvatarThumbnailImageURL string        `json:"currentAvatarThumbnailImageUrl"`
	DeveloperType                  string        `json:"developerType"`
	DisplayName                    string        `json:"displayName"`
	Email                          string        `json:"email"`
	EmailVerified                  bool          `json:"emailVerified"`
	Feature                        struct {
		TwoFactorAuth bool `json:"twoFactorAuth"`
	} `json:"feature"`
	FriendGroupNames       []string `json:"friendGroupNames"`
	FriendKey              string   `json:"friendKey"`
	Friends                []string `json:"friends"`
	HasBirthday            bool     `json:"hasBirthday"`
	HasEmail               bool     `json:"hasEmail"`
	HasLoggedInFromClient  bool     `json:"hasLoggedInFromClient"`
	HasPendingEmail        bool     `json:"hasPendingEmail"`
	HomeLocation           string   `json:"homeLocation"`
	ID                     string   `json:"id"`
	IsFriend               bool     `json:"isFriend"`
	LastLogin              string   `json:"last_login"`
	LastPlatform           string   `json:"last_platform"`
	ObfuscatedEmail        string   `json:"obfuscatedEmail"`
	ObfuscatedPendingEmail string   `json:"obfuscatedPendingEmail"`
	OculusID               string   `json:"oculusId"`
	OfflineFriends         []string `json:"offlineFriends"`
	OnlineFriends          []string `json:"onlineFriends"`
	PastDisplayNames       []struct {
		DisplayName string `json:"displayName"`
		UpdatedAt   string `json:"updated_at"`
	} `json:"pastDisplayNames"`
	State             string `json:"state"`
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
	SteamDetails      struct {
		Avatar                   string `json:"avatar"`
		Avatarfull               string `json:"avatarfull"`
		Avatarmedium             string `json:"avatarmedium"`
		Commentpermission        int    `json:"commentpermission"`
		Communityvisibilitystate int    `json:"communityvisibilitystate"`
		Lastlogoff               int    `json:"lastlogoff"`
		Personaname              string `json:"personaname"`
		Personastate             int    `json:"personastate"`
		Personastateflags        int    `json:"personastateflags"`
		Primaryclanid            string `json:"primaryclanid"`
		Profilestate             int    `json:"profilestate"`
		Profileurl               string `json:"profileurl"`
		Steamid                  string `json:"steamid"`
		Timecreated              int    `json:"timecreated"`
	} `json:"steamDetails"`
	SteamID              string   `json:"steamId"`
	Tags                 []string `json:"tags"`
	TwoFactorAuthEnabled bool     `json:"twoFactorAuthEnabled"`
	Unsubscribe          bool     `json:"unsubscribe"`
	Username             string   `json:"username"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetUserDetails(username string, password string) string {
	url := CurrentUserDetails
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}

func ExportFriends(username string, password string) []string {
	userDetailsData := GetUserDetails(username, password)
	var userDetails UserDetails
	json.Unmarshal([]byte(userDetailsData), &userDetails)
	return userDetails.Friends
}

func SendFriendRequest(username string, password string, usr_id string) {
	url := fmt.Sprintf(FriendRequest, usr_id)
	method := "POST"

	payload := strings.NewReader("{\n	\"apiKey\": \"" + APIKey + "\"\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	res, err := client.Do(req)
	defer res.Body.Close()
}
