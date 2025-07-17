// @ts-check
import { defineConfig } from "astro/config"
import starlight from "@astrojs/starlight"
import solidJs from "@astrojs/solid-js"
import cloudflare from "@astrojs/cloudflare"
import theme from "toolbeam-docs-theme"
import config from "./config.mjs"
import { rehypeHeadingIds } from "@astrojs/markdown-remark"
import rehypeAutolinkHeadings from "rehype-autolink-headings"
import { spawnSync } from "child_process"

// https://astro.build/config
export default defineConfig({
  site: config.url,
  output: "server",
  adapter: cloudflare({
    imageService: "passthrough",
  }),
  devToolbar: {
    enabled: false,
  },
  server: {
    host: "0.0.0.0",
  },
  markdown: {
    rehypePlugins: [rehypeHeadingIds, [rehypeAutolinkHeadings, { behavior: "wrap" }]],
  },
  build: {},
  integrations: [
    solidJs(),
    starlight({
      title: "autoprovisioner",
      lastUpdated: true,
      pagefind: true,
      expressiveCode: { themes: ["github-light", "github-dark"] },
      social: [{ icon: "github", label: "GitHub", href: config.github }],
      head: [
        {
          tag: "link",
          attrs: {
            rel: "icon",
            href: "/favicon.svg",
          },
        },
      ],
      editLink: {
        baseUrl: `https://github.com/zerosync-co/autoprovisioner-cli/edit/dev/packages/web/`,
      },
      markdown: {
        headingLinks: false,
      },
      customCss: ["./src/styles/custom.css"],
      sidebar: ["docs", "docs/cli", "docs/troubleshooting"],
      components: {
        Hero: "./src/components/Hero.astro",
        Head: "./src/components/Head.astro",
        Header: "./src/components/Header.astro",
        Footer: "./src/components/Footer.astro",
      },
      plugins: [
        theme({
          headerLinks: config.headerLinks,
        }),
      ],
    }),
  ],
})

// function configSchema() {
//   return {
//     name: "configSchema",
//     hooks: {
//       "astro:build:done": async () => {
//         console.log("generating config schema")
//         spawnSync("../opencode/script/schema.ts", ["./dist/config.json"])
//       },
//     },
//   }
// }
