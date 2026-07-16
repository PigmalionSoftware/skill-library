---
name: elicit-requirements
description: 'ETAPA 1: Convierte un Project Brief ya validado (Etapa 0) en un SRS reducido: requisitos funcionales y no funcionales verificables, trazables y sin ambigüedad, que sirven de entrada de diseño para la Etapa 2 (Arquitectura). Es la Etapa 1 (Requisitos) del flujo de proyecto, conducida por un facilitador interno (PM técnico o líder técnico) que trabaja sobre el Brief y sube al cliente solo las decisiones que lo requieren. Úsala SIEMPRE que exista un Project Brief cerrado y haya que "sacar los requisitos", "armar el SRS", "detallar qué tiene que hacer el sistema", "resolver el handoff a requisitos" o "pasar del brief a los requisitos", aunque no usen esas palabras exactas. NO es para el encuadre de negocio inicial (eso es discover-project, Etapa 0), NO decide arquitectura ni stack (eso es Etapa 2), y NO cierra el requisito de un ticket dentro de un repo existente (para eso está close-requirement).'
---

# Elicitación de Requisitos

## Propósito

Ayudar a un facilitador interno a transformar un **Project Brief validado** (Etapa 0) en un **SRS reducido**: la especificación de **qué tiene que hacer el sistema y bajo qué condiciones se considera bien hecho**, todavía sin decir _cómo_ se construye.

"SRS reducido" es la versión podada del clásico _Software Requirements Specification_: conserva lo que da valor —requisitos funcionales y no funcionales **no ambiguos, verificables y trazables**— y descarta la ceremonia documental del formato formal. Es el escalón que traduce el problema de negocio del Brief a requisitos accionables, y es la **entrada de diseño** que el gate de arquitectura de la Etapa 2 revisa y aprueba (ISO 9001 cl. 8.3).

El Brief es la raíz de la cadena de trazabilidad; este artefacto es el **siguiente eslabón**: todo requisito nace de algo del Brief. Por eso:

- No optimices para redacción.
- Optimizá para **no ambigüedad, verificabilidad y trazabilidad hacia el Brief**.

---

## Límite de la etapa (importante)

Esta skill vive entre dos fronteras y las respeta en ambos sentidos.

**No mira hacia atrás (Etapa 0).** No re-hace el encuadre de negocio: no re-litiga el problema, los objetivos ni el alcance ya cerrados en el Brief. Si al analizar el Brief detectás que el encuadre está **roto o incompleto** (un objetivo sin métrica que impide derivar un RNF, un alcance contradictorio, una dimensión de negocio que quedó abierta), **no lo arregles acá**: rebotá a la Etapa 0 con una nota de qué falta. Corregir el encuadre en esta etapa convierte Requisitos en un segundo discovery encubierto y rompe la trazabilidad.

**No mira hacia adelante (Etapa 2 y 3).** No produce:

- decisiones de arquitectura, stack, componentes, esquema de datos o _cómo_ se cumple un RNF (eso es Etapa 2),
- historias de usuario, escenarios Given-When-Then ni tickets (eso es Etapa 3).

El guard central de la etapa es **"qué / qué-tan-bien, nunca cómo"**. Cuando aparezca una decisión técnica (elección de tecnología, diseño de una integración no impuesta, forma de cumplir un RNF), no la resuelvas: anotala en el **Handoff a Etapa 2**. Así la cadena no se corta.

---

## Contexto

### Entrada obligatoria: el Project Brief

El Brief **no es degradable**. A diferencia de cómo la Etapa 0 trata `company_capabilities.md` (si falta, sigue igual porque es un ancla "nice to have"), acá el Brief **es la entrada de registro**: la cadena de trazabilidad literalmente pasa por él. Sin Brief no hay requisitos que derivar.

Por eso, si esta skill se invoca sin que se le pase el Brief (ruta del archivo o su contenido), **detenete y pedilo**. No procedas a elicitar sin él, no intentes reconstruirlo de memoria y no lo inventes.

Cuando lo tengas, leelo **completo**, incluidas dos secciones que son el corazón de tu insumo:

- el **Handoff a Etapa 1**: la lista puntual de preguntas abiertas que el Brief te dejó para resolver (catálogos a definir como ABM, formato de integraciones, reglas de cálculo, modelo de historial/eventos, roles y permisos, fallbacks no bloqueantes, etc.);
- el **Anexo — Insumos técnicos recibidos del cliente** (si existe): detalle que el cliente ya entregó, preservado tal cual. Es **evidencia a reconciliar**, no requisito cerrado.

### Esta skill NO usa `company_capabilities.md`

Es una decisión deliberada. El ancla de la Etapa 1 es el **Brief**, no las capacidades de la empresa. Anclar los defaults de requisitos en el stack o las integraciones que la empresa domina empujaría justo hacia el "cómo / con qué" que pertenece a la Etapa 2. Los defaults de esta etapa se anclan en el alcance, los objetivos y las restricciones **del Brief**.

(Beneficio lateral: `elicit-requirements` es una carpeta autocontenida sin `references/`, así que no arrastra el problema de duplicación del doc de capacidades.)

### No hay codebase

Igual que en la Etapa 0, el producto todavía no existe: no hay código que inspeccionar. La única fuente de anclaje es el Brief y el criterio del facilitador.

---

## Rol de quien ejecuta esta skill

La conduce el **facilitador interno** (idealmente el mismo del discovery), no el cliente solo. Sobre el Brief, el facilitador deriva y detalla los requisitos con su criterio; **el cliente entra solo como validador puntual** de las decisiones que de verdad requieren su firma (reglas de negocio, prioridades, restricciones que él impone).

Consecuencia directa: distinguí siempre entre lo que podés cerrar con el Brief + tu criterio y lo que **depende de una decisión del cliente**. Lo que depende del cliente y no está confirmado **no se inventa**: se registra como supuesto con **dueño = cliente** y forma de validación.

**Guard obligatorio.** El perfil técnico tiende a saltar a la arquitectura. No lo hagas. Tu criterio técnico se usa solo para:

- detectar cuándo un requisito quedó ambiguo o no verificable,
- exigir que cada RNF sea medible,
- reconocer una restricción o dependencia externa real temprano.

No se usa para elegir tecnología, diseñar componentes, definir el esquema de datos ni resolver _cómo_ se cumple un requisito. Si te descubrís pensando "esto lo haríamos con X" o "esto va en tal tabla", esa es la señal de que te saliste del espacio del _qué_ y te metiste en la Etapa 2.

---

## Comportamiento

### 1. Cargar y validar el Brief

Si no tenés el Brief (ruta o contenido), **detenete y pedilo** antes de cualquier otra cosa. Con el Brief en mano, leelo completo, prestando especial atención al **Handoff a Etapa 1** y al **Anexo** (si existe).

### 2. Sembrar el trabajo desde el Brief

Antes de hacer preguntas nuevas:

1. **Convertí cada ítem del "Handoff a Etapa 1" en un ítem de trabajo abierto** que hay que cerrar o volver supuesto. Esa lista es el esqueleto inicial de esta etapa, no un punto de partida genérico.
2. **Reconciliá con el Anexo.** Buscá qué insumos técnicos que el cliente ya entregó responden (parcial o totalmente) a esos ítems, para no volver a preguntar lo que ya está. El Anexo es evidencia: se valida y se traduce a requisito, no se copia como requisito. Este cruce es además un **chequeo de completitud**: cada concepto del modelo de datos y cada capacidad que figure en el Anexo tiene que terminar reflejado en algún requisito o marcado explícitamente como fuera de alcance. Un concepto del Anexo que no aparece en ningún requisito ni como no-goal es un hueco silencioso —no lo dejes pasar por el solo hecho de estar en el Anexo—: cerralo o registralo como supuesto.
3. **Verificá la integridad del encuadre.** Si para derivar requisitos te falta algo que debía venir cerrado del Brief (p. ej. un objetivo sin métrica, un alcance contradictorio, un decisor sin definir), **rebotá a la Etapa 0** con una nota de qué falta, en vez de resolverlo acá.

### 3. Triage de granularidad (obligatorio cuando la entrada trae detalle técnico)

En esta etapa el detalle técnico entrante puede ser de dos tipos, y se tratan distinto:

- **Contrato de integración impuesto por el cliente** (una API fija que el sistema debe consumir sí o sí): es una **restricción dura**. No la diseñes: preservala **tal cual** en el `## Anexo — Insumos técnicos` del SRS y marcá qué es decisión ya tomada por el cliente vs. qué es solo una propuesta suya a validar.
- **Diseño, stack o forma de resolver algo**: no se resuelve acá. Va al **Handoff a Etapa 2**.

En ningún caso el detalle técnico se traslada tal cual a las secciones de negocio del SRS (funcionales, datos del dominio, integraciones): esas se redactan a nivel de capacidad y significado. El triage se hace una vez, al principio, y condiciona cómo redactás el resto.

### 4. Elicitar por dimensiones

Hacé las preguntas en el mismo idioma que use el cliente/facilitador. Tu objetivo no es explorar sin rumbo: es **cerrar cada dimensión**. Reglas:

- Tono conversacional; lenguaje de negocio con el cliente, criterio técnico con el facilitador.
- Hacé tantas preguntas como haga falta (sin límite artificial).
- Preferí preguntas de trade-off (A vs B) sobre preguntas abiertas.
- Cuando puedas, incluí un **default sugerido anclado en el Brief** (alcance, objetivo o restricción) y explicá en una línea por qué encaja.

#### Dimensiones obligatorias

Tus preguntas deben cubrir colectivamente:

1. **Requisitos funcionales (capacidades).** Qué debe _poder hacer_ el sistema, como capacidades discretas. Se siembra desde el Alcance del Brief + los ítems del Handoff.
2. **Reglas de negocio y lógica.** Distinto de la capacidad: "el sistema permite registrar un cliente" es funcional; "no se puede registrar sin CUIT válido" o "el descuento no supera 20% salvo aprobación" es regla. Acá viven las **reglas de cálculo**. Es la zona donde más se esconde la ambigüedad.
3. **Datos del dominio (perspectiva de negocio).** Qué información maneja el negocio y qué significa: entidades, **catálogos** (los ABM del Handoff), **modelo de historial/eventos**. _Guard: qué información y qué significa, nunca tipos de dato, tablas ni esquema —eso es Etapa 2._
4. **Roles y permisos.** La matriz **actor → capacidad** (qué rol puede ejecutar qué). El _porqué_ de cada autorización se apoya en una regla de negocio (dimensión 2), pero la matriz vive acá.
5. **Requisitos no funcionales (atributos de calidad).** Performance, disponibilidad, seguridad-como-requisito, usabilidad, escala. _Guard duro: medibles y describen el qué-tan-bien, nunca el cómo._ Ejemplo: "P95 < 2s con 500 usuarios concurrentes" es Etapa 1; "usamos caché + balanceador" es Etapa 2.
6. **Interfaces externas / integraciones (contrato funcional).** Qué sistema, qué información de negocio fluye en cada dirección, y por qué. El contrato técnico exacto es Etapa 2, **salvo** que el cliente lo tenga fijo (restricción dura) → se preserva en el Anexo, no se diseña.
7. **Comportamiento ante fallos y bordes.** Qué pasa cuando algo sale mal, y los **fallbacks de funcionalidades no bloqueantes**. Va como dimensión propia justamente porque el happy path es fácil y los modos de falla se olvidan sistemáticamente.
8. **Restricciones y su traducción a requisitos.** No se trata de descubrir restricciones nuevas (deadline, presupuesto, regulatorio ya vienen del Brief), sino de **traducir las que generan requisitos** a algo concreto: regulatorio "derecho al borrado" → RF "el usuario puede solicitar la baja de sus datos". Las que no generan requisito quedan como límites que acotan la Etapa 2.

La UI **no** es una dimensión de esta etapa: "qué puede hacer el usuario" es funcional (dim. 1); el layout, las pantallas o las columnas de un listado van al Anexo; los escenarios paso a paso son Etapa 3.

#### Alcance emergente (lo que el Brief no tenía)

La elicitación no solo cierra los ítems abiertos del Brief: también **descubre** capacidades o integraciones que el Brief no anticipó (p. ej. "hay que poder dar de alta a los colaboradores antes de asignarles un asset", o "los colaboradores en realidad vienen de tal sistema"). Sumarlas está permitido —para eso se elicita—, pero con tres recaudos, porque es acá donde el alcance se infla sin que nadie lo decida:

1. **Marcalo como alcance agregado en Etapa 1.** En el `traza-a` del requisito, dejá explícito "decisión de elicitación — no figuraba en el Brief", para que se vea que expande el encuadre original en vez de derivarse de él.
2. **No lo metas a must-v1 por inercia.** Algo que agranda el alcance —sobre todo una integración no dominada— no es must-v1 solo porque surgió en la misma charla. Separá lo esencial de lo diferible aunque vengan juntos: si una capacidad manual alcanza para el v1, la automatización o integración que la mejora suele ser diferible.
3. **Si mueve el tamaño o el riesgo del v1, es decisión del decisor, no tuya.** Cuando lo emergente cambia de forma significativa el esfuerzo o el riesgo del v1, ya no es un detalle de requisitos: toca una decisión de alcance que pertenece al Brief. Marcalo para revisión del decisor (confirmación explícita o rebote suave a Etapa 0), no lo absorbas en silencio.

### 5. Exigir las reglas transversales a cada requisito

No son dimensiones: son propiedades que **todo requisito arrastra**, lo produzcas donde lo produzcas.

- **Verificabilidad.** Todo RF y todo RNF testeable lleva **criterios de aceptación a nivel de requisito**: condición → resultado esperado, en lenguaje de verificación. **No** son Given-When-Then ni pasos de UI (eso es la descomposición en historias de la Etapa 3; acá decís _qué condición debe cumplirse_, no _el escenario paso a paso_). Un requisito no verificable se reescribe hasta que lo sea, o se vuelve supuesto.
- **Trazabilidad.** Todo requisito traza hacia atrás a algo del Brief (alcance, objetivo o ítem del Handoff), vía un campo `traza-a`. Un requisito que no traza a nada es una de dos cosas: **scope creep** (sacalo) o un **hueco del Brief** (rebote a Etapa 0).
- **Prioridad v1.** Cada requisito se marca _must-v1_ o _diferido_, consistente con el alcance de v1 del Brief. Alimenta el backlog de la Etapa 3.

### 6. Iterar hasta cerrar (regla de cierre-o-supuesto)

Por cada dimensión, **o se cierra una decisión, o se registra como supuesto con dueño y forma de validación**. Nunca queda una ambigüedad silenciosa. Un supuesto bien registrado incluye: qué se asume, quién lo valida (recordá: lo que depende del cliente lleva **dueño = cliente**) y cómo/cuándo se valida.

- Si una respuesta está incompleta → volvé a preguntar.
- No avances mientras queden dimensiones sin cerrar y sin supuesto registrado.

### 7. Pregunta de cierre (obligatoria, siempre la última antes del gate)

Una vez que las dimensiones están cerradas o registradas como supuesto, y **antes** de pasar al gate, hacé siempre esta pregunta, sin excepción:

"Antes de armar el SRS: ¿hay algún requisito, regla o caso que quieras agregar, corregir o aclarar?"

- Si agrega o cambia algo → volvé al paso 6.
- Si no hay nada más → recién ahí pasá al gate del paso 8.

Esta pregunta existe para capturar lo que las dimensiones no anticiparon. No la reemplaces ni la saltees aunque parezca completo.

### 8. Confirmar antes de escribir (gate)

Cuando todo esté cerrado, preguntá:

"Los requisitos están cerrados. ¿Querés que redacte el SRS reducido?"

No lo escribas todavía.

### 9. Redactar y guardar el archivo

Solo si el facilitador confirma explícitamente en el paso 8, escribí el artefacto final con el schema de abajo **y guardalo como archivo**, no solo en el chat.

**Nombre del archivo:**

    srs-reducido-<slug-del-proyecto>-<YYYY-MM-DD>.md

Donde `<slug-del-proyecto>` es el título del proyecto en minúsculas, sin acentos, con espacios reemplazados por guiones. Usá el **mismo slug del Brief** para mantener la continuidad de la cadena. La fecha evita pisar corridas anteriores.

**Ubicación:** la raíz del proyecto/repositorio donde se ejecuta la skill, salvo que el facilitador pida otra.

Después de guardarlo, confirmá al facilitador el nombre y la ubicación del archivo creado.

### 10. Handoff a la Etapa 2

El SRS termina con una sección de handoff que lista todas las decisiones técnicas que aparecieron y **no** se resolvieron (el _cómo_ de cada RNF, elección de stack, diseño de integraciones no impuestas). Es la entrada del gate de arquitectura y de la skill de Etapa 2. Así la cadena no se corta.

---

## Convención de IDs

Los identificadores incluyen el slug del proyecto para que sean únicos entre proyectos (importante cuando terminen referenciados desde Jira en la Etapa 3):

- Requisito funcional: `RF-<slug>-01`, `RF-<slug>-02`, …
- Requisito no funcional: `RNF-<slug>-01`, …
- Regla de negocio: `RN-<slug>-01`, …

Los RF referencian las RN que aplican; la matriz de roles referencia los RF; los modos de falla referencian el RF afectado.

---

## Output

### Si falta el Brief

Respondé en el idioma del usuario, pidiéndolo, sin avanzar:

## Falta la entrada obligatoria

Esta etapa necesita el Project Brief validado de la Etapa 0 (ruta del archivo o su contenido). Sin él no puedo derivar requisitos con trazabilidad. Pasámelo y arranco.

---

### Si el encuadre del Brief está roto o incompleto

## Rebote a Etapa 0 (Descubrimiento)

El Brief no permite derivar requisitos todavía. Falta cerrar en el encuadre de negocio:

- <qué falta y por qué bloquea la derivación de requisitos>

Sugiero volver a `discover-project` para cerrar esto antes de seguir con Requisitos.

---

### Si todavía hay requisitos por cerrar

## Entendimiento

Requisitos ya sembrados desde el Brief (Handoff + Anexo): <...>
Qué falta cerrar: <dimensiones abiertas>

## Preguntas

1. <pregunta de trade-off>

   Default sugerido:
   <opción anclada en el Brief (alcance / objetivo / restricción) + por qué encaja>

2. <pregunta>

   Default sugerido:
   <...>

---

### Si todo está cerrado pero no confirmado

## Estado

Todos los requisitos están cerrados o registrados como supuestos con dueño y forma de validación. Cada uno es verificable y traza a algo del Brief.

## Confirmación

¿Querés que redacte el SRS reducido?

---

### Si está confirmado

# SRS reducido: <título del proyecto>

## Metadata

- Versión: 1.0
- Dueño (facilitador): <nombre / rol>
- Cliente / stakeholder: <...>
- Fecha: <fecha>
- Estado: Borrador para gate de arquitectura
- **Brief de origen:** <nombre del archivo del Brief> (v<versión>)

## Resumen del alcance

<3-4 líneas que reponen el problema y el alcance del v1 tomados del Brief, para que el SRS se entienda sin tener el Brief al lado. No es re-discovery: es el ancla. No re-litigar el encuadre.>

## Requisitos funcionales

> Cada RF: enunciado como capacidad (no como solución). Referencia las RN que aplica.

- **RF-<slug>-01** — <enunciado de la capacidad>
  - Prioridad: <must-v1 | diferido>
  - Traza-a: <ítem del Brief / Handoff que lo origina>
  - Reglas que aplica: <RN-<slug>-0x, …> (si corresponde)
  - Criterios de aceptación:
    - <condición → resultado esperado>
    - <...>
- **RF-<slug>-02** — <...>

## Reglas de negocio

> Reglas verificables, en términos de negocio. Las de cálculo van acá con su lógica.

- **RN-<slug>-01** — <regla> | Traza-a: <...>
- **RN-<slug>-02** — <...>

## Datos del dominio

> Qué información maneja el negocio y qué significa. NO tipos de dato, tablas ni esquema (eso es Etapa 2).

- Entidades y significado: <...>
- Catálogos (ABM): <catálogo, qué representa, quién lo administra>
- Historial / eventos: <qué se registra y por qué>

## Roles y permisos

> Matriz actor → capacidad. El porqué de cada restricción se apoya en la RN correspondiente.

| Rol   | Capacidades (RF) que puede ejecutar | Regla que lo justifica |
| ----- | ----------------------------------- | ---------------------- |
| <rol> | <RF-<slug>-0x, …>                   | <RN-<slug>-0x>         |

## Requisitos no funcionales

> Cada uno medible y describiendo el qué-tan-bien, nunca el cómo.

- **RNF-<slug>-01** — <atributo de calidad> | Métrica y valor objetivo: <p. ej. "P95 < 2s con 500 usuarios concurrentes"> | Traza-a: <...>
- **RNF-<slug>-02** — <...>

## Interfaces externas / integraciones

> Qué sistema, qué información de negocio fluye y en qué dirección, y por qué. El contrato técnico exacto solo si el cliente lo impuso → referencia al Anexo.

- <sistema>: <info de negocio que fluye, dirección, propósito> | Contrato impuesto: <sí → ver Anexo | no → diseño en Etapa 2>

## Comportamiento ante fallos y bordes

> Modos de falla relevantes y fallbacks de lo no bloqueante. Referencia cruzada al RF afectado.

- <modo de falla / borde> → <comportamiento esperado / fallback> | Afecta: <RF-<slug>-0x>

## Supuestos a validar

- Supuesto: <...> | Dueño: <quién valida — cliente si requiere su decisión> | Validación: <cómo / cuándo>
- <...>

## Riesgos

- <riesgo> | Impacto: <...> | Mitigación propuesta: <...>

## Anexo — Insumos técnicos

> Se propaga el Anexo del Brief (si existía) y se le suma lo recibido en esta etapa (p. ej. contratos de integración impuestos por el cliente). Evidencia preservada tal cual, sin resumir. NO es requisito cerrado. Si no hubo insumos de este tipo, omitir la sección.

<detalle técnico preservado, agrupado por tipo: UI / datos / integraciones / reglas ya especificadas>

## Handoff a Etapa 2 (Arquitectura)

Decisiones técnicas que aparecieron y NO se resolvieron acá:

- <el "cómo" de un RNF concreto>
- <elección de stack / componentes>
- <diseño de una integración no impuesta>
- <...>

---

## Reglas

- Sin Project Brief, no arranques: es entrada obligatoria y no degradable. Pedilo y detenete.
- No re-litigues el encuadre de negocio. Si el Brief está roto o incompleto, rebotá a Etapa 0 con nota; no lo arregles acá.
- No escribas arquitectura, stack, esquema de datos, historias, Given-When-Then ni código. Toda decisión técnica va al Handoff a Etapa 2.
- Guard: "qué / qué-tan-bien, nunca cómo". Si pensás "esto lo haríamos con X", te pasaste a Etapa 2.
- Todo requisito es verificable (criterios de aceptación a nivel requisito, no GWT ni pasos de UI) y traza a algo del Brief vía `traza-a`. El que no traza es scope creep o hueco del Brief.
- Lo que la elicitación descubre y el Brief no tenía se puede sumar, pero se marca como alcance agregado en E1 y no entra a must-v1 por inercia; una integración no dominada que agranda el v1 va, por default, diferida o a revisión del decisor.
- Antes de cerrar, hacé el barrido de completitud contra el Anexo: todo concepto o capacidad del Anexo termina en un requisito o como no-goal explícito; lo que no aparece es hueco a cerrar o supuesto.
- Todo RNF es medible: métrica y valor objetivo, o no es un RNF cerrado.
- Los datos del dominio se describen por significado, nunca por tipo/tabla/esquema. La UI no es dimensión de esta etapa.
- El detalle técnico de la entrada no se traslada a las secciones de negocio: contrato impuesto → Anexo; diseño → Handoff a Etapa 2.
- No dejes ninguna dimensión en ambigüedad: cerrala o registrala como supuesto con dueño y validación. Lo que depende del cliente lleva dueño = cliente.
- No redactes el SRS hasta que el facilitador confirme.
- Siempre hacé la pregunta de cierre ("¿algo más?") antes del gate de confirmación, sin excepción.
- El SRS final se guarda como archivo en la raíz del proyecto, con nombre `srs-reducido-<slug>-<fecha>.md` y el mismo slug del Brief; no alcanza con mostrarlo solo en el chat.
- Esta skill no usa `company_capabilities.md`: los defaults se anclan en el Brief.
- Respondé siempre en el idioma del usuario.
- Optimizá para no ambigüedad, verificabilidad y trazabilidad, no para verbosidad.
