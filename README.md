![Castro AAC](https://i.gyazo.com/f328c60ee8c219b94a521e3e51fa66e7.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/Raggaer/castro)](https://goreportcard.com/report/github.com/Raggaer/castro)

High performance Open Tibia automatic account creator written in **Go** using **Lua** for the subtopics

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.


All pages and widgets (sidebar content) are done using the **Lua** bindings. This provides a great resource for newcomers to learn the Castro bindings. For more information the source code is available at [github](https://github.com/Raggaer/castro).

# Setup

Running Castro for the first time will generate a **config.toml** file. Thats the main configuration file for Castro.

Castro will then get everything it needs from your server folder. Reading your **config.lua** and other files (otbm, xml)

# Contact

Please if you find a bug head over to the [github issues](https://github.com/Raggaer/castro/issues) page. Feel free to ask anything about Castro too.

You can also contact the project maintainer at **nakotoffana@gmail.com**.

# License

**Castro** is made available under the **MIT** license.