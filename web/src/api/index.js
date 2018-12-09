import axios from 'axios';
import Home from './home.js';
import User from './user.js';
import Category from './category.js';
import Topic from './topic.js';
import Comment from './comment.js';

axios.defaults.baseURL = 'https://api.godiscourse.com';
if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:4000';
}
axios.defaults.headers.post['Content-Type'] = 'application/json';

function API() {
  this.home = new Home(this);
  this.user = new User(this);
  this.category = new Category(this);
  this.topic = new Topic(this);
  this.comment = new Comment(this);
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
        // TODO process errors
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
          // TODO
          let githubClientId = '71905afbd6e4541ad62b';
          if (process.env.NODE_ENV === 'development') {
            githubClientId = 'b9b78f343f3a5b0d7c99';
          }
          window.location.href = `https://github.com/login/oauth/authorize?scope=user:email&client_id=${githubClientId}`;
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
