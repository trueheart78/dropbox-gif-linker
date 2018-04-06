# Dropbox Gif Linker

Designed to make working with your Dropbox gifs easier when wanting to share them.

## Usage

Download the respective binary for your system, open a terminal, and execute it.

### `gif-listener`

Listens for gif paths to be entered and checks with the local database before reaching out to 
create a new shareable link via the Dropbox API.

Data displayed includes the name of the id gif, its directories, basename, size, number of times 
used, and the data copied to the clipboard

![listener example](assets/images/listener-example.gif)

![taylor.gif](https://dl.dropboxusercontent.com/s/rhkozj2hwt82bc7/taylor.gif)

### Drag and Drop

Drag a file from your Dropbox gif directory into the terminal that your running this program in, and
press enter to have it present a shareable link. You can also drag and drop multiple gif files at once.

## :warning: Shared Gif Folder :warning:

If you have shared your Dropbox gifs directory at its root, this program will not work as expected.
