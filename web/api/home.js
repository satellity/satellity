import axios from 'axios';

function Home(api) {
  this.api = api;
}

Home.prototype = {
  index: function(callback) {
    axios.get('/_hc').then(function(resp) {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  }
}

export default Home;
