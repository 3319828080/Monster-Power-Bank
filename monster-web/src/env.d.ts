/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface AMapEvent {
  (event: string, handler: (...args: any[]) => void): void
}

interface AMapInstance {
  setCenter: (center: [number, number]) => void
  setZoom: (zoom: number) => void
  setFitView: (markers: any[], immediate?: boolean, margin?: number[]) => void
  add: (item: any) => void
  remove: (item: any) => void
  destroy: () => void
  on: AMapEvent
  off: AMapEvent
}

interface AMapMarker {
  new (opts: {
    position: [number, number]
    title?: string
    label?: { content: string; offset?: [number, number] }
  }): {
    on: AMapEvent
    off: AMapEvent
  }
}

interface AMapClass {
  Map: new (container: HTMLElement | undefined, opts: {
    zoom: number
    center: [number, number]
    mapStyle?: string
  }) => AMapInstance
  Marker: AMapMarker
}

interface Window {
  AMap: AMapClass
}
