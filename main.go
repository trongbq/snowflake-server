package main

import (
  "snowflake-server/snowflake"
)


func main() {
    server, err := snowflake.NewServer(1)
    if err != nil {
        panic(err)
    }
    server.Start()
}

