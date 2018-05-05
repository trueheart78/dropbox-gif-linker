package dropbox

type Client struct {
}

// NewClient create a new Client for interacting with Dropbox
func NewClient() (c Client) {

	return c
}

func (c Client) apiURL() string {
	return "https://api.dropboxapi.com/2"
}

func (c Client) creationPath() string {
	return "sharing/create_shared_link_with_settings"
}

func (c Client) existingPath() string {
	return "sharing/list_shared_links"
}
