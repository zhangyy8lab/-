-- 配置 nvim-cmp
local cmp = require('cmp')
local cmp_autopairs = require("nvim-autopairs.completion.cmp")
-- 插件需要导出的快捷键设置
local pluginKey = {}

cmp.setup({
  snippet = {
    expand = function(args)
      vim.fn["vsnip#anonymous"](args.body) -- For `vsnip` users.
      -- require("luasnip").lsp_expand(args.body) -- For `luasnip` users.
    end,
  },
  mapping = {
    ['<Tab>'] = cmp.mapping.select_next_item({ behavior = cmp.SelectBehavior.Insert }),
    ['<CR>'] = cmp.mapping.confirm({ select = true }),
  },

  formatting = require("lsp.ui").cmpFormatting,
  sources = {
    { name = 'nvim_lsp' },
    -- { name = "luasnip" },
  }

})

-- 代码补全tab选中后回车可用
cmp.event:on("confirm_done", cmp_autopairs.on_confirm_done())

-- 命令行补全的来源和映射方式
-- Use buffer source for `/` and `?` (if you enabled `native_menu`, this won't work anymore).
for _, v in pairs({ "/", "?" }) do
    cmp.setup.cmdline(v, {
        mapping = cmp.mapping.preset.cmdline(),
        sources = {
            { name = "buffer" },
        },
    })
end


-- cmp 代码补全
pluginKey.cmp = function(cmp)
	return {
		["<D-,>"] = cmp.mapping({
			i = cmp.mapping.abort(),
			c = cmp.mapping.close(),
		}),
		["<D-.>"] = cmp.mapping(cmp.mapping.complete(), { "i", "c" }),
		["<Down>"] = cmp.mapping.select_prev_item(),
		["<Up>"] = cmp.mapping.select_next_item(),
		-- 上一个 在一个
		["<C-p>"] = cmp.mapping.select_prev_item(),
		["<C-n>"] = cmp.mapping.select_next_item(),
		-- 确定
		["<CR>"] = cmp.mapping({
			i = function(fallback)
				if cmp.visible() and cmp.get_active_entry() then
					cmp.confirm({
						select = true,
						behavior = cmp.ConfirmBehavior.Replace,
					})
				else
					fallback() -- If you use vim-endwise, this fallback will behave the same as vim-endwise.
				end
			end,
			s = cmp.mapping.confirm({ select = true }),
			c = cmp.mapping.confirm({
				select = true,
				behavior = cmp.ConfirmBehavior.Replace,
			}),
		}),

		-- 如果窗口内容太多，可以滚动
		["<C-u>"] = cmp.mapping(cmp.mapping.scroll_docs(-4), { "i", "c" }),
		["<C-d>"] = cmp.mapping(cmp.mapping.scroll_docs(4), { "i", "c" }),
		-- tab 选择下一个
		["<Tab>"] = function(fallback)
			if cmp.visible() then
				cmp.select_next_item()
			else
				fallback()
			end
		end,
	}
end

return pluginKey

