local status, treesitter = pcall(require, "nvim-treesitter.configs")
if not status then
    vim.notify("not found nvim-treesitter")
    return
end

treesitter.setup({
    ensure_installed = {
        "lua",
        "rust",
        "go",
        "json",
        "yaml",
        "typescript",
        "java",
        "javascript",
        "html",
        "css",
        "vim",
        "thrift",
        "python",
        "proto",
        "gosum",
        "gomod",
        "markdown",
        "markdown_inline",
        "vue",
    },
    highlight = {
        enable = true,
        disable = {},
        additional_vim_regex_highlighting = false,
    },
    indent = {
        enable = true,
    },
    -- 启动折叠
    fold = {
        enable = true,
    },
})

vim.opt.foldtext = "v:lua.custom_fold_text()"

function _G.custom_fold_text()
    -- 获取当前行的内容
    local line = vim.fn.getline(vim.v.foldstart)
    return line
end

-- vim.opt.foldmethod = "manual"
vim.opt.foldmethod = "expr"
vim.opt.foldexpr = "nvim_treesitter#foldexpr()"
-- 默认不要折叠
-- https://stackoverflow.com/questions/8316139/how-to-set-the-default-to-unfolded-when-you-open-a-file
vim.opt.foldlevel = 99
