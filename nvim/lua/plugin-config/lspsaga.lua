vim.cmd([[packadd lspsaga.nvim]])

local status, lspsaga = pcall(require, "lspsaga")
if not status then
    vim.notify("not found lspsaga")
    return
end

lspsaga.setup({
    diagaostic = {
        on_insert = true,
        on_insert_follow = true,
    },
    finder = {
        -- 是否在浮动窗口中显示finder的标题
        show_title = true,
        -- 主题颜色
        border_style = "rounded",
    },
})
