package service

import (
	"github.com/KateGF/Http-Server-Project-SO/core"
	"strconv"
)

// Calcula el n-ésimo número de Fibonacci utilizando recursión con memoización.
func Fibonacci(num int) int {
	// Casos base de la recursión.
	if num <= 0 {
		return 0
	}
	if num == 1 {
		return 1
	}

	a, b := 0, 1
	for i := 2; i <= num; i++ {
		a, b = b, a+b
	}

	return b
}

// Extrae el parámetro 'num' de la consulta, calcula el número de Fibonacci correspondiente y retorna la respuesta HTTP.
func FibonacciHandler(request *core.HttpRequest) (*core.HttpResponse, error) {
	// Obtiene el valor del parámetro 'num' de la URL query.
	numStr := request.Target.Query().Get("num")

	// Valida si el parámetro 'num' está presente.
	if numStr == "" {
		return core.BadRequest().Text("num is required"), nil
	}

	// Convierte el parámetro 'num' de string a entero.
	num, err := strconv.Atoi(numStr)

	// Valida si la conversión fue exitosa.
	if err != nil {
		return core.BadRequest().Text("num must be a number"), nil
	}

	// Valida si el número está dentro del rango permitido (0 a 92).
	// El límite 92 se debe a que Fibonacci(93) excede el máximo valor de int64.
	if num < 0 || num > 92 {
		return core.BadRequest().Text("num must be between 0 and 92"), nil
	}

	// Calcula el número de Fibonacci usando la función optimizada.
	v := Fibonacci(num)

	// Crea una respuesta HTTP 200 OK con el resultado como texto plano.
	return core.Ok().Text(strconv.Itoa(v)), nil
}
