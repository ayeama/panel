<script setup>
const props = defineProps({
  title: {
    type: String,
    default: 'Are you sure?',
  },
})

const emit = defineEmits(['cancel', 'confirm'])

function onCancel() {
  emit('cancel')
}

function onConfirm() {
  emit('confirm')
}
</script>

<template>
  <div class="modal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h1 class="modal-title fs-5">{{ title }}</h1>
        </div>

        <div class="modal-body">
          <div class="row">
            <div class="col">
              <p>This action cannot be undone!</p>
            </div>
          </div>

          <div class="row">
            <div class="col">
              <div class="float-end">
                <button
                  type="button"
                  class="btn btn-secondary"
                  data-bs-dismiss="modal"
                  v-on:click="onCancel"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  class="btn btn-danger ms-2"
                  data-bs-dismiss="modal"
                  v-on:click="onConfirm"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
