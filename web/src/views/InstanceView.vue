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
    <div class="col">
      <InstanceTerminal v-if="instance" v-bind:id="instance.id" />
    </div>
  </div>
  
  <div class="row">
    <div class="col">
      <div class="row">
        <div class="col">
          <div v-if="instance">
            id: {{ instance.id }}
          </div>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <div v-if="instance">
            name: {{ instance.name }}
          </div>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <div v-if="instance">
            image: {{ instance.image }}
          </div>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <div v-if="instance">
            status: {{ instance.status }}
          </div>
        </div>
      </div>
    </div>

    <div class="col">
      <InstanceStatistics v-if="instance" v-bind:id="instance.id" />
    </div>
  </div>
</template>
