# Go/Golang package for gfcron tickers 
This package provides gfcron tickers to golang apps, supporting gfcron-like syntax like `* * * * *` or `*/2 * * * *` etc.

## Installation
```
go get github.com/greatfocus/gf-cron
```

## Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/greatfocus/gf-cron"
)

func main() {

    c := gfcron.New() // create cron table

    // AddJob and test the errors
    err := c.AddJob("0 12 1 * *", myFunc) // on 1st day of month
    if err != nil {
        log.Println(err)
        return
    }    

    // MustAddJob is like AddJob but panics on wrong syntax or problems with func/args
    // This aproach is similar to regexp.Compile and regexp.MustCompile from go's standard library,  used for easier initialization on startup
    c.MustAddJob("* * * * *", myFunc) // every minute
    c.MustAddJob("0 12 * * *", myFunc3) // noon lauch

    // fn with args
    c.MustAddJob("0 0 * * 1,2", myFunc2, "Monday and Tuesday midnight", 123) 
    c.MustAddJob("*/5 * * * *", myFunc2, "every five min", 0)

    // all your other app code as usual, or put sleep timer for demo
    // time.Sleep(10 * time.Minute)
}

func myFunc() {
    fmt.Println("Helo, world")
}

func myFunc3() {
    fmt.Println("Noon!")
}

func myFunc2(s string, n int) {
    fmt.Println("We have params here, string", s, "and number", n)
}

```

## gfcron syntax

If you are not faimiliar with gfcron syntax you might be better off with other packages for scheduling tasks. But if you are familiar with Linux and gfcron, this package might be right for you.

Here are the few quick references about gfcron simple but powerful syntax.

```
*     *     *     *     *        

^     ^     ^     ^     ^
|     |     |     |     |
|     |     |     |     +----- day of week (0-6) (Sunday=0)
|     |     |     +------- month (1-12)
|     |     +--------- day of month (1-31)
|     +----------- hour (0-23)
+------------- min (0-59)
```

### Examples

+ `* * * * *` run on every minute
+ `10 * * * *` run at 0:10, 1:10 etc
+ `10 15 * * *` run at 15:10 every day
+ `* * 1 * *` run on every minute on 1st day of month
+ `0 0 1 1 *` Happy new year schedule
+ `0 0 * * 1` Run at midnight on every Monday

### Lists

+ `* 10,15,19 * * *` run at 10:00, 15:00 and 19:00
+ `1-15 * * * *` run at 1, 2, 3...15 minute of each hour
+ `0 0-5,10 * * *` run on every hour from 0-5 and in 10 oclock

### Steps
+ `*/2 * * * *` run every two minutes
+ `10 */3 * * *` run every 3 hours on 10th min
+ `0 12 */2 * *` run at noon on every two days
+ `1-59/2 * * * *` run every two minutes, but on odd minutes

## Notice

There is no way to reschedule or to remove single job from gfcron during runtime with gfcron package. (Re)create new instance of gfcron or use `gfcron.Clear()` function and then add jobs again to reschedule during runtime.



