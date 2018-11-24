import axios from 'axios';
import Home from './home.js';
import User from './user.js';
import Category from './category.js';

axios.defaults.baseURL = 'https://api.suntin.com';
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
    const self = this;
    axios({
      method: method,
      url: url,
      headers: {'Authorization': 'Bearer ' + self.user.token(method, url, data)},
      data: data
    })
      .then((resp) => {
        console.log(resp.data);
        if (resp.data.error) {
          return Promise.reject(resp.data.error);
        }
        callback(resp.data);
      })
      .catch((error) => {
        if (errCallback === 'function') {
          errCallback(error)
          return
        }
        if (error.code === 401) {
          self.user.clear();
          window.location.href = '/sign_in';
          return
        }
        if (error.response) {
          let data = error.response.data.error;
          // TODO should handle 500 server error
          if (data.code === 500) {
            window.location.href = '/';
            return
          }
        }
        // TODO
        console.info(error);
      })
  }
}

export default API;
