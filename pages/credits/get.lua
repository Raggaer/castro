local contributors_url = "https://api.github.com/repos/Raggaer/castro/contributors"
local commits_url = "https://api.github.com/repos/Raggaer/castro/commits?per_page=7"

function get()
    local data = cache:get("github")
    if not data then
        data = {}
    	data.contributors = http:get(contributors_url)
        data.commits = http:get(commits_url)
        cache:set("github", data, "30m")
    end

    if data.contributors then
        data.contributors = json:unmarshal(data.contributors).object
    end

    if data.commits then
        data.commits = json:unmarshal(data.commits).object
        for i, commit in pairs(data.commits) do
            local message = commit.commit.message:sub(0, commit.commit.message:find("\n\n", 1))
            if message:len() > 45 then
                message = message:sub(0, 45) .. "..."
            end
            data.commits[i].commit.message = message
        end
    end

    http:render("credits.html", data)
end
