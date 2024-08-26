-- 使用<leader>gr 查询方法被引用，选中其中一个后查询结果窗口进行消失
vim.cmd([[
  augroup QuickfixClose
    autocmd!
    autocmd QuickFixCmdPost *grep* cclose
    autocmd BufLeave * if &buftype == 'quickfix' | q | endif
  augroup END
]])
