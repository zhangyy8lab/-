local status, luasnip = pcall(require, "luasnip")
if not status then
	vim.notify("not found luasnip")
	return
end

luasnip.setup({})
