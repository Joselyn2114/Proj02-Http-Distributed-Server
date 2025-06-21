package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sort"
)

// Representa una respuesta HTTP.
type HttpResponse struct {
	StatusCode int               // Código de estado HTTP (ej. 200, 404).
	StatusText string            // Texto del estado HTTP (ej. "OK", "Not Found").
	Headers    map[string]string // Cabeceras HTTP.
	Body       string            // Cuerpo de la respuesta.
}

// Crea una nueva instancia de HttpResponse con los valores proporcionados.
func NewHttpResponse(statusCode int, statusText string, body string) *HttpResponse {
	return &HttpResponse{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    make(map[string]string),
		Body:       body,
	}
}

// Crea una respuesta HTTP 200 OK predeterminada.
func Ok() *HttpResponse {
	return NewHttpResponse(200, "OK", "")
}

// Crea una respuesta HTTP 404 Not Found predeterminada.
func NotFound() *HttpResponse {
	return NewHttpResponse(404, "Not Found", "")
}

// Crea una respuesta HTTP 400 Bad Request predeterminada.
func BadRequest() *HttpResponse {
	return NewHttpResponse(400, "Bad Request", "")
}

// Establece el código de estado de la respuesta.
func (response *HttpResponse) SetStatusCode(code int) *HttpResponse {
	response.StatusCode = code
	return response
}

// Establece el texto del estado de la respuesta.
func (response *HttpResponse) SetStatusText(text string) *HttpResponse {
	response.StatusText = text
	return response
}

// Establece una cabecera HTTP específica.
func (response *HttpResponse) SetHeader(key, value string) *HttpResponse {
	response.Headers[key] = value
	return response
}

// Establece el cuerpo de la respuesta.
func (response *HttpResponse) SetBody(body string) *HttpResponse {
	response.Body = body
	return response
}

// Establece la cabecera Content-Type.
func (response *HttpResponse) SetContentType(contentType string) *HttpResponse {
	response.SetHeader("Content-Type", contentType)
	return response
}

// Establece el cuerpo como texto plano y ajusta la cabecera Content-Type.
func (response *HttpResponse) Text(text string) *HttpResponse {
	response.SetContentType("text/plain")
	response.SetBody(text)
	return response
}

// Establece el cuerpo como JSON y ajusta la cabecera Content-Type.
func (response *HttpResponse) Json(json string) *HttpResponse {
	response.SetContentType("application/json")
	response.SetBody(json)
	return response
}

// Convierte la respuesta HTTP a su representación en formato de cadena HTTP/1.0.
// Calcula automáticamente la cabecera Content-Length.
func (response *HttpResponse) String() string {
	// Calcula y establece la longitud del contenido.
	contentLength := len(response.Body)
	response.SetHeader("Content-Length", fmt.Sprint(contentLength))

	// Formatea las cabeceras.
	keys := make([]string, 0, len(response.Headers))
	for key := range response.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	headersStr := ""
	for _, key := range keys {
		headersStr += fmt.Sprintf("%s: %s\r\n", key, response.Headers[key])
	}

	// Construye la cadena de respuesta HTTP completa.
	return fmt.Sprintf("HTTP/1.0 %d %s\r\n%s\r\n%s", response.StatusCode, response.StatusText, headersStr, response.Body)
}

func (response *HttpResponse) WriteResponse(conn net.Conn) error {
	slog.Info("Response", "address", conn.RemoteAddr().String(), "status_code", response.StatusCode, "status_text", response.StatusText)

	_, err := conn.Write([]byte(response.String()))
	if err != nil {
		return err
	}

	return nil

}

// JsonObj serializa v a JSON y lo pone en el body con application/json.
func (r *HttpResponse) JsonObj(v interface{}) *HttpResponse {
	data, err := json.Marshal(v)
	if err != nil {
		// Construyo la respuesta 500 y luego le pongo el Content-Type.
		resp := NewHttpResponse(500, "Internal Server Error", "json marshal error")
		resp.SetContentType("text/plain")
		return resp
	}
	// En el caso normal, reutilizo 'r'.
	r.SetContentType("application/json")
	r.SetBody(string(data))
	return r
}
