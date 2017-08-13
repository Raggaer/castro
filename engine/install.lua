-- This file is executed at the installation process
-- The returning table will be set inside the Config.Custom value

local custom = {}

-- Online chart defaults
custom.OnlineChart = {
	Enabled = false,
	Interval = "1h",
	Display = 8,
}

-- Menu pages
custom.MenuPages = {
	{
		Name = "Home",
		Dropdown = false,
		Link = "/"
	},
	{
		Name = "Library",
		Dropdown = true,
		DropdownItems = {
			{
				Link = "/subtopic/library/houses",
				Text = "House list"
			},
			{
				Link = "/subtopic/library/serverinfo",
				Text = "Server information"
			},
			{
				Link = "/subtopic/library/support",
				Text = "Support list"
			},
			{
				Link = "/subtopic/library/spells",
				Text = "Spell list"
			}
		}
	},
	{
		Name = "Community",
		Dropdown = true,
		DropdownItems = {
			{
				Link = "/subtopic/community/online",
				Text = "Who is online"
			},
			{
				Link = "/subtopic/community/highscores",
				Text = "Highscores"
			},
			{
				Link = "/subtopic/community/search",
				Text = "Search character"
			},
			{
				Link = "/subtopic/community/guilds/list",
				Text = "Guild list"
			},
			{
				Link = "/subtopic/community/guilds/wars",
				Text = "Guild wars"
			},
			{
				Link = "/subtopic/community/deaths",
				Text = "Latest deaths"
			}
		}
	},
	{
		Name = "Shop",
		Dropdown = true,
		DropdownItems = {
			{
				Link = "/subtopic/shop/view",
				Text = "House list"
			},
			{
				Link = "/subtopic/shop/paypal",
				Text = "Buypoints PayPal"
			},
			{
				Link = "/subtopic/shop/paygol",
				Text = "Buypoints PayGol"
			},
			{
				Link = "/subtopic/shop/fortumo",
				Text = "Buypoints Fortumo"
			}
		}
	}
}

return custom
