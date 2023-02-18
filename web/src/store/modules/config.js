const state = {
    config: {}
}

const mutations = {
    setConfig(state, config) {
        state.config = config;
    }
}

const actions = {
    async setConfig({commit}) {
        const response = await fetch('http://localhost:8080/api/v1/config');
        const config = await response.json();
        commit('setConfig', config);
    },

    async updateConfig({commit}, config) {
        console.log(config);
        const response = await fetch('http://localhost:8080/api/v1/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
        });
        const updatedConfig = await response.json();
        commit('setConfig', updatedConfig);
    }
}

const getters = {
    config: state => state.config
}

export default {
    namespaced: true,
    state,
    mutations,
    actions,
    getters
}
