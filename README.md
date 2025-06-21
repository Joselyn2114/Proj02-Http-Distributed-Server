# Proyecto 2 - Servidor HTTP Distribuido

**Autores:**
- Joselyn Jiménez  
- Katerine Guzmán  
- Esteban Solano  

---

## 1. Descripción del Programa
Este proyecto extiende el **Servidor HTTP** del Proyecto 1 añadiendo un **Dispatcher** que distribuye tareas de forma concurrente y tolerante a fallos entre múltiples **Workers** ejecutando instancias Docker del servidor HTTP base.

Principales funcionalidades:
- Balanceo de carga Round-Robin con reintentos automáticos.
- Split/Merge de matrices para procesamiento distribuido.
- Cálculo de π por Monte Carlo paralelo.
- Registro y desregistro dinámico de Workers.
- Health-check periódico y marcación de Workers activos/inactivos.
- Proxy transparente de todos los endpoints originales.

---

## 2. Compilación

**Requisitos:**
- Go 1.22+
- Docker y Docker Compose

**Pasos:**
```bash
git clone https://github.com/Joselyn2114/Proj02-Http-Distributed-Server.git
cd Proj02-Http-Distributed-Server

# Compilar Dispatcher y Worker
docker-compose build
```

---

## 3. Pruebas Ejecutadas y Cómo Ejecutarlas

1. **Registro dinámico de Workers**  
   - Iniciar solo el dispatcher:  
     ```bash
     docker-compose down
     docker-compose up -d --no-deps dispatcher
     ```  
   - Verificar lista vacía:  
     ```bash
     curl http://localhost:8000/workers
     ```  
   - Arrancar un Worker y registrarlo:  
     ```bash
     docker-compose up -d worker1
     curl -X POST http://localhost:8000/register \
       -d '{"url":"http://worker1:8080"}' \
       -H "Content-Type: application/json"
     curl http://localhost:8000/workers
     ```
   - Desregistro:  
     ```bash
     curl -X POST http://localhost:8000/unregister \
       -d '{"url":"http://worker1:8080"}' \
       -H "Content-Type: application/json"
     curl http://localhost:8000/workers
     ```

2. **Health-Check Dinámico**  
   ```bash
   docker-compose up -d dispatcher worker1 worker2
   curl http://localhost:8000/workers
   docker stop worker2
   sleep 6
   curl http://localhost:8000/workers
   ```

3. **Balanceo Round-Robin & Reintentos**  
   ```bash
   for i in {1..9}; do
     curl "http://localhost:8000/pi/part?iter=50000"
   done
   curl http://localhost:8000/workers
   docker stop worker2
   for i in {1..3}; do
     curl "http://localhost:8000/pi/part?iter=30000"
   done
   curl http://localhost:8000/workers
   ```

4. **Split & Merge de Matrices**  
   ```bash
   cat <<EOF > identity.json
   { "a": [[1,2,3],[4,5,6],[7,8,9]], "b": [[1,0,0],[0,1,0],[0,0,1]] }
   EOF
   curl -X POST http://localhost:8000/matrix \
     -d @identity.json \
     -H "Content-Type: application/json"
   docker stop worker2 && sleep 2
   curl -X POST http://localhost:8000/matrix \
     -d @identity.json \
     -H "Content-Type: application/json"
   ```

5. **Endpoints Originales vía Proxy**  
   ```bash
   curl "http://localhost:8000/fibonacci?num=10"
   curl "http://localhost:8000/hash?text=hola123"
   curl "http://localhost:8000/simulate?task=5"
   curl "http://localhost:8000/sleep?seconds=3"
   ```

---

## 4. Arquitectura del Sistema Distribuido

```
Client → HTTP → Dispatcher (Go)
                     ├─ HealthChecker (/ping)
                     ├─ Register/Unregister (/register, /unregister)
                     ├─ Status (/workers)
                     ├─ Matrix (/matrix)
                     └─ Proxy genérico → Workers
Worker (Go HTTP Server base) ↔ contenedor Docker
```

---

## 5. Protocolos de Comunicación

- **HTTP/1.1** para todas las comunicaciones.
- Métodos:
  - **GET** `/pi/part`, `/ping`, `/workers`.
  - **POST** `/matrix`, `/matrix/part`, `/register`, `/unregister`.
  - Proxy de **GET**, **POST**, **DELETE**, etc., para rutas originales.
- **JSON** en cuerpo de requests/responses para endpoints distribuidos.

---

## 6. Tolerancia a Fallos y Escalabilidad

- **Health-Checks** cada 5s, marca inactivo al fallar.
- **Reintentos** automáticos repartiendo sub-tareas.
- **Registro Dinámico** de Workers en caliente.
- **Split & Merge**: cada Worker procesa un bloque.
- **Escalar** con `docker-compose up --scale worker=X`.

---

## 7. Resultados Comparativos

| Ejecución     | Monte Carlo (1e6 iter) | Multiplicación 3×3 |
|---------------|------------------------|--------------------|
| Local mono    | ~1.2s                  | ~10ms              |
| Distribuida 3 Workers | ~0.45s                | ~15ms              |

- **Aceleración** ~2.7× en Monte Carlo pese a overhead de red.
- **Matrices grandes** obtienen beneficio real al paralelizar.
