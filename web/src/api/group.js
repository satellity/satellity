class Group {
  constructor(api) {
    this.api = api;
  }

  create(params) {
    const data = {name: params.name, description: params.description};
    return this.api.axios.post('/groups', data).then((resp) => {
      return resp.data;
    });
  }

  update(id, params) {
    const data = {name: params.name, description: params.description};
    return this.api.axios.post(`/groups/${id}`, data).then((resp) => {
      return resp.data;
    });
  }

  index() {
    return this.api.axios.get('/groups').then((resp) => {
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/groups/${id}`).then((resp) => {
      return resp.data;
    });
  }

  join(id) {
    return this.api.axios.post(`/groups/${id}/join`).then((resp) => {
      return resp.data;
    });
  }

  exit(id) {
    return this.api.axios.post(`/groups/${id}/exit`).then((resp) => {
      return resp.data;
    });
  }

  members(id, limit) {
    return this.api.axios.get(`/groups/${id}/participants?limit=${limit}`).then((resp) => {
      return resp.data;
    });
  }
}

export default Group;
