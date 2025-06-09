# Proj02-Http-Distributed-Server

## Descripción del Proyecto

Este proyecto implementa un **Dispatcher** y dos **Workers** (π y Matrix) basados en el servidor HTTP del Proyecto 1 (`httpserver-core` v0.1.1).

- **Dispatcher**
  - Recibe peticiones en `/fibonacci`, `/hash`, `/simulate`, `/pi` y `/matrix`.
  - Balancea la carga round-robin (con fail-over) hacia los Workers.
  - Ejecuta health-check periódico (`/ping`) de cada Worker.
  - Expone `/workers` que lista `{url, alive, load, done}` de cada uno.

- **Worker π**
  - Expondrá `/pi?trials=N` para cálculo de π vía Monte Carlo.
  - Responde `/ping` con su estado interno (`status`, `load`, `done`).

- **Worker Matrix**
  - Expondrá `POST /matrix` recibiendo JSON `{A,B}` y devolviendo `{result: A×B}`.
  - Responde `/ping` con su estado interno (`status`, `load`, `done`).

> **Nota**: En esta fase ambos Workers sólo devuelven un stub “not implemented” y `/ping`.

---

## Cómo Compilar

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/Joselyn2114/Proj02-Http-Distributed-Server.git
   cd Proj02-Http-Distributed-Server
   ```

2. (Opcional) Ajustar variables de entorno en un archivo `.env`:
   ```
   PORT=8080
   WORKERS=http://worker-pi:8081,http://worker-matrix:8082
   ```

3. Construir y levantar con Docker Compose:
   ```bash
   docker-compose up --build -d
   ```

---

## Cómo Testear

1. Verificar que los contenedores estén corriendo:
   ```bash
   docker ps
   ```

2. Listar estado de Workers:
   ```bash
   curl http://localhost:8080/workers
   ```

3. Probar rutas stub via Dispatcher:
   ```bash
   curl http://localhost:8080/pi?trials=10        # → "not implemented"
   curl -X POST http://localhost:8080/matrix         -H "Content-Type:application/json"         -d '{"A":[[1]],"B":[[1]]}'               # → Bad Request
   curl http://localhost:8080/fibonacci?n=5      # → 404 Not Found
   ```

4. Simular caída y fail-over:
   ```bash
   docker stop worker-pi
   curl http://localhost:8080/workers            # worker-pi.alive:false
   curl http://localhost:8080/pi?trials=1         # “todos los workers caídos”
   docker start worker-pi
   ```

5. Revisar logs:
   ```bash
   docker-compose logs dispatcher
   docker-compose logs worker-pi
   docker-compose logs worker-matrix
   ```

---

## Próximos Pasos

### Dev B: Worker π

1. Implementar handler `/pi?trials=N`:
   - Parsear `trials` de la URL.
   - Calcular π con Monte Carlo.
   - Medir tiempo de ejecución.
   - Actualizar counters `load` y `done`.
   - Devolver JSON `{trials, pi, duration}`.

2. Escribir tests unitarios para:
   - Cálculo de π (casos `trials=0`, `trials>0`).
   - Correcto reporting en `/ping`.

3. Validar integración:
   ```bash
   docker-compose up --build -d
   curl http://localhost:8080/pi?trials=1000
   ```

### Dev C: Worker Matrix

1. Implementar handler `POST /matrix`:
   - Leer JSON `{A, B}` del body.
   - Validar dimensiones conformes.
   - Calcular producto `C = A × B`.
   - Actualizar counters `load` y `done`.
   - Devolver JSON `{result: C}`.

2. Escribir tests unitarios para:
   - Productos de matrices válidas (2×2, 3×3).
   - Error 400 en dimensiones inválidas.
   - Correcto reporting en `/ping`.

3. Validar integración:
   ```bash
   docker-compose up --build -d
   curl -X POST http://localhost:8080/matrix         -H "Content-Type:application/json"         -d '{"A":[[1,2],[3,4]],"B":[[5,6],[7,8]]}'
   ```

---

Con esto tienes todo documentado para compilar, testear y continuar el desarrollo de los Workers específicos. ¡A programar!
