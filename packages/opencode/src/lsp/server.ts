import { spawn, type ChildProcessWithoutNullStreams } from "child_process"
import type { App } from "../app/app"
import path from "path"
import { Global } from "../global"
import { Log } from "../util/log"
import { BunProc } from "../bun"
import { $ } from "bun"
import fs from "fs/promises"
import { Filesystem } from "../util/filesystem"

export namespace LSPServer {
  const log = Log.create({ service: "lsp.server" })

  export interface Handle {
    process: ChildProcessWithoutNullStreams
    initialization?: Record<string, any>
  }

  type RootFunction = (file: string, app: App.Info) => Promise<string | undefined>

  const NearestRoot = (patterns: string[]): RootFunction => {
    return async (file, app) => {
      const files = Filesystem.up({
        targets: patterns,
        start: path.dirname(file),
        stop: app.path.root,
      })
      const first = await files.next()
      await files.return()
      if (!first.value) return app.path.root
      return path.dirname(first.value)
    }
  }

  export interface Info {
    id: string
    extensions: string[]
    global?: boolean
    root: RootFunction
    spawn(app: App.Info, root: string): Promise<Handle | undefined>
  }

  export const Typescript: Info = {
    id: "typescript",
    root: NearestRoot(["tsconfig.json", "package.json", "jsconfig.json"]),
    extensions: [".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".mts", ".cts"],
    async spawn(app, root) {
      let tsserver = await Bun.resolve("typescript/lib/tsserver.js", app.path.cwd).catch(() => {})
      if (!tsserver) tsserver = await Bun.resolve("typescript/lib/tsserver.js", root).catch(() => {})

      if (!tsserver) return
      const proc = spawn(BunProc.which(), ["x", "typescript-language-server", "--stdio"], {
        cwd: root,
        env: {
          ...process.env,
          BUN_BE_BUN: "1",
        },
      })
      return {
        process: proc,
        initialization: {
          tsserver: {
            path: tsserver,
          },
        },
      }
    },
  }

  export const Gopls: Info = {
    id: "golang",
    root: async (file, app) => {
      const work = await NearestRoot(["go.work"])(file, app)
      if (work) return work
      return NearestRoot(["go.mod", "go.sum"])(file, app)
    },
    extensions: [".go"],
    async spawn(_, root) {
      let bin = Bun.which("gopls", {
        PATH: process.env["PATH"] + ":" + Global.Path.bin,
      })
      if (!bin) {
        if (!Bun.which("go")) return
        log.info("installing gopls")
        const proc = Bun.spawn({
          cmd: ["go", "install", "golang.org/x/tools/gopls@latest"],
          env: { ...process.env, GOBIN: Global.Path.bin },
          stdout: "pipe",
          stderr: "pipe",
          stdin: "pipe",
        })
        const exit = await proc.exited
        if (exit !== 0) {
          log.error("Failed to install gopls")
          return
        }
        bin = path.join(Global.Path.bin, "gopls" + (process.platform === "win32" ? ".exe" : ""))
        log.info(`installed gopls`, {
          bin,
        })
      }
      return {
        process: spawn(bin!, {
          cwd: root,
        }),
      }
    },
  }

  export const RubyLsp: Info = {
    id: "ruby-lsp",
    root: NearestRoot(["Gemfile"]),
    extensions: [".rb", ".rake", ".gemspec", ".ru"],
    async spawn(_, root) {
      let bin = Bun.which("ruby-lsp", {
        PATH: process.env["PATH"] + ":" + Global.Path.bin,
      })
      if (!bin) {
        const ruby = Bun.which("ruby")
        const gem = Bun.which("gem")
        if (!ruby || !gem) {
          log.info("Ruby not found, please install Ruby first")
          return
        }
        log.info("installing ruby-lsp")
        const proc = Bun.spawn({
          cmd: ["gem", "install", "ruby-lsp", "--bindir", Global.Path.bin],
          stdout: "pipe",
          stderr: "pipe",
          stdin: "pipe",
        })
        const exit = await proc.exited
        if (exit !== 0) {
          log.error("Failed to install ruby-lsp")
          return
        }
        bin = path.join(Global.Path.bin, "ruby-lsp" + (process.platform === "win32" ? ".exe" : ""))
        log.info(`installed ruby-lsp`, {
          bin,
        })
      }
      return {
        process: spawn(bin!, ["--stdio"], {
          cwd: root,
        }),
      }
    },
  }

  export const Pyright: Info = {
    id: "pyright",
    extensions: [".py", ".pyi"],
    root: NearestRoot(["pyproject.toml", "setup.py", "setup.cfg", "requirements.txt", "Pipfile", "pyrightconfig.json"]),
    async spawn(_, root) {
      const proc = spawn(BunProc.which(), ["x", "pyright-langserver", "--stdio"], {
        cwd: root,
        env: {
          ...process.env,
          BUN_BE_BUN: "1",
        },
      })
      return {
        process: proc,
      }
    },
  }

  export const ElixirLS: Info = {
    id: "elixir-ls",
    extensions: [".ex", ".exs"],
    root: NearestRoot(["mix.exs", "mix.lock"]),
    async spawn(_, root) {
      let binary = Bun.which("elixir-ls")
      if (!binary) {
        const elixirLsPath = path.join(Global.Path.bin, "elixir-ls")
        binary = path.join(
          Global.Path.bin,
          "elixir-ls-master",
          "release",
          process.platform === "win32" ? "language_server.bar" : "language_server.sh",
        )

        if (!(await Bun.file(binary).exists())) {
          const elixir = Bun.which("elixir")
          if (!elixir) {
            log.error("elixir is required to run elixir-ls")
            return
          }

          log.info("downloading elixir-ls from GitHub releases")

          const response = await fetch("https://github.com/elixir-lsp/elixir-ls/archive/refs/heads/master.zip")
          if (!response.ok) return
          const zipPath = path.join(Global.Path.bin, "elixir-ls.zip")
          await Bun.file(zipPath).write(response)

          if (process.platform === "win32") {
            await $`powershell -Command "Expand-Archive -Path '${zipPath}' -DestinationPath '${Global.Path.bin}' -Force"`
              .cwd(Global.Path.bin)
              .nothrow()
          } else {
            await $`unzip -o -q ${zipPath}`.cwd(Global.Path.bin).nothrow()
          }

          await fs.rm(zipPath, {
            force: true,
            recursive: true,
          })

          await $`mix deps.get && mix compile && mix elixir_ls.release2 -o release`
            .quiet()
            .cwd(path.join(Global.Path.bin, "elixir-ls-master"))
            .env({ MIX_ENV: "prod", ...process.env })

          log.info(`installed elixir-ls`, {
            path: elixirLsPath,
          })
        }
      }

      return {
        process: spawn(binary, {
          cwd: root,
        }),
      }
    },
  }

  export const Zls: Info = {
    id: "zls",
    extensions: [".zig", ".zon"],
    root: NearestRoot(["build.zig"]),
    async spawn(_, root) {
      let bin = Bun.which("zls", {
        PATH: process.env["PATH"] + ":" + Global.Path.bin,
      })

      if (!bin) {
        const zig = Bun.which("zig")
        if (!zig) {
          log.error("Zig is required to use zls. Please install Zig first.")
          return
        }

        log.info("downloading zls from GitHub releases")

        const releaseResponse = await fetch("https://api.github.com/repos/zigtools/zls/releases/latest")
        if (!releaseResponse.ok) {
          log.error("Failed to fetch zls release info")
          return
        }

        const release = await releaseResponse.json()

        const platform = process.platform
        const arch = process.arch
        let assetName = ""

        let zlsArch: string = arch
        if (arch === "arm64") zlsArch = "aarch64"
        else if (arch === "x64") zlsArch = "x86_64"
        else if (arch === "ia32") zlsArch = "x86"

        let zlsPlatform: string = platform
        if (platform === "darwin") zlsPlatform = "macos"
        else if (platform === "win32") zlsPlatform = "windows"

        const ext = platform === "win32" ? "zip" : "tar.xz"

        assetName = `zls-${zlsArch}-${zlsPlatform}.${ext}`

        const supportedCombos = [
          "zls-x86_64-linux.tar.xz",
          "zls-x86_64-macos.tar.xz",
          "zls-x86_64-windows.zip",
          "zls-aarch64-linux.tar.xz",
          "zls-aarch64-macos.tar.xz",
          "zls-aarch64-windows.zip",
          "zls-x86-linux.tar.xz",
          "zls-x86-windows.zip",
        ]

        if (!supportedCombos.includes(assetName)) {
          log.error(`Platform ${platform} and architecture ${arch} is not supported by zls`)
          return
        }

        const asset = release.assets.find((a: any) => a.name === assetName)
        if (!asset) {
          log.error(`Could not find asset ${assetName} in latest zls release`)
          return
        }

        const downloadUrl = asset.browser_download_url
        const downloadResponse = await fetch(downloadUrl)
        if (!downloadResponse.ok) {
          log.error("Failed to download zls")
          return
        }

        const tempPath = path.join(Global.Path.bin, assetName)
        await Bun.file(tempPath).write(downloadResponse)

        if (ext === "zip") {
          if (process.platform === "win32") {
            await $`powershell -Command "Expand-Archive -Path '${tempPath}' -DestinationPath '${Global.Path.bin}' -Force"`
              .cwd(Global.Path.bin)
              .nothrow()
          } else {
            await $`unzip -o -q ${tempPath}`.cwd(Global.Path.bin).nothrow()
          }
        } else {
          await $`tar -xf ${tempPath}`.cwd(Global.Path.bin).nothrow()
        }

        await fs.rm(tempPath, { force: true })

        bin = path.join(Global.Path.bin, "zls" + (platform === "win32" ? ".exe" : ""))

        if (!(await Bun.file(bin).exists())) {
          log.error("Failed to extract zls binary")
          return
        }

        if (platform !== "win32") {
          await $`chmod +x ${bin}`.nothrow()
        }

        log.info(`installed zls`, { bin })
      }

      return {
        process: spawn(bin, {
          cwd: root,
        }),
      }
    },
  }

  export const TofuLS: Info = {
    id: "tofu-ls",
    root: NearestRoot([
      "main.tf",
      "variables.tf",
      "outputs.tf",
      "provider.tf",
      "terraform.tfvars",
      ".terraform.lock.hcl",
    ]),
    extensions: [".tf", ".tfvars", ".tofu"],
    async spawn(_, root) {
      let bin = Bun.which("tofu-ls", {
        PATH: process.env["PATH"] + ":" + Global.Path.bin,
      })

      if (!bin) {
        log.info("downloading tofu-ls from opentofu releases")

        const releaseResponse = await fetch("https://api.github.com/repos/opentofu/tofu-ls/releases/latest")
        if (!releaseResponse.ok) {
          log.error("Failed to fetch tofu-ls release info")
          return
        }

        const release = await releaseResponse.json()

        const platform = process.platform
        const arch = process.arch
        let assetName = ""

        let tofuLsArch: string = arch
        if (arch === "arm64") tofuLsArch = "arm64"
        else if (arch === "x64") tofuLsArch = "x86_64"
        else if (arch === "ia32") tofuLsArch = "i386"

        let tofuLsPlatform: string = platform
        if (platform === "darwin") tofuLsPlatform = "Darwin"
        else if (platform === "win32") tofuLsPlatform = "Windows"
        else if (platform === "linux") tofuLsPlatform = "Linux"

        const ext = "tar.gz"

        assetName = `tofu-ls_${tofuLsPlatform}_${tofuLsArch}.${ext}`

        const supportedCombos = [
          "tofu-ls_Darwin_arm64.tar.gz",
          "tofu-ls_Darwin_x86_64.tar.gz",
          "tofu-ls_Linux_arm64.tar.gz",
          "tofu-ls_Linux_i386.tar.gz",
          "tofu-ls_Linux_x86_64.tar.gz",
          "tofu-ls_Windows_arm64.tar.gz",
          "tofu-ls_Windows_i386.tar.gz",
          "tofu-ls_Windows_x86_64.tar.gz",
        ]

        if (!supportedCombos.includes(assetName)) {
          log.error(`Platform ${platform} and architecture ${arch} is not supported by tofu-ls`)
          return
        }

        const asset = release.assets.find((a: any) => a.name === assetName)
        if (!asset) {
          log.error(`Could not find asset ${assetName} in latest tofu-ls release`)
          return
        }

        const downloadUrl = asset.browser_download_url
        const downloadResponse = await fetch(downloadUrl)
        if (!downloadResponse.ok) {
          log.error("Failed to download tofu-ls")
          return
        }

        const tempPath = path.join(Global.Path.bin, assetName)
        await Bun.file(tempPath).write(downloadResponse)

        await $`tar -xzvf ${tempPath}`.cwd(Global.Path.bin).nothrow()

        await fs.rm(tempPath, { force: true })

        bin = path.join(Global.Path.bin, "tofu-ls" + (platform === "win32" ? ".exe" : ""))

        if (!(await Bun.file(bin).exists())) {
          log.error("Failed to extract tofu-ls binary")
          return
        }

        if (platform !== "win32") {
          await $`chmod +x ${bin}`.nothrow()
        }

        log.info(`installed tofu-ls`, { bin })
      }

      return {
        process: spawn(bin, ["serve"], {
          cwd: root,
        }),
      }
    },
  }
}
