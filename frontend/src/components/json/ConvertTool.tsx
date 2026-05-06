import { useState } from 'react'
import { jsonApi } from '../../api/json'
import { useJsonStore } from '../../store/json'

export function ConvertTool() {
  const { input, setOutput } = useJsonStore()
  const [from, setFrom] = useState('json')
  const [to, setTo] = useState('yaml')
  const [error, setError] = useState('')

  const handleConvert = async () => {
    setError('')
    try {
      if (from === 'json') {
        const res = await jsonApi.convert(input, from, to)
        setOutput(res.converted)
      } else {
        const res = await jsonApi.parse(input, from)
        setOutput(res.parsed)
      }
    } catch {
      setError('转换失败，请检查输入格式')
    }
  }

  return (
    <div className="p-4 space-y-4">
      <div className="flex gap-4 items-center">
        <div>
          <label className="text-sm mr-2">从：</label>
          <select value={from} onChange={(e) => setFrom(e.target.value)} className="px-2 py-1 border border-input rounded bg-background">
            <option value="json">JSON</option>
            <option value="yaml">YAML</option>
            <option value="xml">XML</option>
            <option value="toml">TOML</option>
          </select>
        </div>
        <div>
          <label className="text-sm mr-2">到：</label>
          <select value={to} onChange={(e) => setTo(e.target.value)} className="px-2 py-1 border border-input rounded bg-background" disabled={from !== 'json'}>
            <option value="yaml">YAML</option>
            <option value="xml">XML</option>
            <option value="toml">TOML</option>
          </select>
        </div>
      </div>
      <button onClick={handleConvert} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
        转换
      </button>
      {error && <p className="text-sm text-destructive">{error}</p>}
    </div>
  )
}
