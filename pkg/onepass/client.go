package onepass

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Client allows you to interact with 1password for a subdomain
type Client struct {
	Subdomain string
	session   string
}

// Login provides a username and password from 1password by parsing the
// 1password item and trying to extract the values
type Login struct {
	Username string
	Password string
}

type parsedItem struct {
	UUID    string `json:"uuid"`
	Details struct {
		Fields []struct {
			Designation string `json:"designation"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			Value       string `json:"value"`
		} `json:"fields"`
		Sections []struct {
			Title string `json:"title"`
		} `json:"sections"`
	} `json:"details"`
}

// NewClient helps get a client for a subdomain. This will fail if the user is
// not already signed in to the `op` command and with a valid session already
// established and configured.
func NewClient(subdomain string) (*Client, error) {
	op := &Client{
		Subdomain: subdomain,
		session:   os.Getenv(fmt.Sprintf("OP_SESSION_%s", subdomain)),
	}
	if op.session == "" {
		return nil, errors.New("failed to retrieve 1Password session info, run op signin and try again")
	}
	return op, nil
}

// GetLogin will return a login from the specified vault. The item and vault
// values can be either the name or the UUID. The UUID is recommended as it will
// not change even if the item name or vault name changes.
func (op *Client) GetLogin(item, vault string) (*Login, error) {

	out, err := op.runCmd("get", "item", item, fmt.Sprintf("--vault=%s", vault))
	if err != nil {
		return nil, err
	}
	return parseItemResponse(out)
}

func (op *Client) runCmd(args ...string) ([]byte, error) {
	args = append(args, fmt.Sprintf("--session=%s", op.session))
	cmd := exec.Command("op", args...)
	res, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("failed calling 1Password: %s\n%s", err, res)
	}
	return res, err
}

func parseItemResponse(res []byte) (*Login, error) {
	login := &Login{}
	pItem := parsedItem{}
	if err := json.Unmarshal(res, &pItem); err != nil {
		return nil, err
	}
	for _, field := range pItem.Details.Fields {
		if field.Designation == "username" {
			login.Username = field.Value
		} else if field.Designation == "password" {
			login.Password = field.Value
		}

	}
	return login, nil
}
