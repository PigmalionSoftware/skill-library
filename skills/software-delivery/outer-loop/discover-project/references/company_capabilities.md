# Capacidades de la Empresa

> Este documento describe **qué puede hacer realmente la empresa**. Es la fuente de anclaje que usa la skill `discover-project` (Etapa 0) para sugerir defaults realistas durante el descubrimiento con un cliente, en lugar de proponer cosas genéricas o que la empresa no ofrece.
>
> Regla de oro al completarlo: **escribí solo lo que es verdad hoy**. Si acá dice que dominamos algo que en realidad no, la skill lo va a sugerir con confianza al cliente. Preferí ser honesto y conservador.

---

## Metadata (documento controlado — ISO 9001 cl. 7.5)

- Versión: `0.1 (borrador)`
- Dueño: _[COMPLETAR: quién firma este documento — típicamente CTO o líder técnico senior]_
- Fecha de última revisión: _[COMPLETAR]_
- Frecuencia de revisión sugerida: _[COMPLETAR: ej. cada 3 o 6 meses, o al cerrar cada proyecto nuevo]_

---

## Cómo se usa este documento

Lo lee `discover-project` al inicio de cada descubrimiento. Cuando la skill le propone un default al facilitador (por ejemplo, ante "necesito cobros dentro de la app"), lo ancla en lo que figura acá. No es un catálogo exhaustivo: es una **foto honesta de las capacidades**.

---

## 1. Servicios que ofrecemos

> Qué tipo de trabajo vende la empresa. Sirve para que la skill entienda si la idea del cliente entra dentro de lo que hacemos.

- _[COMPLETAR: ej. desarrollo de productos a medida]_
- _[COMPLETAR: ej. staff augmentation / equipos dedicados]_
- _[COMPLETAR: ej. modernización de sistemas legacy]_
- _[COMPLETAR: ej. consultoría de arquitectura, discovery, etc.]_

---

## 2. Tipos de proyecto

### Hacemos

> Los proyectos donde la empresa tiene experiencia real y se siente cómoda.

- _[COMPLETAR: ej. plataformas web SaaS B2B]_
- _[COMPLETAR: ej. apps móviles con backend propio]_
- _[COMPLETAR: ej. integraciones entre sistemas / APIs]_

### No hacemos (anti-capacidades)

> Igual de importante que lo anterior. Deja explícito lo que la empresa **no** ofrece, para que la skill no sugiera meterse en terreno desconocido. Ser honesto acá evita comprometer al equipo con algo que no domina.

- _[COMPLETAR: ej. desarrollo de hardware / firmware]_
- _[COMPLETAR: ej. modelos de ML/IA a medida desde cero]_
- _[COMPLETAR: ej. juegos / motores gráficos]_

---

## 3. Stack tecnológico

> Lenguajes, frameworks y herramientas con los que el equipo trabaja habitualmente. Cuanto más específico, mejores los defaults.

### Lenguajes y runtimes

- JavaScript / TypeScript (Node.js)
- Java
- Go
- C# / .NET

### Frontend web

- React
- Next.js

### Mobile

- React Native

### Backend / frameworks

- Quarkus (Java)
- .NET
- Node.js — _[COMPLETAR: ¿qué framework? ej. NestJS, Express, Fastify]_
- Go — _[COMPLETAR: ¿algún framework/estándar? ej. net/http, Gin, Echo]_

### Infraestructura y DevOps

- Docker
- _[COMPLETAR: orquestación — ej. Kubernetes, ECS, Docker Compose en prod]_
- _[COMPLETAR: CI/CD — ej. GitHub Actions, GitLab CI, Jenkins]_

### Cloud

- _[COMPLETAR: ¿AWS, GCP, Azure, on-premise? ¿alguno preferido?]_

### Bases de datos

- _[COMPLETAR: relacionales — ej. PostgreSQL, SQL Server, MySQL]_
- _[COMPLETAR: no relacionales / cache — ej. MongoDB, Redis]_
- _[COMPLETAR: vectoriales, si aplica]_

### Mensajería / eventos _(opcional)_

- _[COMPLETAR: ej. RabbitMQ, Kafka, SQS — o "no aplica"]_

### Testing / calidad _(opcional)_

- _[COMPLETAR: ej. Jest, JUnit, Playwright, cobertura mínima acordada]_

---

## 4. Integraciones que ya dominamos

> Servicios de terceros que el equipo ya integró en proyectos reales. Este es el anclaje más valioso para la Etapa 0: permite que la skill sugiera "esto lo resolvemos con X, que ya usamos", en vez de proponer algo genérico.

| Categoría         | Proveedor(es)            | Proyecto(s) donde se usó | Nivel de dominio |
|-------------------|--------------------------|--------------------------|------------------|
| Pagos             | _[COMPLETAR: ej. Mercado Pago, Stripe]_ | _[COMPLETAR]_ | _[alto/medio/bajo]_ |
| Email / mensajería| _[COMPLETAR: ej. SendGrid, SES]_        | _[COMPLETAR]_ | _[...]_ |
| Autenticación     | _[COMPLETAR: ej. Auth0, Keycloak, Cognito]_ | _[COMPLETAR]_ | _[...]_ |
| Mapas / geo       | _[COMPLETAR: ej. Google Maps]_          | _[COMPLETAR]_ | _[...]_ |
| Notificaciones    | _[COMPLETAR: ej. Firebase Cloud Messaging]_ | _[COMPLETAR]_ | _[...]_ |
| Otras             | _[COMPLETAR]_                            | _[COMPLETAR]_ | _[...]_ |

---

## 5. Restricciones regulatorias y compliance

> Marcos que la empresa maneja o debe respetar. La skill los usa para levantar temprano restricciones que el cliente quizás no menciona.

- Gestión de calidad: _[COMPLETAR: ej. ISO 9001 — estado de certificación]_
- Protección de datos: _[COMPLETAR: ej. Ley 25.326 (Argentina); GDPR si hay clientes en la UE]_
- Sectoriales: _[COMPLETAR: ej. PCI-DSS para pagos, HIPAA para salud, si aplica]_
- Otras: _[COMPLETAR]_

---

## 6. Capacidades del equipo _(opcional)_

> Contexto no técnico que igual condiciona qué proyectos son viables.

- Tamaño del equipo: _[COMPLETAR]_
- Seniority / roles disponibles: _[COMPLETAR]_
- Idiomas de trabajo: _[COMPLETAR: ej. español, inglés]_
- Metodología: _[COMPLETAR: ej. Scrum, Kanban]_
- Zonas horarias / disponibilidad para clientes remotos: _[COMPLETAR]_

---

## 7. Notas de mantenimiento

> Registrá acá los cambios relevantes, para que quede el rastro de revisión que pide ISO.

- _[COMPLETAR: fecha — qué cambió — quién]_
