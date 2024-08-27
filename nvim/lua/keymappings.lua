-- 定义 set key map 的快捷函数
local map = vim.api.nvim_set_keymap
local opt = { noremap = true, silent = true }

-- 插件需要导出的快捷键设置
-- local pluginKey = {}

-- 设置常用的快捷的前缀 key 为空格
vim.g.mapleader = " "
vim.g.maplocalleader = " "

-- 取消原来 s 快捷键的功能
map("n", "s", "", opt)
-- 水平分割窗口
map("n", "sl", ":vsplit<CR>", opt)

-- 垂直分割窗口
map("n", "sh", ":split<CR>", opt)
-- 关闭当前窗口
map("n", "sc", "<C-w>c", opt)
-- 关闭其他窗口
-- map("n", "so", "<C-w>o", opt)

-- 切换窗口
map("n", "<C-j>", "<C-w>j", opt)
map("n", "<C-k>", "<C-w>k", opt)
map("n", "<C-h>", "<C-w>h", opt)
map("n", "<C-l>", "<C-w>l", opt)

-- 退出
map("n", "ww", ":w<CR>", opt)
map("n", "rl", ":e!<CR>", opt)
map("n", "q", ":Bdelete!<CR>", opt)
map("n", "qq", ":q<CR>", opt)
map("n", "wq", ":w<CR> | :Bdelete!<CR>", opt)

-- 搜索后， 取消颜色内容
map("n", "<Esc>", ":nohlsearch<CR><Esc>", opt)

-- tree 文件树左侧快捷键设置
map("n", "tr", ":NvimTreeToggle<CR>", opt)

-- 根据文件名查询文件
map("n", "<C-p>", ":Telescope find_files<CR>", opt)
-- 根据字符查询在哪些文件中出现过
map("n", "<C-f>", ":Telescope live_grep<CR>", opt)

-- git commit history
map("n", "gitc", ":Telescope git_commits<CR>", opt)
-- git status
map("n", "gits", ":Telescope git_status<CR>", opt)
-- map("n", "<C-;>", ":Telescope project<CR>", opt)

-- git add . and git commit -m ``
map("n", "gam", ":lua GitAddCommit()", opt)
-- git push
map("n", "gp", ":lua GitPush()", opt)

-- bufferline  左右切换 tab
map("n", "<Tab>h", ":BufferLineCyclePrev<CR>", opt)
map("n", "<Tab>", ":BufferLineCycleNext<CR>", opt)

-- treesitter 代码块折叠
map("n", "zz", ":foldclose<CR>", opt)
map("n", "zf", ":foldopen<CR>", opt)

-- lsp 快捷键
-- 跳转到该方法定义的位置
map("n", "gd", "<cmd>Lspsaga goto_definition<CR>", opt)

-- 列出该方法被哪些调用
map("n", "gr", "<cmd>lua vim.lsp.buf.references()<CR><cmd>cclose<CR>", opt)

-- 代码格式化
map("n", "ff", "<cmd>lua vim.lsp.buf.format({ async = true })<CR><cmd>cclose<CR>", opt)

-- 设置快捷键
map("n", "iw", [[:lua Highlight_word()<CR>]], opt)
map("n", "iwc", [[:lua Clear_highlight()<CR>]], opt)


-- map("n", "gi", "<cmd>lua vim.lsp.buf.implementation()<CR>", opt)
-- 非自定义g开头的快捷键 使用:nmap xx， 查询xx定义的快捷键
-- gc/gcc 注释行，数字+gcc 表示注释以下几行
-- gx Opens filepath or URI under cursor with the system handler (file explorer, web browser, …)

-- pluginKey.lspList = function(bufnr)
-- 	bufmap(bufnr, "n", "gi", "<cmd>lua vim.lsp.buf.implementation()<CR>", opt)
-- end

-- go to
-- bufmap(bufnr, "n", "gr", "<cmd>lua vim.lsp.buf.references()<CR><cmd>cclose<CR>", opt)
-- bufmap(bufnr, "n", "gr", "<cmd>Lspsaga lsp_finder<CR>", opt)
--
-- bufmap(bufnr, "n", "grn", "<cmd>Lspsaga rename<CR>", opt)
-- bufmap(bufnr, "n", "<leader>ca", "<cmd>Lspsaga cade_action<CR>", opt)
--
--
-- bufmap(bufnr, "n", "gD", "<cmd>lua vim.lsp.buf.declaration()<CR>", opt)

-- bufmap(bufnr, "n", "gh", "<cmd>lua vim.lsp.buf.hover()<CR>", opt)
-- bufmap(bufnr, "n", "gh", "<cmd>Lspsaga hover_doc<CR>", opt)


-- map("n", "gt", "<cmd>Lspsaga goto_type_definition<CR>", opt)
-- diagnostic
-- bufmap(bufnr, "n", "go", "<cmd>lua vim.diagnostic.open_float()<CR>", opt)
-- bufmap(bufnr, "n", "go", "<cmd>Lspsaga show_line_diagnostics<CR>", opt)

-- bufmap(bufnr, "n", "gn", "<cmd>lua vim.diagnostic.goto_next()<CR>", opt)
-- bufmap(bufnr, "n", "gn", "<cmd>Lspsaga diagnostic_jump_next<CR>", opt)

-- bufmap(bufnr, "n", "gp", "<cmd>lua vim.diagnostic.goto_prev()<CR>", opt)
-- bufmap(bufnr, "n", "gp", "<cmd>Lspsaga diagnostic_jump_prev<CR>", opt)
-- bufmap(bufnr, "n", "gf", "<cmd>lua vim.lsp.buf.format({ async = true })<CR>", opt)
--
-- bufmap(bufnr, "n", [[<M-\>]], "<cmd>Lspsaga term_toggle<CR>", opt)
-- bufmap(bufnr, "t", [[<M-\>]], "<cmd>Lspsaga term_toggle<CR>", opt)
-- end


-- 修改窗口大小
-- map("n", "<M-Left>", ":vertical resize -2<CR>", opt)
-- map("n", "<M-Right>", ":vertical resize +2<CR>", opt)
-- map("n", "<M-Up>", ":resize -2<CR>", opt)
-- map("n", "<M-Down>", ":resize +2<CR>", opt)

-- map("n", "s,", ":vertical res+ze -20<CR>", opt)
-- map("n", "s.", ":vertical resize +20<CR>", opt)
-- map("n", "sk", ":resize +20<CR>", opt)
-- map("n", "sj", ":resize -20<CR>", opt)
-- map("n", "s=", "<C-w>=", opt)

-- insert mormal 模式下的快捷键 跳转到行首和行尾
-- map("i", "<C-a>", "<ESC>I", opt)
-- map("i", "<C-e>", "<ESC>A", opt)
-- map("n", "<C-a>", "<ESC>I", opt)
-- map("n", "<C-e>", "<ESC>A", opt)

-- terminal 终端快捷键设置
-- map("t", "li", ":LspInstallInfo<CR>", opt)
-- map("n", "<leader>t", ":sp | terminal<CR>", opt)
-- map("n", "<leader>vt", ":vsp | terminal<CR>", opt)

-- visual 模式下的快捷键
-- map("v", "<", "<gv", opt)
-- map("v", ">", ">gv", opt)
--
-- }
--     { key = "<F5>",                           action = "refresh" },
--     { key = "a",                              action = "create" },
--     { key = "d",                              action = "remove" },
--     { key = "r",                              action = "rename" },
--     { key = "x",                              action = "cut" },
--     { key = "c",                              action = "copy" },
--     { key = "p",                              action = "paste" },
--     { key = "<leader>",                       action = "system_open" },
-- }
