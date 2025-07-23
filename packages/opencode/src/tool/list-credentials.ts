import { z } from "zod"
import { Tool } from "./tool"

export const ListCredentialsTool = Tool.define({
  id: "list-credentials",
  description:
    "Lists available credential environment variables. Use this tool BEFORE attempting operations that may require authentication (deployments, API calls, cloud operations) to verify credentials are available.",
  parameters: z.object({}),
  async execute(_params, _ctx) {
    // Get all environment variables that match credential patterns
    const credentialPatterns = [
      /.*TOKEN.*/i,
      /.*KEY.*/i,
      /.*SECRET.*/i,
      /.*PASSWORD.*/i,
      /.*CREDENTIAL.*/i,
      /.*AUTH.*/i,
    ]

    const availableCredentials = Object.keys(process.env)
      .filter((key) => credentialPatterns.some((pattern) => pattern.test(key)))
      .sort()

    const output =
      availableCredentials.length > 0 ? availableCredentials.join("\n") : "No credential environment variables found"

    return {
      title: "Available Credentials",
      metadata: {
        count: availableCredentials.length,
        credentials: availableCredentials,
      },
      output,
    }
  },
})
