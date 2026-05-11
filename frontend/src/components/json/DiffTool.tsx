import { useState } from 'react'
import { jsonApi } from '../../api/json'
import { EditorPanel } from './EditorPanel'

export function DiffTool() {
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [result, setResult] = useState('')
  const [error, setError] = useState('')

  const handleDiff = async () => {
    setError('')
    try {
      const res = await jsonApi.diff(left, right)
      setResult(JSON.stringify(res.diff, null, 2))
    } catch {
      setError('比较失败，请检查 JSON 格式')
    }
  }

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b border-border flex gap-2">
        <button onClick={handleDiff} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
          比较差异
        </button>
        {error && <p className="text-sm text-destructive self-center">{error}</p>}
      </div>
      <div className="flex flex-1 min-h-0">
        <div className="flex-1 border-r border-border">
          <div className="px-3 py-1 text-xs text-muted-foreground border-b border-border">左侧 JSON</div>
          <div className="h-[calc(100%-28px)]">
            <EditorPanel value={left} onChange={setLeft} />
          </div>
        </div>
        <div className="flex-1 border-r border-border">
          <div className="px-3 py-1 text-xs text-muted-foreground border-b border-border">右侧 JSON</div>
          <div className="h-[calc(100%-28px)]">
            <EditorPanel value={right} onChange={setRight} />
          </div>
        </div>
        <div className="flex-1">
          <div className="px-3 py-1 text-xs text-muted-foreground border-b border-border">差异结果</div>
          <div className="h-[calc(100%-28px)]">
            <EditorPanel value={result} readOnly />
          </div>
        </div>
      </div>
    </div>
  )
}
