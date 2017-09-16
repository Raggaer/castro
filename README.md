# Castro AAC

[![Go Report Card](https://goreportcard.com/badge/github.com/Raggaer/castro)](https://goreportcard.com/report/github.com/Raggaer/castro)
[![lICENSE](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/Raggaer/castro/blob/master/LICENSE)
[![Build status](https://ci.appveyor.com/api/projects/status/yhrx9l6jrbvxhw5p?svg=true)](https://ci.appveyor.com/project/Raggaer/castro)

High performance Open Tibia content management system written in **Go** using **Lua** scripting

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.

## Documentation

Everything you might need [is here](https://docs.castroaac.org/). The official documentation site of Castro

## License

**Castro** is available under the **MIT** license.
