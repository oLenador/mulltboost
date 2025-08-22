import clsx from "clsx"
import * as React from "react"
import { createContext, useContext, useState, useCallback, useRef, useEffect } from "react"

type LayerType = 'toast' | 'dialog' | 'popover' | 'dropdown' | 'tooltip' | 'notification' | 'custom'

interface LayerConfig {
  id: string
  type: LayerType
  zIndex: number
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'center'
  priority: number // Maior número = maior prioridade
}

interface FloatManagerContextType {
  registerLayer: (config: Omit<LayerConfig, 'zIndex'>) => number
  unregisterLayer: (id: string) => void
  getLayerZIndex: (id: string) => number | undefined
  isLayerActive: (id: string) => boolean
  getActiveLayersByType: (type: LayerType) => LayerConfig[]
}

const BASE_Z_INDEXES: Record<LayerType, number> = {
  tooltip: 1000,
  dropdown: 1010,
  popover: 1020,
  notification: 1030,
  toast: 1040,
  dialog: 1050,
  custom: 1060
}

const FloatManagerContext = createContext<FloatManagerContextType | null>(null)

interface FloatManagerProviderProps {
  children: React.ReactNode
}

export const FloatManagerProvider: React.FC<FloatManagerProviderProps> = ({ children }) => {
  const [layers, setLayers] = useState<Map<string, LayerConfig>>(new Map())
  const layerCounterRef = useRef<Record<LayerType, number>>({
    toast: 0,
    dialog: 0,
    popover: 0,
    dropdown: 0,
    tooltip: 0,
    notification: 0,
    custom: 0
  })

  const registerLayer = useCallback((config: Omit<LayerConfig, 'zIndex'>) => {
    const { id, type, priority = 0 } = config

    // Incrementa o contador para o tipo
    layerCounterRef.current[type] += 1

    // Calcula z-index baseado no tipo + contador + prioridade
    const zIndex = BASE_Z_INDEXES[type] + layerCounterRef.current[type] + priority

    const layerConfig: LayerConfig = {
      ...config,
      zIndex
    }

    setLayers(prev => new Map(prev).set(id, layerConfig))

    return zIndex
  }, [])

  const unregisterLayer = useCallback((id: string) => {
    setLayers(prev => {
      const newMap = new Map(prev)
      const layer = newMap.get(id)

      if (layer) {
        // Decrementa o contador
        layerCounterRef.current[layer.type] = Math.max(0, layerCounterRef.current[layer.type] - 1)
        newMap.delete(id)
      }

      return newMap
    })
  }, [])

  const getLayerZIndex = useCallback((id: string) => {
    return layers.get(id)?.zIndex
  }, [layers])

  const isLayerActive = useCallback((id: string) => {
    return layers.has(id)
  }, [layers])

  const getActiveLayersByType = useCallback((type: LayerType) => {
    return Array.from(layers.values()).filter(layer => layer.type === type)
  }, [layers])

  const value = {
    registerLayer,
    unregisterLayer,
    getLayerZIndex,
    isLayerActive,
    getActiveLayersByType
  }

  return (
    <FloatManagerContext.Provider value={value}>
      {children}
    </FloatManagerContext.Provider>
  )
}

export const useFloatManager = () => {
  const context = useContext(FloatManagerContext)
  if (!context) {
    throw new Error('useFloatManager deve ser usado dentro de FloatManagerProvider')
  }
  return context
}

export const useFloatLayer = (
  config: Omit<LayerConfig, 'zIndex'>,
  active: boolean = true
) => {
  const { registerLayer, unregisterLayer, getLayerZIndex } = useFloatManager()
  const [zIndex, setZIndex] = useState<number | undefined>()

  useEffect(() => {
    if (active) {
      const z = registerLayer(config)
      setZIndex(z)

      return () => {
        unregisterLayer(config.id)
        setZIndex(undefined)
      }
    } else {
      unregisterLayer(config.id)
      setZIndex(undefined)
    }
  }, [active, config.id, config.type, config.priority, registerLayer, unregisterLayer])

  return zIndex
}

interface FloatElementProps {
  id: string
  type: LayerType
  priority?: number
  position?: LayerConfig['position']
  active?: boolean
  className?: string
  style?: React.CSSProperties
  children: React.ReactNode
}

export const FloatElement: React.FC<FloatElementProps> = ({
  id,
  type,
  priority = 0,
  position = 'center',
  active = true,
  className,
  style,
  children
}) => {
  const zIndex = useFloatLayer({ id, type, priority, position }, active)

  if (!active || zIndex === undefined) {
    return null
  }

  const positionClasses = {
    'top-right': 'fixed top-10 right-10',
    'top-left': 'fixed top-10 left-10',
    'bottom-right': 'fixed bottom-10 right-8',
    'bottom-left': 'fixed bottom-10 left-10',
    'center': 'fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2'
  }

  return (
    <div
      className={clsx(positionClasses[position], className)}
      style={{
        zIndex,
        ...style
      }}
    >
      {children}
    </div>
  )
}

export const useLayerConflicts = (type: LayerType) => {
  const { getActiveLayersByType } = useFloatManager()

  return useCallback(() => {
    const activeLayers = getActiveLayersByType(type)
    return activeLayers.length > 1
  }, [getActiveLayersByType, type])
}
