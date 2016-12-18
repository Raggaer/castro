# Castro

High performance Open Tibia automatic account creator written in **Go** using **LUA** for the custom pages

# How

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool


All pages are done using the LUA bindings. This provides a great resource for newcomers to learn the Castro bindings. For more information the source code is available at [github](https://github.com/Raggaer/castro)

# Contact

Please if you find a bug head over to the [github issues](https://github.com/Raggaer/castro/issues) page. Feel free to ask anything about Castro too

You can also contact the project maintainer at **nakotoffana@gmail.com**