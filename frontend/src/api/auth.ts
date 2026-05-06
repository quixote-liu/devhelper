import { client } from './client'

export const authApi = {
  login: (username: string, password: string) =>
    client.post('/auth/login', { username, password }).then((r) => r.data.data),
  register: (username: string, email: string, password: string) =>
    client.post('/auth/register', { username, email, password }).then((r) => r.data.data),
  me: () => client.get('/auth/me').then((r) => r.data.data),
}
