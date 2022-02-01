# TimeWheel

An implementation of Simple Timing Wheels.
![TimeWheel](./timewheel.jpeg)

## Features

- Set a unit of timer
- Remove a unit of timer

## Installing

To start using `tw`, install Go and run `go get`:

```sh
$ go get -u https://github.com/x-debug/tw
```

This will retrieve the library.

## Example
Set timer in time-wheel
```go
wheel := NewTimeWheel(1, 60)
defer wheel.StopTimer()

_ = wheel.SetTimer("timer1", 2*time.Second, func() {
    //Do anything
})
```

Remove timer in time-wheel
```go
wheel := NewTimeWheel(1, 60)
defer wheel.StopTimer()

//Do anything
_ = wheel.RemoveTimer("timer1")
```
## License

`tw` source code is available under the MIT [License](/LICENSE).