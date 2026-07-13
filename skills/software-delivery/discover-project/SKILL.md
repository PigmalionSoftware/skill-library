---
name: discover-project
description: 'Convierte la idea vaga de un cliente en un Project Brief cerrado en el encuadre de negocio: separa el problema de la solución, fuerza objetivos medibles y deja lo desconocido como supuestos a validar en vez de ambigüedad silenciosa. Es la Etapa 0 (descubrimiento) del flujo de proyecto, pensada para un facilitador interno (PM técnico o líder técnico) que conduce la conversación con el cliente. Úsala SIEMPRE que arranque un proyecto nuevo, cuando un cliente llegue con una idea poco definida, o cuando pidan "armar el brief", "encuadrar el proyecto", "entender qué quiere el cliente" o "hacer el descubrimiento", aunque no usen esas palabras exactas. NO es para cerrar el requisito de un ticket dentro de un repo existente — para eso está close-requirement.'
---

# Project Discovery

## Propósito

Ayudar a un facilitador interno a transformar la idea vaga de un cliente en un **Project Brief cerrado en el encuadre de negocio**.

El único trabajo de esta etapa es **separar el problema de la solución** y dejar registrado el problema real, medible y acotado. El cliente casi siempre llega con una solución ya pensada ("quiero una app que…"); tu tarea es extraer el problema de negocio que hay debajo.

Este artefacto es la **raíz de la cadena de trazabilidad** del proyecto y la entrada de diseño de la Etapa 1 (Requisitos). Por eso:

- No optimices para redacción.
- Optimiza para **claridad del problema, objetivos medibles y decisiones/supuestos explícitos**.

---

## Límite de la etapa (importante)

Esta skill se detiene en el **encuadre de negocio**. No produce:

- requisitos funcionales o criterios de aceptación detallados (eso es Etapa 1),
- decisiones de arquitectura, stack o diseño técnico (eso es Etapa 2),
- historias de usuario ni tickets (eso es Etapa 3).

Si empujás la conversación hacia esas cosas, las etapas se fusionan y se pierde el escalón que da la trazabilidad. Cuando aparezca una decisión técnica, anotala como pregunta abierta para la etapa correspondiente y seguí.

---

## Contexto

Antes de hacer preguntas, leé el documento de capacidades de la empresa, que viaja junto con esta skill:

    references/company_capabilities.md

Sirve para el mismo rol que la arquitectura cumple en el loop interno: **anclar los defaults sugeridos en la realidad de la empresa** (qué servicios ofrece, con qué stacks trabaja, qué restricciones regulatorias/ISO maneja, qué integraciones ya domina).

- Usalo siempre para que tus defaults sean específicos y no genéricos.
- Si por algún motivo no podés leerlo (por ejemplo, la skill se instaló sin bundled resources), no te detengas: procedé igual, avisá una sola vez que los defaults estarán menos anclados y registralo como un supuesto en el brief final.
- Este documento es específico de la empresa y se mantiene por PR dentro del repo de esta skill, igual que el resto del contenido de `discover-project`. Si en el futuro otra skill del flujo también lo necesita, va a requerir su propia copia en `references/`, porque `npx skills add` instala cada skill como unidad autocontenida y no comparte recursos entre carpetas de skills distintas.

A diferencia del loop interno, **acá no hay codebase que inspeccionar**: el producto todavía no existe. La única fuente de anclaje es el doc de capacidades y el criterio del facilitador.

---

## Rol de quien ejecuta esta skill

La conduce un **facilitador interno** (PM técnico o líder técnico), no el cliente solo. Podés hablar en lenguaje de negocio con el cliente y, a la vez, aprovechar el criterio técnico del facilitador para proponer defaults realistas.

**Guard obligatorio.** El perfil técnico tiende a saltar a la solución. No lo hagas. Tu criterio técnico se usa solo para:

- detectar cuándo una idea del cliente esconde un problema distinto al que enuncia,
- proponer defaults de *negocio* anclados en lo que la empresa puede hacer,
- reconocer una restricción o riesgo real temprano.

No se usa para elegir tecnología, diseñar la solución ni prometer una implementación. Si te descubrís pensando "esto lo haríamos con X", esa es la señal de que te saliste del espacio del problema.

---

## Comportamiento

### 1. Entender la solicitud

Leé la entrada y separá explícitamente:

- **Problema**: qué duele hoy, en términos de negocio.
- **Solución propuesta por el cliente**: lo que pidió (guardala, pero no la trates como un requisito).
- **Qué no está claro**: los huecos que hay que cerrar o convertir en supuestos.

### 2. Triage de granularidad (obligatorio cuando la entrada trae detalle técnico)

Algunos clientes —sobre todo con perfil técnico— llegan con documentación que ya incluye specs de UI, endpoints, tablas/atributos de datos, o catálogos exhaustivos. Ese detalle es valioso, pero **no pertenece al brief de negocio**: pertenece a Requisitos (Etapa 1) o Arquitectura (Etapa 2).

Antes de armar cualquier pregunta, recorré la entrada e identificá todo lo que sea de este tipo:

- specs de UI (columnas de listados, filtros, layout, campos de formulario),
- modelos de datos (atributos, tipos de dato, catálogos, entidades),
- endpoints, contratos de API o integraciones ya especificadas a nivel técnico,
- reglas de cálculo o lógica de negocio ya detallada a nivel de implementación.

Con eso:

1. **No lo traslades tal cual al brief.** El brief resume esto en términos de capacidad de negocio (qué proceso cubre, qué decisión habilita), nunca como lista de columnas/campos/endpoints.
2. **No lo descartes.** Preservalo íntegro en una sección aparte del brief: `## Anexo — Insumos técnicos recibidos del cliente`. Ahí sí va tal cual vino, sin resumir, para que Requisitos/Arquitectura lo aprovechen sin tener que pedirlo de nuevo.
3. **Marcá qué quedó como decisión ya tomada por el cliente vs. qué es solo una propuesta a validar.** Que el cliente haya escrito un endpoint no significa que esa forma de resolverlo esté aprobada — a menos que el propio cliente lo presente como una restricción dura (ej. "tiene que integrar con este sistema que ya tiene esta API").

Este triage se hace una sola vez, al principio, y no se repite en cada pregunta — pero condiciona cómo redactás "Alcance" en el paso final: ese detalle no debe reaparecer ahí.

### 3. Hacer preguntas de descubrimiento

Hacé las preguntas en el mismo idioma que use el cliente/facilitador.

Tu objetivo no es explorar sin rumbo: es **cerrar el encuadre**. Reglas:

- Tono conversacional, lenguaje de negocio (no jerga técnica con el cliente).
- Hacé tantas preguntas como haga falta para cerrar el encuadre (sin límite artificial).
- Cada pregunta resuelve una decisión concreta o valida un supuesto.
- Preferí preguntas de trade-off (A vs B) sobre preguntas abiertas.
- Cuando puedas, incluí un **default sugerido** anclado en el doc de capacidades y explicá en una línea por qué encaja.

#### Dimensiones de encuadre obligatorias

Tus preguntas deben cubrir colectivamente:

1. **Problema y job-to-be-done** — qué intenta lograr el cliente, no qué solución pidió.
2. **Actor y estado actual** — quién tiene el problema y cómo lo resuelve hoy (herramienta, proceso manual, workaround).
3. **Objetivos medibles** — qué resultado de negocio se busca, expresado como algo observable.
4. **Alcance del v1** — qué entra, qué queda afuera, y **no-goals** explícitos.
5. **Restricciones** — deadline, presupuesto, regulatorio/compliance, integraciones obligatorias, plataforma.
6. **Stakeholders y decisor** — quiénes intervienen y, sobre todo, **quién firma / decide**.
7. **Supuestos y riesgos** — qué se está dando por sentado y qué podría hacer fracasar el proyecto.

Si alguna de estas no está clara, tenés que preguntar por ella.

#### Regla de cierre-o-supuesto (adaptación propia del descubrimiento)

En descubrimiento no siempre se puede cerrar todo: hay cosas genuinamente desconocidas que dependen del cliente o del mercado. La regla no es "cerrá todo", sino:

> Por cada dimensión, **o se cierra una decisión, o se registra como supuesto con dueño y forma de validación**. Nunca queda una ambigüedad silenciosa.

Un supuesto bien registrado incluye: qué se asume, quién lo valida y cómo/ cuándo se valida.

#### Objetivos medibles (obligatorio)

Toda meta vaga se empuja a algo medible antes de aceptarla.

**Ejemplo:**
Cliente dice: "quiero mejorar la experiencia del usuario"
Cerrás a: "reducir el tiempo de alta de un cliente nuevo de ~15 min a menos de 5 min"

Si el cliente no puede dar una métrica todavía, registralo como supuesto a validar, no como objetivo cerrado.

#### Defaults anclados (crítico)

Antes de sugerir cualquier default:

- revisá el doc de capacidades cuando sea relevante,
- alineá con lo que la empresa realmente ofrece y con sus restricciones.

Evitá defaults genéricos si hay evidencia en el doc de capacidades. Cuando sugieras uno, explicá brevemente por qué encaja con la empresa y, si aplica, referenciá el servicio o capacidad concreta.

### 4. Iterar hasta cerrar el encuadre

- Si una respuesta está incompleta → volvé a preguntar.
- Si algo es ambiguo → cerralo o convertilo en supuesto explícito.
- No avances mientras queden dimensiones sin cerrar y sin supuesto registrado.

### 5. Pregunta de cierre (obligatoria, siempre la última antes del gate)

Una vez que las 7 dimensiones están cerradas o registradas como supuesto, y **antes** de pasar al gate de confirmación, hacé siempre esta pregunta, sin excepción:

"Antes de armar el brief: ¿hay algo más que quieras agregar, corregir o aclarar?"

- Si la respuesta agrega o cambia algo → volvé al paso 4 (puede reabrir una dimensión o agregar un supuesto nuevo).
- Si la respuesta es que no hay nada más → recién ahí pasá al gate del paso 6.

Esta pregunta existe para capturar lo que las 7 dimensiones no anticiparon. No la reemplaces por ninguna de las preguntas de encuadre ni la saltees aunque el encuadre parezca completo.

### 6. Confirmar antes de escribir (gate)

Cuando el encuadre esté cerrado, preguntá:

"El encuadre está cerrado. ¿Querés que redacte el Project Brief?"

No lo escribas todavía.

### 7. Redactar y guardar el archivo

Solo si el facilitador confirma explícitamente en el paso 6, escribí el artefacto final con el schema de abajo **y guardalo como archivo**, no solo en el chat.

**Nombre del archivo:**

    project-brief-<slug-del-proyecto>-<YYYY-MM-DD>.md

Donde `<slug-del-proyecto>` es el título del proyecto en minúsculas, sin acentos, con espacios reemplazados por guiones (ej. `project-brief-portal-clientes-acme-2026-07-13.md`). La fecha evita que se pise con briefs de otros proyectos o con una corrida anterior del mismo día.

**Ubicación:** la raíz del proyecto/repositorio donde se está ejecutando esta skill (no dentro de una subcarpeta, salvo que el facilitador pida explícitamente otra ubicación).

Si el paso 2 identificó insumos técnicos del cliente, el archivo incluye la sección `## Anexo — Insumos técnicos recibidos del cliente` con ese detalle preservado tal cual.

Después de guardarlo, confirmá al facilitador el nombre y la ubicación del archivo creado.

### 8. Handoff a la Etapa 1

El brief termina con una sección de handoff que lista las preguntas abiertas que la Etapa 1 (Requisitos) deberá resolver. Así la cadena no se corta.

---

## Output

### Si todavía hay encuadre abierto

Respondé en el idioma del usuario:

## Entendimiento

Problema (no la solución): <...>
Solución que pidió el cliente: <...>
Qué falta cerrar: <...>

## Preguntas

1. <pregunta de trade-off>

   Default sugerido:
   <opción anclada en el doc de capacidades + por qué encaja>

2. <pregunta>

   Default sugerido:
   <...>

---

### Si el encuadre está cerrado pero no confirmado

## Estado

El encuadre de negocio está cerrado. Todas las dimensiones están resueltas o registradas como supuestos con dueño y forma de validación.

## Confirmación

¿Querés que redacte el Project Brief?

---

### Si está confirmado

# Project Brief: <título claro del proyecto>

## Metadata

- Versión: 1.0
- Dueño (facilitador): <nombre / rol>
- Cliente / stakeholder: <...>
- Fecha: <fecha>
- Estado: Borrador para validación del cliente

## Problema

<el problema de negocio, en términos del cliente, NO la solución>

## Job-to-be-done y estado actual

- Qué intenta lograr el cliente: <...>
- Cómo lo resuelve hoy: <herramienta / proceso / workaround>

## Objetivos medibles

- <resultado de negocio observable, con métrica y valor objetivo>
- <...>

## Actores y stakeholders

- Usuario(s) afectado(s): <...>
- Decisor / quien firma: <...>
- Otros stakeholders: <...>

## Alcance del v1

### En alcance

- <ítem, a nivel de capacidad de negocio — NO columnas de listado, ni atributos de datos, ni endpoints>

### Fuera de alcance / no-goals

- <ítem que explícitamente NO se hace en el v1>

## Restricciones

- Deadline: <...>
- Presupuesto: <...>
- Regulatorio / compliance: <...>
- Integraciones obligatorias: <...>
- Plataforma: <...>

## Supuestos a validar

- Supuesto: <...> | Dueño: <quién valida> | Validación: <cómo / cuándo>
- <...>

## Riesgos

- <riesgo> | Impacto: <...> | Mitigación propuesta: <...>

## Criterios de éxito

- <condición de negocio observable que indica que el v1 resolvió el problema>
- <...>

## Anexo — Insumos técnicos recibidos del cliente

> Solo si el paso 2 detectó detalle técnico en la entrada (specs de UI, modelos de datos, endpoints, catálogos). Se preserva tal cual vino, sin resumir, para que Requisitos/Arquitectura lo reciban directo. Si no hubo insumos de este tipo, omitir esta sección.

<detalle técnico preservado, agrupado por tipo: UI / datos / integraciones / reglas de cálculo>

## Handoff a Etapa 1 (Requisitos)

Preguntas abiertas que Requisitos debe resolver:

- <pregunta que quedó fuera del encuadre de negocio>
- <...>

---

## Reglas

- No escribas requisitos técnicos, arquitectura ni código.
- No trates la solución que pidió el cliente como un requisito cerrado.
- Si la entrada trae detalle técnico (UI, datos, endpoints, catálogos), no lo traslades a "Alcance": resumilo a nivel de negocio y preservá el detalle íntegro en el Anexo.
- No dejes ninguna dimensión en ambigüedad: cerrala o registrala como supuesto con dueño y validación.
- Toda meta vaga se convierte en objetivo medible o en supuesto a validar.
- No redactes el brief hasta que el facilitador confirme.
- Siempre hacé la pregunta de cierre ("¿algo más?") antes del gate de confirmación, sin excepción.
- El brief final se guarda como archivo en la raíz del proyecto, con el nombre `project-brief-<slug>-<fecha>.md`; no alcanza con mostrarlo solo en el chat.
- Anclá los defaults en el doc de capacidades; si no existe, avisalo y registralo como supuesto.
- Respondé siempre en el idioma del usuario.
- Optimizá para claridad del problema, no para verbosidad.
