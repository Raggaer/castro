-- Reverse the given table
function reverse(tbl)
    for i=1, math.floor(#tbl / 2) do
        tbl[i], tbl[#tbl - i + 1] = tbl[#tbl - i + 1], tbl[i]
    end
end

-- Return application running absolute URL
function runningURL()
    local mode = "http"

    if app.SSL.Enabled then
        mode = "https"
    end

    return string.format("%s://%s:%s", mode, app.URL, app.Port)
end

-- Return plugin type string
function pluginTypeToString(plugin)
    if plugin == 1 then
        return "Page"
    elseif plugin == 2 then
        return "Widget"
    elseif plugin == 3 then
        return "Engine"
    elseif plugin == 4 then
        return "Template"
    end
end

function __genOrderedIndex( t )
    local orderedIndex = {}
    for key in pairs(t) do
        table.insert( orderedIndex, key )
    end
    table.sort( orderedIndex )
    return orderedIndex
end

function orderedNext(t, state)
    -- Rturns the keys in the alphabetic
    -- order. We use a temporary ordered key table that is stored in the
    -- table being iterated.

    local key = nil
    --print("orderedNext: state = "..tostring(state) )
    if state == nil then
        -- the first time, generate the index
        t.__orderedIndex = __genOrderedIndex( t )
        key = t.__orderedIndex[1]
    else
        -- fetch the next value
        for i = 1,table.getn(t.__orderedIndex) do
            if t.__orderedIndex[i] == state then
                key = t.__orderedIndex[i+1]
            end
        end
    end

    if key then
        return key, t[key]
    end

    -- no more value to return, cleanup
    t.__orderedIndex = nil
    return
end

function orderedPairs(t)
    -- Equivalent of the pairs() function on tables. Allows to iterate
    -- in order
    return orderedNext, t, nil
end

-- Trim leading and trailing whitespace from string
function string:trim()
   return self:match("^()%s*$") and "" or self:match("^%s*(.*%S)" )
end

-- Split a string into a table on separator
function string:split(sep)
   local fields, pattern = {}, string.format("([^%s]+)", sep)
   self:gsub(pattern, function(c) fields[#fields+1] = c:trim() end)
   return fields
end

-- Unmarshal json from file
function json:unmarshalFile(path)
    local f, e = io.open(path, "r")
    if not f then
        return nil, e
    end
    local t = self:unmarshal(f:read("*a"))
    f:close()
    return t
end

-- Explode a string into a table
function explode(d,p)
   local t, ll
   t={}
   ll=0
   if(#p == 1) then
      return {p}
   end
   while true do
      l = string.find(p, d, ll, true) -- find the next d in the string
      if l ~= nil then -- if "not not" found then..
        if string.sub(p,ll,l-1) ~= "" and string.sub(p,ll,l-1) ~= " " then
            table.insert(t, trim(string.sub(p,ll,l-1))) -- Save it in our array.
        end
        ll = l + 1 -- save just after where we found it for searching next time.
      else
        table.insert(t, trim(string.sub(p,ll))) -- Save what's left in our array.
        break -- Break at end, as it should be, according to the lua manual.
      end
   end
   return t
end

function trim(s)
  return s:match "^%s*(.-)%s*$"
end

-- A function in Lua similar to PHP's print_r, from http://luanet.net/lua/function/print_r

function print_r ( t ) 
    local print_r_cache={}
    local function sub_print_r(t,indent)
        if (print_r_cache[tostring(t)]) then
            print(indent.."*"..tostring(t))
        else
            print_r_cache[tostring(t)]=true
            if (type(t)=="table") then
                for pos,val in pairs(t) do
                    if (type(val)=="table") then
                        print(indent.."["..pos.."] => "..tostring(t).." {")
                        sub_print_r(val,indent..string.rep(" ",string.len(pos)+8))
                        print(indent..string.rep(" ",string.len(pos)+6).."}")
                    else
                        print(indent.."["..pos.."] => "..tostring(val))
                    end
                end
            else
                print(indent..tostring(t))
            end
        end
    end
    sub_print_r(t,"  ")
end
