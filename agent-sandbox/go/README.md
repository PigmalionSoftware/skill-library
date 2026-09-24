# agent-sandbox

`agent-sandbox` ejecuta Codex, Claude Code, opencode o pi dentro de un
contenedor Docker y sobre un *worktree* de Git propio. El agente trabaja en un
branch y directorio separados, sin modificar la copia de trabajo desde la que
se invocó el comando.

## Instalación

La release publicada actualmente es para Linux x86_64. Instalala con:

```bash
curl -fsSL https://raw.githubusercontent.com/PigmalionSoftware/skill-library/master/agent-sandbox/go/install.sh | bash
```

El instalador descarga la última release de `agent-sandbox` y deja el binario
en `~/.local/bin`. Si ese directorio todavía no está en `PATH`, el instalador
indica la línea que hay que agregar al perfil de la shell.

Para compilarlo desde este directorio:

```bash
go build -o agent-sandbox ./cmd/agent-sandbox
```

## Requisitos

- Docker en ejecución.
- Git y un repositorio Git: ejecutá el comando desde cualquier directorio
  dentro del repositorio.
- Una sesión ya autenticada del agente elegido en el host. El contenedor recibe
  esa configuración mediante montajes; no recibe credenciales como variables
  de entorno.

| Agente | Configuración requerida en el host |
| --- | --- |
| Codex | `~/.codex/` |
| Claude Code | `~/.claude/` y `~/.claude.json` |
| opencode | `~/.config/opencode/`, `~/.local/share/opencode/` y `~/.local/state/opencode/` |
| pi | `~/.pi/agent/` |

Si falta una de esas rutas, el comando falla sin crearla. Ejecutá y autenticá
primero el agente correspondiente en el host.

En cada ejecución, `agent-sandbox` comprueba la versión del agente seleccionado
dentro de la imagen contra npm. Construye la imagen si no existe y la reconstruye
si encuentra una versión más reciente; por eso la primera ejecución puede tardar
y necesita acceso a npm.

## Uso

```text
agent-sandbox run [-b <branch>] -a <codex|claude|opencode|pi> [-m <modelo>] [-i <imagen>] [--image <archivo>]... [-p] [-c <mensaje-commit>] (-q <consulta> | -f <archivo-prompt>)
agent-sandbox resume -b <branch> -a <codex|claude|opencode|pi> [opciones] (-q <consulta> | -f <archivo-prompt>)
```
| Parámetro | Descripción |
| --- | --- |
| `-b`, `--branch` | Opcional. Branch para el worktree aislado; si se omite, se genera uno. |
| `-a`, `--agent` | Obligatorio. Uno de `codex`, `claude`, `opencode` o `pi`. |
| `-m`, `--model` | Opcional. Sobrescribe el modelo que resuelve el agente. |
| `-i`, `--base-image` | Opcional. Deriva una imagen desde una base compatible, para disponer de su toolchain dentro del sandbox. |
| `-p`, `--push` | Al finalizar, agrega todos los cambios, crea un commit y hace `git push --set-upstream origin <branch>`. |
| `-q`, `--query` | Instrucción para el agente. |
| `-c`, `--commit-message` | Opcional. Mensaje del commit creado por `-p` o `--push`; si se omite, usa el prompt resuelto. Sin `-p` o `--push`, no tiene efecto. |
| `-f`, `--file-prompt` | Archivo cuyo contenido se usa como instrucción para el agente, en lugar de `-q` o `--query`. |
| `--image <archivo>` | Opcional y repetible. Adjunta imágenes al prompt inicial de Codex o Claude Code. Cada ruta debe ser un archivo regular existente en el host; opencode y pi la ignoran. Usá `--` antes del prompt de texto para que Codex no lo interprete como otra imagen. |

Sin `-p` o `--push`, los cambios quedan sin commitear en el worktree. Con `-p`
o `--push`, si el agente no produjo cambios, no se crea ningún commit. Hay que
proporcionar exactamente una fuente de prompt: `-q`/`--query`, `-f`/`--file-prompt`
o el valor `"query"`/`"file-prompt"` del JSON. No se pueden usar ambas fuentes
a la vez ni se aceptan instrucciones posicionales. Si se usa `--push` sin
`--commit-message`, el contenido
del archivo se convierte en el mensaje de commit por defecto cuando se eligió
`--file-prompt`.

Para continuar un worktree registrado, usá `agent-sandbox resume -b <branch>`.
El branch es el nombre mostrado por `worktree-list`. Si se combina con
`--push`, se commitean todos los cambios pendientes.

`run` y `resume` leen `./agent-sandbox.json` si existe en el directorio desde
el que se ejecutan. Cada sección admite los nombres largos de las opciones:
`branch`, `agent`, `model`, `base-image`, `query`, `push`, `commit-message`,
`file-prompt` e `image`. Las opciones explícitas de la línea de comandos
prevalecen sobre el JSON.

```json
{
  "run": {
    "agent": "codex",
    "base-image": "golang:1.26-alpine",
    "query": "run go version and do not change any files",
    "push": false
  },
  "resume": {
    "branch": "fix-login",
    "agent": "codex",
    "query": "add a regression test"
  }
}
```

Con este archivo, `agent-sandbox run` usa la sección `run`; `resume` usa la
sección `resume`. También podés pasar `-q "otra tarea"` para reemplazar la
consulta del JSON. Los comandos de gestión de worktrees no leen el archivo.

### Imagen base externa

`--base-image` usa una imagen que ya trae el runtime del proyecto. Se admiten
Alpine, Debian, Ubuntu, Fedora, RHEL 8/9, UBI 8/9 y Amazon Linux 2023. Por
ejemplo, `golang:1.26-alpine` deja disponibles Go, `gofmt` y `go test` dentro
del contenedor:

```bash
agent-sandbox run -b fix-go-tests -a codex -i golang:1.26-alpine -q "run gofmt and go test ./..., then fix failures"
```

La primera ejecución crea una imagen local derivada e instala lo necesario para
ejecutar los agentes: Node.js, npm, Bash, Codex, Claude Code, opencode, pi,
ripgrep, certificados CA, curl y Git. Las siguientes reutilizan esa imagen para
la misma base.

La base debe ofrecer `apk` (Alpine), `apt-get` (Debian/Ubuntu), `dnf`
(Fedora/RHEL/UBI/Amazon Linux) o `microdnf` (UBI minimal). Las demás fallan
durante el build. Alpine instala Node.js desde sus paquetes; las otras bases
usan Node.js 22 de NodeSource.

En RHEL, la imagen debe tener repositorios habilitados y, si corresponde, una
suscripción válida: el sandbox no monta credenciales del host. No habilita
EPEL ni CRB/CodeReady Builder, ni soporta `yum`; si falta un paquete como
`ripgrep`, el build falla con el error del gestor de paquetes.

La imagen derivada no se actualiza sola. Para reconstruirla, eliminá la imagen
correspondiente:

```bash
docker image ls 'agent-sandbox-base-*'
docker image rm <IMAGE_ID>
```

### Modelos

Al omitir `--model`, Codex usa `gpt-5.6-terra` y Claude Code usa `opus`.
opencode y pi dejan que su propia configuración elija el modelo. Los valores de
`--model` se pasan directamente al agente: por ejemplo, `gpt-5.6-sol` para
Codex, `sonnet` para Claude Code y `proveedor/modelo` para opencode o pi.

## Ejemplos

Crear un worktree para Codex y dejar sus cambios listos para revisar:

```bash
agent-sandbox run -a codex -m gpt-5.6-sol -q "fix the login redirect loop"
```

Ejecutar Claude Code y publicar el branch al terminar:

```bash
agent-sandbox run -b add-test -a claude -m sonnet -p -q "add a regression test for the login redirect"
```

Continuar un worktree registrado por su nombre:

```bash
agent-sandbox resume -b fix-login -a codex -q "add a regression test for the login redirect"
```

Dejar que opencode resuelva su modelo configurado:

```bash
agent-sandbox run -b update-copy -a opencode -q "update the empty-state copy"
```

Ejecutar Codex con imagen de golang

```bash
agent-sandbox run -b fix-go-tests -a codex -i golang:1.26-alpine -q "run go test ./... and fix failures"
```

Obtener prompt de archivo
```bash
agent-sandbox run -b prompt-file-test -a codex -f prompt.md
```

Adjuntar una o más imágenes al prompt inicial de Codex o Claude Code
```bash
agent-sandbox run -a claude \
  --image "/home/user/Pictures/mockup.png" \
  --image "/home/user/Pictures/reference.png" \
  -q "compare these screenshots and implement the resulting UI"
```

## Gestionar worktrees creados por el sandbox

El comando crea los worktrees en
`~/.config/agent-sandbox/worktrees/<repo>-<hash>/<branch>` y registra los que
creó en `~/.config/agent-sandbox/worktrees.jsonl`. El identificador del
repositorio evita que branches con el mismo nombre en repositorios distintos
colisionen. Los siguientes subcomandos solo ven y eliminan esas entradas; no
afectan worktrees creados manualmente.

```bash
# Mostrar los worktrees registrados para el repositorio actual.
agent-sandbox worktree-list

# Abrir un worktree registrado en VS Code.
agent-sandbox worktree-editor -b fix-login

# Eliminar un worktree limpio y su branch.
agent-sandbox worktree-delete --branch fix-login

# Permitir eliminar también sus cambios sin commitear.
agent-sandbox worktree-delete -b fix-login --force

# Mostrar todos, pedir confirmación y eliminarlos junto con sus branches.
agent-sandbox worktree-delete-all

# Omitir la confirmación de la eliminación masiva.
agent-sandbox worktree-delete-all --yes
```

`worktree-delete` se niega a descartar cambios sin `--force`.
`worktree-delete-all` siempre elimina los cambios sin commitear después de la
confirmación (o inmediatamente con `--yes`). No ejecutes los comandos de
eliminación desde el worktree que querés borrar.
`worktree-editor` requiere `-b` y abre en vscode el branch worktree
