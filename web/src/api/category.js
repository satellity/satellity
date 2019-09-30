import Base64 from '../components/base64.js';

class Category {
  constructor(api) {
    this.api = api;
    this.base64 = new Base64();
    this.admin = new Admin(api);
  }

  index() {
    return this.api.axios.get('/categories').then((resp) => {
      const categories = resp.data.map((o) => {
        return {category_id: o.category_id, name: o.name, alias: o.alias}
      });
      window.localStorage.setItem('categories', this.base64.encode(JSON.stringify(categories)));
      return resp;
    }).catch((error) => {
      return error;
    });
  }

  topics(id, offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B')
    }
    return this.api.axios.get(`/categories/${id}/topics?offset=${offset}`).then((resp) => {
      return resp.data;
    });
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/categories').then((resp) => {
      return resp.data;
    })
  }

  create(params) {
    params['position'] = params['position'] === '' ? 0 : parseInt(params['position']);
    const data = {name: params.name, alias: params.alias, description: params.description, position: params.position}
    return this.api.axios.post('/admin/categories', data).then((resp) => {
      return resp.data;
    });
  }

  update(id, params) {
    params['position'] = params['position'] === '' ? 0 : parseInt(params['position']);
    const data = {name: params.name, alias: params.alias, description: params.description, position: params.position}
    return this.api.axios.post(`/admin/categories/${id}`, data).then((resp) => {
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/admin/categories/${id}`).then((resp) => {
      return resp.data;
    });
  }
}

export default Category;
