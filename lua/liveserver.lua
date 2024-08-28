local ffi = require("ffi")
local example = ffi.load(script_path() .. "example.so")

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

local function serverStart()
  if not M.active then
    example.StartServer()
    M.active = true
    print("Running Live Server")
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

function script_path()
  local str = debug.getinfo(2, "S").source:sub(2)
  return str:match("(.*/)")
end

return M
