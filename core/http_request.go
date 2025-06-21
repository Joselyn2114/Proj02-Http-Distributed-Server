package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

var HTTP_METHODS = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
	"PATCH":   true,
}

// Representa una solicitud HTTP recibida.
// Contiene el método, el objetivo (URL), las cabeceras y el cuerpo de la solicitud.
type HttpRequest struct {
	Method  string            // Método HTTP (GET, POST, etc.)
	Target  *url.URL          // URL objetivo de la solicitud
	Headers map[string]string // Cabeceras HTTP como un mapa de clave-valor
	Body    string            // Cuerpo de la solicitud (si existe)
}

// Crea una nueva instancia de HttpRequest.
func NewHttpRequest(method string, target *url.URL, header map[string]string, body string) *HttpRequest {
	return &HttpRequest{
		Method:  method,
		Target:  target,
		Headers: header,
		Body:    body,
	}
}

// Lee una solicitud HTTP completa desde una conexión de red.
// Devuelve un puntero a HttpRequest o un error si ocurre algún problema.
func ReadRequest(conn net.Conn) (*HttpRequest, error) {
	lines := make([]string, 0)

	reader := bufio.NewReader(conn)

	// Lee las líneas de la cabecera hasta encontrar una línea vacía
	for {
		line, err := reader.ReadString('\n')

		// Si es fin de archivo (EOF), puede ser normal si la conexión se cierra
		if errors.Is(err, io.EOF) {
			break
		}

		// Otro error durante la lectura
		if err != nil {
			return nil, err
		}

		// Elimina los sufijos de retorno de carro y nueva línea
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")

		// Una línea vacía indica el final de las cabeceras
		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	// Si no se leyeron líneas, la solicitud está vacía
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	// Une las líneas de la cabecera y añade el doble salto de línea final
	headersPart := strings.Join(lines, "\r\n") + "\r\n\r\n"

	// Parsea la parte de las cabeceras para obtener la estructura HttpRequest inicial
	request, err := ParseRequest(headersPart)
	if err != nil {
		return nil, fmt.Errorf("can't parse request: %w", err)
	}

	// Si el método es POST, exige Content-Length
	if request.Method == "POST" {
		if _, ok := request.Headers["Content-Length"]; !ok {
			return nil, fmt.Errorf("post request without content length")
		}
	}
	
	// Parsea el cuerpo de la solicitud si Content-Length existe (o body vacío en otro caso)
	if err := ParseBody(request, reader); err != nil {
		return nil, err
	}

	return request, nil
}

// Analiza la parte de las cabeceras de una solicitud HTTP (como string).
// Devuelve un puntero a HttpRequest (sin el cuerpo) o un error.
func ParseRequest(headersPart string) (*HttpRequest, error) {
	// Encuentra el final de las cabeceras (doble salto de línea CRLF)
	headerEndIndex := strings.Index(headersPart, "\r\n\r\n")
	if headerEndIndex == -1 {
		return nil, fmt.Errorf("no header end")
	}

	// Extrae la parte de las cabeceras
	headerPart := headersPart[:headerEndIndex]

	// Divide las cabeceras en líneas individuales usando CRLF
	lines := strings.Split(headerPart, "\r\n")

	// Parsea la línea de inicio (ej: "GET /path HTTP/1.0")
	start := strings.SplitN(lines[0], " ", 3)
	if len(start) < 2 {
		return nil, fmt.Errorf("no method or target")
	}

	// Extrae el método
	method := start[0]

	// Comprueba si el método es válido
	if _, ok := HTTP_METHODS[method]; !ok {
		return nil, fmt.Errorf("bad method: %s", method)
	}

	// Extrae el target (URL)
	targetStr := start[1]

	// Parsea el target como una URL
	target, err := url.Parse(targetStr)
	if err != nil {
		return nil, fmt.Errorf("bad target format: %w", err)
	}

	// Extrae la versión
	version := start[2]

	// Comprueba si la versión es válida
	if version != "HTTP/1.0" && version != "HTTP/1.1" {
		return nil, fmt.Errorf("bad version: %s", version)
	}

	// Crea un mapa para almacenar las cabeceras
	headers := make(map[string]string)

	// Procesa cada línea de cabecera (a partir de la segunda línea)
	for _, line := range lines[1:] {
		// Divide la línea en clave y valor por el primer ":"
		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			// Ignora líneas mal formadas
			continue
		}

		// Limpia espacios en blanco de clave y valor
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			// Ignora cabeceras con clave vacía
			continue
		}

		headers[k] = v
	}

	body := ""

	// Crea y devuelve el objeto HttpRequest con los datos parseados
	return NewHttpRequest(method, target, headers, body), nil
}

// Parsea el cuerpo de la solicitud HTTP si existe.
func ParseBody(request *HttpRequest, reader *bufio.Reader) error {
	// Comprueba si existe la cabecera Content-Length para leer el cuerpo
	contentLengthStr, ok := request.Headers["Content-Length"]
	if !ok {
		return nil
	}

	// Convierte el valor de Content-Length a entero
	var contentLength int
	_, err := fmt.Sscan(contentLengthStr, &contentLength)
	if err != nil {
		return fmt.Errorf("bad content length format: %w", err)
	}

	// Si hay longitud de contenido, lee el cuerpo
	if contentLength <= 0 {
		return nil
	}

	body := make([]byte, contentLength)

	// Lee exactamente contentLength bytes desde el reader
	_, err = io.ReadAtLeast(reader, body, contentLength)
	if err != nil {
		return fmt.Errorf("can't read body: %w", err)
	}

	request.Body = string(body)

	return nil
}
