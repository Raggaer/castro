package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func compileWidgetList(req *http.Request, w http.ResponseWriter, sess map[string]interface{}) ([]template.HTML, error) {
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

		// Get state
		state, err := WidgetList.Get(filepath.Join("widgets", widget.Name, widget.Name+".lua"))

		if err != nil {
			return nil, err
		}

		defer WidgetList.Put(state, filepath.Join("widgets", widget.Name, widget.Name+".lua"))

		// Set session user data
		SetSessionMetaTableUserData(state, sess)

		// Execute widget
		tbl, useCache, err := widget.Execute(state)

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
