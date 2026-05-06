import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { authApi } from '../api/auth'
import { client } from '../api/client'
import { useAuthStore } from '../store/auth'

export function ProfilePage() {
  const { user, setUser } = useAuthStore()
  const [pwd, setPwd] = useState({ current: '', newPwd: '', confirm: '' })
  const [msg, setMsg] = useState('')

  useQuery({ queryKey: ['me'], queryFn: async () => {
    const u = await authApi.me(); setUser(u); return u
  }})

  const updatePwd = async (e: React.FormEvent) => {
    e.preventDefault()
    if (pwd.newPwd !== pwd.confirm) { setMsg('两次密码不一致'); return }
    try {
      await client.put('/user/password', { current_password: pwd.current, new_password: pwd.newPwd })
      setMsg('密码修改成功')
      setPwd({ current: '', newPwd: '', confirm: '' })
    } catch { setMsg('修改失败') }
  }

  return (
    <div className="p-6 max-w-lg space-y-6">
      <h2 className="text-xl font-semibold">个人设置</h2>
      <div className="border border-border rounded-lg p-4 space-y-2">
        <div className="text-sm"><span className="text-muted-foreground">用户名：</span>{user?.username}</div>
        <div className="text-sm"><span className="text-muted-foreground">邮箱：</span>{user?.email}</div>
        <div className="text-sm"><span className="text-muted-foreground">角色：</span>{user?.role}</div>
      </div>
      <div className="border border-border rounded-lg p-4 space-y-3">
        <h3 className="font-medium">修改密码</h3>
        <form onSubmit={updatePwd} className="space-y-3">
          <input type="password" placeholder="当前密码" value={pwd.current} onChange={(e) => setPwd({ ...pwd, current: e.target.value })}
            className="w-full px-3 py-2 border border-input rounded bg-background" required />
          <input type="password" placeholder="新密码" value={pwd.newPwd} onChange={(e) => setPwd({ ...pwd, newPwd: e.target.value })}
            className="w-full px-3 py-2 border border-input rounded bg-background" required />
          <input type="password" placeholder="确认新密码" value={pwd.confirm} onChange={(e) => setPwd({ ...pwd, confirm: e.target.value })}
            className="w-full px-3 py-2 border border-input rounded bg-background" required />
          {msg && <p className="text-sm text-muted-foreground">{msg}</p>}
          <button type="submit" className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">保存</button>
        </form>
      </div>
    </div>
  )
}
