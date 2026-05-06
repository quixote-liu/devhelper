import { client } from './client'

export const jsonApi = {
  validate: (content: string) =>
    client.post('/json/validate', { content }).then((r) => r.data.data),
  format: (content: string, indent = 2) =>
    client.post('/json/format', { content, indent }).then((r) => r.data.data),
  minify: (content: string) =>
    client.post('/json/minify', { content }).then((r) => r.data.data),
  convert: (content: string, from: string, to: string) =>
    client.post('/json/convert', { content, from, to }).then((r) => r.data.data),
  parse: (content: string, from: string) =>
    client.post('/json/parse', { content, from }).then((r) => r.data.data),
  generateSchema: (content: string) =>
    client.post('/json/schema/generate', { content }).then((r) => r.data.data),
  validateSchema: (content: string, schema: string) =>
    client.post('/json/schema/validate', { content, schema }).then((r) => r.data.data),
  diff: (left: string, right: string) =>
    client.post('/json/diff', { left, right }).then((r) => r.data.data),
  query: (content: string, path: string) =>
    client.post('/json/query', { content, path }).then((r) => r.data.data),
}
