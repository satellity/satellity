class Ratio {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get(`/ratios`);
  }
}

export default Ratio;
