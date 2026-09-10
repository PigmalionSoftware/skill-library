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
agent-sandbox <branch> --agent <codex|claude|opencode|pi> [--model <modelo>] [--base-image <imagen>] [--push] <prompt...>
```

El branch debe ser el primer argumento. El worktree se crea como directorio
hermano del directorio actual, en `../<branch>`. Si esa ruta ya existe, la
ejecución se rechaza. Un branch existente puede reutilizarse solo si no está
activo en otro worktree.

| Parámetro | Descripción |
| --- | --- |
| `<branch>` | Branch para el worktree aislado. |
| `--agent` | Obligatorio. Uno de `codex`, `claude`, `opencode` o `pi`. |
| `--model` | Opcional. Sobrescribe el modelo que resuelve el agente. |
| `--base-image` | Opcional. Deriva una imagen desde una base compatible, para disponer de su toolchain dentro del sandbox. |
| `--push` | Al finalizar, agrega todos los cambios, crea un commit cuyo mensaje es el prompt y hace `git push --set-upstream origin <branch>`. |
| `<prompt...>` | Instrucción para el agente. Usá comillas para conservarla como una sola cadena. |

Sin `--push`, los cambios quedan sin commitear en el worktree. Con `--push`, si
el agente no produjo cambios, no se crea ningún commit.

### Imagen base externa

`--base-image` permite que el agente elegido use una imagen Alpine, Debian,
Ubuntu, Fedora, RHEL 8/9, UBI 8/9 o Amazon Linux 2023 que ya tiene el runtime
del proyecto. Por ejemplo,
`golang:1.26-alpine` deja disponibles Go, `gofmt` y `go test` dentro del
contenedor:

```bash
agent-sandbox fix-go-tests --agent codex --base-image golang:1.26-alpine "run gofmt and go test ./..., then fix failures"
```

La primera ejecución construye una imagen local derivada e instala Node.js,
npm, Bash, Codex, Claude Code, opencode, pi, ripgrep, certificados CA, curl y
Git.
Las ejecuciones siguientes reutilizan esa imagen para la misma referencia base.
La imagen detecta `apk` para Alpine, `apt-get` para Debian y Ubuntu, `dnf` para
Fedora/RHEL/UBI/Amazon Linux, o `microdnf` para UBI minimal. Las otras familias
fallan durante el build con un mensaje explícito. Alpine instala Node.js y npm
desde sus paquetes; Debian, Ubuntu y las bases RPM usan el repositorio firmado
de NodeSource para instalar Node.js 22, que incluye npm y cumple el mínimo de
Claude Code.

Una imagen RHEL debe incluir repositorios habilitados y una suscripción válida
si la requiere: el sandbox no monta credenciales ni entitlement certificates del
host. Esta primera versión tampoco habilita EPEL, CRB/CodeReady Builder ni
soporta `yum`; si los repositorios activos no contienen un paquete requerido
como `ripgrep`, el build falla mostrando el error del gestor de paquetes.

La imagen derivada no se actualiza automáticamente. Para forzar un nuevo build,
listá las imágenes `agent-sandbox-base-*` y eliminá explícitamente la
que corresponde a la base elegida:

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
agent-sandbox fix-login --agent codex --model gpt-5.6-sol "fix the login redirect loop"
```

Ejecutar Claude Code y publicar el branch al terminar:

```bash
agent-sandbox add-test --agent claude --model sonnet --push "add a regression test for the login redirect"
```

Dejar que opencode resuelva su modelo configurado:

```bash
agent-sandbox update-copy --agent opencode "update the empty-state copy"
```

Ejecutar Codex con el toolchain Go de una base Alpine:

```bash
agent-sandbox fix-go-tests --agent codex --base-image golang:1.26-alpine "run go test ./... and fix failures"
```

## Gestionar worktrees creados por el sandbox

El comando registra los worktrees que creó en
`~/.config/agent-sandbox/worktrees.jsonl`. Los siguientes subcomandos solo ven
y eliminan esas entradas; no afectan worktrees creados manualmente.

```bash
# Mostrar los worktrees registrados para el repositorio actual.
agent-sandbox worktree-list

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
