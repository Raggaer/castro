![Castro AAC](https://i.gyazo.com/7a9f5b334221456a6c9f80c3274b917c.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/Raggaer/castro)](https://goreportcard.com/report/github.com/Raggaer/castro)
[![lICENSE](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/Raggaer/castro/blob/master/LICENSE)
[![Build status](https://ci.appveyor.com/api/projects/status/yhrx9l6jrbvxhw5p?svg=true)](https://ci.appveyor.com/project/Raggaer/castro)

High performance Open Tibia automatic account creator written in **Go** using **Lua** scripting

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.


All pages and widgets (sidebar content) are done using the **Lua** bindings. This provides a great resource for newcomers to learn the Castro bindings. For more information the source code is available at [github](https://github.com/Raggaer/castro).

# Pages and widgets

To create a custom page head to the `pages` folder and create a new directory with your page name. The name is equivalent to `/subtopic/:name`. GET requests are mapped to the `get.lua` file and POST request to `post.lua`
 
To create a custom widget head to the `widgets` folder and create a new directory with your widget name. Castro will look for `:name.html` and `:name.lua` files inside your widget directory

**For detailed information check out the project wiki**

# Setup

Running Castro for the first time will generate a **config.toml** file. Thats the main configuration file for Castro.

Castro will then get everything it needs from your server folder. Reading your **config.lua** and other files (otbm, xml)

# Compiling

Check out our wiki to see how to compile Castro. You can however get the latest commit build at [appveyor](https://ci.appveyor.com/project/Raggaer/castro/build/artifacts) for Linux-amd64 and Windows-amd64 only

# License

**Castro** is made available under the **MIT** license.