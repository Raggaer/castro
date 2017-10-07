---
name: Custom pages
---

# Custom pages

Castro uses lua to handle all your pages. To create a new page navigate to pages directory and create a new folder with your subtopic name. The page will be accessed as `/subtopic/:name`. You can use multiple levels of directories

`pages/community/view` can be accessed as `/subtopic/community/view`

On your new folder you can then create a `get.lua` (to handle GET requests) or `post.lua` (to handle POST requests)

Each file must contain a function with the method they correspond to. `get.lua` files should have the `function get()` and `post.lua` files should have the `function post()`. Below is an example on how a simple page will look like:

```lua
require "paginator"
require "bbcode"

function get()
    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/index")
        return
    end

    local articleCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_articles", true)
    local pg = paginator(page, 5, tonumber(articleCount.total))
    local data = {}

    data.articles, cache = db:query("SELECT title, text, created_at FROM castro_articles ORDER BY id DESC LIMIT ?, ?", pg.limit, pg.offset, true)
    data.paginator = pg

    if data.articles == nil and page > 0 then
        http:redirect("/subtopic/index")
        return
    end

    if data.articles ~= nil then
        if not cache then
            for _, article in pairs(data.articles) do
                article.text = article.text:parseBBCode()
            end
        end
    end

    http:render("home.html", data)
end
```

This example loads the server latest news from database, we are also using the paginator function to paginate the results. This function takes the current page, number of results per page and the total count of items. Then we load the news from database or from cache, you can modify the table and it will modify the cache value also.

Lastly we render our template and pass the data we want.

```html
{{ template "header.html" . }}
{{ if .articles }}
    {{ range $index, $element := .articles }}
    <div class="news-box">
        <h3>
            <a href="#">{{ $element.title }}</a>
            <small class="small-info">{{ unixToDate $element.created_at }}</small>
        </h3>
        <hr>
        <p>
            {{ nl2br $element.text }}
        </p>
    </div>
    {{ end }}
    <ul class="pagination pagination-sm">
        {{ if .paginator.prev }}
        <li><a href="{{ url "subtopic" "index" }}?page={{ .paginator.firstpage.num }}">First</a></li>
        <li><a href="{{ url "subtopic" "index" }}?page={{ .paginator.prevnumber }}">&lt;</a></li>
        {{ end }}
        </li>
        {{ if $.paginator.last }}
        <li><a href="{{ url "subtopic" "index" }}?page={{ .paginator.lastnumber }}">&gt;</a></li>
        <li><a href="{{ url "subtopic" "index" }}?page={{ .paginator.lastpage.num }}">Last</a></li>
        {{ end }}
    </ul>
{{ else }}
<div class="news-box">
    <p>There are no articles</p>
</div>
{{ end }}
{{ template "footer.html" . }}
```

We first include our header and our footer, check if the articles variable is not nil then loop that variable to show the news. We also add pagination buttons at the end using the paginator variable.