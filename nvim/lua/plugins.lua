-- 插件管理
local packer = require("packer")

packer.startup({
    function(use)
        -- 包管理器
        use("wbthomason/packer.nvim")

        -- dashboard nvim 启动页配
        use({ "glepnir/dashboard-nvim", requires = { "nvim-tree/nvim-web-devicons" } })

        -- nvim 主题
        use({ "ellisonleao/gruvbox.nvim" })
        use("sheerun/vim-polyglot")

        -- project.nvim
        -- https://github.com/ahmedkhalf/project.nvim
        use("ahmedkhalf/project.nvim")

        -- 用于模糊查找和快速导航文件、缓冲区、Git 分支、符号、命令等内容
        -- telescope.nvim
        -- https://github.com/nvim-telescope/telescope.nvim
        use({ "nvim-telescope/telescope.nvim" })

        -- https://github.com/nvim-treesitter/nvim-treesitter/wiki/Installation
        use({ "nvim-treesitter/nvim-treesitter", run = ":TSUpdate" })

        -- folder tree
        use({ "kyazdani42/nvim-tree.lua" })
        use({ "kyazdani42/nvim-web-devicons" })

        -- 成对出现
        -- https://github.com/windwp/nvim-autopairs
        use("windwp/nvim-autopairs")

        -- bufferline
        -- using packer.nvim
        use({ "akinsho/bufferline.nvim", requires = { "moll/vim-bbye" } })

        -- comment.vim gcc 注释行
        -- https://github.com/numToStr/Comment.nvim
        use("numToStr/Comment.nvim")

        -- git blame 提示文件在什么时候被谁修改的， 看着不得劲
        -- https://github.com/f-person/git-blame.nvim
        -- use"f-person/git-blame.nvim")

        --  gitsigns.nvim 修改行 有对应提示
        -- https://github.com/lewis6991/gitsigns.nvim
        use("lewis6991/gitsigns.nvim")

        -- ui 图标
        use("onsails/lspkind-nvim")

        -- 下方系统图标
        use("nvim-lualine/lualine.nvim")

        -- 复制单词时高亮
        use("machakann/vim-highlightedyank")

        -- 指定单词高亮
        use("RRethy/vim-illuminate")

        use("lfv89/vim-interestingwords")
        -- https://github.com/jose-elias-alvarez/null-ls.nvim
        use("jose-elias-alvarez/null-ls.nvim")
        use("nvim-lua/plenary.nvim")

        -- go/python... lsp插件
        use("neovim/nvim-lspconfig")
        use("williamboman/mason.nvim")
        -- lsp installer
        use("williamboman/nvim-lsp-installer")
        use("williamboman/mason-lspconfig.nvim")
        use({ "L3MON4D3/LuaSnip", requires = { "saadparwaiz1/cmp_luasnip" } })
        -- nvim-cmp 补全引擎
        -- https://github.com/hrsh7th/nvim-cmp
        use("hrsh7th/nvim-cmp")
        -- 补全源
        use("hrsh7th/vim-vsnip")
        use("hrsh7th/cmp-nvim-lsp")          -- name = nvim_lsp
        use("hrsh7th/cmp-buffer")            -- name = buffer
        use("hrsh7th/cmp-path")              -- name = path
        use("hrsh7th/cmp-cmdline")           -- name = cmdline
        use("hrsh7th/cmp-vsnip")             -- 代码片段提示
        use("f3fora/cmp-spell")              -- 单词拼写
        use("hrsh7th/cmp-nvim-lsp-signature-help") -- 函数签名

        -- lspsaga
        -- https://github.com/glepnir/lspsaga.nvim
        use("glepnir/lspsaga.nvim")

        use("b0o/schemastore.nvim")
    end,

    config = {
        display = {
            open_fn = require("packer.util").float,
        },
        max_jobs = nil,
    },
})

-- 当该文件变化时，自动进行编译
vim.cmd([[	
    augroup packer_user_config
        autocmd!
        autocmd BufWritePost plugins.lua source <afile> | PackerSync
    augroup end
]])
