import axios from 'axios';
import Home from './home.js';

axios.defaults.baseURL = 'https://localhost:4000';
if (process.env.NODE_ENV === 'development') {
  axios.defaults.baseURL = 'http://localhost:4000';
}
axios.defaults.headers.post['Content-Type'] = 'application/json';

function API() {
  this.home = new Home(this);
}

export default API;
