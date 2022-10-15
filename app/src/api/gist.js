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

  genres(name, offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B').replace(' ', '%2B');
    }
    return this.api.axios.get(`/genres/${name}?offset=${offset}`);
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/gists');
  }

  genres(name, offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B').replace(' ', '%2B');
    }
    return this.api.axios.get(`/genres/${name}?offset=${offset}`);
  }

  delete(id) {
    return this.api.axios.delete(`/admin/gists/${id}`);
  }
}

export default Gist;
