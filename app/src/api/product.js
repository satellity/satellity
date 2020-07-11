class Product {
  constructor(api) {
    this.api = api;
    this.admin = new Admin(api);
  }

  index() {
    return this.api.get('/products');
  }

  show(id) {
    return this.api.get(`/products/${id}`);
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  create(params) {
    const tags = params.tags.map((e) => {
      if (typeof e !== 'string') return e;
      e = e.trim();
      return e.charAt(0).toUpperCase() + e.slice(1);
    });
    const data = {name: params.name, body: params.body, cover: params.cover, source: params.source, tags: tags};
    return this.api.post('/admin/products', data);
  }

  update(params) {
    const tags = params.tags.map((e) => {
      if (typeof e !== 'string') return e;
      e = e.trim();
      return e.charAt(0).toUpperCase() + e.slice(1);
    });
    const data = {name: params.name, body: params.body, cover: params.cover, source: params.source, tags: tags};
    return this.api.post(`/admin/products/${params.product_id}`, data);
  }
}

export default Product;
