import path from "path"
import { $ } from "bun"
import { z } from "zod"
import { NamedError } from "../util/error"
import { Bus } from "../bus"
import { Log } from "../util/log"

declare global {
  const OPENCODE_VERSION: string
}

const GITHUB_OWNER = "zerosync-co"
const GITHUB_REPOSITORY = "autoprovisioner-cli"

export namespace Installation {
  const log = Log.create({ service: "installation" })

  export type Method = Awaited<ReturnType<typeof method>>

  export const Event = {
    Updated: Bus.event(
      "installation.updated",
      z.object({
        version: z.string(),
      }),
    ),
  }

  export const Info = z
    .object({
      version: z.string(),
      latest: z.string(),
    })
    .openapi({
      ref: "InstallationInfo",
    })
  export type Info = z.infer<typeof Info>

  export async function info() {
    return {
      version: VERSION,
      latest: await latest(),
    }
  }

  export function isSnapshot() {
    return VERSION.startsWith("0.0.0")
  }

  export function isDev() {
    return VERSION === "dev"
  }

  export async function method() {
    const exec = process.execPath.toLowerCase()

    const checks = [
      {
        name: "npm" as const,
        command: () => $`npm list -g --depth=0`.throws(false).text(),
      },
      {
        name: "yarn" as const,
        command: () => $`yarn global list`.throws(false).text(),
      },
      {
        name: "pnpm" as const,
        command: () => $`pnpm list -g --depth=0`.throws(false).text(),
      },
      {
        name: "bun" as const,
        command: () => $`bun pm ls -g`.throws(false).text(),
      },
      {
        name: "brew" as const,
        command: () => $`brew list --formula autoprovisioner-ai`.throws(false).text(),
      },
    ]

    checks.sort((a, b) => {
      const aMatches = exec.includes(a.name)
      const bMatches = exec.includes(b.name)
      if (aMatches && !bMatches) return -1
      if (!aMatches && bMatches) return 1
      return 0
    })

    for (const check of checks) {
      const output = await check.command()
      if (output.includes("autoprovisioner-ai")) {
        return check.name
      }
    }

    // Only fall back to curl/script installation if no package manager installation found
    if (process.execPath.includes(path.join(".autoprovisioner", "bin"))) return "curl"

    return "unknown"
  }

  export const UpgradeFailedError = NamedError.create(
    "UpgradeFailedError",
    z.object({
      stderr: z.string(),
    }),
  )

  export async function upgrade(method: Method, target: string) {
    const cmd = (() => {
      switch (method) {
        case "curl":
          if (process.platform === "win32") {
            throw new Error(
              'Windows auto-updates are currently disabled due to persistent issues. Please update manually using: powershell -c "irm https://cli.autoprovisioner.ai/install.ps1 | iex"',
            )
          }
          return $`curl -fsSL https://cli.autoprovisioner.ai/install | bash`.env({
            ...process.env,
            VERSION: target,
          })
        case "npm":
          return $`npm install -g autoprovisioner-ai@${target}`
        case "pnpm":
          return $`pnpm install -g autoprovisioner-ai@${target}`
        case "bun":
          return $`bun install -g autoprovisioner-ai@${target}`
        // case "brew":
        //   return $`brew install sst/tap/autoprovisioner`.env({
        //     HOMEBREW_NO_AUTO_UPDATE: "1",
        //   })
        default:
          throw new Error(`Unknown method: ${method}`)
      }
    })()
    const result = await cmd.quiet().throws(false)
    log.info("upgraded", {
      method,
      target,
      stdout: result.stdout.toString(),
      stderr: result.stderr.toString(),
    })
    if (result.exitCode !== 0)
      throw new UpgradeFailedError({
        stderr: result.stderr.toString("utf8"),
      })
  }

  export const VERSION = typeof OPENCODE_VERSION === "string" ? OPENCODE_VERSION : "dev"

  export async function latest() {
    return fetch(`https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPOSITORY}/releases/latest`)
      .then((res) => res.json())
      .then((data) => {
        if (typeof data.tag_name !== "string") {
          log.error("GitHub API error", data)
          throw new Error("failed to fetch latest version")
        }
        return data.tag_name.slice(1) as string
      })
  }
}
