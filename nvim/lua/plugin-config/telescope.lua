local status, telescope = pcall(require, "telescope")
if not status then
	vim.notify("telescope.nvim not installed", vim.log.levels.WARN)
	return
end

telescope.setup({
	defaults = {
		initial_mode = "insert",
		layout_strategy = "horizontal",
		-- mappings = require("keymappings").telescopeList,
	},
	pickers = {
		find_files = {
			theme = "ivy",
		},
	},
	extensions = {
		fzf = {
			fuzzy = true,
			override_generic_sorter = true,
			overrode_file_sorter = true,
			case_mode = "smart_case",
		},
		media_files = {
			filetypes = { "png", "webp", "jpg", "jpeg" },
			find_cmd = "rg",
		},
	},
})

-- telescope extensions
pcall(telescope.load_extension, "ui-select")
pcall(telescope.load_extension, "env")
pcall(telescope.load_extension, "projects")
pcall(telescope.load_extension, "fzf")
pcall(telescope.load_extension, "media_files")

-- 插件需要导出的快捷键设置
local pluginKey = {}
pluginKey.telescopeList = {
	i = {
		["<Esc>"] = "close",
		-- 预览窗口上下滚动
		["<C-u>"] = "preview_scrolling_up",
		["<C-n>"] = "preview_scrolling_down",
		["<C-h>"] = "which_key",
	},
}
return pluginKey
