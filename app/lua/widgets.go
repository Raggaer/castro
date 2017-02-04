package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"html/template"
	"net/http"
	"time"
)

func compileWidgetList(luaState *lua.LState, req *http.Request) ([]template.HTML, error) {
	// Data holder
	results := []template.HTML{}

	// Loop widget list
	for _, widget := range util.Widgets.List {

		// Check if widget is on the cache
		buff, cache := widget.IsCached()

		// If widget is cached append
		if cache {

			// Append to result slice
			results = append(results, buff)

			continue
		}

		// Get widget source
		source, err := Widgets.Get("widgets", widget.Name, widget.Name)

		if err != nil {
			return nil, err
		}

		// Execute widget
		tbl, useCache, err := widget.Execute(luaState, source)

		if err != nil {
			return nil, err
		}

		// Execute widget template
		buffer, err := util.WidgetTemplate.RenderWidget(req, widget.Name+".html", TableToMap(tbl))

		if err != nil {
			return nil, err
		}

		// Transform result to template HTML
		widgetResult := template.HTML(buffer.String())

		// Append result
		results = append(results, widgetResult)

		if useCache {

			// Add widget to cache
			util.Cache.Add(
				fmt.Sprintf("widget_%v", widget.Name),
				widgetResult,
				time.Minute*3,
			)
		}
	}

	return results, nil
}
