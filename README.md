# Dropbox Gif Linker

Designed to make working with your Dropbox gifs easier when wanting to share them.

## Dropbox Integration

First, you need to create a new [Dropbox app][dropbox-new-app], using the **Dropbox API** (not the 
business option), with **Full Dropbox** access. Once you have that setup, you will need to click 
the _Generate_ button beneath the **Generate Access Token** header of the **OAuth2** section. This 
is the token that will be used for interacting with your Dropbox account.

## Configuration

In your home directory, make sure to create `.dgl.json` file, and fill in the details:

```json
{
	"dropbox_path" : "~/Dropbox",
	"dropbox_gif_dir" : "gifs/",
	"dropbox_api_token" : "YOUR_API_TOKEN"
}
```

⚠️ The program will not load if you do not have this file setup correctly. All details are required.

## Usage

Download the respective binary for your system, open a terminal, and execute it.

### `dropbox-gif-linker`

Handles talking to the DropBox API for you.

![listener example](assets/images/listener-example.gif)

![taylor.gif](https://dl.dropboxusercontent.com/s/rhkozj2hwt82bc7/taylor.gif)

### Drag and Drop

Drag a file from your Dropbox gif directory into the terminal that your running this program in, and
press enter to have it present a shareable link. You can also drag and drop multiple gif files at once.

## :warning: Shared Gif Folder :warning:

If you have shared your Dropbox gifs directory at its root, this program will not work as expected.

[dropbox-new-app]: https://www.dropbox.com/developers/apps
