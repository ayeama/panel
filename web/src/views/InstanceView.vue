<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useInstance } from '@/composables/useInstance'
import InstanceStatistics from '@/components/InstanceStatistics.vue'
import InstanceTerminal from '@/components/InstanceTerminal.vue'

const route = useRoute()
const router = useRouter()

const id = route.params.id

const { instance, instanceRead, instanceDelete, instanceStart, instanceStop } = useInstance()

onMounted(() => {
  instanceRead(id)
})

async function instanceDeleteRedirect(id) {
  await instanceDelete(id)
  router.push("/")
}

function running(status) {
  switch (status) {
    case 'running':
      return true
    default:
      return false
  }
}
</script>

<template>
  <div v-if="instance" class="row g-2 py-2">
    <div class="col">
      <div class="row">
        <div class="col d-flex">
          <h1 class="mb-0">{{ instance.name }}</h1>
        </div>

        <div class="col col-auto d-flex align-items-end">
          <div class="btn-group">
            <button v-if="!running(instance.status)" type="button" class="btn btn-secondary" v-on:click="instanceStart(id)">Start</button>
            <button v-else type="button" class="btn btn-secondary" v-on:click="instanceStop(id)">Stop</button>
            
            <button type="button" class="btn btn-secondary dropdown-toggle dropdown-toggle-split" data-bs-toggle="dropdown" aria-expanded="false">
              <span class="visually-hidden">Toggle Dropdown</span>
            </button>
            
            <ul class="dropdown-menu">
              <li><a class="dropdown-item">Backup</a></li>
              <li><a class="dropdown-item">Restore</a></li>
              <li><hr class="dropdown-divider"></li>
              <li><a class="dropdown-item text-danger" v-on:click="instanceDeleteRedirect(id)">Delete</a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="col col-12">
      <InstanceTerminal v-if="instance" :id="instance.id" />
    </div>
  
    <div class="col">
      <div v-if="instance">
        <input id="instanceName" class="form-control-plaintext" type="text" :value="instance.name" readonly>
        <input id="instanceImage" class="form-control-plaintext" type="text" :value="instance.image" readonly>
        <input id="instanceStatus" class="form-control-plaintext" type="text" :value="instance.status" readonly>
        
        <input v-for="port in instance.ports" id="instancePort" class="form-control-plaintext" type="text" :value="port" readonly>
        
        <input id="instanceWebhook" class="form-control-plaintext" type="text" :value="instance.webhook" readonly>
      </div>
    </div>
  
    <div class="col">
      <InstanceStatistics v-if="instance" :id="instance.id" />
    </div>
  </div>
</template>
