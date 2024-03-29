class Gist {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }

  index(offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B').replace(' ', '%2B');
    }
    return this.api.axios.get(`/gists?offset=${offset}`);
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/gists');
  }

  delete(id) {
    return this.api.axios.delete(`/admin/gists/${id}`);
  }
}

export default Gist;
