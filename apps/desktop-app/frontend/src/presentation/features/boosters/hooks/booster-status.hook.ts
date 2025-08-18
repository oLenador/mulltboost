import { useState, useEffect, useCallback } from 'react'

export interface LoadingItem {
  id: string
  title: string
  status: "idle" | "error" | "applied" | "applying"
  progress?: number
}

interface BoosterStatusState {
  items: LoadingItem[]
  totalItems: number
  completedItems: number
  isLoading: boolean
  progress: number
}

function useBoosterStatus(initialItems: LoadingItem[] = []) {
  const [state, setState] = useState<BoosterStatusState>({
    items: initialItems,
    totalItems: initialItems.length,
    completedItems: initialItems.filter(item => item.status === 'applied').length,
    isLoading: false,
    progress: 0
  })

  // Calcula o progresso automaticamente
  useEffect(() => {
    const completed = state.items.filter(item => item.status === 'applied').length
    const applying = state.items.filter(item => item.status === 'applying').length
    const isLoading = applying > 0
    const progress = state.totalItems > 0 ? (completed / state.totalItems) * 100 : 0

    setState(prev => ({
      ...prev,
      completedItems: completed,
      isLoading,
      progress
    }))
  }, [state.items, state.totalItems])

  const addItem = useCallback((item: LoadingItem) => {
    setState(prev => ({
      ...prev,
      items: [...prev.items, item],
      totalItems: prev.totalItems + 1
    }))
  }, [])

  const updateItem = useCallback((id: string, updates: Partial<LoadingItem>) => {
    setState(prev => ({
      ...prev,
      items: prev.items.map(item =>
        item.id === id ? { ...item, ...updates } : item
      )
    }))
  }, [])

  const removeItem = useCallback((id: string) => {
    setState(prev => ({
      ...prev,
      items: prev.items.filter(item => item.id !== id),
      totalItems: Math.max(0, prev.totalItems - 1)
    }))
  }, [])

  const startProcess = useCallback(async (itemId: string, processFunction: () => Promise<void>) => {
    updateItem(itemId, { status: 'applying' })
    try {
      await processFunction()
      updateItem(itemId, { status: 'applied' })
    } catch (error) {
      updateItem(itemId, { status: 'error' })
      throw error
    }
  }, [updateItem])

  const reset = useCallback(() => {
    setState(prev => ({
      ...prev,
      items: prev.items.map(item => ({ ...item, status: 'idle' as const }))
    }))
  }, [])

  return {
    ...state,
    addItem,
    updateItem,
    removeItem,
    startProcess,
    reset
  }
}