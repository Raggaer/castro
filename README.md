# Castro AAC

[![Join the chat at https://gitter.im/castro-aac/Lobby](https://badges.gitter.im/castro-aac/Lobby.svg)](https://gitter.im/castro-aac/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/Raggaer/castro)](https://goreportcard.com/report/github.com/Raggaer/castro)
[![LICENSE](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/Raggaer/castro/blob/master/LICENSE)
[![Build status](https://ci.appveyor.com/api/projects/status/yhrx9l6jrbvxhw5p?svg=true)](https://ci.appveyor.com/project/Raggaer/castro)

High performance Open Tibia content management system written in **Go** using **Lua** for the scripting part.

Castro provides lua bindings. Using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.

## Documentation

Castro also ships with the documentation so you can view it offline. Everything is located under **/en-US/** directory.

The documentation is also hosted at [castroaac.org](https://castroaac.org).

## Extensions

Castro ships with a very solid extension system. You can read more about it on the extensions part of the documentation.

There is a public plugin list hosted at [plugins.castroaac.org](https://plugins.castroaac.org)

## License

**Castro** is licensed under the **MIT** license.
