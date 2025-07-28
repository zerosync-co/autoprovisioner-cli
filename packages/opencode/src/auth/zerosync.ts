import { Auth } from "./index"
import open from "open"
import * as prompts from "@clack/prompts"
import { UI } from "../cli/ui"

export namespace AuthZerosync {
  export async function access() {
    const info = await Auth.get("zerosync")
    if (info?.type === "oauth" && info.access && info.expires > Date.now()) return info.access

    const token = await login()
    await Auth.set("zerosync", {
      type: "oauth",

      // @ts-ignore FIXME
      access: token,
      expires: Date.now() + 1000 * 60 * 60 * 24 * 7,
      // refresh
    })
    process.exit()

    return token
  }

  export class ExchangeFailed extends Error {
    constructor() {
      super("Exchange failed")
    }
  }
}

async function login() {
  const urlToOpen = "https://autoprovisioner.zerosync.co/auth/cli/sign-in"

  let server: Bun.Server
  let loginTimeoutHandle: ReturnType<typeof setTimeout>

  const timerPromise = new Promise<void>((_, reject) => {
    loginTimeoutHandle = setTimeout(async () => {
      await server.stop(true)
      clearTimeout(loginTimeoutHandle)
      reject("Timed out waiting for authorization code, please try again.")
    }, 120000) // wait for 120 seconds for the user to authorize
  })

  const loginPromise = new Promise<string>((resolve, reject) => {
    server = Bun.serve({
      port: 8976,
      hostname: "localhost",
      fetch: (req) => {
        async function finish(token: string | null, error?: Error) {
          clearTimeout(loginTimeoutHandle)

          setTimeout(async () => {
            await server.stop(true)
          }, 100)

          if (error) {
            reject(error)
          } else {
            resolve(token as string)
          }
        }

        if (req.method !== "GET") {
          return new Response("OK")
        }

        const { searchParams } = new URL(req.url)
        const token = searchParams.get("token")

        if (!token?.length) {
          finish(null, new Error("failed to resolve token"))
          return new Response("failed to resolve token", { status: 400 })
        }

        const response = new Response()

        finish(token as string)

        return response
      },
    })
  })

  prompts.note("This will open a browser window to sign in to ZeroSync")
  const shouldOpen = await prompts.confirm({
    message: "Open signin link in browser?",
  })

  if (prompts.isCancel(shouldOpen)) throw new UI.CancelledError()

  if (!shouldOpen) {
    prompts.log.info("You can manually visit: https://autoprovisioner.zerosync.co/auth/cli/sign-in")
  } else {
    open(urlToOpen)
  }

  const token = await Promise.race([timerPromise, loginPromise])
  return token
}
