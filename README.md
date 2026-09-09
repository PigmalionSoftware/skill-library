# Skills

Colección de **skills** (habilidades) para clientes de agentes de IA (como Claude Code, entre otros). Cada skill extiende al agente con conocimiento y flujos de trabajo especializados que se activan automáticamente según la petición del usuario o de forma manual mediante un comando (`/nombre-skill`).

Cada skill vive en su propia carpeta dentro de `skills/` y se define en un archivo `SKILL.md`.

## Skills disponibles

| Skill | Descripción |
|-------|-------------|
| **caveman** | Modo de comunicación ultracomprimido. Reduce el uso de tokens (~75%), manteniendo toda la precisión técnica. Soporta distintos niveles de intensidad. |
| **commit-push** | Hace commit de los cambios en staging y los sube a la rama actual. |
| **deep-research** | Investigación técnica rigurosa. Combina búsqueda web con documentación autorizada (Context7 MCP), verifica afirmaciones entre fuentes y genera un informe estructurado y citado guardado en `research/`. |
| **explain-changes** | Explica en detalle cambios de código y configuración, con diagramas ASCII en la terminal. Explica que cambio, por qué importa, etc |
| **plan-from-spec** | Convierte una especificación o definición de feature en un plan de implementación riguroso y revisado, guardado en `plans/`. Interroga la spec con preguntas hasta no dejar nada librado a suposiciones. |
| **postman-collection-generator** | Genera una colección de Postman (v2.1) importable a partir de una base de código. Escanea las rutas, extrae métodos/paths/params/bodies/headers, agrupa endpoints en carpetas y configura variables de entorno. Agnóstico al lenguaje y framework. |
| **supabase-postgres-best-practices** | Guía de optimización de rendimiento y buenas prácticas de Postgres mantenida por Supabase. Reglas en 8 categorías priorizadas por impacto, para escribir, revisar u optimizar consultas y esquemas. |
| **update-readme** | Revisa los cambios del repositorio y mantiene el README.md alineado con el comportamiento, instalación, uso y compatibilidad verificados. |
| **push-github-tag** | Determina el último tag semántico estable, incrementa su versión patch, crea y publica el nuevo tag anotado en GitHub. |

## Uso

Las skills se activan de dos maneras:

- **Automática**: el agente detecta cuándo una skill aplica según su descripción y la petición del usuario.
- **Manual**: invocando el comando correspondiente, por ejemplo `/commit-push` o `/deep-research`.

## agent-sandbox

Además de las skills, este repositorio incluye **agent-sandbox** (`agent-sandbox/go/`): una
herramienta de línea de comandos que ejecuta un agente de código (**Codex**, **Claude Code**,
**opencode** o **pi**) dentro de un contenedor Docker, sobre un *worktree* de Git aislado. El
agente trabaja en un branch propio, en un directorio separado, sin tocar tu copia de trabajo
actual.

### Instalación

```bash
curl -fsSL https://raw.githubusercontent.com/PigmalionSoftware/skill-library/master/agent-sandbox/go/install.sh | bash
```

El script descarga el binario de la última release, lo instala en `~/.local/bin/agent-sandbox`
y te avisa si ese directorio no está en tu `PATH`. Solo Linux x86_64 por ahora.

### Modo de uso

```
agent-sandbox <branch> --agent <codex|claude|opencode|pi> [--model <modelo>] [--push] <prompt...>
```

**El nombre del branch va siempre primero**, antes de cualquier opción.

| Parámetro | Obligatorio | Descripción |
| --- | --- | --- |
| `<branch>` | Sí | Nombre del branch. El worktree se crea en `../<branch>`. |
| `--agent <codex\|claude\|opencode\|pi>` | Sí | Agente a ejecutar. |
| `--model <modelo>` | No | Modelo a usar. Por defecto: `gpt-5.6-terra` (codex) y `opus` (claude); opencode y pi resuelven el suyo. |
| `--push` | No | Al terminar: `git add -A`, commit con el prompt como mensaje y `git push --set-upstream`. |
| `<prompt...>` | Sí | Instrucción para el agente, entre comillas. |

Requiere Docker, ejecutarse dentro de un repositorio Git, y tener la configuración del agente
ya autenticada en el host (`~/.codex`, `~/.claude`, `~/.config/opencode` o `~/.pi/agent`
según corresponda): las credenciales salen de ahí, nunca de variables de entorno.

