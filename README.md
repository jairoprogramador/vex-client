<div align="center">
  <h1>Vex Client</h1>
  <p><strong>Tu asistente personal para desplegar aplicaciones sin complicaciones.</strong></p>
  <p><i>Orquesta despliegues complejos con comandos sencillos.</i></p>
  
  <p>
    <a href="https://github.com/jairoprogramador/vex-client/releases">
      <img src="https://img.shields.io/github/v/release/jairoprogramador/vex-client?style=for-the-badge" alt="Latest Release">
    </a>
    <a href="https://github.com/jairoprogramador/vex-client/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/jairoprogramador/vex-client?style=for-the-badge" alt="License">
    </a>
  </p>
</div>

---

**`vex`** es una herramienta de línea de comandos (CLI) que actúa como un cliente inteligente para `vexc`. Su misión es simplificar al máximo el proceso de despliegue, permitiéndote inicializar y ejecutar el despliegue de tus proyectos en un entorno contenerizado con una configuración mínima y comandos intuitivos.

Olvídate de la complejidad de Docker y los detalles de bajo nivel, `vex` es el puente que te conecta con un motor de despliegue potente, haciendo que el proceso sea simple y repetible.

## ✨ Características Principales

*   **🚀 Inicialización Rápida**: Con `vex init`, la herramienta genera un archivo `vexconfig.yaml` adaptado a tus necesidades.
*   **📄 Configuración Declarativa**: Define tu configuracion de despliegue en un único archivo `vexconfig.yaml`. Fácil de leer, modificar y versionar.
*   **🐳 Abstracción de Docker**: `vex [step] [environment]` se encarga de construir la imagen de Docker y ejecutar comando en el contenedor que aloja a `vex`. No necesitas ser un experto.
*   **🔌 Orquestación Transparente**: Actúa como un punto de entrada único para `vex`, pasándole tus instrucciones y gestionando el ciclo de vida del contenedor por ti.

## 🚀 Instalación

Instala `vex` en segundos.

*(Nota: Las siguientes instrucciones son un ejemplo. Ajústalas según tu método de distribución final).*

### macOS (Homebrew)
```sh
brew install --cask jairoprogramador/vex-client/vex
```
Si macOS indica que no puede verificar el desarrollador, puedes permitir la ejecución en **Ajustes del sistema → Privacidad y seguridad → "Abrir de todos modos"**, o en Terminal: `xattr -cr $(which vex)`.

### Linux
Puedes descargar el paquete `.deb` o `.rpm` desde la [página de Releases](https://github.com/jairoprogramador/vex-client/releases) y usar tu gestor de paquetes.

```sh
# Para sistemas basados en Debian/Ubuntu
sudo dpkg -i vex-client_*.deb

# Para sistemas basados en Red Hat/Fedora
sudo rpm -i vex-client_*.rpm
```
Alternativamente, puedes descargar el binario directamente:
```sh
curl -sL https://github.com/jairoprogramador/vex-client/releases/latest/download/vex-client_linux_amd64.tar.gz | tar xz

sudo mv vex /usr/local/bin/
```

### Windows
1.  Descarga el archivo `vex-client_*_windows_a*64.zip` desde la [página de Releases](https://github.com/jairoprogramador/vex-client/releases).
2.  Descomprime el archivo.
3.  Añade el ejecutable `vex.exe` a tu variable de entorno `PATH`.

## 🏁 Guía de Inicio Rápido

Este es el flujo de trabajo típico con `vex`.

### Paso 1: Inicializa tu Proyecto

Navega al directorio raíz de tu proyecto y ejecuta:
```sh
vex init
```
La herramienta te guiará con unas sencillas preguntas para generar el archivo `vexconfig.yaml`, que conecta tu proyecto con la plantilla de despliegue de `vex`.

### Paso 2: Ejecuta los Pasos de Despliegue

Una vez configurado, usa el comando `vex` para enviar instrucciones directamente a `vex`. Los `steps` como `test`, `supply`, `package`  o `deploy` son gestionados por el motor de `vex`, no por esta CLI.

Por ejemplo, para ejecutar las pruebas en el entorno de `sand`:
```sh
vex test sand
```
Para desplegar en el mismo entorno:
```sh
vex deploy sand
```
`vex` se encargará de iniciar el contenedor con el core de `vex` y le pasará estos comandos para que los ejecute.

## 📚 Comandos Básicos

| Comando | Descripción |
| :--- | :--- |
| `vex init` | Inicializa un proyecto creando el archivo de configuración `vexconfig.yaml`. |
| `vex [step] [env]` | Ejecuta un comando en `vex`. Los `steps` (`test`, `supply`, `deploy`, etc.) dependen de la plantilla utilizada. |
| `vex version` | Muestra la versión de la CLI. |

**Flags comunes:**
*   `--yes` o `-y`: Salta las confirmaciones interactivas para `vex init`.

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Si tienes ideas, sugerencias o encuentras un error, por favor abre un [issue](https://github.com/jairoprogramador/vex-client/issues) o envía un [pull request](https://github.com/jairoprogramador/vex-client/pulls).

## 📄 Licencia

`vex` está distribuido bajo la [Apache License 2.0](https://github.com/jairoprogramador/vex-client/blob/main/LICENSE).
