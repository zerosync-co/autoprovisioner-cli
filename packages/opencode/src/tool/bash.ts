import { z } from "zod"
import { Tool } from "./tool"
import DESCRIPTION from "./bash.txt"
import { App } from "../app/app"

const MAX_OUTPUT_LENGTH = 30000
const DEFAULT_TIMEOUT = 1 * 60 * 1000
const MAX_TIMEOUT = 10 * 60 * 1000

export const BashTool = Tool.define({
  id: "bash",
  description: DESCRIPTION,
  parameters: z.object({
    command: z.string().describe("The command to execute"),
    timeout: z.number().min(0).max(MAX_TIMEOUT).describe("Optional timeout in milliseconds").optional(),
    description: z
      .string()
      .describe(
        "Clear, concise description of what this command does in 5-10 words. Examples:\nInput: ls\nOutput: Lists files in current directory\n\nInput: git status\nOutput: Shows working tree status\n\nInput: npm install\nOutput: Installs package dependencies\n\nInput: mkdir foo\nOutput: Creates directory 'foo'",
      ),
  }),
  async execute(params, ctx) {
    const timeout = Math.min(params.timeout ?? DEFAULT_TIMEOUT, MAX_TIMEOUT)

    const process = Bun.spawn({
      cmd: ["bash", "-c", params.command],
      cwd: App.info().path.cwd,
      maxBuffer: MAX_OUTPUT_LENGTH,
      signal: ctx.abort,
      timeout: timeout,
      stdout: "pipe",
      stderr: "pipe",
    })

    let stdoutBuffer = ""
    let stderrBuffer = ""

    if (ctx.stream) {
      const streamOutput = async () => {
        const stdoutReader = process.stdout.getReader()
        const stderrReader = process.stderr.getReader()
        const decoder = new TextDecoder()

        const readStream = async (reader: ReadableStreamDefaultReader<Uint8Array>, type: "stdout" | "stderr") => {
          try {
            while (true) {
              const { done, value } = await reader.read()
              if (done) break

              const text = decoder.decode(value, { stream: true })

              if (type === "stdout") {
                stdoutBuffer += text
              } else {
                stderrBuffer += text
              }

              if (text.trim() && typeof ctx.stream === "function") {
                ctx.stream({
                  type,
                  data: text,
                  timestamp: Date.now(),
                })
              }
            }
          } catch (error) {
            if (!ctx.abort.aborted) {
              throw error
            }
          } finally {
            reader.releaseLock()
          }
        }

        await Promise.all([readStream(stdoutReader, "stdout"), readStream(stderrReader, "stderr")])
      }

      streamOutput().catch(() => {
        // Ignore streaming errors, we'll still get final output
      })
    }

    await process.exited

    if (!ctx.stream) {
      stdoutBuffer = await new Response(process.stdout).text()
      stderrBuffer = await new Response(process.stderr).text()
    }

    return {
      title: params.command,
      metadata: {
        stderr: stderrBuffer,
        stdout: stdoutBuffer,
        exit: process.exitCode,
        description: params.description,
      },
      output: [`<stdout>`, stdoutBuffer ?? "", `</stdout>`, `<stderr>`, stderrBuffer ?? "", `</stderr>`].join("\n"),
    }
  },
})
