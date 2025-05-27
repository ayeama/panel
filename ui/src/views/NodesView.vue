<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HOST } from '@/config'
import ArrowClockwise from '@/components/icons/ArrowClockwise.vue'

const router = useRouter()

onMounted(async () => {
  getNodes()
})

const paginationLimits = [5, 10, 20, 50, 100] // TODO handle if there are many pages...
const paginationLimit = ref(paginationLimits[1])

async function paginationLimitSelect(value) {
  paginationLimit.value = value
  await getNodes()
}

const nodesPaginated = ref({ limit: 10, offset: 0, total: 1, items: [] }) // TODO remvoe defaults?

async function getNodes(page = 1) {
  var limit = paginationLimit.value
  var offset = (page - 1) * paginationLimit.value

  try {
    const response = await fetch(`https://${HOST}/nodes?limit=${limit}&offset=${offset}`)
    const data = await response.json()
    nodesPaginated.value = data
  } catch (error) {
    console.log('Failed to fetch nodes', error)
  }
}

function nodeView(id) {
  router.push(`/nodes/${id}`)
}
</script>

<template>
  <div>
    <div class="row">
      <div class="col">
        <h1>Nodes</h1>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <RouterLink to="/nodes/create" class="btn btn-primary">Create</RouterLink>
        </div>
      </div>
    </div>

    <div>
      <div class="card">
        <div class="card-header py-3">
          <a
            class="icon-link"
            href=""
            v-on:click.prevent="getNodes(nodesPaginated.offset / nodesPaginated.limit + 1)"
            >Refresh
            <ArrowClockwise />
          </a>
        </div>

        <div class="card-body p-0">
          <table class="table table-hover m-0">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th>URI</th>
              </tr>
            </thead>

            <tbody>
              <tr
                v-for="node in nodesPaginated.items"
                v-bind:key="node.id"
                v-on:click="nodeView(node.id)"
              >
                <td scope="row">{{ node.name }}</td>
                <td>{{ node.uri }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="card-footer pt-3 pb-0 border-0">
          <div class="row">
            <div class="col my-auto">
              <p>
                Showing {{ nodesPaginated.offset + 1 }} to
                {{
                  nodesPaginated.offset + nodesPaginated.limit > nodesPaginated.total
                    ? nodesPaginated.total
                    : nodesPaginated.offset + nodesPaginated.limit
                }}
                of {{ nodesPaginated.total }} results.</p>
            </div>

            <div class="col d-flex justify-content-center">
              <nav aria-label="Page navigation">
                <ul class="pagination">
                  <li class="page-item">
                    <a
                      class="page-link"
                      href="#"
                      aria-label="Previous"
                      v-bind:class="{ disabled: !(nodesPaginated.offset > 0) }"
                      v-on:click.prevent="
                        getNodes(nodesPaginated.offset / nodesPaginated.limit + 1 - 1)
                      "
                    >
                      <span aria-hidden="true">&laquo;</span>
                    </a>
                  </li>

                  <li
                    class="page-item"
                    v-for="(item, index) in Math.ceil(nodesPaginated.total / nodesPaginated.limit)"
                    v-bind:key="item"
                  >
                    <a
                      class="page-link"
                      href="#"
                      v-bind:class="{
                        active: index + 1 == nodesPaginated.offset / nodesPaginated.limit + 1,
                      }"
                      v-on:click.prevent="getNodes(index + 1)"
                      >{{ index + 1 }}</a
                    >
                  </li>

                  <li class="page-item">
                    <a
                      class="page-link"
                      href="#"
                      aria-label="Next"
                      v-bind:class="{
                        disabled: !(
                          nodesPaginated.offset + nodesPaginated.limit <
                          nodesPaginated.total
                        ),
                      }"
                      v-on:click.prevent="
                        getNodes(nodesPaginated.offset / nodesPaginated.limit + 1 + 1)
                      "
                    >
                      <span aria-hidden="true">&raquo;</span>
                    </a>
                  </li>
                </ul>
              </nav>
            </div>

            <div class="col">
              <div class="float-end">
                <div class="input-group">
                  <span class="input-group-text">Per page</span>

                  <button
                    class="btn btn-outline-secondary dropdown-toggle"
                    type="button"
                    data-bs-toggle="dropdown"
                    aria-expanded="false"
                  >
                    {{ paginationLimit }}
                  </button>
                  <ul class="dropdown-menu">
                    <li v-for="pl in paginationLimits" :key="pl">
                      <a class="dropdown-item" v-on:click="paginationLimitSelect(pl)">{{ pl }}</a>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
