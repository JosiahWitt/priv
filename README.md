# priv
> Privatize data in your Go (Golang) application

[![Build Status](https://travis-ci.com/JosiahWitt/priv.svg?branch=master)](https://travis-ci.com/JosiahWitt/priv)
[![codecov](https://codecov.io/gh/JosiahWitt/priv/branch/master/graph/badge.svg)](https://codecov.io/gh/JosiahWitt/priv)
[![Maintainability](https://api.codeclimate.com/v1/badges/a8d6b159cbfb58d2e4a3/maintainability)](https://codeclimate.com/github/JosiahWitt/priv/maintainability)
[![GoDoc](https://godoc.org/github.com/JosiahWitt/priv?status.svg)](https://godoc.org/github.com/JosiahWitt/priv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Under development; API may change.**

> go get github.com/JosiahWitt/priv

Tested with Go 1.9, 1.10, and 1.11.

## About

Currently `priv` exports a simple API, which is composed of `priv.ToMap` and `priv.ToMapErr`. `priv.ToMap` panics if a field that does not exist is requested, whereas `priv.ToMapErr` returns an error. Most of the time panicing is not a big deal, since it would indicate a typo in the code.

These functions allow mapping structs and arrays/slices of structs to `map[string]interface{}` and `[]map[string]interface{}{}`, respectively.
One example of using `priv` is if you only want to expose certain fields at an API endpoint.

## Example

```go
users := [User{ID: "123", Some: {Nested: {Field: "abc"}}, SomethingElse: true}]
v := priv.ToMap(users, "ID", "Some.Nested.Field->Renamed.Location")
// Thus, v = [{ID: "123", Renamed: {Location: "abc"}}]
```
