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
pluginKey.config = function(config)
	return {
		["<D-,>"] = config.mapping({
			i = config.mapping.abort(),
			c = config.mapping.close(),
		}),
		["<D-.>"] = config.mapping(config.mapping.complete(), { "i", "c" }),
		["<Down>"] = config.mapping.select_prev_item(),
		["<Up>"] = config.mapping.select_next_item(),
		-- 上一个 在一个
		["<C-p>"] = config.mapping.select_prev_item(),
		["<C-n>"] = config.mapping.select_next_item(),
		-- 确定
		["<CR>"] = config.mapping({
			i = function(fallback)
				if config.visible() and config.get_active_entry() then
					config.confirm({
						select = true,
						behavior = config.ConfirmBehavior.Replace,
					})
				else
					fallback() -- If you use vim-endwise, this fallback will behave the same as vim-endwise.
				end
			end,
			s = config.mapping.confirm({ select = true }),
			c = config.mapping.confirm({
				select = true,
				behavior = config.ConfirmBehavior.Replace,
			}),
		}),

		-- 如果窗口内容太多，可以滚动
		["<C-u>"] = config.mapping(config.mapping.scroll_docs(-4), { "i", "c" }),
		["<C-d>"] = config.mapping(config.mapping.scroll_docs(4), { "i", "c" }),
		-- tab 选择下一个
		["<Tab>"] = function(fallback)
			if config.visible() then
				config.select_next_item()
			else
				fallback()
			end
		end,
	}
end

return pluginKey

