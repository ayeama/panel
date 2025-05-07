<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
    url: {
        required: true,
    }
})

let socket = null

const status = ref('disconnected')

const input = ref('')
const output = ref('')

const textarea = ref(null)

onMounted(() => {
    status.value = 'connecting'

    socket = new WebSocket(props.url)
    socket.onopen = () => {
        status.vaule = 'connected'
    }
    socket.onclose = () => {
        status.value = 'disconnected'
    }
    socket.onerror = () => {
        status.value = 'error'
    }
    socket.onmessage = (event) => {
        output.value += event.data
        textarea.value.scrollTop = textarea.value.scrollHeight - 1
    }
})

onBeforeUnmount(() => {
    if (socket) {
        socket.close()
    }
})

async function send() {
    if (socket) {
        socket.send(input.value + "\n")
        input.value = ''
    }
}
</script>

<template>
    <div>
        <div class="card">
            <div class="card-header">
                <h1 class="fs-5 position-relative">Terminal <span class="position-absolute top-0 translate-middle badge text-bg-secondary">{{ status }}</span></h1>
                <!-- <p>Terminal <span class="badge text-bg-success">{{ status }}</span></p> -->
            </div>

            <div class="card-body">
                <form @submit.prevent="send">
                    <div class="row">
                        <div class="col">
                            <textarea ref="textarea" v-model="output" class="form-control" rows="18" style="resize: none; font-family: monospace;" readonly></textarea>
                        </div>
                    </div>

                    <div class="row mt-2">
                        <div class="col">
                            <div class="input-group">
                                <input v-model="input" class="form-control" style="font-family: monospace;">
                                <button type="submit" class="btn btn-secondary">Send</button>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>

    </div>
</template>
