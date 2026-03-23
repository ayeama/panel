<script setup lang="ts">
import { computed } from 'vue'

type Point = {
  t: number
  value: number
}

type Series = {
  key: string
  label: string
  color: string
  points: Point[]
}

const props = withDefaults(defineProps<{
  title: string
  series: Series[]
  min?: number
  max?: number
  height?: number
  yTicks?: number[]
  fillArea?: boolean
}>(), {
  min: 0,
  max: 100,
  height: 160,
  yTicks: () => [0, 25, 50, 75, 100],
  fillArea: false
})

const width = 1000
const paddingX = 8
const paddingY = 4
const windowMs = 5 * 60 * 1000

const latestTimestamp = computed(() => {
  let latest = Date.now()

  for (const s of props.series) {
    const last = s.points[s.points.length - 1]
    if (last && last.t > latest) {
      latest = last.t
    }
  }

  return latest
})

const startTime = computed(() => latestTimestamp.value - windowMs)

const visibleSeries = computed(() => {
  return props.series.map((s) => ({
    ...s,
    points: s.points.filter((p) => p.t >= startTime.value)
  }))
})

function xFor(t: number) {
  const ratio = (t - startTime.value) / windowMs
  return paddingX + ratio * (width - paddingX * 2)
}

function yFor(value: number) {
  const min = props.min
  const max = props.max
  const clamped = Math.max(min, Math.min(max, value))
  const ratio = (clamped - min) / (max - min || 1)

  return props.height - paddingY - ratio * (props.height - paddingY * 2)
}

function buildLinePath(points: Point[]) {
  if (points.length === 0) {
    return ''
  }

  return points
    .map((p, i) => `${i === 0 ? 'M' : 'L'} ${xFor(p.t)} ${yFor(p.value)}`)
    .join(' ')
}

function buildAreaPath(points: Point[]) {
  if (points.length === 0) {
    return ''
  }

  const bottom = props.height - paddingY
  const first = points[0]
  const last = points[points.length - 1]

  const top = points
    .map((p, i) => `${i === 0 ? 'M' : 'L'} ${xFor(p.t)} ${yFor(p.value)}`)
    .join(' ')

  return `${top} L ${xFor(last.t)} ${bottom} L ${xFor(first.t)} ${bottom} Z`
}

const renderedSeries = computed(() => {
  return visibleSeries.value.map((s) => ({
    ...s,
    linePath: buildLinePath(s.points),
    areaPath: buildAreaPath(s.points),
    lastPoint: s.points[s.points.length - 1] ?? null
  }))
})

const legend = computed(() => {
  return renderedSeries.value.map((s) => ({
    key: s.key,
    label: s.label,
    color: s.color,
    latest: s.lastPoint ? s.lastPoint.value : null
  }))
})

function gradientId(key: string) {
  return `chart-fill-${key}`
}
</script>

<template>
  <UCard
    class="overflow-hidden"
    :ui="{
      body: 'p-3 sm:p-3'
    }"
  >
    <template #header>
      <div class="flex items-center justify-between gap-4">
        <div class="font-medium">
          {{ title }}
        </div>

        <div class="flex items-center gap-4 text-xs">
          <div
            v-for="item in legend"
            :key="item.key"
            class="flex items-center gap-2"
          >
            <span
              class="inline-block h-2.5 w-2.5 rounded-full"
              :style="{ backgroundColor: item.color }"
            />
            <span class="text-muted">
              {{ item.label }}
            </span>
            <span class="tabular-nums">
              {{ item.latest !== null ? item.latest.toFixed(2) : '--' }}
            </span>
          </div>
        </div>
      </div>
    </template>

    <svg
      :viewBox="`0 0 ${width} ${height}`"
      class="h-40 w-full"
      preserveAspectRatio="none"
    >
      <defs>
        <linearGradient
          v-for="s in renderedSeries"
          :id="gradientId(s.key)"
          :key="gradientId(s.key)"
          x1="0"
          y1="0"
          x2="0"
          y2="1"
        >
          <stop offset="0%" :stop-color="s.color" stop-opacity="0.28" />
          <stop offset="100%" :stop-color="s.color" stop-opacity="0.03" />
        </linearGradient>
      </defs>

      <g>
        <line
          v-for="tick in yTicks"
          :key="tick"
          :x1="paddingX"
          :x2="width - paddingX"
          :y1="yFor(tick)"
          :y2="yFor(tick)"
          stroke="currentColor"
          stroke-opacity="0.12"
          stroke-width="1"
        />
      </g>

      <g v-for="s in renderedSeries" :key="s.key">
        <path
          v-if="fillArea && s.areaPath"
          :d="s.areaPath"
          :fill="`url(#${gradientId(s.key)})`"
        />

        <path
          v-if="s.linePath"
          :d="s.linePath"
          fill="none"
          :stroke="s.color"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          vector-effect="non-scaling-stroke"
        />

        <circle
          v-if="s.lastPoint"
          :cx="xFor(s.lastPoint.t)"
          :cy="yFor(s.lastPoint.value)"
          r="3"
          :fill="s.color"
        />
      </g>
    </svg>
  </UCard>
</template>