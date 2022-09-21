class Gist {
  constructor(api) {
    this.api = api;
  }

  index(offset) {
    if (!!offset) {
      offset = offset.replace('+', '%2B').replace(' ', '%2B');
    }
    return this.api.axios.get(`/gists?offset=${offset}`);
  }
}

export default Gist;
