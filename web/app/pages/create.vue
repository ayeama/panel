<script setup lang="ts">
import type { Image } from '~/types/image'
import type { ServerCreateRequest } from '~/types/server'

useHead({
  title: 'Panel - Create'
})

definePageMeta({
  layout: 'default'
})

const { create } = useServers()
const { read } = useImages()

const data = ref<Image[]>()
const selected = ref<Image>()

async function createServer() {
    if (!selected.value) {
      return
    }

    const request: ServerCreateRequest = {
      image: selected.value.reference,
      env: formState.env
    }
    const server = await create(request)

    navigateTo(`/servers/${server.name}`)
}

const serverEnv = computed(() => {
  return Object.keys(selected.value?.env ?? {})
})

const formState = reactive({
  env: {} as Record<string, string>
})

watch(selected, (image) => {
  Object.keys(formState.env).forEach((key) => {
    delete formState.env[key]
  })

  Object.assign(formState.env, image?.env ?? {})
})

function envLabel(label: string) {
  return label.charAt(0).toUpperCase() + label.slice(1).toLowerCase().replaceAll('_', ' ')
}

onMounted(async () => {
  data.value = (await read()).sort((a, b) => a.reference.localeCompare(b.reference))
})
</script>

<template>
  <UContainer>
    <div class="mt-5">
      <h1 class="text-2xl font-bold">Create server</h1>
    </div>

    <div>
      <UForm :state="formState">
        <div class="mt-4">
          <UFormField label="Server">
            <USelectMenu v-model="selected" :items="data" class="w-100" label-key="reference" placeholder="Select server" color="secondary" />
          </UFormField>
        </div>

        <div class="mt-4 grid sm:grid-cols-3 gap-4">
          <UFormField v-for="k in serverEnv" :key="k" :label="envLabel(k)" :name="`env.${k}`">
            <UInput v-model="formState.env[k]" class="w-full" />
          </UFormField>
        </div>

        <div class="mt-4" v-if="selected">
          <UButton variant="outline" color="primary" v-on:click="createServer">Create</UButton>
        </div>
      </UForm>
    </div>
  </UContainer>
</template>
