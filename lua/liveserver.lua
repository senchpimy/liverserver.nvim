local ffi = require("ffi")
local example = ffi.load("/home/plof/Documents/lua/liveserver/lua/example.so")

--local example = ffi.load("./example.so")
ffi.cdef([[
extern void StartServer();
extern void SendUpdate();
]]);

local M = { active = false }

function M.on_save()
  if M.active then
    example.SendUpdate()
  end
end

local buffer_to_string = function()
  local content = vim.api.nvim_buf_get_lines(0, 0, vim.api.nvim_buf_line_count(0), false)
  return table.concat(content, "\n")
end

local function serverStart()
  if not M.active then
    example.StartServer()
    M.active = true
    print("Live server started")
    -- Create the global autocommand to trigger on every buffer save
    vim.api.nvim_create_autocmd('BufWritePost', {
      group = vim.api.nvim_create_augroup("liveServerGroup", { clear = true }),
      pattern = '*', -- Matches all file patterns
      callback = M.on_save,
    })
  else
    print("Live server is already running")
  end
end

vim.api.nvim_create_user_command('Liveserver', function()
  serverStart()
end, {})

return M
