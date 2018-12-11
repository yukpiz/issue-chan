package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var GITHUB_ISSUE_URL = "https://api.github.com/repos/yukpiz/private/issues"
var SLACK_POST_URL = "https://slack.com/api/chat.postMessage"

type GithubIssue struct {
	URL           string `json:"url"`
	RepositoryURL string `json:"repository_url"`
	LabelsURL     string `json:"labels_url"`
	CommentsURL   string `json:"comments_url"`
	EventsURL     string `json:"events_url"`
	HTMLURL       string `json:"html_url"`
	ID            int    `json:"id"`
	NodeID        string `json:"node_id"`
	Number        int    `json:"number"`
	Title         string `json:"title"`
	User          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
	Labels   []interface{} `json:"labels"`
	State    string        `json:"state"`
	Locked   bool          `json:"locked"`
	Assignee struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"assignee"`
	Assignees []struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"assignees"`
	Milestone         interface{} `json:"milestone"`
	Comments          int         `json:"comments"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ClosedAt          interface{} `json:"closed_at"`
	AuthorAssociation string      `json:"author_association"`
	Body              string      `json:"body"`
}

type SlackAttachment struct {
	Color   string `json:"color,omitempty"`
	PreText string `json:"pretext,omitempty"`
	Text    string `json:"text,omitempty"`
}

func Handler() error {
	issues, err := GetGithubIssues()
	if err != nil {
		return err
	}
	fmt.Printf("Issues: %+v\n", issues)

	atts, err := CreateSlackMessage(issues)
	if err != nil {
		return err
	}

	fmt.Printf("Slack Request: %+v\n", atts)
	err = PostSlackMessage(atts)
	if err != nil {
		return err
	}

	return nil
}

func GetGithubIssues() (*[]GithubIssue, error) {
	req, err := http.NewRequest(http.MethodGet, GITHUB_ISSUE_URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_AUTH_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	issues := []GithubIssue{}
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, err
	}
	return &issues, nil
}

func CreateSlackMessage(issues *[]GithubIssue) (*[]SlackAttachment, error) {
	atts := []SlackAttachment{}

	for _, issue := range *issues {
		atts = append(atts, SlackAttachment{
			Color: "good",
			Text:  issue.Title,
		})
	}
	return &atts, nil
}

func PostSlackMessage(atts *[]SlackAttachment) error {
	params, err := json.Marshal(atts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, SLACK_POST_URL, nil)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Add("token", os.Getenv("SLACK_AUTH_TOKEN"))
	v.Add("channel", "rikka-chan")
	v.Add("text", "@yukpiz\n今日のタスク一覧だよ〜\n今日も1日がんばれ〜！\nhttps://github.com/yukpiz/private#boards")
	v.Add("icon_url", "https://i.gyazo.com/874720c0b9fc4f5b05688714c68e1f1d.jpg")
	v.Add("link_names", "true")
	v.Add("username", "六花")
	v.Add("attachments", string(params))

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = v.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Resp: %s\n", string(b))
	return nil
}

func main() {
	lambda.Start(Handler)
}
