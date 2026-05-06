import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { schemasApi } from '../api/schemas'

export function SchemasPage() {
  const qc = useQueryClient()
  const { data: schemas = [] } = useQuery({ queryKey: ['schemas'], queryFn: schemasApi.list })
  const [editing, setEditing] = useState<any>(null)
  const [form, setForm] = useState({ name: '', description: '', schema: '' })

  const createMut = useMutation({
    mutationFn: () => schemasApi.create(form.name, form.description, form.schema),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['schemas'] }); setForm({ name: '', description: '', schema: '' }) },
  })
  const updateMut = useMutation({
    mutationFn: () => schemasApi.update(editing.id, form.name, form.description, form.schema),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['schemas'] }); setEditing(null) },
  })
  const deleteMut = useMutation({
    mutationFn: (id: number) => schemasApi.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['schemas'] }),
  })

  const startEdit = (s: any) => { setEditing(s); setForm({ name: s.name, description: s.description, schema: s.schema }) }

  return (
    <div className="p-6 space-y-6">
      <h2 className="text-xl font-semibold">Schema 管理</h2>
      <div className="border border-border rounded-lg p-4 space-y-3">
        <h3 className="font-medium">{editing ? '编辑 Schema' : '新建 Schema'}</h3>
        <input placeholder="名称" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
          className="w-full px-3 py-2 border border-input rounded bg-background" />
        <input placeholder="描述" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })}
          className="w-full px-3 py-2 border border-input rounded bg-background" />
        <textarea placeholder="Schema JSON" value={form.schema} onChange={(e) => setForm({ ...form, schema: e.target.value })}
          className="w-full h-32 px-3 py-2 border border-input rounded bg-background font-mono text-sm" />
        <div className="flex gap-2">
          <button onClick={() => editing ? updateMut.mutate() : createMut.mutate()}
            className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
            {editing ? '保存' : '创建'}
          </button>
          {editing && <button onClick={() => setEditing(null)} className="px-4 py-2 bg-secondary text-secondary-foreground rounded">取消</button>}
        </div>
      </div>
      <div className="space-y-2">
        {schemas.map((s: any) => (
          <div key={s.id} className="border border-border rounded-lg p-4 flex justify-between items-start">
            <div>
              <div className="font-medium">{s.name}</div>
              {s.description && <div className="text-sm text-muted-foreground">{s.description}</div>}
            </div>
            <div className="flex gap-2">
              <button onClick={() => startEdit(s)} className="text-sm text-primary hover:underline">编辑</button>
              <button onClick={() => deleteMut.mutate(s.id)} className="text-sm text-destructive hover:underline">删除</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
