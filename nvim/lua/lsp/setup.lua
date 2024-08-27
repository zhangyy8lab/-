local lspconfig = require("lspconfig")

require("mason").setup()
require("mason-lspconfig").setup({
    ensure_installed = { "pyright", "gopls", "jsonls", "yamlls", "lua_ls", "bashls" }, -- 确保这几个 LSP 服务器安装
})

-- 配置 Python LSP (pyright)
lspconfig.pyright.setup({})

-- 配置 Go LSP (gopls)
lspconfig.gopls.setup({
    settings = {
        gopls = {
            staticcheck = true
        }
    }
})

-- 配置 JSON LSP (jsonls)
lspconfig.jsonls.setup({})

-- 将补全插件连接到 LSP 服务器
local capabilities = require("cmp_nvim_lsp").default_capabilities()
lspconfig.pyright.setup({ capabilities = capabilities })
lspconfig.gopls.setup({ capabilities = capabilities })
lspconfig.jsonls.setup({ capabilities = capabilities })
lspconfig.yamlls.setup({ capabilities = capabilities })
lspconfig.lua_ls.setup({ capabilities = capabilities })
lspconfig.bashls.setup({ capabilities = capabilities })
