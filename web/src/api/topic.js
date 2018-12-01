import axios from 'axios';

function Topic(api) {
  this.api = api;
}

Topic.prototype = {
  index: function(callback) {
    this.api.request('get', '/topics', {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  create: function(params, callback) {
    this.api.request('post', '/topics', data, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  update: function(id, params, callback) {
    const data = {title: params.title, body: params.body, category_id: params.category_id};
    this.api.request('post', '/topics', data, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  update: function(id, params, callback) {
    const data = {title: params['title'], body: params['body'], category_id: params['category_id']};
    this.api.request('post', `/topics/${id}`, params, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  show: function(id, callback) {
    this.api.request('get', `/topics/${id}`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },
}

export default Topic;
