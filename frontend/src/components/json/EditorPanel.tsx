import Editor from '@monaco-editor/react'

interface EditorPanelProps {
  value: string
  onChange?: (value: string) => void
  language?: string
  readOnly?: boolean
}

export function EditorPanel({ value, onChange, language = 'json', readOnly = false }: EditorPanelProps) {
  const theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'vs-dark' : 'light'

  return (
    <Editor
      height="100%"
      language={language}
      value={value}
      onChange={(v) => onChange?.(v || '')}
      theme={theme}
      options={{
        readOnly,
        minimap: { enabled: false },
        fontSize: 14,
        lineNumbers: 'on',
        scrollBeyondLastLine: false,
        automaticLayout: true,
      }}
    />
  )
}
