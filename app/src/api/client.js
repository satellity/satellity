class Client {
  constructor(api) {
    this.api = api;
  }

  configs() {
    return this.api.axios.get(`/client`);
  }
}

export default Client;
