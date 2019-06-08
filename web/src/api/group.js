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

  show(id) {
    return this.api.axios.get(`/groups/${id}`).then((resp) => {
      return resp.data;
    });
  }
}

export default Group;
