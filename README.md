# Pipeline - Нелинейный асинхронный конвейер выполнения

Пакет для создания и выполнения нелинейных конвейеров обработки данных на основе JSON конфигурации.

## Особенности

-  **Передача данных в JSON формате** между блоками
-  **Нелинейные конвейеры** с поддержкой ветвления и циклов
-  **Асинхронное выполнение** параллельных веток
-  **Условные переходы** для динамического управления потоком
-  **Простая архитектура** без сложных зависимостей

## Структура блока

```json
{
  "name": "имя_блока",
  "function": "имя_функции",
  "input": {
    "параметр1": "значение_или_ссылка_на_блок",
    "параметр2": 42
  },
  "output": ["следующий_блок1", "следующий_блок2"]
}
```

## Примеры использования

### 1. Простое условное ветвление (`pipeline1.json`)

```json
[
  {
    "name": "start",
    "function": "add",
    "input": { "a": 10, "b": 5 },
    "output": ["check"]
  },
  {
    "name": "check", 
    "function": "conditional",
    "input": { "value": "start" },
    "output": ["small_path", "big_path"]
  }
]
```

**Схема выполнения:**
```
start -> check -> small_path (если < 100)
               -> big_path   (если >= 100)
                     |
                   final
```

### 2. Параллельные ветки (`pipeline2.json`)

**Схема выполнения:**
```
input -> branch1 (async) \
      -> branch2 (async) -> merge
```

### 3. Циклический конвейер (`pipeline_cycle.json`)

**Схема выполнения:**
```
init -> counter -> multiplier -> condition
         ^                         |
         |-- continue_loop <-------|
                   |
                finish (конец)
```

## API

### Основные функции

```go
// Парсинг JSON в структуры блоков
func Parse(input string) ([]Block, error)

// Выполнение конвейера
func Run(blocks []Block, funcs map[string]Function)
```

### Типы данных

```go
type Block struct {
    Name     string                 `json:"name"`
    Function string                 `json:"function"`
    Input    map[string]interface{} `json:"input"`
    Output   []string               `json:"output"`
}

type Function func(input map[string]interface{}) interface{}
```

## Запуск демонстрации

```bash
go run main.go
```

## Архитектура

- **Рекурсивное выполнение** с проверкой уже выполненных блоков
- **sync.Map** для безопасного concurrent доступа к результатам
- **Условные блоки** выбирают одну из веток на основе boolean результата  
- **Обычные блоки** запускают все output-блоки параллельно
- **Предотвращение дедлоков** через неблокирующее разрешение зависимостей

## Ограничения

- Циклические конвейеры могут выполняться бесконечно без внешних ограничений
- Строковые значения в `input` интерпретируются как ссылки на блоки
- Функции должны быть зарегистрированы заранее

## Примеры функций

```go
funcs := map[string]pipeline.Function{
    "add": func(input map[string]interface{}) interface{} {
        a := input["a"].(float64)
        b := input["b"].(float64)
        return a + b
    },
    "conditional": func(input map[string]interface{}) interface{} {
        value := input["value"].(float64)
        return value < 100
    },
}
```