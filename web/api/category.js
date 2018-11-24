import axios from 'axios';

function Category(api) {
  this.api = api;
}

Category.prototype = {
  adminIndex: function(callback) {
    this.api.request('get', '/admin/categories', {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  create: function(params, callback) {
    this.api.request('post', '/admin/categories', params, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    })
  },

  show: function(id, callback) {
    this.api.request('get', `/admin/categories/${id}`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    })
  }
}

export default Category;
