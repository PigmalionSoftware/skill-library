---
name: design-architecture
description: 'ETAPA 2: Convierte un SRS reducido ya cerrado (Etapa 1) en un documento de arquitectura único —cuerpo afirmativo del diseño resuelto, más un anexo de decisiones firmables— que sirve al gate humano del arquitecto (ISO 9001 cl. 8.3) y, sobre todo, como contexto sin ambigüedad para que un modelo de IA implemente después sin alucinar. Es la Etapa 2 (Arquitectura) del flujo de proyecto, conducida por el arquitecto o líder técnico. Úsala SIEMPRE que exista un SRS cerrado y haya que "definir la arquitectura", "elegir el stack", "armar los ADRs", "diseñar el sistema", "resolver el handoff a arquitectura" o "pasar de los requisitos al diseño", aunque no usen esas palabras exactas. NO es el encuadre de negocio (eso es discover-project, Etapa 0), NO elicita ni cambia requisitos (eso es elicit-requirements, Etapa 1), NO descompone en historias ni tickets (eso es Etapa 3), y NO implementa ni planifica el trabajo de un ticket dentro de un repo existente (para eso están close-requirement y plan-implementation).'
---

# Diseño de Arquitectura

## Propósito

Ayudar al arquitecto a transformar un **SRS reducido cerrado** (Etapa 1) en un **documento de arquitectura**: el diseño técnico resuelto, con el _cómo_ de cada requisito, listo para ser firmado en el gate y consumido por quien implementa.

El artefacto tiene **dos consumidores y eso condiciona toda su forma**:

- **El gate humano.** El arquitecto (u otra persona que revise después) tiene que poder ver qué se decidió, por qué, y qué alternativas se descartaron. Es la revisión y validación de diseño que pide ISO 9001 cl. 8.3.
- **El implementador, que normalmente es un modelo de IA.** Lee el documento **en cada ticket** y necesita lo contrario que el gate: el estado resuelto, afirmativo, sin alternativas ni deliberación.

Esa tensión se resuelve con **un solo archivo de dos partes**: el **cuerpo**, escrito en afirmativo y sin argumentar, y el **anexo de decisiones**, donde vive el porqué. El cuerpo referencia al anexo por ID; no lo repite.

Por eso, al redactar:

- No optimices para brevedad ni para elegancia argumental.
- Optimizá para **estado resuelto sin ambigüedad, complejidad mínima y trazabilidad hacia el SRS**.

### Por qué el cuerpo no argumenta

No es una preferencia de estilo. Si el cuerpo dice "consideramos A, B y C; elegimos B", **A y C entran igual al contexto del implementador como opciones plausibles** y el resultado es un híbrido que nadie decidió. Peor todavía las consecuencias con hedges ("habría que revisarlo si la escala crece"): para el arquitecto son honestidad, para el implementador son permiso.

La deliberación no se pierde —va al anexo, completa— pero no se mezcla con el estado.

---

## Límite de la etapa (importante)

### No mira hacia atrás (Etapa 1)

No re-litiga requisitos: no cambia un target, no relaja una regla de negocio, no agrega una capacidad, no reinterpreta un criterio de aceptación. El SRS es la entrada de diseño ya revisada.

Hay **dos motivos legítimos de rebote a Etapa 1**, y conviene distinguirlos porque se resuelven distinto:

- **Rebote por hueco.** El SRS no alcanza para diseñar: un RNF sin métrica, una regla que no cubre un caso que el diseño necesita, un contrato de integración que nadie definió. Falta información.
- **Rebote por insatisfacibilidad.** El requisito está bien escrito y el diseño muestra que **no se puede cumplir al costo aceptable**, o que dos requisitos juntos son incompatibles en cualquier diseño posible. El caso típico: un RNF de disponibilidad alto sobre un sistema cuya única fuente de verdad es un servicio externo — la disponibilidad del producto queda acotada por la del tercero, y ningún diseño lo arregla.

El segundo es el hallazgo más valioso que produce esta etapa, y es la razón por la que existe el rebote: **la Etapa 1 no tiene forma de descubrirlo sola**, porque recién aparece cuando alguien intenta diseñar. No lo absorbas bajando la vara en silencio ni lo archives como riesgo: volvé a Etapa 1 con la nota, o escalalo al decisor como decisión explícita de cambiar el requisito.

### No mira hacia adelante (Etapas 3 y 4)

No produce:

- historias de usuario, escenarios Given-When-Then ni tickets (eso es Etapa 3),
- código, nombres de funciones, migraciones concretas ni pseudocódigo de la lógica (eso es Etapa 4).

**El triage de granularidad se invierte respecto de las etapas anteriores.** En Etapa 0 y 1 el detalle técnico era contaminación y se mandaba al Anexo. Acá el detalle técnico **es el entregable**: el modelo de datos, los contratos y las convenciones son justamente lo que esta etapa debe cerrar. Lo que hay que interceptar es un nivel más abajo: el detalle de **implementación**. La regla práctica: definí la forma y el contrato, no el cuerpo de la función.

### Guard central: "cómo, nunca qué ni cuánto"

Es el espejo del guard de la Etapa 1. Si te descubrís pensando "este target es demasiado exigente, lo bajo a…", "esta regla no tiene sentido, la cambio" o "conviene agregar también la capacidad de…", esa es la señal de que te saliste del espacio del diseño. Ninguna de esas tres cosas es una decisión de arquitectura: son rebote a Etapa 1 o decisión del decisor.

### Guard de complejidad (propio de esta etapa)

Este flujo acepta deliberadamente verbosidad documental, porque con un implementador automático la prosa de más casi no tiene costo y la ambigüedad tiene un costo altísimo. **Ese razonamiento vale para la prosa y no vale para la estructura.** Un documento de 400 líneas no cuesta nada; un componente de más cuesta código, tests y superficie de falla para siempre.

La regla es entonces asimétrica: **detalle documental ilimitado, complejidad estructural mínima demostrable contra los requisitos.** Todo componente, capa, cola o servicio que no puedas justificar contra un requisito concreto, sacalo.

---

## Contexto

### Entrada obligatoria: el SRS reducido

El SRS **no es degradable**. Es la entrada de registro y la cadena de trazabilidad pasa por él. Si esta skill se invoca sin el SRS (ruta del archivo o su contenido), **detenete y pedilo**. No lo reconstruyas de memoria ni lo inventes.

Cuando lo tengas, leelo **completo**. Es un error tomar solo su sección `Handoff a Etapa 2`: ese handoff es el índice más visible del trabajo que te toca, pero **no es el único lugar donde el SRS te asigna trabajo** (ver paso 2). Necesitás además los requisitos funcionales y sus criterios de aceptación, las reglas de negocio, los datos del dominio, la matriz de roles, los RNF con sus métricas, las integraciones, el comportamiento ante fallos, el Anexo de insumos técnicos y la sección de fuera de alcance.

### Esta skill SÍ usa `company_capabilities.md`

Es la primera etapa del flujo que lo necesita, porque es la primera que elige tecnología. Viaja junto con la skill:

    references/company_capabilities.md

Su función es **anclar el diseño en la realidad de la empresa**: con qué stacks trabaja el equipo, dónde despliega, qué integraciones ya domina, qué restricciones regulatorias maneja, y qué explícitamente **no** hace.

A diferencia del SRS, este documento es degradable — pero degrada de una forma particular: **no te detiene, te cambia el modo de operar**. Donde el doc responde, se confirma un anclaje; donde no responde, se pregunta o se decide. Cómo tratar cada hueco depende de **por qué** está en blanco, y tratarlos igual es donde se rompe:

1. **En blanco, pero la empresa sí tiene respuesta.** Es el caso más frecuente y el más peligroso: cloud, CI/CD, orquestación, base de datos. La empresa hoy despliega en algún lado. **Acá preguntá, no propongas.** Si proponés un default, inventás una decisión que ya estaba tomada, alguien la firma por inercia y queda un documento que contradice la operación real. Un anclaje adivinado y un anclaje real se ven exactamente igual de prolijos en el documento final; la única diferencia la hace haber preguntado.
2. **En blanco porque la empresa no tiene esa capacidad y el proyecto la necesita.** Acá sí proponé, pero la salida no es "usamos X": es una decisión marcada como **capacidad nueva a adquirir**, con la curva de aprendizaje como consecuencia aceptada. Eso es lo que el gate necesita ver y lo que alimenta cualquier reestimación. Caso límite: si lo que el SRS pide cae dentro de las **anti-capacidades** declaradas del doc ("no hacemos hardware", "no hacemos modelos de ML a medida"), no es un default a proponer — es un escalamiento al decisor.
3. **Lleno pero no alcanza para este proyecto.** Tienen contenedores pero no orquestación, y un RNF de disponibilidad pide más. Es una decisión real: pregunta de trade-off y decisión registrada en el anexo.

Si el documento no existe (por ejemplo, la skill se instaló sin bundled resources), no te detengas: avisá una sola vez que vas a tener que preguntar todo lo que el doc habría anclado, y seguí.

### Subproducto: devolverle al doc lo que aprendiste

Todo hueco que se cierre preguntando tiene que volver al documento de capacidades, o el doc nunca se completa: la skill lo suple en cada proyecto, el siguiente vuelve a preguntar lo mismo y `company_capabilities.md` queda para siempre como un formulario que nadie llena porque nadie lo necesita.

Al cerrar, entregá una lista de **actualizaciones sugeridas al doc de capacidades** con lo que se respondió. Va **en el chat, no en el artefacto**: no es parte del diseño de este proyecto. Es además el ciclo de mejora continua de ISO 9001 cl. 10 funcionando de verdad.

### No hay codebase

Igual que en las etapas anteriores, el producto todavía no existe: no hay código que inspeccionar. Las fuentes de anclaje son el SRS, el doc de capacidades y el criterio del arquitecto.

---

## Rol de quien ejecuta esta skill

La conduce el **arquitecto o líder técnico**. Puede ser la misma persona que facilitó las Etapas 0 y 1; eso está bien y es lo habitual en equipos chicos. La firma del gate la puede revisar otra persona después, y el artefacto tiene que estar escrito para que eso sea posible: quien lo revise sin haber estado en la conversación tiene que poder entender qué se decidió y contra qué requisito.

**Guard obligatorio.** Acá el sesgo del perfil técnico no es saltar a la solución —esta etapa *es* la solución— sino **sobre-diseñar**: anticipar escala que nadie pidió, agregar capas por completitud estética, elegir la herramienta interesante sobre la adecuada. Tu criterio se usa para:

- elegir el diseño más simple que satisface los requisitos escritos,
- detectar cuándo un requisito no es satisfacible al costo aceptable,
- cerrar las micro-decisiones que, si quedan abiertas, el implementador va a inventar distinto en cada ticket.

No se usa para mejorar los requisitos, ampliar el alcance ni preparar el sistema para un futuro que el SRS no declaró. Lo que el SRS puso fuera de alcance, se queda afuera; la única concesión legítima al futuro es **no cerrarle la puerta**, y solo cuando el SRS o su handoff lo pidan explícitamente.

---

## Comportamiento

### 1. Cargar y validar las entradas

Si no tenés el SRS (ruta o contenido), **detenete y pedilo** antes de cualquier otra cosa. Con el SRS en mano, leelo completo. Después leé `references/company_capabilities.md`.

### 2. Barrido del trabajo a Etapa 2 (obligatorio)

El SRS te asigna trabajo desde **tres secciones distintas**, y la que se llama handoff es habitualmente la menos completa de las tres. Recorré las tres antes de armar nada:

1. **`Handoff a Etapa 2`** — la lista explícita de decisiones técnicas que la Etapa 1 no resolvió.
2. **`Supuestos a validar`** — buscá los que tienen forma de validación "confirmar en el gate de arquitectura" o equivalente. Son típicamente los **números que el facilitador propuso y el negocio no confirmó**: targets de performance, RPO/RTO, escala, umbrales, límites. Importan enormemente porque **determinan el costo de tu diseño** y ninguno figura en el handoff. Un RPO de 24 h y uno de 1 h producen arquitecturas distintas, y si el número es un supuesto sin confirmar estás diseñando sobre arena.
3. **`Riesgos`** — buscá los que tienen mitigación asignada a la Etapa 2 ("dimensionar en Etapa 2", "cuantificar el impacto sobre tal RNF", "espiga técnica temprana"). Es trabajo de arquitectura asignado por nombre desde una sección que no se llama handoff.

Sumá lo que salga de leer el cuerpo del SRS aunque nadie lo haya señalado: toda regla de negocio, criterio de aceptación o modo de falla implica algo del diseño.

#### Clasificá cada ítem por naturaleza

Los ítems que salen del barrido son de cuatro tipos y **se tratan distinto**. Mezclarlos en una lista plana es lo que hace que un handoff se lea como una pila de temas en vez de un plan:

| Naturaleza | Qué es | Destino |
|---|---|---|
| **Decisión** | Hay que elegir entre opciones y registrar la elección | Anexo de decisiones (`AD-…`) y su reflejo en el cuerpo |
| **Espiga** | No se puede decidir sin investigar o probar antes | `ESP-…` con criterio de salida y qué decisión desbloquea |
| **Restricción de evolución** | No es una decisión a tomar sino una fuerza a respetar ("no cerrarle la puerta a X") | Califica varias decisiones; se declara una vez y se referencia |
| **Gestión** | Reestimar esfuerzo, decidir un rollout escalonado, fechas | **No es arquitectura.** Señalalo al facilitador y sacalo del alcance del documento |

### 3. Verificar la integridad del SRS

Si para diseñar te falta algo que debía venir cerrado (un RNF sin métrica, una regla que no cubre un caso que el diseño necesita, un contrato que nadie definió), **rebotá a Etapa 1** con la nota de qué falta y por qué bloquea el diseño, en vez de resolverlo acá. Rellenar un hueco de requisitos con una decisión de arquitectura rompe la trazabilidad: queda una regla de negocio inventada por el arquitecto y sin firma de nadie.

### 4. Armar el grafo de decisiones y mostrarlo

Acá esta etapa se aparta de la Etapa 1, y por un motivo estructural: **las decisiones de arquitectura están acopladas.** Las dimensiones de la elicitación son casi independientes entre sí —se puede cerrar la matriz de roles sin haber tocado los RNF—, mientras que acá el stack condiciona la persistencia, que condiciona cómo se garantiza la unicidad de un campo opcional, que condiciona la forma de un identificador. Preguntar en orden arbitrario obliga a reabrir lo ya cerrado.

Antes de preguntar nada, armá el grafo:

- **Decisiones raíz** — las que condicionan a otras. Típicamente stack y plataforma, el modelo de persistencia, y el diseño de la integración central cuando el sistema depende de un tercero.
- **Decisiones derivadas** — las que solo tienen sentido una vez resuelta su raíz.

Y clasificá cada una **por origen**, contra el doc de capacidades:

- **Anclada** — el doc la responde. No se pregunta: se enuncia para confirmar, en lote.
- **Decisión real** — el doc no la responde o no alcanza (los tres casos de hueco de la sección Contexto).
- **Bloqueada por espiga** — no se decide en esta corrida.

**Mostrale el mapa al facilitador antes de arrancar** con las preguntas: cuántas decisiones salieron, cuáles son raíz, cuáles quedan ancladas y cuáles esperan una espiga. Sirve para dos cosas: le da la forma de la conversación que viene, y hace visible un error de clasificación tuyo antes de que cueste caro.

### 5. Confirmar los anclajes en un solo lote

Las decisiones ancladas van juntas, en un solo mensaje, como confirmación y no como pregunta de trade-off. Nadie quiere deliberar sobre algo que la empresa ya tiene decidido, y montar una pregunta abierta sobre un anclaje invita a cambiarlo sin motivo.

Si el facilitador corrige un anclaje, actualizalo y anotalo para el subproducto de capacidades: significa que el doc está desactualizado.

### 6. Recorrer las decisiones

Hacé las preguntas en el mismo idioma que use el facilitador.

#### Ritmo

- **Decisiones raíz: de a una.** Cada respuesta reordena o elimina decisiones derivadas, así que agruparlas obliga a reabrir. Después de cerrar cada raíz, **recalculá qué queda pendiente** antes de preguntar lo siguiente.
- **Decisiones derivadas: en tandas de hasta 4, solo si son independientes entre sí.** Si dos derivadas se condicionan, no van en la misma tanda.

Señal de diagnóstico: si una tanda produce respuestas que se contradicen entre sí, clasificaste mal y alguna de esas era en realidad una decisión raíz. Reordená el grafo y seguí.

#### Guard contra el trade-off falso

En la Etapa 1 las opciones A y B eran las dos legítimas porque decidía el cliente. Acá, muchas veces **el SRS ya determina la respuesta**: si un RNF exige recuperación ante pérdida de datos, no hay elección sobre si existe backup. Montar una pregunta de trade-off igual invita a elegir arbitrariamente algo que los requisitos ya decidieron.

La regla: **si el SRS determina la respuesta, enunciá la derivación y pedí confirmación** ("RNF-06 pide RPO ≤ 24 h, así que el diseño incluye X; confirmame"). Reservá la pregunta de trade-off para cuando los requisitos genuinamente no deciden. Un trade-off inventado es una alucinación con formato de diligencia.

Cuando sí preguntes, incluí un **default anclado** en el SRS o en el doc de capacidades, y explicá en una línea por qué encaja.

#### Las tres salidas de una decisión

Esta etapa cambia la regla de cierre-o-supuesto de las anteriores, porque cambia la naturaleza de lo desconocido: en descubrimiento y requisitos lo que faltaba dependía del cliente y por eso se suponía; **acá lo desconocido es técnico y se resuelve investigando, no suponiendo**. Toda decisión termina en una de tres:

1. **Decidida** → se registra como `AD-<slug>-NN` en el anexo y su estado resuelto va al cuerpo.
2. **Diferida a espiga** → `ESP-<slug>-NN`, con **criterio de salida explícito** (qué hay que saber para decidir) y **qué decisión desbloquea**. Una espiga sin criterio de salida es una decisión postergada sin fecha.
3. **Rebote a Etapa 1** → cuando el requisito resulta un hueco o no es satisfacible al costo aceptable.

El supuesto sigue existiendo, pero **solo para lo que depende de una confirmación del decisor**: típicamente los números heredados del SRS que este gate tenía que validar. No lo uses para tapar una decisión técnica que no quisiste tomar.

#### Dimensiones obligatorias

Tus decisiones deben cubrir colectivamente:

1. **Stack y plataforma.** Lenguajes, frameworks, runtime, tipo de aplicación. Anclado en el doc de capacidades siempre que se pueda.
2. **Modelo de datos.** Entidades, atributos, tipos, relaciones, restricciones de integridad, esquema de persistencia. La Etapa 1 excluyó esto a propósito ("qué significa, nunca tipos ni tablas"), y por eso es **probablemente la sección de mayor rendimiento del artefacto**: cada entidad y atributo que no esté escrito acá, el implementador lo inventa, y lo inventa distinto en cada ticket.
3. **Componentes y sus límites.** Qué partes tiene el sistema, de qué responde cada una, y qué requisitos atiende. Mínimo demostrable: ver el guard de complejidad.
4. **Contratos entre componentes.** Todo límite que cruza un requisito funcional necesita su contrato escrito: qué entra, qué sale, qué errores.
5. **Convenciones transversales.** Las micro-decisiones que, si no están cerradas, el implementador resuelve distinto en cada ticket. La lista concreta la dicta el SRS, no un checklist fijo, pero estas aparecen casi siempre y valen como sonda: **representación de dinero** (monto, moneda, precisión) cuando hay precios; **fecha, hora y zonas horarias** cuando hay cálculos temporales o conversión; **generación de identificadores** cuando el SRS pide un identificador propio; **unicidad de campos opcionales**; **soft delete**, cuando el SRS repite "no existe operación que elimine" sobre varias entidades —eso es un patrón transversal, se declara una vez y no se redescubre entidad por entidad—; **forma de los errores**; **dónde valida qué**; **estructura del repositorio**.
6. **Mecanismo de cada requisito no funcional.** Uno por uno, con el mecanismo nombrado. No aceptes un "en particular los más críticos": un RNF sin mecanismo es un RNF que nadie va a cumplir.
7. **Integraciones.** El contrato técnico de cada una. Si el cliente lo impuso, ya está preservado en el Anexo del SRS: referencialo, no lo rediseñes.
8. **Entornos y despliegue.** Entornos, cómo se aplican los cambios de esquema, dónde vive la configuración y cómo se manejan los secretos. Rara vez aparece en el SRS —no es un requisito de negocio— y es justamente donde el implementador inventa con más libertad porque nadie le dijo nada.

### 7. Chequeo de consistencia (antes de la pregunta de cierre)

Cada decisión puede ser razonable por separado y el conjunto igual estar mal. Los pasos anteriores revisan cada decisión contra su requisito; este revisa el conjunto contra sí mismo. Recorré los cinco controles:

1. **Cobertura bidireccional requisitos ↔ componentes.** Hacia adelante: **todo RF must-v1 mapea a un conjunto nombrado de componentes que existen en el diseño, y todo límite que ese RF cruza tiene su contrato definido**. Un RF que no mapea a nada es diseño incompleto, y es exactamente lo que después impide que la Etapa 3 sepa si ese requisito es un ticket o cuatro. Hacia atrás: **todo componente traza a algún requisito**. Un componente al que no traza nada es arquitectura especulativa y se saca — es el análogo del "requisito que no traza a nada es scope creep" de la Etapa 1, y es el sentido que nadie corre.
2. **Todo RNF con mecanismo, y los mecanismos entre sí.** Verificá primero que ningún RNF quedó sin mecanismo nombrado. Después cruzalos: **el mecanismo de un RNF no puede violar el target de otro**. El caso típico es un mecanismo de frescura de datos que agrega latencia y compite con un target de tiempo de respuesta. Si compiten, la resolución se decide y se registra; no queda como tensión implícita.
3. **Fuente única por convención.** Cada convención se enuncia en **un solo lugar** del documento; el resto la referencia. Si el manejo de fechas está dicho en el modelo de datos y otra vez en convenciones transversales, no son dos secciones prolijas: son dos verdades que van a divergir, y el implementador no tiene forma de saber cuál manda.
4. **Ningún contrato que viva solo en el diagrama.** Los diagramas sirven al gate humano y casi nada al implementador. Todo límite entre componentes que cruza un requisito está escrito en prosa afirmativa en el cuerpo, además de dibujado.
5. **Coherencia cuerpo ↔ anexo.** Ninguna afirmación del cuerpo contradice la decisión que la sostiene, y toda decisión del anexo se refleja en el cuerpo. Este control es **propio de haber unificado los dos documentos en uno**: en un esquema de archivos separados esta grieta no existe, acá sí, y es silenciosa porque las dos partes se leen bien por separado.

Lo que encuentres se resuelve acá. Si la resolución requiere una decisión del decisor, escalala; si requiere cambiar un requisito, es rebote a Etapa 1. Una contradicción entre decisiones no es un riesgo a documentar: es un defecto.

### 8. Pregunta de cierre (obligatoria, siempre la última antes del gate)

Con todas las decisiones cerradas, diferidas a espiga o rebotadas, y **antes** del gate, hacé siempre esta pregunta, sin excepción:

"Antes de armar el documento de arquitectura: ¿hay alguna decisión, restricción o convención que quieras agregar, corregir o aclarar?"

- Si agrega o cambia algo → volvé al paso 6 y **volvé a correr el chequeo de consistencia** del paso 7.
- Si no hay nada más → recién ahí pasá al gate del paso 9.

### 9. Confirmar antes de escribir (gate)

Cuando todo esté cerrado, preguntá:

"La arquitectura está cerrada. ¿Querés que redacte el documento?"

No lo escribas todavía.

### 10. Redactar y guardar el archivo

Solo si el facilitador confirma explícitamente en el paso 9, escribí el artefacto con el schema de abajo **y guardalo como archivo**, no solo en el chat.

**Nombre del archivo:**

    arquitectura-<slug-del-proyecto>-<YYYY-MM-DD>.md

Usá el **mismo slug del Brief y del SRS** para mantener la continuidad de la cadena. **Ubicación:** la raíz del proyecto/repositorio donde se ejecuta la skill, salvo que el facilitador pida otra.

Al redactar, sostené la separación que justifica todo el diseño del artefacto: **el cuerpo afirma, el anexo argumenta**. En el cuerpo no escribas "elegimos", "consideramos", "podría", "en el futuro quizás". Escribí el estado.

Después de guardarlo, confirmá el nombre y la ubicación.

### 11. Handoff a Etapa 3 y subproducto de capacidades

El documento termina con el handoff a la Etapa 3: qué necesita la descomposición en tickets (orden de construcción, dependencias entre componentes, qué RF conviene partir y por dónde). Así la cadena no se corta.

Y **en el chat, fuera del artefacto**, entregá la lista de actualizaciones sugeridas a `company_capabilities.md`.

---

## Convención de IDs

Igual que en las etapas anteriores, los identificadores llevan el slug del proyecto para que sean únicos entre proyectos:

- Decisión de arquitectura: `AD-<slug>-01`, `AD-<slug>-02`, …
- Componente: `C-<slug>-01`, …
- Espiga: `ESP-<slug>-01`, …

Las afirmaciones del cuerpo referencian su decisión entre paréntesis —`(AD-<slug>-04)`— igual que los RF del SRS referencian sus RN. Las decisiones referencian los requisitos que las obligan (`RF-…`, `RN-…`, `RNF-…`), nunca al revés: el SRS no se modifica desde acá.

Las decisiones **no se editan**: si una decisión posterior la reemplaza, la nueva declara a cuál supersede y la anterior queda con estado `supersedida por AD-<slug>-NN`. Es lo que conserva el rastro de por qué la arquitectura es como es, y lo que hace auditable la revisión de diseño.

---

## Output

### Si falta el SRS

Respondé en el idioma del usuario, pidiéndolo, sin avanzar:

## Falta la entrada obligatoria

Esta etapa necesita el SRS reducido cerrado de la Etapa 1 (ruta del archivo o su contenido). Sin él no puedo diseñar con trazabilidad. Pasámelo y arranco.

---

### Si el SRS tiene un hueco o un requisito no satisfacible

## Rebote a Etapa 1 (Requisitos)

El SRS no permite cerrar el diseño todavía:

- **<qué falta o qué no es satisfacible>** — Requisito afectado: <ID>. Por qué bloquea: <...>. Tipo: <hueco | no satisfacible al costo aceptable>.

Sugiero volver a `elicit-requirements` para resolver esto, o escalarlo al decisor si implica cambiar lo que el v1 se comprometió a lograr.

---

### Mapa de decisiones (antes de preguntar)

## Trabajo a resolver en esta etapa

Del SRS salieron <N> decisiones (Handoff a Etapa 2: <n>; supuestos a confirmar en este gate: <n>; riesgos con mitigación en Etapa 2: <n>; derivadas del cuerpo: <n>).

- **Raíz (<n>):** <...>
- **Derivadas (<n>):** <...>
- **Ancladas en el doc de capacidades (<n>):** <...>
- **Bloqueadas por espiga (<n>):** <...>
- **No son arquitectura (<n>):** <ítem de gestión> → para el facilitador

Arranco confirmando los anclajes y sigo con las raíz de a una.

---

### Anclajes a confirmar

## Anclajes heredados del estándar de la empresa

Estos no los pregunto como decisión porque el doc de capacidades ya los responde; confirmame que siguen vigentes:

- <decisión> — <valor del doc de capacidades>
- <...>

---

### Si todavía hay decisiones abiertas

## Estado

Decisiones cerradas: <...>
Abiertas: <...>

## Decisión <raíz | derivada>: <nombre>

<pregunta de trade-off, o derivación a confirmar si el SRS ya la determina>

   Requisitos que la obligan: <IDs>
   Default sugerido:
   <opción anclada en el SRS o en el doc de capacidades + por qué encaja>

---

### Si todo está cerrado pero no confirmado

## Estado

Todas las decisiones están cerradas, diferidas a espiga con criterio de salida, o rebotadas a Etapa 1. El chequeo de consistencia pasó: todo RF must-v1 mapea a componentes, ningún componente quedó sin requisito, todo RNF tiene mecanismo y cada convención se enuncia una sola vez.

## Confirmación

¿Querés que redacte el documento de arquitectura?

---

### Si está confirmado

# Arquitectura: <título del proyecto>

## Metadata

- Versión: 1.0
- Dueño (arquitecto): <nombre / rol>
- Fecha: <fecha>
- Estado: Borrador para gate de arquitectura
- **Firma del gate:** _pendiente_ → al aprobarse: `aprobado por <nombre>, <fecha>`
- **SRS de origen:** <nombre del archivo> (v<versión>)
- **Brief de origen:** <nombre del archivo> (v<versión>)

## Resumen del enfoque

<3-4 líneas: qué forma tiene el sistema y qué lo condiciona. Para que el documento se entienda sin tener el SRS al lado. No re-litiga requisitos.>

## Anclajes heredados del estándar de la empresa

> Lo que no se decidió en este proyecto porque ya estaba decidido. Existe para que el gate sepa qué le están pidiendo firmar: confirmar un estándar es distinto de aprobar una decisión nueva.

- <decisión> — <valor> | Fuente: `company_capabilities.md`

## Decisiones propias de este proyecto

> Índice de lo que sí se decidió acá. El detalle está en el anexo.

- **AD-<slug>-01** — <título de la decisión> | Requisitos que la obligan: <IDs>
- <...>

## Supuestos del SRS que este gate debía confirmar

> Los números que el SRS dejó como supuesto con validación en el gate de arquitectura, y que condicionan este diseño. Cada uno con destino explícito: el diseño no se puede firmar dejándolos en silencio.

| Supuesto (SRS) | Valor usado en este diseño | Estado | Consecuencia si cambia |
|---|---|---|---|
| <S-NN> | <valor> | <confirmado por <quién>, <fecha> \| sigue abierto> | <qué parte del diseño se rehace> |

## Stack y plataforma

<afirmativo>

## Modelo de datos

> Entidades, atributos, tipos, relaciones y restricciones de integridad. Es lo que la Etapa 1 excluyó a propósito y lo que el implementador necesita con más precisión.

<afirmativo>

## Mapa de componentes

| Componente | Responsabilidad | Requisitos que atiende |
|---|---|---|
| **C-<slug>-01** — <nombre> | <de qué responde> | <RF-…, RNF-…> |

## Contratos entre componentes

> Todo límite que cruza un requisito funcional. Ninguno vive solo en el diagrama.

- **<C-…> → <C-…>** — <qué entra, qué sale, qué errores> | Afecta: <RF-…>

## Convenciones transversales

> Las micro-decisiones que, si no están escritas, se resuelven distinto en cada ticket. Cada una enunciada una sola vez.

- **<convención>** — <regla afirmativa> `(AD-<slug>-NN)`

## Cumplimiento de requisitos no funcionales

| RNF | Métrica del SRS | Mecanismo | Decisión |
|---|---|---|---|
| <RNF-…> | <target> | <cómo se cumple> | `(AD-<slug>-NN)` |

## Integraciones

- **<sistema>** — <contrato técnico> | Contrato impuesto por el cliente: <sí → ver Anexo del SRS \| no → definido acá, `(AD-…)`>

## Entornos y despliegue

<entornos, cambios de esquema, configuración y secretos — afirmativo>

## Trazabilidad requisitos → componentes

> Cierra el test de salida de esta etapa en las dos direcciones y es el insumo que le permite a la Etapa 3 decidir la granularidad de los tickets.

| Requisito | Componentes | Contratos que cruza |
|---|---|---|
| <RF-…> | <C-…, C-…> | <...> |

## Espigas pendientes

- **ESP-<slug>-01** — <qué hay que averiguar> | Criterio de salida: <qué hay que saber para decidir> | Desbloquea: <decisión> | Requisito en juego: <ID>

## Rebote a Etapa 1

> Solo si el diseño encontró un hueco o un requisito no satisfacible al costo aceptable. Si no hubo ninguno, omitir la sección.

- <requisito> — Tipo: <hueco \| no satisfacible> | Qué mostró el diseño: <...> | Destino: <vuelve a Etapa 1 \| decisión del decisor, registrada acá>

## Diagramas

> Contexto y contenedores/componentes, en mermaid. Sirven al gate humano; no son la fuente de ningún contrato.

<diagramas>

## Handoff a Etapa 3 (Historias + backlog)

- Orden de construcción sugerido y por qué: <...>
- Dependencias entre componentes que condicionan la secuencia: <...>
- Requisitos que conviene partir en más de un ticket, y por dónde: <...>
- Espigas que deben resolverse antes de tocar cierto requisito: <...>

## Anexo — Decisiones de arquitectura

> Acá vive la deliberación completa. Las decisiones no se editan: se supersedan.

### AD-<slug>-01 — <título>

- **Estado:** <propuesta | aceptada | supersedida por AD-<slug>-NN>
- **Requisitos que la obligan:** <RF-…, RN-…, RNF-…>
- **Contexto:** <qué fuerzas obligan a decidir, en términos de los requisitos>
- **Decisión:** <qué se eligió>
- **Alternativas consideradas:** <opción — por qué se descartó>
- **Consecuencias:** <qué se gana, qué se acepta como costo>
- **Tipo:** <anclada en el estándar | decisión de este proyecto | capacidad nueva a adquirir>

---

## Reglas

- Sin SRS reducido, no arranques: es entrada obligatoria y no degradable. Pedilo y detenete. Leelo completo, no solo su `Handoff a Etapa 2`.
- Barré el trabajo a Etapa 2 desde las **tres** secciones del SRS: el handoff, los supuestos con validación en este gate y los riesgos con mitigación asignada acá. La sección que se llama handoff es habitualmente la menos completa.
- Clasificá cada ítem por naturaleza: decisión, espiga, restricción de evolución o gestión. Lo de gestión no es arquitectura: señalalo y sacalo del documento.
- No re-litigues requisitos. Si falta algo, es rebote a Etapa 1 por hueco; si el requisito no es satisfacible al costo aceptable, es rebote por insatisfacibilidad o escalamiento al decisor. Nunca bajes la vara en silencio ni lo archives como riesgo.
- Guard: "cómo, nunca qué ni cuánto". Si te descubrís cambiando un target, agregando una capacidad o relajando una regla, te salió del alcance de esta etapa.
- Guard de complejidad: detalle documental ilimitado, complejidad estructural mínima demostrable. Todo componente sin requisito que lo justifique se saca.
- Armá el grafo de decisiones antes de preguntar y mostraselo al facilitador. Raíz de a una, recalculando lo pendiente después de cada una; derivadas en tandas de hasta 4 y solo si son independientes. Si una tanda se contradice a sí misma, alguna de esas era raíz.
- Las decisiones ancladas en `company_capabilities.md` se confirman en un solo lote, no se preguntan como trade-off.
- Un hueco del doc de capacidades **se pregunta cuando la empresa tiene la respuesta** (cloud, CI/CD, orquestación, base de datos): proponer un default ahí inventa una decisión que ya estaba tomada y queda indistinguible de un anclaje real. Se propone solo cuando la capacidad no existe, y entonces se marca como capacidad nueva a adquirir. Si cae en las anti-capacidades declaradas, se escala al decisor.
- Si el SRS ya determina la respuesta, enunciá la derivación y pedí confirmación. Reservá el trade-off para cuando los requisitos genuinamente no deciden: un trade-off inventado es una alucinación con formato de diligencia.
- Toda decisión termina decidida (`AD-…`), diferida a espiga (`ESP-…` con criterio de salida y qué desbloquea) o rebotada a Etapa 1. El supuesto queda solo para lo que depende de una confirmación del decisor, no para tapar una decisión técnica.
- El cuerpo del documento afirma; el anexo argumenta. En el cuerpo no van alternativas descartadas ni hedges: entran al contexto del implementador como opciones plausibles y como permiso.
- Las decisiones no se editan: se supersedan, declarando a cuál reemplazan.
- Cerrá las convenciones transversales que el SRS haga necesarias (dinero, fecha/hora y zonas, generación de identificadores, unicidad de campos opcionales, soft delete, forma de los errores, dónde valida qué, estructura del repo). Cada una enunciada una sola vez; lo que quede abierto se inventa distinto en cada ticket.
- Cerrá siempre entornos y despliegue, aunque el SRS no diga nada: no es un requisito de negocio y es donde el implementador inventa con más libertad.
- Antes de la pregunta de cierre corré el chequeo de consistencia: cobertura bidireccional requisitos ↔ componentes, todo RNF con mecanismo y los mecanismos sin competir entre sí, fuente única por convención, ningún contrato que viva solo en el diagrama, y coherencia cuerpo ↔ anexo. Una contradicción entre decisiones es un defecto, no un riesgo.
- Todo RF must-v1 mapea a componentes que existen y a los contratos que cruza; todo componente traza a un requisito. Sin ese mapa la Etapa 3 no puede decidir la granularidad de los tickets.
- Los diagramas sirven al gate humano. Ningún contrato existe solo ahí.
- No escribas historias, tickets, código ni pseudocódigo de la lógica. Definí la forma y el contrato, no el cuerpo de la función.
- No redactes el documento hasta que el facilitador confirme.
- Siempre hacé la pregunta de cierre ("¿algo más?") antes del gate de confirmación, sin excepción.
- El documento final se guarda como archivo en la raíz del proyecto, con nombre `arquitectura-<slug>-<fecha>.md` y el mismo slug del Brief y del SRS; no alcanza con mostrarlo solo en el chat.
- Al cerrar, entregá **en el chat** las actualizaciones sugeridas a `company_capabilities.md` con lo que se respondió en la conversación. Sin eso el doc nunca se completa.
- Respondé siempre en el idioma del usuario.
- Optimizá para estado resuelto sin ambigüedad y complejidad mínima, no para verbosidad argumental.
