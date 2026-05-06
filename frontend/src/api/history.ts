import { client } from './client'

export const historyApi = {
  list: (sessionId: string) =>
    client.get('/history', { params: { session_id: sessionId } }).then((r) => r.data.data),
  save: (sessionId: string, content: string, note?: string) =>
    client.post('/history', { session_id: sessionId, content, note }).then((r) => r.data.data),
  delete: (id: number) =>
    client.delete(`/history/${id}`).then((r) => r.data.data),
}
