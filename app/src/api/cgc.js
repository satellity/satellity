import axios from 'axios';

class CoinGecko {
  constructor() {
    this.instance = axios.create({
      baseURL: 'https://api.coingecko.com',
      timeout: 5000,
      headers: {'Content-Type': 'application/json'},
    });
  }

  coins() {
    return this.instance.get('/api/v3/simple/price?ids=bitcoin,ethereum,binancecoin,dogecoin&vs_currencies=usd&include_24hr_change=true');
  }
}

export default CoinGecko;
