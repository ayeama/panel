<script setup lang="ts">
import type { TableColumn, TableRow } from '@nuxt/ui'
import type { Server } from '~/types/server'

definePageMeta({
  layout: 'default'
})

const ServerImageBadge = resolveComponent('ServerImageBadge') as Component
const ServerStatusBadge = resolveComponent('ServerStatusBadge') as Component

const { read } = useServers()

const filter = ref('')

const data = ref<Server[]>([])
const columns: TableColumn<Server>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
  },
  {
    accessorKey: 'image',
    header: 'Image',
    cell: ({ row }) => {
      return h(ServerImageBadge, {
        image: row.original.image
      })
    }
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      return h(ServerStatusBadge, {
        status: row.original.status
      })
    }
  }
]

function handleSelect(e: Event, row: TableRow<Server>) {
  navigateTo(`/servers/${row.original.name}`)
}

onMounted(async () => {
  data.value = await read()
})
</script>

<template>
  <UContainer>

    <div class="mt-5">
      <UCard :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
          <div class="flex justify-between">
            <UButton label="Create" color="primary" variant="outline" to="/create" />
            <UInput icon="i-lucide-search" size="md" color="secondary" variant="outline" placeholder="Search..." v-model="filter" />
          </div>
        </template>
  
        <UTable v-model:global-filter="filter" :data="data" :columns="columns" @select="handleSelect" />
      </UCard>
    </div>
  </UContainer>
</template>
