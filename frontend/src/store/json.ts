import { create } from 'zustand'

interface JsonState {
  input: string
  output: string
  sessionId: string
  setInput: (input: string) => void
  setOutput: (output: string) => void
  newSession: () => void
}

export const useJsonStore = create<JsonState>((set) => ({
  input: '',
  output: '',
  sessionId: crypto.randomUUID(),
  setInput: (input) => set({ input }),
  setOutput: (output) => set({ output }),
  newSession: () => set({ sessionId: crypto.randomUUID(), input: '', output: '' }),
}))
