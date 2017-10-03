package lua

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// setWidgetMetaTable sets the widget metatable to the given state
func setWidgetMetaTable(luaState *lua.LState) {
	// Create and set the widget metatable
	widgetMetaTable := luaState.NewTypeMetatable(WidgetMetaTableName)
	luaState.SetGlobal(WidgetMetaTableName, widgetMetaTable)

	// Set all captcha metatable functions
	luaState.SetFuncs(widgetMetaTable, widgetMethods)
}

// RenderWidgetTemplate renders the given widget template
func RenderWidgetTemplate(L *lua.LState) int {
	// Get template name
	templateName := L.Get(2)

	// Check valid template type
	if templateName.Type() != lua.LTString {
		L.ArgError(1, "Invalid widget template. Expected string")
		return 0
	}

	// Get template args
	tableArgs := L.ToTable(3)

	// Get http fields
	req, _ := getRequestAndResponseWriter(L)

	// Render widget template
	buff, err := util.WidgetTemplate.RenderWidget(req, templateName.String(), TableToMap(tableArgs))

	if err != nil {
		L.RaiseError("Cannot parse widget template: %v", err)
		return 0
	}

	// Set widget template user data
	templateData := L.NewUserData()
	templateData.Value = template.HTML(buff.String())

	// Get main metatable
	tbl := L.GetTypeMetatable(WidgetMetaTableName)

	// Set data field
	L.SetField(tbl, "__data", templateData)

	return 0
}

func compileWidgetList(req *http.Request, w http.ResponseWriter, sess map[string]interface{}) ([]template.HTML, error) {
	// Data holder
	results := []template.HTML{}

	// Loop widget list
	for _, widget := range util.Widgets.List {

		// Get widget state from the pool
		state, err := WidgetList.Get(filepath.Join("widgets", widget.Name, widget.Name+".lua"))

		if err != nil {
			return nil, err
		}

		// Return state
		defer WidgetList.Put(state, filepath.Join("widgets", widget.Name, widget.Name+".lua"))

		// Set widget metatable
		setWidgetMetaTable(state)

		// Set widget HTTP metatable
		SetWidgetHTTPMetaTable(state)

		// Set HTTP user data
		SetHTTPUserData(state, w, req)

		// Set session user data
		SetSessionMetaTableUserData(state, sess)

		// Call widget function
		if err := ExecuteControllerPage(state, "widget"); err != nil {
			return nil, err
		}

		// Get widget metatable
		tbl := state.GetTypeMetatable(WidgetMetaTableName)

		// Get widget user data
		widgetData := state.GetField(tbl, "__data")

		// Convert data to user data
		data, ok := widgetData.(*lua.LUserData)

		if !ok {
			return nil, errors.New("Cannot convert widget data to user data")
		}

		// Convert data value to template HTML
		templateData, ok := data.Value.(template.HTML)

		if !ok {
			return nil, errors.New("Cannot convert widget user data to template HTML")
		}

		// Append data
		results = append(results, templateData)
	}

	return results, nil
}
