local colorscheme = "gruvbox"                                                                        -- 设置颜色方案名称
local ok, err = pcall(vim.cmd, string.format("colorscheme %s", colorscheme))                         -- 尝试加载颜色方案
if not ok then                                                                                       -- 如果加载失败
    vim.notify(string.format("Colorscheme '%s' not found: %s", colorscheme, err), vim.log.levels.ERROR) -- 发出错误通知并附带错误信息
    return
end
