function string.parseBBCode(str)
	-- Remove all html tags
	str = str:gsub("<", "&lt;")
	str = str:gsub(">", "&gt;")

	-- Store noparse content and remove from processed string
	local np, i = {}, 0
	local function noParse(s)
		i = i + 1
		np[i] = s
		return string.format("NOPARSE-%s", i)
	end

	-- Parse BB tags
	str = str:gsub("%[right%](.-)%[/right%]", [[<div style="float:right;">%1</div>]])
	str = str:gsub("%[noparse%](.-)%[/noparse%]", noParse)
	str = str:gsub("%[url=(.-)%](.-)%[/url%]", [[<a href="%1" target="_BLANK">%2</a>]])
	str = str:gsub("%[url=(.-)%]%[/url%]", [[<a href="%1">%1</a>]])
	str = str:gsub("%[center%](.-)%[/center%]", [[<div style="text-align:center;">%1</div>]])
	str = str:gsub("%[quote%](.-)%[/quote%]", [[<blockquote><p>%1</p></blockquote>]])
	str = str:gsub("%[quote=(.-)%](.-)%[/quote%]", [[<blockquote><p>%2</p><footer><cite title="%1">%1</cite></footer></blockquote>]])
	str = str:gsub("%[size=(.-)%](.-)%[/size%]", [[<span style="font-size: %1;">%2</span>]])
	str = str:gsub("%[color=(.-)%](.-)%[/color%]", [[<span style="color:%1;">%2</span>]])
	str = str:gsub("%[font=(.-)%](.-)%[/font%]", [[<span style="font-family:'%1';">%2</span>]])
	str = str:gsub("%[del%](.-)%[/del%]", [[<del>%1</del>]])
	str = str:gsub("%[small%](.-)%[/small%]", [[<span style="font-size: x-small;">%1</span>]])
	str = str:gsub("%[large%](.-)%[/large%]", [[<span style="font-size: large;">%1</span>]])
	str = str:gsub("%[highlight%](.-)%[/highlight%]", [[<mark>%1</mark>]])
	str = str:gsub("%[[Bb]%](.-)%[/[Bb]%]", [[<strong>%1</strong>]])
	str = str:gsub("%[[Ii]%](.-)%[/[Ii]%]", [[<em>%1</em>]])
	str = str:gsub("%[[Uu]%](.-)%[/[Uu]%]", [[<span style="text-decoration:underline;">%1</span>]])
	str = str:gsub("%[img%](.-)%[/img%]", [[<img src="%1" class="img-responsive" alt="Image" />]])
	str = str:gsub("%[youtube%](.-)%[/youtube%]", [[<div class="embed-responsive embed-responsive-16by9"><iframe class="embed-responsive-item" src="https://www.youtube.com/embed/%1" allowfullscreen></iframe></div>]])

	--Replace noparse content
	local function replace(s)
		return np[tonumber(s)]
	end
	str = str:gsub("NOPARSE%-(%d+)", replace)

	return str
end