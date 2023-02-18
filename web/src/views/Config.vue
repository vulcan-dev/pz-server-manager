<template>
  <div class="bg-gray-100">
    <div class="container">
    <div class="flex justify-between">
        <h1 class="text-2xl font-bold">Config</h1>
        <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" @click="refresh">Refresh</button>
        <button :class="{'bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded': changed, 'bg-gray-500 text-white font-bold py-2 px-4 rounded': !changed}" @click="save" :disabled="!changed">Save</button>
      </div>

      <!--   The config is modifiable   -->
      <table class="table-auto">
        <thead>
          <tr>
            <th class="px-4 py-2">Key</th>
            <th class="px-4 py-2">Type</th>
            <th class="px-4 py-2">Value</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(value, key) in config" :key="key">
            <td class="border px-4 py-2">{{ key }}</td>
            <td class="border px-4 py-2">{{ typeof value }}</td>
            <td class="border px-4 py-2">
              <input
                v-if="typeof value === 'string'"
                class="border rounded w-full"
                type="text"
                v-model="config[key]"
              />
              <input
                v-if="typeof value === 'number'"
                class="border rounded w-full"
                type="number"
                v-model.number="config[key]"
              />
              <input
                v-if="typeof value === 'boolean'"
                class="border rounded w-full"
                type="checkbox"
                v-model="config[key]"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ConfigComponent',
  methods: {
    refresh() {
      this.$store.dispatch('config/setConfig');
      this.config = this.$store.state.config.config;
    },
    save() {
      // make all values strings
      const config = {};
      for (const key in this.config) {
        config[key] = this.config[key].toString();
      }
      this.$store.dispatch('config/updateConfig', config);
    },
  },
  computed: {
    config() {
      return this.$store.state.config.config;
    },
  },
  data() {
    return {
      // eslint-disable-next-line vue/no-dupe-keys
      config: {},
      changed: false,
    };
  },
  watch: {
    config: {
      handler() {
        for (const key in this.config) {
          if (this.config[key] !== this.$store.state.config.config[key]) {
            this.changed = true;
            return;
          }
        }
        this.changed = false;
      },
      deep: true,
    },
  },
}
</script>

<style scoped>
.container {
  margin: 0 auto;
  max-width: 800px;
}

table {
  border-collapse: collapse;
  width: 100%;
}
</style>
