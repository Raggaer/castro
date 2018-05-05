package lua

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/clbanning/mxj"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/html/charset"
)

// SetXMLMetaTable sets the xml metatable of the given lua state
func SetXMLMetaTable(luaState *lua.LState) {
	// Create and set the xml metatable
	xmlMetaTable := luaState.NewTypeMetatable(XMLMetaTableName)
	luaState.SetGlobal(XMLMetaTableName, xmlMetaTable)

	// Set all xml metatable functions
	luaState.SetFuncs(xmlMetaTable, xmlMethods)
}

// MarshalXML marshals the given lua table
func MarshalXML(L *lua.LState) int {
	// Get table
	tbl := L.Get(2)

	// Check for valid table type
	if tbl.Type() != lua.LTTable {
		L.ArgError(1, "Invalid marshal object. Expected table")
		return 0
	}

	// Convert table to map
	mxj.XmlCharsetReader = charset.NewReaderLabel
	r := mxj.Map(TableToMap(L.ToTable(2)))

	// Marshal converted table
	buff, err := r.Xml()

	if err != nil {
		L.RaiseError("Cannot marshal the given table: %v", err)
		return 0
	}

	// Push result as string
	L.Push(lua.LString(string(buff)))

	return 1
}

// UnmarshalXMLFile unmarshals the given fle
func UnmarshalXMLFile(L *lua.LState) int {
	// Get path
	src := L.Get(2)

	// Check for valid string type
	if src.Type() != lua.LTString {
		L.ArgError(1, "Invalid unmarshal source. Expected string")
		return 0
	}

	// Check if file is already parsed
	xmlResult, found := util.Cache.Get(
		fmt.Sprintf("xml_table_%v", src.String()),
	)

	if found {

		// Push result as table
		L.Push(xmlResult.(*lua.LTable))

		return 1
	}

	// Read the whole file
	buff, err := ioutil.ReadFile(src.String())

	if err != nil {
		L.RaiseError("Cannot unmarshal file. File not found: %v", err)
		return 0
	}

	// Unmarshal string
	mxj.XmlCharsetReader = charset.NewReaderLabel
	result, err := mxj.NewMapXml(buff)

	if err != nil {
		L.RaiseError("Cannot unmarshal the given file: %v", err)
		return 0
	}

	// Convert result to table
	r := MapToTable(result)

	// Save result to cache
	util.Cache.Add(
		fmt.Sprintf("xml_table_%v", src.String()),
		r,
		util.Config.Configuration.Cache.Default.Duration,
	)

	// Push result as table
	L.Push(r)

	return 1
}

// UnmarshalXML unmarshals the given string to a lua table
func UnmarshalXML(L *lua.LState) int {
	// Get string
	src := L.Get(2)

	// Check for valid string type
	if src.Type() != lua.LTString {
		L.ArgError(1, "Invalid unmarshal source. Expected string")
		return 0
	}

	// Unmarshal string
	mxj.XmlCharsetReader = charset.NewReaderLabel
	result, err := mxj.NewMapXml([]byte(src.String()))

	if err != nil {
		L.RaiseError("Cannot unmarshal the given string: %v", err)
		return 0
	}

	// Push result as table
	L.Push(MapToTable(result))

	return 1
}

// GetVocationByName gets a vocation by the given name
func GetVocationByName(L *lua.LState) int {
	// Get name
	name := L.Get(2)

	// Check for valid name type
	if name.Type() != lua.LTString {

		L.ArgError(1, "Invalid name format. Expected string")
		return 0
	}

	// Check if vocation is on cache
	v, found := util.Cache.Get(
		fmt.Sprintf("vocation_%v", name.String()),
	)

	if found {

		// Push vocation table
		L.Push(v.(*lua.LTable))

		return 1
	}

	// Get vocation
	for _, voc := range util.ServerVocationList.List.Vocations {

		// If it is the vocation we are looking for
		if voc.Name == name.String() {

			// Convert vocation to table
			vocation := StructToTable(voc)

			// Save vocation on the cache
			util.Cache.Add(
				fmt.Sprintf("vocation_%v", voc.Name),
				vocation,
				time.Minute*3,
			)

			// Push vocation as lua table
			L.Push(vocation)

			return 1
		}
	}

	return 0
}

// GetVocationByID gets a vocation by the given id
func GetVocationByID(L *lua.LState) int {
	// Get ID
	id := L.Get(2)

	// Check for valid name type
	if id.Type() != lua.LTNumber {

		L.ArgError(1, "Invalid ID format. Expected number")
		return 0
	}

	// Get id as int
	idn := L.ToInt(2)

	// Get vocation
	for _, voc := range util.ServerVocationList.List.Vocations {

		// If it is the vocation we are looking for
		if voc.ID == idn {

			// Check if vocation is on the cache
			vocation, found := util.Cache.Get(
				fmt.Sprintf("vocation_%v", voc.Name),
			)

			if found {

				// Push vocation table
				L.Push(vocation.(*lua.LTable))

				return 1
			}

			// Convert vocation to table
			v := StructToTable(voc)

			// Save vocation to cache
			util.Cache.Add(
				fmt.Sprintf("vocation_%v", voc.Name),
				v,
				time.Minute*3,
			)

			// Push vocation as lua table
			L.Push(v)

			return 1
		}
	}

	return 0
}

// VocationList returns the server vocations xml file
func VocationList(L *lua.LState) int {
	// Check if user wants base vocations
	base := L.ToBool(2)

	// Build cache key
	cacheKeyName := "vocation_list"

	if base {
		cacheKeyName = "vocation_list_base"
	}

	// Check if base vocation list is on the cache
	b, found := util.Cache.Get(cacheKeyName)

	if found {

		// Push list table
		L.Push(b.(*lua.LTable))

		return 1
	}

	// Data holder
	result := &lua.LTable{}

	// Loop vocation list
	for _, vocation := range util.ServerVocationList.List.Vocations {

		// Convert vocation to table
		v := StructToTable(vocation)

		// Check if user wants base vocations
		if base {

			// Check if base vocation
			if vocation.ID == vocation.FromVoc {

				// Push vocation to table
				result.Append(v)
			}

			continue
		}

		// Push vocation to table
		result.Append(v)
	}

	// Add table to cache
	util.Cache.Add(cacheKeyName, result, time.Minute*3)

	// Push table to stack
	L.Push(result)

	return 1
}
