import { client } from './client'

export const schemasApi = {
  list: () => client.get('/schemas').then((r) => r.data.data),
  create: (name: string, description: string, schema: string, isPublic = false) =>
    client.post('/schemas', { name, description, schema, is_public: isPublic }).then((r) => r.data.data),
  get: (id: number) => client.get(`/schemas/${id}`).then((r) => r.data.data),
  update: (id: number, name: string, description: string, schema: string, isPublic = false) =>
    client.put(`/schemas/${id}`, { name, description, schema, is_public: isPublic }).then((r) => r.data.data),
  delete: (id: number) => client.delete(`/schemas/${id}`).then((r) => r.data.data),
}
