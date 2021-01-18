# go-monkey

Go言語で作るインタプリタ

## Getting Started

```sh
$ go run main.go
Hello atrn0!! This is Monkey REPL!
>> 3 + 4
7
>> 78 - 39 == 78 / 2
true
>> 888 == 39
false
>> (true == false) != (4 < 100)
true
>> (120 + 433) * 320
176960
>> if (1 > 2) { 1 + 2 } else { 1 - 2 }
-1
>> let add = fn(x, y) { x + y }
>> let sub = fn(x, y) { x - y }
>> let apply = fn(x, y, func) { func(x, y) }
>> apply(3, 4, add)
7
>> apply(10, 5, sub)
5
```


## Test

```sh
go test ./...
```
