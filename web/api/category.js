import axios from 'axios';

function Category(api) {
  this.api = api;
}

Category.prototype = {
  adminIndex: function(callback) {
    this.api.request('get', '/admin/categories', {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp.data);
      }
    });
  }
}

export default Category;
