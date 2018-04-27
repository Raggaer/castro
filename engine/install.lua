-- This file is executed at the installation process
-- The returning table will be set inside the Config.Custom value

local custom = {}

-- Online chart defaults
custom.OnlineChart = {
	Enabled = false,
	Interval = "1h",
	Display = 8,
}

return custom
