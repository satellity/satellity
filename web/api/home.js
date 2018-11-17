function Home(api) {
  this.api = api;
}

Home.prototype = {
  index: function(callback) {
    this.api('get', '/_hc', {}, function(resp) {
      if (typeof callback === 'function') {
        callback(resp);
      }
    });
  }
}

export default Home;
