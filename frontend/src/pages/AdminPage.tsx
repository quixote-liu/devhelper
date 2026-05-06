import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { client } from '../api/client'

export function AdminPage() {
  const qc = useQueryClient()
  const { data: users = [] } = useQuery({
    queryKey: ['admin-users'],
    queryFn: () => client.get('/admin/users').then((r) => r.data.data),
  })
  const deleteMut = useMutation({
    mutationFn: (id: number) => client.delete(`/admin/users/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-users'] }),
  })

  return (
    <div className="p-6 space-y-4">
      <h2 className="text-xl font-semibold">用户管理</h2>
      <div className="border border-border rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-muted">
            <tr>
              <th className="px-4 py-2 text-left">ID</th>
              <th className="px-4 py-2 text-left">用户名</th>
              <th className="px-4 py-2 text-left">邮箱</th>
              <th className="px-4 py-2 text-left">角色</th>
              <th className="px-4 py-2 text-left">操作</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u: any) => (
              <tr key={u.id} className="border-t border-border">
                <td className="px-4 py-2">{u.id}</td>
                <td className="px-4 py-2">{u.username}</td>
                <td className="px-4 py-2">{u.email}</td>
                <td className="px-4 py-2">{u.role}</td>
                <td className="px-4 py-2">
                  <button onClick={() => deleteMut.mutate(u.id)} className="text-destructive hover:underline">删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
