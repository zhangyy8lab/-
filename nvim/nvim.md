# nvim 



## Install node oh-my-zsh

```bash
# node  
brew install node 

# oh-my-zsh 
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

```



## oh-my-zsh setting

```bash
export ZSH="$HOME/.oh-my-zsh"

export PATH=/opt/homebrew/bin:$PATH
export PATH=$PATH:$(go env GOPATH)/bin

ZSH_THEME="agnoster"

plugins=(
     zsh-autosuggestions
     autojump
     git
     jsontools
)

# 
unsetopt inc_append_history
unsetopt share_history
```



## Download packer Manager

```bash
# https://github.com/wbthomason/packer.nvim
git clone --depth 1 https://github.com/wbthomason/packer.nvim\
  ~/.local/share/nvim/site/pack/packer/start/packer.nvim
```



### font&icon 

```bash
https://www.nerdfonts.com/font-downloads
```



## use 

```bsh
 ~/.config/nvim# tree
.
├── init.lua
├── lua
│   ├── basic.lua
│   ├── colorscheme.lua
│   ├── keymappings.lua
│   ├── lsp
│   │   ├── cmp.lua
│   │   ├── setup.lua
│   │   └── ui.lua
│   ├── plugin-config
│   │   ├── autopairs.lua
│   │   ├── bufferline.lua
│   │   ├── dashboard.lua
│   │   ├── git-blame.lua
│   │   ├── git-signs.lua
│   │   ├── highlight.lua
│   │   ├── lspsaga.lua
│   │   ├── lualine.lua
│   │   ├── luasnip.lua
│   │   ├── mason.lua
│   │   ├── myset.lua
│   │   ├── null-ls.lua
│   │   ├── nvim-tree.lua
│   │   ├── nvim-treesitter.lua
│   │   ├── project.lua
│   │   └── telescope.lua
│   └── plugins.lua
└── testgo
    ├── go.mod
    └── main.go

4 directories, 26 files
```



