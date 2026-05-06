import { useState } from 'react'
import { EditorPanel } from '../components/json/EditorPanel'
import { FormatTool } from '../components/json/FormatTool'
import { ConvertTool } from '../components/json/ConvertTool'
import { SchemaTool } from '../components/json/SchemaTool'
import { DiffTool } from '../components/json/DiffTool'
import { QueryTool } from '../components/json/QueryTool'
import { HistoryPanel } from '../components/json/HistoryPanel'
import { useJsonStore } from '../store/json'

const tabs = [
  { id: 'format', label: '格式化' },
  { id: 'convert', label: '转换' },
  { id: 'schema', label: 'Schema' },
  { id: 'diff', label: 'Diff' },
  { id: 'query', label: '查询' },
]

export function JsonPage() {
  const [tab, setTab] = useState('format')
  const { input, output, setInput } = useJsonStore()

  return (
    <div className="flex h-full">
      <div className="flex-1 flex flex-col">
        <div className="border-b border-border flex">
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`px-4 py-2 text-sm ${tab === t.id ? 'border-b-2 border-primary text-primary' : 'text-muted-foreground'}`}
            >
              {t.label}
            </button>
          ))}
        </div>
        {tab === 'diff' ? (
          <DiffTool />
        ) : (
          <>
            <div className="p-4 border-b border-border">
              {tab === 'format' && <FormatTool />}
              {tab === 'convert' && <ConvertTool />}
              {tab === 'schema' && <SchemaTool />}
              {tab === 'query' && <QueryTool />}
            </div>
            <div className="flex flex-1 min-h-0">
              <div className="flex-1 border-r border-border">
                <div className="px-3 py-1 text-xs text-muted-foreground border-b border-border">输入</div>
                <div className="h-[calc(100%-28px)]">
                  <EditorPanel value={input} onChange={setInput} />
                </div>
              </div>
              <div className="flex-1">
                <div className="px-3 py-1 text-xs text-muted-foreground border-b border-border">输出</div>
                <div className="h-[calc(100%-28px)]">
                  <EditorPanel value={output} readOnly />
                </div>
              </div>
            </div>
          </>
        )}
      </div>
      <HistoryPanel />
    </div>
  )
}
