import { NavLink, useNavigate } from 'react-router-dom'
import { Braces, Database, User, Shield, LogOut } from 'lucide-react'
import { useAuthStore } from '../../store/auth'
import { cn } from '../../lib/utils'

const nav = [
  { to: '/json', icon: Braces, label: 'JSON 工具' },
  { to: '/schemas', icon: Database, label: 'Schema 管理' },
  { to: '/profile', icon: User, label: '个人设置' },
]

export function Sidebar() {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  return (
    <aside className="w-52 shrink-0 border-r border-border flex flex-col bg-muted/30">
      <div className="px-4 py-4 border-b border-border">
        <span className="font-semibold text-foreground">DevHelper</span>
      </div>
      <nav className="flex-1 p-2 space-y-1">
        {nav.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              cn('flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors',
                isActive
                  ? 'bg-primary/10 text-primary font-medium'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground')
            }
          >
            <Icon size={16} />
            {label}
          </NavLink>
        ))}
        {user?.role === 'admin' && (
          <NavLink
            to="/admin"
            className={({ isActive }) =>
              cn('flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors',
                isActive
                  ? 'bg-primary/10 text-primary font-medium'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground')
            }
          >
            <Shield size={16} />
            管理员
          </NavLink>
        )}
      </nav>
      <div className="p-2 border-t border-border">
        <div className="px-3 py-1 text-xs text-muted-foreground truncate">{user?.username}</div>
        <button
          onClick={() => { logout(); navigate('/login') }}
          className="flex items-center gap-2 px-3 py-2 rounded-md text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground w-full transition-colors"
        >
          <LogOut size={16} />
          退出登录
        </button>
      </div>
    </aside>
  )
}
