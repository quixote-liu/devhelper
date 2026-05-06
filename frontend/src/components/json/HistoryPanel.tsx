import { useQuery } from '@tanstack/react-query'
import { historyApi } from '../../api/history'
import { useJsonStore } from '../../store/json'

export function HistoryPanel() {
  const { sessionId, setInput } = useJsonStore()
  const { data: history = [] } = useQuery({
    queryKey: ['history', sessionId],
    queryFn: () => historyApi.list(sessionId),
  })

  return (
    <div className="w-64 border-l border-border flex flex-col">
      <div className="px-4 py-3 border-b border-border font-medium">历史记录</div>
      <div className="flex-1 overflow-auto p-2 space-y-1">
        {history.map((item: any) => (
          <button
            key={item.id}
            onClick={() => setInput(item.content)}
            className="w-full text-left px-3 py-2 text-sm rounded hover:bg-accent"
          >
            <div className="text-xs text-muted-foreground">#{item.seq_num}</div>
            {item.note && <div className="truncate">{item.note}</div>}
          </button>
        ))}
      </div>
    </div>
  )
}
