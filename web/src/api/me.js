class Me {
  constructor(api) {
    this.api = api;
  }

  groups(limit) {
    return this.api.axios.get(`/user/groups?limit=${limit}`).then((resp) => {
      return resp.data;
    });
  }
}

export default Me;
