import { client } from './client'

export const jsonApi = {
  validate: (json: string) =>
    client.post('/json/validate', { json }).then((r) => r.data.data),
  format: (json: string, indent = 2) =>
    client.post('/json/format', { json, indent }).then((r) => r.data.data),
  minify: (json: string) =>
    client.post('/json/minify', { json }).then((r) => r.data.data),
  convert: (json: string, target: string) =>
    client.post('/json/convert', { json, target }).then((r) => r.data.data),
  parse: (content: string, source: string) =>
    client.post('/json/parse', { content, source }).then((r) => r.data.data),
  generateSchema: (json: string) =>
    client.post('/json/schema/generate', { json }).then((r) => r.data.data),
  validateSchema: (schema: string, data: string) =>
    client.post('/json/schema/validate', { schema, data }).then((r) => r.data.data),
  diff: (a: string, b: string) =>
    client.post('/json/diff', { a, b }).then((r) => r.data.data),
  query: (json: string, path: string) =>
    client.post('/json/query', { json, path }).then((r) => r.data.data),
}
