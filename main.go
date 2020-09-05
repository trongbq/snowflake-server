package main

import (
    "os"
    "strconv"
    "snowflake-server/snowflake"
)

func main() {
    mID := os.Getenv("MACHINE_ID")
    if len(mID) == 0 {
        panic("Machine ID is required")
    }
    id, err := strconv.Atoi(mID)
    if err != nil {
        panic(err)
    }
    server, err := snowflake.NewServer(id)
    if err != nil {
        panic(err)
    }
    server.Start()
}
