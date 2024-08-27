-- 使用<leader>gr 查询方法被引用，选中其中一个后查询结果窗口进行消失
vim.cmd([[
  augroup QuickfixClose
    autocmd!
    autocmd QuickFixCmdPost *grep* cclose
    autocmd BufLeave * if &buftype == 'quickfix' | q | endif
  augroup END
]])

-- 定义Git add . 和 git commit -m 函数
function GitAddCommit()
  -- 获取当前文件的路径
  local filepath = vim.fn.expand("%:p")
  -- 执行 git add 命令
  vim.cmd("!git add " .. filepath)
  -- 提示用户输入 commit 信息
  local commit_message = vim.fn.input("Commit message: ")
  -- 执行 git commit 命令
  vim.cmd('!git commit -m "' .. commit_message .. '"')
end

function GitPush()
  -- 执行 git push 命令
  vim.cmd("!git push")
end
