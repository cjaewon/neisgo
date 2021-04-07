# Neisgo
🎓 Neisgo는 나이스 api를 기반으로 하는 전국 초중고 급식, 학사일정, 시간표 파싱 모듈 입니다.

## 설치
```sh
go get -u github.com/cjaewon/neisgo
```

## 예제
먼저 나이스 관련 데이터를 가져오기 위해 neis 인스턴스를 생성해야합니다.
```go
func main() {
  neis := neisgo.New(apiKey)
}
```
만약 apiKey가 공백으로 주어지면 open neis api가 제공하는 데이터가 제한될 수 있습니다. 

### 학교 설정하기
급식과 일정을 가져오기 위해서는 원하는 학교를 설정해야합니다.
첫번째 인자로는 시도교육청코드, 두번째 인자로는 표준학교코드를 전달해줍니다.
```go
func main() {
  neis := neisgo.New(apiKey)
  neis.Set("C10", "7150144")
}
```

### 급식 가져오기
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

meals의 타입은 `[]Meal` 와 같습니다.
```go
type mealTime struct {
  Breakfast string
  Lunch     string
  Dinner    string
}

type Meal struct {
  Date        time.Time
  // 급식 원산지
  Origin      mealTime
  // 급식 성분, 영양소
  Ingredients mealTime
  mealTime
}
```

## 일정 가져오기
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

calendars의 타입은 `[]Calendar` 와 같습니다.
```go
type Calendar struct {
  Date time.Time
  Name string
  // 일정에 대한 자세한 내용을 나타냅니다. (주로 공백이 주워집니다)
  Content string
  // 주야과정명 ("주간" 혹은 "야간")
  ClassTime string
  // 수업공제일명 ("휴업일", "공휴일" ...)
  Deduction string

  // 일정이 대상으로 학년을 나타냅니다.
  // 고등학교 또는 중학교 일 경우에는 3개만 사용됩니다. 
  Target [6]bool
}
```