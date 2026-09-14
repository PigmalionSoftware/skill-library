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
agent-sandbox [-b <branch>] -a <codex|claude|opencode|pi> [-m <modelo>] [-i <imagen>] [-p] [-c <mensaje-commit>] (<prompt...> | -f <archivo-prompt>)
```
| Parámetro | Descripción |
| --- | --- |
| `-b`, `--branch` | Opcional. Branch para el worktree aislado; si se omite, se genera uno. |
| `-a`, `--agent` | Obligatorio. Uno de `codex`, `claude`, `opencode` o `pi`. |
| `-m`, `--model` | Opcional. Sobrescribe el modelo que resuelve el agente. |
| `-i`, `--base-image` | Opcional. Deriva una imagen desde una base compatible, para disponer de su toolchain dentro del sandbox. |
| `-p`, `--push` | Al finalizar, agrega todos los cambios, crea un commit y hace `git push --set-upstream origin <branch>`. |
| `-c`, `--commit-message` | Opcional. Mensaje del commit creado por `-p` o `--push`; si se omite, usa el prompt resuelto. Sin `-p` o `--push`, no tiene efecto. |
| `-f`, `--file-prompt` | Archivo cuyo contenido se usa como instrucción para el agente, en lugar de `<prompt...>`. |
| `<prompt...>` | Instrucción para el agente, en lugar de `-f` o `--file-prompt`. Usá comillas para conservarla como una sola cadena. |

Sin `-p` o `--push`, los cambios quedan sin commitear en el worktree. Con `-p`
o `--push`, si el agente no produjo cambios, no se crea ningún commit. Hay que
proporcionar exactamente una fuente de prompt: `<prompt...>` o `-f`/`--file-prompt`.
No se pueden usar juntas. Si se usa `--push` sin `--commit-message`, el contenido
del archivo se convierte en el mensaje de commit por defecto cuando se eligió
`--file-prompt`.

### Imagen base externa

`--base-image` usa una imagen que ya trae el runtime del proyecto. Se admiten
Alpine, Debian, Ubuntu, Fedora, RHEL 8/9, UBI 8/9 y Amazon Linux 2023. Por
ejemplo, `golang:1.26-alpine` deja disponibles Go, `gofmt` y `go test` dentro
del contenedor:

```bash
agent-sandbox -b fix-go-tests -a codex -i golang:1.26-alpine "run gofmt and go test ./..., then fix failures"
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
agent-sandbox -a codex -m gpt-5.6-sol "fix the login redirect loop"
```

Ejecutar Claude Code y publicar el branch al terminar:

```bash
agent-sandbox -b add-test -a claude -m sonnet -p "add a regression test for the login redirect"
```

Dejar que opencode resuelva su modelo configurado:

```bash
agent-sandbox -b update-copy -a opencode "update the empty-state copy"
```

Ejecutar Codex con imagen de golang

```bash
agent-sandbox -b fix-go-tests -a codex -i golang:1.26-alpine "run go test ./... and fix failures"
```

Obtener prompt de archivo
```bash
agent-sandbox -b prompt-file-test -a codex -f prompt.md
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
