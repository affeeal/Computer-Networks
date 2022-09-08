package proto

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать три значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "add" - передача нового прямоугольника на сервер;
	// * "find" - просьба найти прямоугольник с наибольшей площадью.
	Command string `json:"command"`

	// Если Command == "add", в поле Data должен лежать прямоугольник
	// в виде структуры Rect.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit" или "add";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - прямоугольник с наибольшей площадью найден.
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	// Если Status == "result", в поле Data должен лежать прямоугольник
	// в виде структуры Rect.
	// В противном случае, поле Data пустое.
	Data *json.RawMessage `json:"data"`
}

// Rect -- прямоугольник
type Rect struct {
	// Координаты текущей вершины прямоугольника.
	X1 string `json:"x1"`
	Y1 string `json:"y1"`

	// Координаты противоположной вершины прямоугольника.
	X2 string `json:"x2"`
	Y2 string `json:"y2"`

    // Площадь прямоугольника.
	Area float64
}