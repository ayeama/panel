<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

import { useInstance } from '@/composables/useInstance'
import InstanceStatistics from '@/components/InstanceStatistics.vue'
import InstanceTerminal from '@/components/InstanceTerminal.vue'

const route = useRoute()
const id = route.params.id

const { instance, instanceRead } = useInstance()

onMounted(() => {
  instanceRead(id)
})
</script>

<template>
  <div class="row">
    <div class="col col-12">
      <InstanceTerminal v-if="instance" v-bind:id="instance.id" />
    </div>
  
    <div class="col">
      <div v-if="instance">
        <input id="instanceName" class="form-control-plaintext" type="text" :value="instance.name" readonly>
        <input id="instanceImage" class="form-control-plaintext" type="text" :value="instance.image" readonly>
        <input id="instanceStatus" class="form-control-plaintext" type="text" :value="instance.status" readonly>
      </div>
    </div>
  
    <div class="col">
      <InstanceStatistics v-if="instance" v-bind:id="instance.id" />
    </div>
  </div>

  <div class="row">
    <div class="col">
    </div>
  </div>
</template>
