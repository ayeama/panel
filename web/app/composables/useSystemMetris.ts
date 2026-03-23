import { ref, onBeforeUnmount } from 'vue'

type Point = {
  t: number
  value: number
}

type NetPoint = {
  t: number
  rx: number // bytes/sec
  tx: number // bytes/sec
}

const WINDOW_MS = 5 * 60 * 1000

function parseMetricLine(line: string): {
  cpu: number
  mem: number
  rx: number
  tx: number
} | null {
  const match = line
    .trim()
    .match(/^cpu\s+(\d+(?:\.\d+)?)\s+mem\s+(\d+(?:\.\d+)?)\s+net\s+(\d+)\s+(\d+)$/i)

  if (!match) return null

  const cpu = Number(match[1])
  const mem = Number(match[2])
  const rx = Number(match[3])
  const tx = Number(match[4])

  if ([cpu, mem, rx, tx].some(Number.isNaN)) return null

  return {
    cpu,
    mem,
    rx,
    tx
  }
}

export function useSystemMetrics(wsUrl: string) {
  const cpuPoints = ref<Point[]>([])
  const memPoints = ref<Point[]>([])
  const netPoints = ref<NetPoint[]>([])

  const status = ref<'connecting' | 'open' | 'closed' | 'error'>('connecting')
  const error = ref<string | null>(null)

  let socket: WebSocket | null = null

  // track previous network sample for rate calculation
  let lastNetSample: { t: number; rx: number; tx: number } | null = null

  function prune(now: number) {
    const cutoff = now - WINDOW_MS

    cpuPoints.value = cpuPoints.value.filter((p) => p.t >= cutoff)
    memPoints.value = memPoints.value.filter((p) => p.t >= cutoff)
    netPoints.value = netPoints.value.filter((p) => p.t >= cutoff)
  }

  function pushPoint(cpu: number, mem: number, rxTotal: number, txTotal: number) {
    const t = Date.now()

    cpuPoints.value.push({ t, value: cpu })
    memPoints.value.push({ t, value: mem })

    let rxRate = 0
    let txRate = 0

    if (lastNetSample) {
      const dt = (t - lastNetSample.t) / 1000 // seconds

      if (dt > 0) {
        rxRate = (rxTotal - lastNetSample.rx) / dt
        txRate = (txTotal - lastNetSample.tx) / dt
      }
    }

    // handle counter resets / container restarts
    rxRate = Math.max(0, rxRate)
    txRate = Math.max(0, txRate)

    netPoints.value.push({
      t,
      rx: rxRate,
      tx: txRate
    })

    lastNetSample = {
      t,
      rx: rxTotal,
      tx: txTotal
    }

    prune(t)
  }

  function connect() {
    socket = new WebSocket(wsUrl)
    status.value = 'connecting'
    error.value = null

    socket.addEventListener('open', () => {
      status.value = 'open'
    })

    socket.addEventListener('message', (event) => {
      const raw = String(event.data)
      const lines = raw.split('\n')

      for (const line of lines) {
        if (!line.trim()) continue

        const parsed = parseMetricLine(line)
        if (!parsed) continue

        pushPoint(parsed.cpu, parsed.mem, parsed.rx, parsed.tx)
      }
    })

    socket.addEventListener('close', () => {
      status.value = 'closed'
    })

    socket.addEventListener('error', () => {
      status.value = 'error'
      error.value = 'WebSocket error'
    })
  }

  function disconnect() {
    socket?.close()
    socket = null
  }

  function reconnect() {
    disconnect()
    connect()
  }

  connect()
  onBeforeUnmount(disconnect)

  return {
    cpuPoints,
    memPoints,
    netPoints,
    status,
    error,
    reconnect
  }
}