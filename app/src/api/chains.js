import axios from 'axios';

class Chains {
  constructor() {
    this.instance = axios.create({
      baseURL: 'https://chainid.network',
      timeout: 5000,
      headers: {'Content-Type': 'application/json'},
    });
  }

  list() {
    return this.instance.get('/chains.json');
  }
}

export default Chains;
