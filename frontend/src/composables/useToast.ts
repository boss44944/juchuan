import { readonly, ref } from 'vue'

export type ToastTone = 'success' | 'warning' | 'error' | 'info'

export interface ToastItem {
  id: number
  message: string
  tone: ToastTone
}

const items = ref<ToastItem[]>([])
let nextID = 1

function dismiss(id: number) {
  items.value = items.value.filter((item) => item.id !== id)
}

function show(message: string, tone: ToastTone = 'info') {
  const id = nextID++
  items.value.push({ id, message, tone })
  window.setTimeout(() => dismiss(id), 3600)
}

export function useToast() {
  return {
    items: readonly(items),
    dismiss,
    success: (message: string) => show(message, 'success'),
    warning: (message: string) => show(message, 'warning'),
    error: (message: string) => show(message, 'error'),
    info: (message: string) => show(message, 'info'),
  }
}
