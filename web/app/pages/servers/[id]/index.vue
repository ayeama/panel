<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useClipboard } from '@vueuse/core'
import { AttachAddon } from '@xterm/addon-attach';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import ServerStatusBadge from '~/components/ServerStatusBadge.vue';

definePageMeta({
  layout: 'default'
})

let resizeObserver: ResizeObserver | null = null

const route = useRoute()
const id = computed(() => String(route.params.id))

useHead({
  title: `Panel - ${id.value}`
})

const { copy, copied } = useClipboard()

const { readOne, deleteOne, start, stop, restore } = useServers()
const { data: server, status, error, refresh } = await useAsyncData(
  () => `server-${id.value}`,
  () => readOne(id.value),
  {
    watch: [id]
  }
)

const serverAddress = computed(() => {
  return `${window.location.hostname}:${server.value?.ports[0]}`  // TODO fix
})

const terminalEl = ref<HTMLElement | null>(null)

let termFitAddon: FitAddon | null = null

let wsAttach: WebSocket | null = null

const api = useApi()

function termResize() {
    if (termFitAddon) {
        termFitAddon.fit()
    }
}

onMounted(() => {
    if (!terminalEl.value) {
        return
    }

    wsAttach = new WebSocket(api.ws(`/servers/${id.value}/attach`))
    
    const term = new Terminal();
    
    termFitAddon = new FitAddon()
    term.loadAddon(termFitAddon)
    
    term.open(terminalEl.value)
    termResize()
    resizeObserver = new ResizeObserver(() => {
        termResize()
    })
    // resizeObserver.observe(terminalEl.value)

    const termAttachAddon = new AttachAddon(wsAttach)
    term.loadAddon(termAttachAddon)
})

onBeforeUnmount(() => {
    if (resizeObserver) {
        resizeObserver.disconnect()
    }

    if (wsAttach) {
        wsAttach.close()
    }
})

async function deleteServer() {
  await deleteOne(id.value)
  navigateTo('/')
}

async function startServer() {
  await start(id.value)
  await refresh()
}

async function stopServer() {
  await stop(id.value)
  await refresh()
}

function backupURL(): string {
  return api.http(`/servers/${id.value}/backup`)
}

const restoreFileInput = ref<HTMLInputElement | null>(null)
const restoreFile = ref<File | null>(null)

function restoreServer() {
  restoreFileInput.value?.click()
}

async function restoreServerFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  restoreFile.value = input.files?.[0] ?? null

  if (!restoreFile.value) {
    return
  }

  const form = new FormData()
  form.append('file', restoreFile.value)

  await restore(id.value, form)
  restoreFile.value = null
}

const {
  cpuPoints,
  memPoints,
  netPoints,
  status: metricsStatus,
  error: metricsError,
  reconnect
} = useSystemMetrics(api.ws(`/servers/${id.value}/stats`))

const netRxPoints = computed(() =>
  netPoints.value.map((p) => ({
    t: p.t,
    value: p.rx
  }))
)

const netTxPoints = computed(() =>
  netPoints.value.map((p) => ({
    t: p.t,
    value: p.tx
  }))
)

const netMax = computed(() => {
  let max = 0

  for (const p of netPoints.value) {
    if (p.rx > max) max = p.rx
    if (p.tx > max) max = p.tx
  }

  if (max <= 0) return 1

  // add headroom
  const scaled = max * 1.2

  // round to nice numbers (important for readability)
  const magnitude = Math.pow(10, Math.floor(Math.log10(scaled)))
  return Math.ceil(scaled / magnitude) * magnitude
})
</script>

<template>
  <UContainer>
    <div class="mt-5">
      <h1 class="text-2xl font-bold">{{ server.name }} <ServerStatusBadge :status="server.status" /></h1>
    </div>

    <div class="mt-4">
      <div class="flex h-full min-h-0 flex-col gap-3">
          <div class="flex gap-2">
            <UButton variant="outline" color="neutral" :to="backupURL()">Backup</UButton> <!-- TODO hover icon -->
            <input ref="restoreFileInput" type="file" class="hidden" accept=".tar.gz" @change="restoreServerFileChange" />
            <UButton variant="outline" color="neutral" v-on:click="restoreServer">Restore</UButton>
            <UButton variant="outline" color="neutral" v-on:click="startServer">Start</UButton>
            <UButton variant="outline" color="error" v-on:click="stopServer">Stop</UButton>
            <UButton variant="outline" color="error" v-on:click="deleteServer">Delete</UButton>
            
            <UInput v-model="serverAddress" :ui="{ trailing: 'pr-0.5' }" readonly="true">
              <template v-if="serverAddress?.length" #trailing>
                <UTooltip text="Copy to clipboard" :content="{ side: 'right' }">
                  <UButton color="neutral" variant="link" :icon="copied ? 'i-lucide-copy-check' : 'i-lucide-copy'" @click="copy(serverAddress)" />
                </UTooltip>
              </template>
            </UInput>
          </div>
  
          <UCard
          class="flex h-[60vh] flex-col overflow-hidden"
          :ui="{
              body: 'flex-1 min-h-0 p-0 sm:p-0'
          }"
          >
          <template #header>
              <div class="font-semibold">
              Terminal
              </div>
          </template>
  
          <div ref="terminalEl" class="h-full w-full min-h-0"></div>
          </UCard>
  
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div class="text-lg font-semibold">
                System Metrics
              </div>
  
              <div class="flex items-center gap-3">
                <UBadge
                  :color="metricsStatus === 'open' ? 'success' : metricsStatus === 'error' ? 'error' : 'neutral'"
                  variant="soft"
                >
                  {{ metricsStatus }}
                </UBadge>
  
                <UButton
                  size="sm"
                  variant="soft"
                  v-on:click="reconnect"
                >
                  Reconnect
                </UButton>
              </div>
            </div>
  
            <UAlert
              v-if="metricsError"
              color="error"
              variant="soft"
              :description="metricsError"
            />
  
            <div class="grid gap-4 lg:grid-cols-3">
              <ServerStatChart
                title="CPU Usage"
                :series="[
                  {
                    key: 'cpu',
                    label: 'CPU',
                    color: '#22c55e',
                    points: cpuPoints
                  }
                ]"
                :min="0"
                :max="100"
                :fill-area="true"
              />
  
              <ServerStatChart
                title="Memory Usage"
                :series="[
                  {
                    key: 'mem',
                    label: 'Memory',
                    color: '#3b82f6',
                    points: memPoints
                  }
                ]"
                :min="0"
                :max="100"
                :fill-area="true"
              />
  
              <ServerStatChart
                title="Network"
                :series="[
                  {
                    key: 'rx',
                    label: 'RX',
                    color: '#22c55e',
                    points: netRxPoints
                  },
                  {
                    key: 'tx',
                    label: 'TX',
                    color: '#a855f7',
                    points: netTxPoints
                  }
                ]"
                :min="0"
                :max="netMax"
              />
            </div>
          </div>
        </div>
    </div>
  </UContainer>

</template>
