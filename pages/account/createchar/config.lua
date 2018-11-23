-- Valid town id list
app.Custom.ValidTownList = {1,2,3,4}

-- Valid vocation id list
app.Custom.ValidVocationList = {1, 2, 3, 4}

-- New character values
app.Custom.NewCharacterValues = {
	[1] = { -- no voc (vocation ID -1)
		level = 1,
	    experience = 0,
	    health = 150,
	    healthmax = 150,
	    mana = 0,
	    manamax = 0,
	    cap = 435,
	    soul = 0
	},
	[2] = { -- sorc
		level = 8,
	    experience = 4200,
	    health = 185,
	    healthmax = 185,
	    mana = 35,
	    manamax = 35,
	    cap = 470,
	    soul = 100
	},
	[3] = { -- druid
		inherit_from = 1, -- remove this line if you want to use custom values
	},
	[4] = { -- paladin
		inherit_from = 1
	},
	[5] = { -- knight
		inherit_from = 1
	}
}
