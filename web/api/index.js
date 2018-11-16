import axios from 'axios';
import Home from './home.js';
import User from './user.js';

// TODO should replace with funyeah.com
axios.defaults.baseURL = 'https://localhost:4000';
if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:4000';
}
axios.defaults.headers.post['Content-Type'] = 'application/json';

function API() {
  this.home = new Home(this);
  this.user = new User(this);
}

export default API;
