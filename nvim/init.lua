-- 初始化 packer.nvim 插件管理器

require("plugins")

require("basic")
require("colorscheme")

-- 快捷键映射
require("keymappings")

-- 面板
require("plugin-config.dashboard")

-- 用于模糊查找和快速导航文件、缓冲区、Git 分支、符号、命令等内容
require("plugin-config.telescope")

-- 成对字符 
require("plugin-config.autopairs")

-- 左侧文件管理器,文件树
require("plugin-config.nvim-tree")
require("plugin-config.nvim-treesitter")

-- 自定义配置
require("plugin-config.myset")

-- 复制单词时高亮
require("plugin-config.highlight")

-- 编辑器上方显示打开的buffer
require("plugin-config.bufferline")

-- 注释用 常规模式 gcc / 视觉模式下gc
-- require("plugin-config.comment")

-- buffer中新增/修改/删除 对应提示
require("plugin-config.git-signs")

-- 显示处在什么模式及下方图标信息
require("plugin-config.lualine")

-- 格式化代码
require("plugin-config.null-ls")

-- lsp
require("plugin-config.lspsaga")
require("lsp.cmp")
require("lsp.setup")
require("lsp.ui")

