# data
Пакет содержит алгоритмы и структуры, связанные с хранением и обработкой данных (временные ряды, инфомрмация об активах и т.п.)


## candlestick.go
### `Period`
Структура данных, представляющия, логично, период времени (от какого-то до какого-то).

Структура очень легковесная: она весит всего **48 байт** и предстовляет собой массив из двух `time.Time`

- `ShiftedStart(shift time.Duration) Period` - Возвразает тот-же самый период, но со сдвинутым вперёд временем начала. Исходный период при этом не изменяется.


### `TimeStamp`

```go
type TimeStamp struct {
    timeframe TimeFrame
    Timestamp []time.Time
}
```

`TimeStamp` представляет временные метки для временных рядов.

#### Методы

- `Timeframe() TimeFrame`: Возвращает временной интервал, связанный с `TimeStamp`.
- `At(index int) time.Time`: Возвращает временную метку по индексу.
- `Extend(n int)`: Расширяет временные метки на указанное количество.
- `Append(moments ...time.Time)`: Добавляет временные метки.
- `sliceIdx(start, stop int) TimeStamp`: Возвращает подмножество `TimeStamp` с индексами от `start` до `stop`.
- `End() time.Time`: Возвращает последнюю временную метку.
- `Start() time.Time`: Возвращает первую временную метку.

### `Candle`

```go
type Candle struct {
    Open      float64
    High      float64
    Low       float64
    Close     float64
    Volume    float64
    TimeClose time.Time
}
```

`Candle` представляет свечу с открытием, максимумом, минимумом, закрытием и объемом торгов.

#### Конструктор

- `NewCandle(open, high, low, c, volume float64, timeClose time.Time) *Candle`: Создает новую свечу.

### `InstrumentCandle`

```go
type InstrumentCandle struct {
    Candle
    Instrument
}
```

`InstrumentCandle` представляет свечу с указанием инструмента, к которому она относится.

### `Chart`

```go
type Chart struct {
    Open      []float64
    High      []float64
    Low       []float64
    Close     []float64
    Volume    []float64
    Timestamp TimeStamp
}
```

`Chart` представляет временной ряд свечей и соответствующих временных меток.

#### Конструктор

- `RawChart(timeframe TimeFrame, capacity int) Chart`: Создает новый временной ряд.

#### Методы

- `Add(candle Candle)`: Добавляет свечу в временной ряд.
- `Slice(period Period) Chart`: Возвращает подмножество временного ряда в заданном периоде включительно с началом и концом.
- `Len() int`: Возвращает количество свечей во временном ряде.
- `CandleByIndex(index int) (*Candle, error)`: Возвращает свечу по индексу.


#### Возможные ошибки:

`errors.OutOfIndexError`: Индекс находится за пределами доступного диапазона свечей.

### `ChartContainer`

```go
type ChartContainer map[Instrument]Chart
```

`ChartContainer` представляет коллекцию временных рядов для различных инструментов.

#### Методы

- `ChartsByPeriod(period Period) ChartContainer`: Возвращает временные ряды в заданном периоде.
- `Candles() []InstrumentCandle`: Возвращает все свечи из всех временных рядов в порядке возрастания времени.
Сложность: $O(NK)$, где $K$ - количество инструментов.


## Функции

### `findIndexBeforeOrAtTime`

```go
func findIndexBeforeOrAtTime(series TimeStamp, moment time.Time) (int, error)
```

Представляет функцию для поиска индекса временной метки, ближайшей к указанному моменту времени или находящейся до него в пределах временного ряда.

Сложность: $O(\log_2N)$

#### Параметры

- `series TimeStamp`: слайс меток времени.
- `moment time.Time`: Момент времени, для которого ищется ближайшая временная метка.

#### Возвращаемое значение

- `int`: Индекс временной метки, ближайшей к указанному моменту времени.

#### Возможные ошибки

1. `errors.NewZeroLengthError("timestamp series")`: Возвращается, если временной ряд `TimeStamp` пуст. Это указывает на отсутствие временных меток для поиска.

2. `errors.ValueNotFoundError{}`: Возвращается, если момент времени находится перед началом временного ряда. Это указывает на то, что момент времени раньше, чем самый ранний момент во временном ряде.

## Примечание

Пакет `data` предоставляет удобные структуры данных для работы с временными рядами и свечами в финансовых данных. Код может использоваться в различных приложениях для анализа и обработки данных финансовых рынков.
