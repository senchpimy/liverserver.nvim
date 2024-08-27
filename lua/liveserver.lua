local ffi = require("ffi")
local example = ffi.load("/home/plof/Documents/lua/liveserver/lua/example.so")
ffi.cdef([[
extern void StartServer();
extern void SendUpdate();
]]);

function on_save()
  example.SendUpdate()
end

local buffer_to_string = function()
  local content = vim.api.nvim_buf_get_lines(0, 0, vim.api.nvim_buf_line_count(0), false)
  return table.concat(content, "\n")
end

local function main()
  print("Hello from our plugin")
  example.StartServer()
  local str = buffer_to_string()
  file = io.open("/home/plof/Documents/lua/liveserver/lua/test", "w")
  file:write(str)
  file:close()
end

local function setup()
  vim.api.nvim_create_autocmd("VimEnter",
    { group = augroup, desc = "Set a fennel scratch buffer on load", once = true, callback = main })
end

local group = vim.api.nvim_create_augroup("myGroup", { clear = true })
vim.api.nvim_create_autocmd('BufWritePost', {
  group = group,
  pattern = '*', -- You can specify a particular file pattern if needed
  callback = on_save,
})

return { setup = setup }
