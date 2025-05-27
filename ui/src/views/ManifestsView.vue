<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HOST } from '@/config'
import ArrowClockwise from '@/components/icons/ArrowClockwise.vue'

const router = useRouter()

onMounted(async () => {
  getManifests()
})

const paginationLimits = [5, 10, 20, 50, 100] // TODO handle if there are many pages...
const paginationLimit = ref(paginationLimits[1])

async function paginationLimitSelect(value) {
  paginationLimit.value = value
  await getManifests()
}

const manifestsPaginated = ref({ limit: 10, offset: 0, total: 1, items: [] }) // TODO remvoe defaults?

async function getManifests(page = 1) {
  var limit = paginationLimit.value
  var offset = (page - 1) * paginationLimit.value

  try {
    const response = await fetch(`https://${HOST}/manifests?limit=${limit}&offset=${offset}`)
    const data = await response.json()
    manifestsPaginated.value = data
  } catch (err) {
    console.log('Failed to fetch manifests', err)
  }
}

function manifestView(id) {
  router.push(`/manifests/${id}`)
}
</script>

<template>
  <div>
    <div class="row">
      <div class="col">
        <h1>Manifests</h1>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <!-- <RouterLink to="/manifests/create" class="btn btn-primary">Create</RouterLink> -->
        </div>
      </div>
    </div>

    <div>
      <div class="card">
        <div class="card-header py-3">
          <a
            class="icon-link"
            href=""
            v-on:click.prevent="
              getManifests(manifestsPaginated.offset / manifestsPaginated.limit + 1)
            "
            >Refresh
            <ArrowClockwise />
          </a>
        </div>

        <div class="card-body p-0">
          <table class="table table-hover m-0">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th>Variant</th>
                <th>Version</th>
              </tr>
            </thead>

            <tbody>
              <tr
                v-for="manifest in manifestsPaginated.items"
                v-bind:key="manifest.id"
                v-on:click="manifestView(manifest.id)"
              >
                <td scope="row">{{ manifest.name }}</td>
                <td>{{ manifest.variant }}</td>
                <td>{{ manifest.version }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="card-footer pt-3 pb-0 border-0">
          <div class="row">
            <div class="col my-auto">
              <p>
                Showing {{ manifestsPaginated.offset + 1 }} to
                {{
                  manifestsPaginated.offset + manifestsPaginated.limit > manifestsPaginated.total
                    ? manifestsPaginated.total
                    : manifestsPaginated.offset + manifestsPaginated.limit
                }}
                of {{ manifestsPaginated.total }} results.</p>
            </div>

            <div class="col d-flex justify-content-center">
              <nav aria-label="Page navigation">
                <ul class="pagination">
                  <li class="page-item">
                    <a
                      class="page-link"
                      href="#"
                      aria-label="Previous"
                      v-bind:class="{ disabled: !(manifestsPaginated.offset > 0) }"
                      v-on:click.prevent="
                        getManifests(manifestsPaginated.offset / manifestsPaginated.limit + 1 - 1)
                      "
                    >
                      <span aria-hidden="true">&laquo;</span>
                    </a>
                  </li>

                  <li
                    class="page-item"
                    v-for="(item, index) in Math.ceil(
                      manifestsPaginated.total / manifestsPaginated.limit,
                    )"
                    v-bind:key="item"
                  >
                    <a
                      class="page-link"
                      href="#"
                      v-bind:class="{
                        active:
                          index + 1 == manifestsPaginated.offset / manifestsPaginated.limit + 1,
                      }"
                      v-on:click.prevent="getManifests(index + 1)"
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
                          manifestsPaginated.offset + manifestsPaginated.limit <
                          manifestsPaginated.total
                        ),
                      }"
                      v-on:click.prevent="
                        getManifests(manifestsPaginated.offset / manifestsPaginated.limit + 1 + 1)
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
