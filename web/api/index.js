import axios from 'axios';
import Home from './home.js';
import User from './user.js';
import Category from './category.js';

// TODO should replace with funyeah.com
axios.defaults.baseURL = 'https://localhost:4000';
if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:4000';
}
axios.defaults.headers.post['Content-Type'] = 'application/json';

function API() {
  this.home = new Home(this);
  this.user = new User(this);
  this.category = new Category(this);
}

API.prototype = {
  request: function(method, url, data, callback, errCallback) {
    axios({
      method: method,
      url: url,
      data: data
    })
      .then(function(resp) {
        callback(resp);
      })
      .catch(function(error) {
        if (errCallback === 'function') {
          errCallback(error)
          return
        }
        if (error.code === 401) {
          window.history.pushState(null, null, '/');
          return
        }
        // TODO
        console.info(error);
      })
  }
}

export default API;
