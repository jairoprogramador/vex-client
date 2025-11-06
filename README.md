<div align="center">
  <h1>FastDeploy CLI</h1>
  <p><strong>Tu asistente personal para desplegar aplicaciones sin complicaciones.</strong></p>
  <p><i>Orquesta despliegues complejos con comandos sencillos.</i></p>
  
  <p>
    <a href="https://github.com/jairoprogramador/fastdeploy/releases">
      <img src="https://img.shields.io/github/v/release/jairoprogramador/fastdeploy?style=for-the-badge" alt="Latest Release">
    </a>
    <a href="https://github.com/jairoprogramador/fastdeploy/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/jairoprogramador/fastdeploy?style=for-the-badge" alt="License">
    </a>
  </p>
</div>

---

**`fastdeploy`** es una herramienta de línea de comandos (CLI) que actúa como un cliente inteligente para `fastdeploy-core`. Su misión es simplificar al máximo el proceso de despliegue, permitiéndote inicializar y ejecutar tus proyectos en un entorno contenerizado con una configuración mínima y comandos intuitivos.

Olvídate de la complejidad de Docker y los detalles de bajo nivel. `fastdeploy` es el puente que te conecta con un motor de despliegue potente, haciendo que el proceso sea simple y repetible.

## ✨ Características Principales

*   **🚀 Inicialización Rápida**: Con `fastdeploy init`, la herramienta genera un archivo `fdconfig.yaml` adaptado a tus necesidades.
*   **📄 Configuración Declarativa**: Define todo tu despliegue en un único archivo `fdconfig.yaml`. Fácil de leer, modificar y versionar.
*   **🐳 Abstracción de Docker**: `fastdeploy` se encarga de construir la imagen de Docker y ejecutar el contenedor que aloja a `fastdeploy-core`. No necesitas ser un experto.
*   **🔌 Orquestación Transparente**: Actúa como un punto de entrada único para `fastdeploy-core`, pasándole tus instrucciones y gestionando el ciclo de vida del contenedor por ti.

## 🚀 Instalación

Instala `fastdeploy` en segundos.

*(Nota: Las siguientes instrucciones son un ejemplo. Ajústalas según tu método de distribución final).*

### macOS (Homebrew)
```sh
# brew tap jairoprogramador/fastdeploy
# brew install fastdeploy
```

### Linux
Puedes descargar el paquete `.deb` o `.rpm` desde la [página de Releases](https://github.com/jairoprogramador/fastdeploy/releases) y usar tu gestor de paquetes.

```sh
# Para sistemas basados en Debian/Ubuntu
# sudo dpkg -i fastdeploy_*.deb

# Para sistemas basados en Red Hat/Fedora
# sudo rpm -i fastdeploy_*.rpm
```
Alternativamente, puedes descargar el binario directamente:
```sh
# curl -sL https://github.com/jairoprogramador/fastdeploy/releases/download/v0.0.5/fastdeploy_0.0.5_linux_amd64.tar.gz | tar xz
# sudo mv fastdeploy /usr/local/bin/
```

### Windows
1.  Descarga el archivo `fastdeploy_*_windows_a*64.zip` desde la [página de Releases](https://github.com/jairoprogramador/fastdeploy/releases).
2.  Descomprime el archivo.
3.  Añade el ejecutable `fastdeploy.exe` a tu `PATH`.

## 🏁 Guía de Inicio Rápido

Este es el flujo de trabajo típico con `fastdeploy`.

### Paso 1: Inicializa tu Proyecto

Navega al directorio raíz de tu proyecto y ejecuta:
```sh
fd init
```
La herramienta te guiará con unas sencillas preguntas para generar el archivo `fdconfig.yaml`, que conecta tu proyecto con la plantilla de despliegue de `fastdeploy-core`.

### Paso 2: Ejecuta los Pasos de Despliegue

Una vez configurado, usa el comando `fd` para enviar instrucciones directamente a `fastdeploy-core`. Los `steps` como `test`, `supply` o `deploy` son gestionados por el motor `core`, no por esta CLI.

Por ejemplo, para ejecutar las pruebas en el entorno de `sand`:
```sh
fd test sand
```
Para desplegar en el mismo entorno:
```sh
fd deploy sand
```
`fastdeploy` se encargará de iniciar el contenedor con `fastdeploy-core` y le pasará estos comandos para que los ejecute.

## 📚 Comandos Básicos

| Comando | Descripción |
| :--- | :--- |
| `fd init` | Inicializa un proyecto creando el archivo de configuración `fdconfig.yaml`. |
| `fd [step] [env]` | Ejecuta un comando en `fastdeploy-core`. Los `steps` (`test`, `supply`, `deploy`, etc.) dependen de la plantilla utilizada por el core. |
| `fd version` | Muestra la versión de la CLI. |

**Flags comunes:**
*   `--yes` o `-y`: Salta las confirmaciones interactivas para `fd init`.

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Si tienes ideas, sugerencias o encuentras un error, por favor abre un [issue](https://github.com/jairoprogramador/fastdeploy/issues) o envía un [pull request](https://github.com/jairoprogramador/fastdeploy/pulls).

## 📄 Licencia

`fastdeploy` está distribuido bajo la [Licencia MIT](https://github.com/jairoprogramador/fastdeploy/blob/main/LICENSE).
