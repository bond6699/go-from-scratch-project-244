# Вычислитель отличий на Go

[![CI](https://github.com/bond6699/go-from-scratch-project-244/actions/workflows/CI.yml/badge.svg)](https://github.com/bond6699/go-from-scratch-project-244/actions/workflows/CI.yml)
[![Hexlet Check](https://github.com/bond6699/go-from-scratch-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/bond6699/go-from-scratch-project-244/actions/workflows/hexlet-check.yml)

**GenDiff** — консольная утилита на Go для сравнения двух конфигурационных файлов в форматах JSON и YAML. Программа находит добавленные, удалённые и изменённые свойства и выводит результат в удобном для чтения формате.

Проект выполнен в рамках программы [«Разработчик Go»](https://ru.hexlet.io/programs/go-from-scratch) школы Хекслет.

## Возможности

- Сравнение файлов JSON и YAML.
- Рекурсивное сравнение вложенных объектов.
- Определение добавленных, удалённых и изменённых свойств.
- Три формата вывода:
  - `stylish` — древовидное представление различий;
  - `plain` — список изменений в виде текста;
  - `json` — структурированный JSON.
- Обработка ошибок чтения и разбора файлов.

## Требования

- Go версии, указанной в `go.mod`, или совместимая более новая версия.
- Git — для клонирования репозитория.

## Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/bond6699/go-from-scratch-project-244.git
cd go-from-scratch-project-244
```

Соберите приложение:

```bash
make build
```

После успешной сборки исполняемый файл будет находиться в `bin/gendiff`.

При необходимости зависимости можно загрузить отдельно:

```bash
go mod download
```

## Использование

Общий синтаксис:

```text
gendiff [options] <filepath1> <filepath2>
```

Сравните два файла:

```bash
./bin/gendiff testdata/file1.json testdata/file2.json
```

По умолчанию используется формат `stylish`.

### Форматы вывода

**Stylish**

```bash
./bin/gendiff testdata/file1.json testdata/file2.json
```

Пример вывода:

```text
{
    additional: {
        one: 1
      - three: [
            1, 2, 4
            2
            hello from, GO
        ]
      + three: [
            1, 2
            2
            hello from, GO
        ]
        two: 2
    }
    follow: false
    host: hexlet.io
    proxy: 123.234.53.22
  - secret: null
  + secret: null
    timeout: 50
}
```

**Plain**

```bash
./bin/gendiff --format plain testdata/file1.json testdata/file2.json
```

Пример вывода:

```text
Property 'additional.three' was updated. From [1, 2, 4 2 hello from, GO] to [1, 2 2 hello from, GO]
Property 'secret' was updated. From null to 'null'
```

**JSON**

```bash
./bin/gendiff --format json testdata/file1.json testdata/file2.json
```

Пример вывода:

```json
{
   "child": {
      "additional": {
         "child": {
            "one": {
               "old": 1,
               "key": "one",
               "status": "Unchanged"
            },
            "three": {
               "old": [
                  "1, 2, 4",
                  2,
                  "hello from, GO"
               ],
               "new": [
                  "1, 2",
                  2,
                  "hello from, GO"
               ],
               "key": "three",
               "status": "Changed"
            },
            "two": {
               "old": 2,
               "key": "two",
               "status": "Unchanged"
            }
         },
         "key": "additional",
         "status": "Container"
      },
      "follow": {
         "old": false,
         "key": "follow",
         "status": "Unchanged"
      },
      "host": {
         "old": "hexlet.io",
         "key": "host",
         "status": "Unchanged"
      },
      "proxy": {
         "old": "123.234.53.22",
         "key": "proxy",
         "status": "Unchanged"
      },
      "secret": {
         "old": null,
         "new": "null",
         "key": "secret",
         "status": "Changed"
      },
      "timeout": {
         "old": 50,
         "key": "timeout",
         "status": "Unchanged"
      }
   },
   "key": "root",
   "status": "Container"
}
```

Для выбора формата также доступен короткий флаг `-f`:

```bash
./bin/gendiff -f plain testdata/file1.json testdata/file2.json
```

### Запуск через Make

Для запуска без предварительной сборки используйте `make run`:

```bash
make run ARGS="testdata/file1.json testdata/file2.json"
```

С указанием формата:

```bash
make run ARGS="-f plain testdata/file1.json testdata/file2.json"
```

## Разработка

Запуск тестов:

```bash
make test
```

Запуск линтера:

```bash
make lint
```

Исправление поддерживаемых линтером проблем:

```bash
make lint-fix
```

Сборка приложения:

```bash
make build
```

## Технологии

- **Go** — основной язык разработки.
- **urfave/cli** — создание CLI и обработка аргументов.
- **testify** — вспомогательные инструменты для тестирования.
- **gopkg.in/yaml.v3** — разбор YAML.
- **encoding/json** — работа с JSON.
- **golangci-lint** — статический анализ и проверка качества кода.
- **GitHub Actions** — автоматизация проверок.

## Структура проекта

```text
.
├── cmd/
│   └── gendiff/       # Точка входа приложения
├── internal/
│   ├── analyzer/      # Анализ различий между объектами
│   ├── app/           # Настройка CLI
│   ├── formatter/     # Форматирование результата
│   └── parser/        # Чтение и нормализация файлов
├── gendiff.go         # Публичная функция GenDiff
├── Makefile           # Команды разработки
├── go.mod
└── README.md
```

## Демонстрация

[![asciicast](https://asciinema.org/a/B10iyrpBMxRNYLEX.svg)](https://asciinema.org/a/B10iyrpBMxRNYLEX)