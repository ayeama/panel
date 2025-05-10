<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HOST } from '@/config'
import ArrowClockwise from '@/components/icons/ArrowClockwise.vue'
import ServerStatusBadge from '@/components/ServerStatusBadge.vue'

const router = useRouter()

onMounted(async () => {
  getServers()
})

// pagination
const paginationLimits = [5, 10, 20, 50, 100] // TODO handle if there are many pages...
const paginationLimit = ref(paginationLimits[1])

async function paginationLimitSelect(value) {
  paginationLimit.value = value
  await getServers()
}

const loading = ref(true)
const serversPaginated = ref({ limit: 10, offset: 0, total: 1, items: [] }) // TODO remvoe defaults?

async function getServers(page = 1) {
  var limit = paginationLimit.value
  var offset = (page - 1) * paginationLimit.value

  try {
    const response = await fetch(`https://${HOST}/servers?limit=${limit}&offset=${offset}`)
    const data = await response.json()
    serversPaginated.value = data
  } catch (err) {
    console.log('Failed to fetch servers', err)
  } finally {
    loading.value = false
  }
}

function serverView(id) {
  router.push(`/servers/${id}`)
}
</script>

<template>
  <div>
    <div class="row">
      <div class="col">
        <h1>Servers</h1>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <RouterLink to="/servers/create" class="btn btn-primary">Create</RouterLink>
        </div>
      </div>
    </div>

    <div>
      <div class="card">
        <div class="card-header py-3">
          <a
            class="icon-link"
            href=""
            v-on:click.prevent="getServers(serversPaginated.offset / serversPaginated.limit + 1)"
            >Refresh
            <ArrowClockwise />
          </a>
        </div>

        <div class="card-body p-0">
          <table class="table table-hover m-0">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th>Status</th>
              </tr>
            </thead>

            <tbody>
              <tr v-if="loading">
                <td scope="role">
                  <div class="placeholder-glow"><span class="placeholder col-8"></span></div>
                </td>
                <td>
                  <div class="placeholder-glow"><span class="placeholder col-4"></span></div>
                </td>
              </tr>

              <tr
                v-else
                v-for="server in serversPaginated.items"
                v-bind:key="server.id"
                v-on:click="serverView(server.id)"
              >
                <td scope="row">{{ server.name }}</td>
                <td>
                  <ServerStatusBadge v-bind:status="server.status" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="card-footer pt-3 pb-0">
          <div class="row">
            <div class="col my-auto">
              <p>
                Showing {{ serversPaginated.offset + 1 }} to
                {{
                  serversPaginated.offset + serversPaginated.limit > serversPaginated.total
                    ? serversPaginated.total
                    : serversPaginated.offset + serversPaginated.limit
                }}
                of {{ serversPaginated.total }} results.
              </p>
            </div>

            <div class="col d-flex justify-content-center">
              <nav aria-label="Page navigation">
                <ul class="pagination">
                  <li class="page-item">
                    <a
                      class="page-link"
                      href="#"
                      aria-label="Previous"
                      v-bind:class="{ disabled: !(serversPaginated.offset > 0) }"
                      v-on:click.prevent="
                        getServers(serversPaginated.offset / serversPaginated.limit + 1 - 1)
                      "
                    >
                      <span aria-hidden="true">&laquo;</span>
                    </a>
                  </li>

                  <!-- <li class="page-item"><a class="page-link active" href="#">1</a></li>
                  <li class="page-item"><a class="page-link" href="#">2</a></li>
                  <li class="page-item"><a class="page-link" href="#">3</a></li> -->
                  <li
                    class="page-item"
                    v-for="(item, index) in Math.ceil(
                      serversPaginated.total / serversPaginated.limit,
                    )"
                    v-bind:key="item"
                  >
                    <a
                      class="page-link"
                      href="#"
                      v-bind:class="{
                        active: index + 1 == serversPaginated.offset / serversPaginated.limit + 1,
                      }"
                      v-on:click.prevent="getServers(index + 1)"
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
                          serversPaginated.offset + serversPaginated.limit <
                          serversPaginated.total
                        ),
                      }"
                      v-on:click.prevent="
                        getServers(serversPaginated.offset / serversPaginated.limit + 1 + 1)
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
