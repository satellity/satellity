import axios from 'axios';

function Comment(api) {
  this.api = api;
}

Comment.prototype = {
  index: function(id, callback) {
    this.api.request('get', `topics/${id}/comments`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  create: function(params, callback) {
    this.api.request('post', '/comments', params, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  update: function(id, params, callback) {
    this.api.request('post', `/comments/${id}`, params, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  show: function(id, callback) {
    this.api.request('get', `/comments/${id}`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  },

  delete: function(id, callback) {
    this.api.request('post', `/comments/${id}/delete`, {}, (resp) => {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  }
}

export default Comment;
