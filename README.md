# Neisgo
๐ Neisgo๋ ๋์ด์ค api๋ฅผ ๊ธฐ๋ฐ์ผ๋ก ํ๋ ์ ๊ตญ ์ด์ค๊ณ  ๊ธ์, ํ์ฌ์ผ์ , ์๊ฐํ ํ์ฑ ๋ชจ๋ ์๋๋ค.

## ์ค์น
```sh
go get -u github.com/cjaewon/neisgo
```

## ์์ 
๋จผ์  ๋์ด์ค ๊ด๋ จ ๋ฐ์ดํฐ๋ฅผ ๊ฐ์ ธ์ค๊ธฐ ์ํด neis ์ธ์คํด์ค๋ฅผ ์์ฑํด์ผํฉ๋๋ค.
```go
func main() {
  neis := neisgo.New(apiKey)
}
```
๋ง์ฝ apiKey๊ฐ ๊ณต๋ฐฑ์ผ๋ก ์ฃผ์ด์ง๋ฉด open neis api๊ฐ ์ ๊ณตํ๋ ๋ฐ์ดํฐ๊ฐ ์ ํ๋  ์ ์์ต๋๋ค. 

### ํ๊ต ์ค์ ํ๊ธฐ
๊ธ์๊ณผ ์ผ์ ์ ๊ฐ์ ธ์ค๊ธฐ ์ํด์๋ ์ํ๋ ํ๊ต๋ฅผ ์ค์ ํด์ผํฉ๋๋ค.
์ฒซ๋ฒ์งธ ์ธ์๋ก๋ ์๋๊ต์ก์ฒญ์ฝ๋, ๋๋ฒ์งธ ์ธ์๋ก๋ ํ์คํ๊ต์ฝ๋๋ฅผ ์ ๋ฌํด์ค๋๋ค.
```go
func main() {
  neis := neisgo.New(apiKey)
  neis.Set("C10", "7150144")
}
```

### ๊ธ์ ๊ฐ์ ธ์ค๊ธฐ
```go
func main() {
  neis := neisgo.New(apiKey)
  neis.Set("C10", "7150144")

  meals, err := neis.GetMeal(2021, 1)
  if err != nil {
    panic(err)
  }
}
```

meals์ ํ์์ `[]Meal` ์ ๊ฐ์ต๋๋ค.
```go
type mealTime struct {
  Breakfast string
  Lunch     string
  Dinner    string
}

type Meal struct {
  Date        time.Time
  // ๊ธ์ ์์ฐ์ง
  Origin      mealTime
  // ๊ธ์ ์ฑ๋ถ, ์์์
  Ingredients mealTime
  mealTime
}
```

## ์ผ์  ๊ฐ์ ธ์ค๊ธฐ
```go
func main() {
  neis := neisgo.New(apiKey)
  neis.Set("C10", "7150144")

  calendars, err := neis.GetCalendar(2021, 1)
  if err != nil {
    panic(err)
  }
}
```

calendars์ ํ์์ `[]Calendar` ์ ๊ฐ์ต๋๋ค.
```go
type Calendar struct {
  Date time.Time
  Name string
  // ์ผ์ ์ ๋ํ ์์ธํ ๋ด์ฉ์ ๋ํ๋๋๋ค. (์ฃผ๋ก ๊ณต๋ฐฑ์ด ์ฃผ์์ง๋๋ค)
  Content string
  // ์ฃผ์ผ๊ณผ์ ๋ช ("์ฃผ๊ฐ" ํน์ "์ผ๊ฐ")
  ClassTime string
  // ์์๊ณต์ ์ผ๋ช ("ํด์์ผ", "๊ณตํด์ผ" ...)
  Deduction string

  // ์ผ์ ์ด ๋์์ผ๋ก ํ๋์ ๋ํ๋๋๋ค.
  // ๊ณ ๋ฑํ๊ต ๋๋ ์คํ๊ต ์ผ ๊ฒฝ์ฐ์๋ 3๊ฐ๋ง ์ฌ์ฉ๋ฉ๋๋ค. 
  Target [6]bool
}
```