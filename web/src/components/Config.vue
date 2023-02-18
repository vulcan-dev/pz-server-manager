<template>
  <div class="container">
    <h1>Config</h1>
    <button @click="refresh">Refresh</button>
    <pre>{{ config }}</pre>
  </div>
</template>

<script>
export default {
  name: 'ConfigComponent',
  methods: {
    // refresh looks at ":8080/api/v1/config" and updates the config object
    refresh() {
      fetch('http://localhost:8080/api/v1/config')
          .then(response => response.json())
          .then(data => {
            this.$data.config = data;
            console.log(data);
          });
    }
  },
  data() {
    return {
      config: {},
    };
  },
}
</script>

