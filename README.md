# LiveServer.nvim

**LiveServer.nvim** is a Neovim plugin that allows you to serve files from your current directory using a Go server. The server automatically reloads the contents of the page whenever you save your file, making it ideal for web development.

## Features

- **Instant Live Preview**: Automatically serves and reloads files in your current directory on every save.
- **Lightweight and Fast**: Powered by a Go server, ensuring high performance.
- **Easy to Use**: Start the server with a simple `:Liveserver` command.
- **Cross-Platform**: Works seamlessly on all platforms supported by Neovim.

## Requirements

- Neovim 0.5.0 or higher.
- Go installed on your machine.

## Installation

### Using [lazy.nvim](https://github.com/folke/lazy.nvim)

Add the following line to your Neovim configuration:

```lua
{
  'senchpimy/liveserver.nvim',
  config = function()
      require('liveserver')
  end,
}
```

Then, sync your plugins with:

```vim
:Lazy sync
```

## Usage

Once installed, you can start the server with the `:Liveserver` command:

```vim
:Liveserver
```

This will start a Go server in the current directory, serving all files and automatically reloading the content when you save any file.

### Default Behavior

- The server runs on `http://localhost:2324`.
- Files in the current working directory are served.
- The browser automatically reloads when any file is saved.

## License

This plugin is licensed under the GPLv3 License.
