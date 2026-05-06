import { useState } from 'react'
import { jsonApi } from '../../api/json'
import { useJsonStore } from '../../store/json'

export function QueryTool() {
  const { input, setOutput } = useJsonStore()
  const [path, setPath] = useState('$')
  const [error, setError] = useState('')

  const handleQuery = async () => {
    setError('')
    try {
      const res = await jsonApi.query(input, path)
      setOutput(JSON.stringify(res.result, null, 2))
    } catch {
      setError('查询失败，请检查路径格式')
    }
  }

  return (
    <div className="p-4 space-y-4">
      <div>
        <label className="block text-sm mb-1">JSONPath：</label>
        <input
          type="text"
          value={path}
          onChange={(e) => setPath(e.target.value)}
          className="w-full px-3 py-2 border border-input rounded bg-background"
          placeholder="例如：$.store.book[0].title"
        />
      </div>
      <button onClick={handleQuery} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
        查询
      </button>
      {error && <p className="text-sm text-destructive">{error}</p>}
    </div>
  )
}
