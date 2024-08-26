-- local status, hight = pcall(require, "vim-illuminate")
-- if not status then
--     vim.notify("not found vim-illuminate")
--     return
-- end
--
-- hight.setup({})

-- 定义高亮组
local function define_highlights()
    local colors = { "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White" }
    for i, color in ipairs(colors) do
        vim.cmd(string.format("highlight WordHighlight%d guifg=%s guibg=NONE", i, color))
    end
end

-- 高亮光标下的单词
local function highlight_word()
    local word = vim.fn.expand("<cword>")
    if word == "" then
        return
    end

    -- 随机选择一个高亮组
    local group_id = math.random(1, 7)
    local highlight_group = string.format("WordHighlight%d", group_id)

    -- 清除以前的高亮
    -- vim.cmd("match none")

    -- 应用新的高亮
    vim.cmd(string.format("match %s /\\<%s\\>/", highlight_group, word))
end

-- 绑定快捷键
-- vim.api.nvim_set_keymap("n", "<leader>hw", [[:lua highlight_word()<CR>]], { noremap = true, silent = true })

-- 在启动时定义高亮组
define_highlights()
-- highlight_word()
