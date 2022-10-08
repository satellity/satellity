class Source {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/sources');
  }

  index() {
    return this.api.axios.get('/admin/sources');
  }

  delete(id) {
    return this.api.axios.delete(`/admin/sources/${id}`);
  }
}

export default Source;
