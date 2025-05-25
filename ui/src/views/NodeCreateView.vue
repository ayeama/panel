<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { HOST } from '@/config'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'

const router = useRouter()

const name = ref('')
const uri = ref('')

async function createNode() {
  try {
    const response = await fetch(`https://${HOST}/nodes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: name.value,
        uri: uri.value,
      }),
    })
  } catch (err) {
    console.log('Failed create node', err)
  } finally {
    router.push('/nodes')
  }
}
</script>

<template>
  <div>
    <div>
      <RouterLink to="/nodes" class="icon-link">
        <ArrowLeft />
        Nodes
      </RouterLink>
    </div>

    <div class="row">
      <div class="col">
        <h1>Create Node</h1>
      </div>

      <div class="col">
        <div class="float-end">
          <button form="createNode" type="submit" class="btn btn-primary">Create</button>
        </div>
      </div>
    </div>

    <form v-on:submit.prevent="createNode" id="createNode">
      <div class="row">
        <div class="col">
          <div>
            <label for="nameInput" class="form-label">Name</label>
            <input type="text" class="form-control" id="nameInput" v-model="name" />
          </div>
        </div>

        <div class="col">
          <div>
            <label for="uriInput" class="form-label">URI</label>
            <input type="text" class="form-control" id="uriInput" v-model="uri" />
          </div>
        </div>
      </div>
    </form>
  </div>
</template>
