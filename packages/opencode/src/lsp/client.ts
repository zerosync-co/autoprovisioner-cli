import path from "path"
import { createMessageConnection, StreamMessageReader, StreamMessageWriter } from "vscode-jsonrpc/node"
import type { Diagnostic as VSCodeDiagnostic } from "vscode-languageserver-types"
import { App } from "../app/app"
import { Log } from "../util/log"
import { LANGUAGE_EXTENSIONS } from "./language"
import { Bus } from "../bus"
import z from "zod"
import type { LSPServer } from "./server"
import { NamedError } from "../util/error"
import { withTimeout } from "../util/timeout"

export namespace LSPClient {
  const log = Log.create({ service: "lsp.client" })

  export type Info = NonNullable<Awaited<ReturnType<typeof create>>>

  export type Diagnostic = VSCodeDiagnostic

  export const InitializeError = NamedError.create(
    "LSPInitializeError",
    z.object({
      serverID: z.string(),
    }),
  )

  export const Event = {
    Diagnostics: Bus.event(
      "lsp.client.diagnostics",
      z.object({
        serverID: z.string(),
        path: z.string(),
      }),
    ),
  }

  export async function create(input: { serverID: string; server: LSPServer.Handle; root: string }) {
    const app = App.info()
    const l = log.clone().tag("serverID", input.serverID)
    l.info("starting client")

    const connection = createMessageConnection(
      new StreamMessageReader(input.server.process.stdout),
      new StreamMessageWriter(input.server.process.stdin),
    )

    const diagnostics = new Map<string, Diagnostic[]>()
    connection.onNotification("textDocument/publishDiagnostics", (params) => {
      const path = new URL(params.uri).pathname
      l.info("textDocument/publishDiagnostics", {
        path,
      })
      const exists = diagnostics.has(path)
      diagnostics.set(path, params.diagnostics)
      if (!exists && input.serverID === "typescript") return
      Bus.publish(Event.Diagnostics, { path, serverID: input.serverID })
    })
    connection.onRequest("window/workDoneProgress/create", (params) => {
      l.info("window/workDoneProgress/create", params)
      return null
    })
    connection.onRequest("workspace/configuration", async () => {
      return [{}]
    })
    connection.listen()

    l.info("sending initialize")
    await withTimeout(
      connection.sendRequest("initialize", {
        rootUri: "file://" + input.root,
        processId: input.server.process.pid,
        workspaceFolders: [
          {
            name: "workspace",
            uri: "file://" + input.root,
          },
        ],
        initializationOptions: {
          ...input.server.initialization,
        },
        capabilities: {
          window: {
            workDoneProgress: true,
          },
          workspace: {
            configuration: true,
          },
          textDocument: {
            synchronization: {
              didOpen: true,
              didChange: true,
            },
            publishDiagnostics: {
              versionSupport: true,
            },
          },
        },
      }),
      5_000,
    ).catch((err) => {
      l.error("initialize error", { error: err })
      throw new InitializeError(
        { serverID: input.serverID },
        {
          cause: err,
        },
      )
    })

    await connection.sendNotification("initialized", {})

    const files: {
      [path: string]: number
    } = {}

    const result = {
      root: input.root,
      get serverID() {
        return input.serverID
      },
      get connection() {
        return connection
      },
      notify: {
        async open(input: { path: string }) {
          input.path = path.isAbsolute(input.path) ? input.path : path.resolve(app.path.cwd, input.path)
          const file = Bun.file(input.path)
          const text = await file.text()
          const version = files[input.path]
          if (version !== undefined) {
            // File is already open, send didChange instead of close/open
            const newVersion = version + 1
            log.info("textDocument/didChange", input)
            await connection.sendNotification("textDocument/didChange", {
              textDocument: {
                uri: `file://` + input.path,
                version: newVersion,
              },
              contentChanges: [
                {
                  text,
                },
              ],
            })
            files[input.path] = newVersion
            return
          }
          log.info("textDocument/didOpen", input)
          diagnostics.delete(input.path)
          const extension = path.extname(input.path)
          const languageId = LANGUAGE_EXTENSIONS[extension] ?? "plaintext"
          await connection.sendNotification("textDocument/didOpen", {
            textDocument: {
              uri: `file://` + input.path,
              languageId,
              version: 0,
              text,
            },
          })
          files[input.path] = 0
          return
        },
      },
      get diagnostics() {
        return diagnostics
      },
      async waitForDiagnostics(input: { path: string }) {
        input.path = path.isAbsolute(input.path) ? input.path : path.resolve(app.path.cwd, input.path)
        log.info("waiting for diagnostics", input)
        let unsub: () => void
        return await withTimeout(
          new Promise<void>((resolve) => {
            unsub = Bus.subscribe(Event.Diagnostics, (event) => {
              if (event.properties.path === input.path && event.properties.serverID === result.serverID) {
                log.info("got diagnostics", input)
                unsub?.()
                resolve()
              }
            })
          }),
          3000,
        )
          .catch(() => {})
          .finally(() => {
            unsub?.()
          })
      },
      async shutdown() {
        l.info("shutting down")
        connection.end()
        connection.dispose()
        input.server.process.kill()
        l.info("shutdown")
      },
    }

    l.info("initialized")

    return result
  }
}
