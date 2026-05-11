import { useState } from 'react'
import { jsonApi } from '../../api/json'
import { useJsonStore } from '../../store/json'

export function SchemaTool() {
  const { input, setOutput } = useJsonStore()
  const [schema, setSchema] = useState('')
  const [error, setError] = useState('')

  const handleGenerate = async () => {
    setError('')
    try {
      const res = await jsonApi.generateSchema(input)
      setOutput(res.schema)
    } catch {
      setError('生成失败，请检查 JSON 格式')
    }
  }

  const handleValidate = async () => {
    setError('')
    try {
      await jsonApi.validateSchema(schema, input)
      setOutput('✓ 数据符合 Schema')
    } catch {
      setError('验证失败，数据不符合 Schema')
    }
  }

  return (
    <div className="p-4 space-y-4">
      <div className="flex gap-2">
        <button onClick={handleGenerate} className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
          生成 Schema
        </button>
      </div>
      <div>
        <label className="block text-sm mb-1">Schema（用于验证）：</label>
        <textarea
          value={schema}
          onChange={(e) => setSchema(e.target.value)}
          className="w-full h-32 px-3 py-2 border border-input rounded bg-background font-mono text-sm"
          placeholder="粘贴 JSON Schema..."
        />
      </div>
      <button onClick={handleValidate} className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/90">
        验证数据
      </button>
      {error && <p className="text-sm text-destructive">{error}</p>}
    </div>
  )
}
