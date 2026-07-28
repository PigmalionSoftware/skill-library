---
name: elicit-requirements
description: 'Convierte un Project Brief ya validado (Etapa 0) en un SRS reducido: requisitos funcionales y no funcionales verificables, trazables y sin ambigüedad, que sirven de entrada de diseño para la Etapa 2 (Arquitectura). Es la Etapa 1 (Requisitos) del flujo de proyecto, conducida por un facilitador interno (PM técnico o líder técnico) que trabaja sobre el Brief y sube al cliente solo las decisiones que lo requieren. Úsala SIEMPRE que exista un Project Brief cerrado y haya que "sacar los requisitos", "armar el SRS", "detallar qué tiene que hacer el sistema", "resolver el handoff a requisitos" o "pasar del brief a los requisitos", aunque no usen esas palabras exactas. NO es para el encuadre de negocio inicial (eso es discover-project, Etapa 0), NO decide arquitectura ni stack (eso es Etapa 2), y NO cierra el requisito de un ticket dentro de un repo existente (para eso está close-requirement).'
---

# Elicitación de Requisitos

## Propósito

Ayudar a un facilitador interno a transformar un **Project Brief validado** (Etapa 0) en un **SRS reducido**: la especificación de **qué tiene que hacer el sistema y bajo qué condiciones se considera bien hecho**, todavía sin decir *cómo* se construye.

"SRS reducido" es la versión podada del clásico *Software Requirements Specification*: conserva lo que da valor —requisitos funcionales y no funcionales **no ambiguos, verificables y trazables**— y descarta la ceremonia documental del formato formal. Es el escalón que traduce el problema de negocio del Brief a requisitos accionables, y es la **entrada de diseño** que el gate de arquitectura de la Etapa 2 revisa y aprueba (ISO 9001 cl. 8.3).

El Brief es la raíz de la cadena de trazabilidad; este artefacto es el **siguiente eslabón**: todo requisito nace de algo del Brief. Por eso:

- No optimices para redacción.
- Optimizá para **no ambigüedad, verificabilidad y trazabilidad hacia el Brief**.

---

## Límite de la etapa (importante)

Esta skill vive entre dos fronteras y las respeta en ambos sentidos.

**No mira hacia atrás (Etapa 0).** No re-hace el encuadre de negocio: no re-litiga el problema, los objetivos ni el alcance ya cerrados en el Brief. Si al analizar el Brief detectás que el encuadre está **roto o incompleto** (un objetivo sin métrica que impide derivar un RNF, un alcance contradictorio, una dimensión de negocio que quedó abierta), **no lo arregles acá**: rebotá a la Etapa 0 con una nota de qué falta. Corregir el encuadre en esta etapa convierte Requisitos en un segundo discovery encubierto y rompe la trazabilidad.

**No mira hacia adelante (Etapa 2 y 3).** No produce:

- decisiones de arquitectura, stack, componentes, esquema de datos o *cómo* se cumple un RNF (eso es Etapa 2),
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

No se usa para elegir tecnología, diseñar componentes, definir el esquema de datos ni resolver *cómo* se cumple un requisito. Si te descubrís pensando "esto lo haríamos con X" o "esto va en tal tabla", esa es la señal de que te saliste del espacio del *qué* y te metiste en la Etapa 2.

---

## Comportamiento

### 1. Cargar y validar el Brief

Si no tenés el Brief (ruta o contenido), **detenete y pedilo** antes de cualquier otra cosa. Con el Brief en mano, leelo completo, prestando especial atención al **Handoff a Etapa 1** y al **Anexo** (si existe).

### 2. Sembrar el trabajo desde el Brief

Antes de hacer preguntas nuevas:

1. **Convertí cada ítem del "Handoff a Etapa 1" en un ítem de trabajo abierto** que hay que cerrar o volver supuesto. Esa lista es el esqueleto inicial de esta etapa, no un punto de partida genérico.
2. **Reconciliá con el Anexo.** Buscá qué insumos técnicos que el cliente ya entregó responden (parcial o totalmente) a esos ítems, para no volver a preguntar lo que ya está. El Anexo es evidencia: se valida y se traduce a requisito, no se copia como requisito. Este cruce es además un **chequeo de completitud** en dos niveles:
   - **Conceptos y capacidades.** Cada entidad, atributo, catálogo y capacidad que figure en el Anexo tiene que terminar reflejado en algún requisito o marcado explícitamente como fuera de alcance.
   - **Calificadores y restricciones.** No basta con recoger el concepto: los adjetivos que lo califican son requisitos en sí mismos y se pierden con facilidad porque parecen decoración del atributo. Si el Anexo dice que algo es *único*, *obligatorio*, *configurable*, *derivado*, *inmutable* o que solo admite ciertos valores, ese calificador tiene que aterrizar en una **regla de negocio** explícita, con su borde correspondiente en "Comportamiento ante fallos". Un atributo recogido sin su calificador es una restricción de integridad perdida en silencio (p. ej. tomar "identificador" y perder "único").

   Lo que no aparece ni en un requisito ni como no-goal es un hueco silencioso —no lo dejes pasar por el solo hecho de estar en el Anexo—: cerralo o registralo como supuesto.
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

1. **Requisitos funcionales (capacidades).** Qué debe *poder hacer* el sistema, como capacidades discretas. Se siembra desde el Alcance del Brief + los ítems del Handoff.
2. **Reglas de negocio y lógica.** Distinto de la capacidad: "el sistema permite registrar un cliente" es funcional; "no se puede registrar sin CUIT válido" o "el descuento no supera 20% salvo aprobación" es regla. Acá viven las **reglas de cálculo**. Es la zona donde más se esconde la ambigüedad. *Guard para todo cálculo que dependa del tiempo: además de la fórmula, tiene que declarar su **corte** —hasta cuándo cuenta— para las entidades que salieron del ciclo de vida activo (vendidas, cerradas, dadas de baja, archivadas). Una fórmula "desde tal fecha hasta hoy" sigue corriendo para siempre sobre algo que ya dejó de estar vivo, y eso casi nunca es lo que el negocio quiere; preguntá el corte en vez de asumirlo. Y no alcanza con definir el corte para el caso obvio: recorré **uno por uno todos los estados terminales** y verificá que cada uno tenga una fecha de corte que exista de verdad. Es fácil escribir "el corte es la fecha de venta o de baja" y no notar que un estado terminal como "rota" no tiene ninguna de las dos; ese estado queda sin corte definido y el cálculo vuelve a correr para siempre justo donde nadie mira. Si un estado terminal no tiene una fecha propia, hay que decidir cuál usa.*
3. **Datos del dominio (perspectiva de negocio).** Qué información maneja el negocio y qué significa: entidades, **catálogos** (los ABM del Handoff), **modelo de historial/eventos**. *Guard: qué información y qué significa, nunca tipos de dato, tablas ni esquema —eso es Etapa 2.*
4. **Roles y permisos.** La matriz **actor → capacidad** (qué rol puede ejecutar qué). El *porqué* de cada autorización se apoya en una regla de negocio (dimensión 2), pero la matriz vive acá.
5. **Requisitos no funcionales (atributos de calidad).** Performance, disponibilidad, seguridad-como-requisito, usabilidad, escala. *Guard duro: medibles y describen el qué-tan-bien, nunca el cómo.* Ejemplo: "P95 < 2s con 500 usuarios concurrentes" es Etapa 1; "usamos caché + balanceador" es Etapa 2.
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

- **Verificabilidad.** Todo RF y todo RNF testeable lleva **criterios de aceptación a nivel de requisito**: condición → resultado esperado, en lenguaje de verificación. **No** son Given-When-Then ni pasos de UI (eso es la descomposición en historias de la Etapa 3; acá decís *qué condición debe cumplirse*, no *el escenario paso a paso*). Un requisito no verificable se reescribe hasta que lo sea, o se vuelve supuesto. **Corolario, sin excepción: un requisito sin criterios de aceptación no es un requisito todavía.** Si no podés escribirle criterios porque nadie lo especificó, no le pongas un ID y lo dejes vacío ni con un "a definir": no pertenece al cuerpo del SRS. Va a la lista de fuera de alcance del v1 (ver abajo) o a un supuesto.
- **Trazabilidad.** Todo requisito traza hacia atrás a algo del Brief (alcance, objetivo o ítem del Handoff), vía un campo `traza-a`. Un requisito que no traza a nada es una de dos cosas: **scope creep** (sacalo) o un **hueco del Brief** (rebote a Etapa 0).
- **Prioridad v1.** Cada requisito se marca *must-v1* o *diferido*, consistente con el alcance de v1 del Brief. Alimenta el backlog de la Etapa 3.

#### Fuera de alcance: no lo disfraces de requisito

El barrido de completitud del paso 2 pide que nada quede sin destino, pero **"tener un destino" no significa "ser un requisito"**. Distinguí dos cosas:

- **Requisito diferido:** está especificado —tiene enunciado y criterios de aceptación verificables— y simplemente no entra al v1. Vive en el cuerpo del SRS con su ID y prioridad *diferido*.
- **Fuera de alcance:** los no-goals que el Brief ya declaró y lo que se decide no hacer durante esta etapa. **No recibe ID de requisito ni entra al cuerpo.** Va a la sección `## Fuera de alcance del v1`, como lista breve con el motivo y el destino (v1.x, v2, descartado).

Convertir un no-goal del Brief en un requisito con ID vacío rompe dos cosas: contradice una decisión de alcance que el Brief ya tomó, y ensucia el insumo de la Etapa 3, que podría generar tickets para algo que nadie pidió construir. Lo que el Brief ya declaró fuera de alcance se **propaga como fuera de alcance**, no se reabre acá.

### 6. Iterar hasta cerrar (regla de cierre-o-supuesto)

Por cada dimensión, **o se cierra una decisión, o se registra como supuesto con dueño y forma de validación**. Nunca queda una ambigüedad silenciosa. Un supuesto bien registrado incluye: qué se asume, quién lo valida (recordá: lo que depende del cliente lleva **dueño = cliente**) y cómo/cuándo se valida.

- Si una respuesta está incompleta → volvé a preguntar.
- No avances mientras queden dimensiones sin cerrar y sin supuesto registrado.

#### Todo valor que no vino del Brief ni del cliente se registra como supuesto

Especificar bien exige poner números concretos, y muchos de esos números los vas a proponer vos: targets de performance, volúmenes objetivo, umbrales, límites de tamaño y cantidad, monedas o unidades admitidas, cantidad de reintentos. **Ponerlos está bien y es necesario** —un requisito sin número no es verificable— pero cada valor que no salió del Brief ni de una respuesta del cliente lleva su supuesto con dueño y forma de validación.

El supuesto **no saca el valor del cuerpo del SRS**: la regla sigue diciendo el número exacto, con la misma precisión, para que quien implemente no tenga nada que adivinar. Lo único que agrega es el rastro de su origen. Sin esa marca, a las semanas nadie distingue lo que el cliente pidió de lo que vos decidiste, y cualquiera de los dos se defiende o se cambia con la misma (in)seguridad.

Cuidado especialmente con el margen de holgura: multiplicar la escala declarada del Brief "por si crece" es una decisión de diseño, no un dato del cliente. Si el Brief dice 50 y especificás para 500, el 500 va con supuesto.

### 7. Chequeo de consistencia (antes de la pregunta de cierre)

Un requisito puede ser verificable, trazable y estar completo, y el conjunto igual estar mal: **dos requisitos que se contradicen entre sí**. Los pasos anteriores revisan cada requisito contra su origen; este revisa el conjunto contra sí mismo, que es lo único que detecta una contradicción.

Importa especialmente porque este SRS suele ser el insumo de una implementación asistida por IA. Frente a una contradicción o un caso sin definir, una persona pregunta; un implementador automático elige una opción plausible en silencio, no deja rastro de que eligió, y puede elegir distinto la próxima vez. Una ambigüedad que a un humano le cuesta un mensaje de Slack acá se convierte en comportamiento arbitrario.

Antes de la pregunta de cierre, recorré el conjunto y verificá:

1. **Pares incompatibles.** Cruzá las reglas de negocio entre sí y contra los criterios de aceptación de los RF que las invocan. Buscá dos enunciados que no puedan ser ciertos a la vez. El caso típico: una regla define un valor derivado en función del estado actual, y el criterio de un RF afirma que editar el registro no altera ese valor — si el estado se puede editar, las dos cosas no se sostienen.
2. **Cobertura de valores enumerados.** Por cada conjunto cerrado de valores que declares (estados, categorías, tipos de comprador, roles), verificá que **cada valor** tenga comportamiento definido en **todas** las reglas que dependen de ese conjunto. Alcanza con que un solo valor quede afuera para que haya un caso sin definir; y suele quedar afuera justo el valor menos obvio, no el principal.
3. **Fuente única por decisión.** Cada regla o valor se enuncia en **un solo lugar** del SRS; el resto lo referencia por ID. Si la misma decisión aparece dicha dos veces (p. ej. un límite en una regla de negocio y otra vez como métrica de un RNF), no son dos requisitos: son dos fuentes de verdad que pueden divergir, y quien implemente no tiene forma de saber cuál manda. Dejá una y referenciá.
4. **Consistencia de prioridad.** Un requisito must-v1 no puede depender de otro diferido o fuera de alcance. Si depende, o sube el dependido o baja el dependiente.
5. **Conflicto con lo que el Brief se comprometió a lograr.** Cruzá el conjunto contra los **objetivos medibles y los criterios de éxito** del Brief. Un requisito puede trazar correctamente a algo del Brief y sin embargo dejar sin cumplir lo que el Brief declaró que el v1 iba a conseguir. Ejemplo del patrón: el criterio de éxito dice que todo cambio de cierto tipo queda registrado y consultable, y una regla decide que corregir uno de esos datos es una edición normal que no genera registro. Cada pieza es defendible por separado; juntas incumplen la promesa.

Lo que encuentres se resuelve acá: se elige cuál enunciado gana y se corrige el otro, o se define el caso faltante. Si la resolución requiere una decisión del cliente, va como supuesto —pero la contradicción no puede quedar en pie: dos requisitos que se contradicen no son un supuesto, son un defecto.

**Un conflicto con un objetivo o criterio de éxito del Brief no se archiva como riesgo.** Esta es la salida fácil y hay que resistirla: escribir un riesgo bien redactado se lee como diligencia, deja el problema descrito con precisión y no lo resuelve. Un riesgo sirve para lo que *podría* salir mal; no para blanquear un requisito que incumple lo que el proyecto declaró que iba a lograr. Frente a uno de estos, hay dos salidas legítimas y ninguna más: **corregir el requisito** para que la promesa se cumpla, o **escalarlo al decisor** como decisión explícita de bajar la vara que el Brief había fijado. Si el decisor la baja, eso se registra como tal —decisión de alcance, con quién la tomó— no como un riesgo que quedó anotado.

El motivo es el mismo que sostiene toda la etapa: si un implementador automático lee el SRS, no lee los riesgos. Implementa la regla tal como está escrita, y el incumplimiento entra al producto con el defecto documentado a un costado.

### 8. Pregunta de cierre (obligatoria, siempre la última antes del gate)

Una vez que las dimensiones están cerradas o registradas como supuesto, y **antes** de pasar al gate, hacé siempre esta pregunta, sin excepción:

"Antes de armar el SRS: ¿hay algún requisito, regla o caso que quieras agregar, corregir o aclarar?"

- Si agrega o cambia algo → volvé al paso 6 y volvé a correr el chequeo de consistencia del paso 7.
- Si no hay nada más → recién ahí pasá al gate del paso 9.

Esta pregunta existe para capturar lo que las dimensiones no anticiparon. No la reemplaces ni la saltees aunque parezca completo.

### 9. Confirmar antes de escribir (gate)

Cuando todo esté cerrado, preguntá:

"Los requisitos están cerrados. ¿Querés que redacte el SRS reducido?"

No lo escribas todavía.

### 10. Redactar y guardar el archivo

Solo si el facilitador confirma explícitamente en el paso 9, escribí el artefacto final con el schema de abajo **y guardalo como archivo**, no solo en el chat.

**Nombre del archivo:**

    srs-reducido-<slug-del-proyecto>-<YYYY-MM-DD>.md

Donde `<slug-del-proyecto>` es el título del proyecto en minúsculas, sin acentos, con espacios reemplazados por guiones. Usá el **mismo slug del Brief** para mantener la continuidad de la cadena. La fecha evita pisar corridas anteriores.

**Ubicación:** la raíz del proyecto/repositorio donde se ejecuta la skill, salvo que el facilitador pida otra.

Después de guardarlo, confirmá al facilitador el nombre y la ubicación del archivo creado.

### 11. Handoff a la Etapa 2

El SRS termina con una sección de handoff que lista todas las decisiones técnicas que aparecieron y **no** se resolvieron (el *cómo* de cada RNF, elección de stack, diseño de integraciones no impuestas). Es la entrada del gate de arquitectura y de la skill de Etapa 2. Así la cadena no se corta.

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

| Rol | Capacidades (RF) que puede ejecutar | Regla que lo justifica |
|-----|-------------------------------------|------------------------|
| <rol> | <RF-<slug>-0x, …> | <RN-<slug>-0x> |

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

## Fuera de alcance del v1

> Los no-goals que el Brief ya declaró, más lo que se decidió no hacer en esta etapa. **Sin ID de requisito y sin criterios de aceptación**: no son requisitos, son límites. Existe para que la Etapa 3 no genere tickets de esto y para que el barrido de completitud tenga dónde cerrar lo que no se construye.

- <capacidad / integración> — Motivo: <no-goal del Brief | decidido en E1> | Destino: <v1.x | v2 | descartado>

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
- Un requisito sin criterios de aceptación no es un requisito: no le pongas ID con los criterios vacíos ni "a definir". Va a `Fuera de alcance del v1` o a un supuesto.
- Los no-goals que el Brief ya declaró se propagan a `Fuera de alcance del v1`, no se convierten en requisitos con ID. Diferido con criterios escritos sí es requisito; diferido sin especificar, no.
- Lo que la elicitación descubre y el Brief no tenía se puede sumar, pero se marca como alcance agregado en E1 y no entra a must-v1 por inercia; una integración no dominada que agranda el v1 va, por default, diferida o a revisión del decisor.
- Antes de cerrar, hacé el barrido de completitud contra el Anexo en dos niveles: (a) todo concepto o capacidad termina en un requisito o en `Fuera de alcance del v1`; (b) todo calificador declarado —único, obligatorio, configurable, derivado, inmutable, valores admitidos— aterriza en una regla de negocio con su borde. Lo que no aparece es hueco a cerrar o supuesto.
- Toda regla de cálculo que dependa del tiempo declara su corte para entidades fuera del ciclo de vida activo (vendidas, cerradas, de baja), y el corte se verifica **estado terminal por estado terminal**: si alguno no tiene una fecha propia que exista, hay que decidir cuál usa. Preguntá el corte, no lo asumas.
- Antes de la pregunta de cierre, corré el chequeo de consistencia: pares incompatibles entre reglas y criterios, cobertura de todos los valores de cada conjunto enumerado, una sola fuente por decisión, ningún must-v1 que dependa de algo diferido, y ningún requisito que incumpla un objetivo o criterio de éxito del Brief. Dos requisitos que se contradicen no son un supuesto: son un defecto y se resuelven acá.
- Un requisito que incumple un objetivo o criterio de éxito del Brief no se archiva como riesgo: se corrige, o se escala al decisor como decisión explícita de bajar la vara del Brief. Un riesgo describe lo que podría salir mal, no blanquea una promesa incumplida.
- Cada decisión se enuncia en un solo lugar y el resto la referencia por ID. La misma regla dicha dos veces son dos fuentes de verdad que divergen.
- Todo valor que no vino del Brief ni de una respuesta del cliente (targets, umbrales, límites, unidades, márgenes de holgura sobre la escala declarada) va con supuesto y dueño. El supuesto no saca el número del cuerpo: el requisito sigue siendo igual de preciso, solo queda el rastro de su origen.
- Todo RNF es medible: métrica y valor objetivo, o no es un RNF cerrado. Un RNF cuya métrica reenuncia un RF o una regla de negocio no es un RNF: es duplicación.
- Los datos del dominio se describen por significado, nunca por tipo/tabla/esquema. La UI no es dimensión de esta etapa.
- El detalle técnico de la entrada no se traslada a las secciones de negocio: contrato impuesto → Anexo; diseño → Handoff a Etapa 2.
- No dejes ninguna dimensión en ambigüedad: cerrala o registrala como supuesto con dueño y validación. Lo que depende del cliente lleva dueño = cliente.
- No redactes el SRS hasta que el facilitador confirme.
- Siempre hacé la pregunta de cierre ("¿algo más?") antes del gate de confirmación, sin excepción.
- El SRS final se guarda como archivo en la raíz del proyecto, con nombre `srs-reducido-<slug>-<fecha>.md` y el mismo slug del Brief; no alcanza con mostrarlo solo en el chat.
- Esta skill no usa `company_capabilities.md`: los defaults se anclan en el Brief.
- Respondé siempre en el idioma del usuario.
- Optimizá para no ambigüedad, verificabilidad y trazabilidad, no para verbosidad.
