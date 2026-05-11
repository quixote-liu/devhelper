import { useState } from 'react'
import { jsonApi } from '../../api/json'
import { useJsonStore } from '../../store/json'

export function FormatTool() {
  const { input, setOutput } = useJsonStore()
  const [indent, setIndent] = useState(2)
  const [error, setError] = useState('')

  const handleFormat = async () => {
    setError('')
    try {
      const res = await jsonApi.format(input, indent)
      setOutput(res.result)
    } catch {
      setError('格式化失败，请检查 JSON 格式')
    }
  }

  const handleMinify = async () => {
    setError('')
    try {
      const res = await jsonApi.minify(input)
      setOutput(res.result)
    } catch {
      setError('压缩失败，请检查 JSON 格式')
    }
  }

  const handleValidate = async () => {
    setError('')
    try {
      await jsonApi.validate(input)
      setOutput('✓ JSON 格式有效')
    } catch {
      setError('JSON 格式无效')
    }
  }

  return (
    <div className="p-4 space-y-4">
      <div className="flex gap-2 items-center">
        <label className="text-sm">缩进空格数：</label>
        <input
          type="number"
          value={indent}
          onChange={(e) => setIndent(Number(e.target.value))}
          className="w-16 px-2 py-1 border border-input rounded bg-background"
          min={1}
          max={8}
        />
      </div>
      <div className="flex gap-2">
        <button onClick={handleFormat} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
          格式化
        </button>
        <button onClick={handleMinify} className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/90">
          压缩
        </button>
        <button onClick={handleValidate} className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/90">
          验证
        </button>
      </div>
      {error && <p className="text-sm text-destructive">{error}</p>}
    </div>
  )
}
