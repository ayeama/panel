<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()

const users = [
  { id: 1, first: 'Mark', last: 'Otto', handle: '@mdo' },
  { id: 2, first: 'Jacob', last: 'Thornton', handle: '@fat' },
  { id: 3, first: 'John', last: 'Doe', handle: '@social' },
]

function userEdit(id) {
  router.push(`/users/${id}/edit`)
}
</script>

<template>
  <div>
    <h1>Users</h1>

    <table class="table table-bordered table-hover">
      <thead>
        <tr>
          <th scope="col">ID</th>
          <th scope="col">First</th>
          <th scope="col">Last</th>
          <th scope="col">Handle</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" v-bind:key="user.id" v-on:click="userEdit(user.id)">
          <td scope="row">{{ user.id }}</td>
          <td>{{ user.first }}</td>
          <td>{{ user.last }}</td>
          <td>{{ user.handle }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style>
/* TODO this affects all tables?? */
.table-hover tbody tr:hover {
  cursor: pointer;
}
</style>
