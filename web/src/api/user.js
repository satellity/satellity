import KJUR from 'jsrsasign';
import uuid from 'uuid/v4';
import Cookies from 'js-cookie';
import Base64 from '../components/base64.js';

class User {
  constructor(api) {
    this.api = api;
    this.base64 = new Base64();
    this.admin = new Admin(api);
    this.fixed_schema_header = '3059301306072a8648ce3d020106082a8648ce3d030107034200';
  }

  get ecdsa() {
    let priv = window.localStorage.getItem('token');
    let pwd = Cookies.get('sid');
    if (!priv || !pwd) {
      return "";
    }
    let ec = KJUR.KEYUTIL.getKey(priv, pwd);
    return KJUR.KEYUTIL.getPEM(ec, 'PKCS1PRV');
  }

  signIn(code, email, password) {
    let pwd = uuid().toLowerCase();
    let ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    let pub = ec.generateKeyPairHex().ecpubhex;
    let priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    let params = {'session_secret': this.fixed_schema_header + pub, 'code': code, email: email, password: password};
    let request;
    if (code != '') {
      request = this.api.axios.post('/oauth/github', params);
    } else {
      request = this.api.axios.post('/sessions', params);
    }
    return request.then((resp) => {
      const data = resp.data;
      Cookies.set('sid', pwd, { expires: 365 });
      window.localStorage.setItem('token', priv);
      window.localStorage.setItem('uid', data.user_id);
      window.localStorage.setItem('sid', data.session_id);
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(data)));
      return data;
    });
  }

  verify(params) {
    let pwd = uuid().toLowerCase();
    let ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    let pub = ec.generateKeyPairHex().ecpubhex;
    let priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    params = {'verification_id': params.verification_id, 'code': params.code, 'username': params.username, 'password': params.password, 'session_secret': this.fixed_schema_header + pub};
    return this.api.axios.post(`/email_verifications/${params.verification_id}`, params).then((resp) => {
      const data = resp.data;
      Cookies.set('sid', pwd, { expires: 365 });
      window.localStorage.setItem('token', priv);
      window.localStorage.setItem('uid', data.user_id);
      window.localStorage.setItem('sid', data.session_id);
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(data)));
      return resp.data;
    });
  }

  update(params) {
    return this.api.axios.post('/me', params).then((resp) => {
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(resp.data)));
      return resp.data;
    });
  }

  show(id) {
    return this.api.axios.get(`/users/${id}`).then((resp) => {
      return resp.data;
    });
  }

  remote() {
    return this.api.axios.get('/me').then((resp) => {
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(resp.data)));
      return resp.data;
    })
  }

  local() {
    const user = window.localStorage.getItem('user');
    if (!user) {
      return {};
    }
    return JSON.parse(this.base64.decode(user));
  }

  loggedIn() {
    const user = this.local();
    return user.user_id !== undefined || user.username !== undefined || user.nickname !== undefined;
  }

  isAdmin() {
    return this.local().role === 'admin';
  }

  topics(id) {
    return this.api.axios.get(`/users/${id}/topics`).then((resp) => {
      return resp.data;
    })
  }

  clear() {
    window.localStorage.clear();
  }
}

class Admin {
  constructor(api) {
    this.api = api;
  }

  index() {
    return this.api.axios.get('/admin/users').then((resp) => {
      return resp.data;
    })
  }
}

export default User;
