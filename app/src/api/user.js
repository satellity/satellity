import KJUR from 'jsrsasign';
import {encode as encodeBase64, decode as decodeBase64} from '@stablelib/base64';
import {encode as encodeUTF8, decode as decodeUTF8} from '@stablelib/utf8';
import {v4 as uuidv4} from 'uuid';
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
    const priv = window.localStorage.getItem('token');
    const pwd = Cookies.get('sid');
    if (!priv || !pwd) {
      return '';
    }
    const ec = KJUR.KEYUTIL.getKey(priv, pwd);
    return KJUR.KEYUTIL.getPEM(ec, 'PKCS1PRV');
  }

  signIn(email, password, provider, code) {
    const pwd = uuidv4().toLowerCase();
    const ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    const pub = ec.generateKeyPairHex().ecpubhex;
    const priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    const data = {session_secret: this.fixed_schema_header + pub, code: code, email: email, password: password};
    let request;
    if (code) {
      request = this.api.axios.post(`/oauth/${provider}`, data);
    } else {
      request = this.api.axios.post('/sessions', data);
    }
    return request.then((resp) => {
      if (resp.error) {
        return resp;
      }
      const data = resp.data;
      Cookies.set('sid', pwd, {expires: 365});
      window.localStorage.setItem('token', priv);
      window.localStorage.setItem('uid', data.user_id);
      window.localStorage.setItem('sid', data.session_id);
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(data)));
      return resp;
    });
  }

  verify(params) {
    if (params.purpose === 'PASSWORD') {
      const data = {purpose: params.purpose, verification_id: params.verification_id, code: params.code, password: params.password};
      return this.api.axios.post(`/email_verifications/${params.verification_id}`, data);
    }
    const pwd = uuidv4().toLowerCase();
    const ec = new KJUR.crypto.ECDSA({'curve': 'secp256r1'});
    const pub = ec.generateKeyPairHex().ecpubhex;
    const priv = KJUR.KEYUTIL.getPEM(ec, 'PKCS8PRV', pwd);
    const data = {
      purpose: params.purpose,
      verification_id: params.verification_id,
      code: params.code,
      username: params.username,
      password: params.password,
      session_secret: this.fixed_schema_header + pub,
    };
    return this.api.axios.post(`/email_verifications/${params.verification_id}`, data).then((resp) => {
      if (resp.error) {
        return resp;
      }
      const data = resp.data;
      Cookies.set('sid', pwd, {expires: 365});
      window.localStorage.setItem('token', priv);
      window.localStorage.setItem('uid', data.user_id);
      window.localStorage.setItem('sid', data.session_id);
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(data)));
      return resp;
    });
  }

  create(publicKey, sessionPublic, sessionPrivate, sig) {
    const data = {public_key: publicKey, session_public: sessionPublic, signature: sig};
    return this.api.axios.post('/sessions', data).then((resp) => {
      if (resp.error) {
        return resp;
      }
      const user = resp.data;
      user.private = sessionPrivate;
      window.localStorage.setItem('user', encodeBase64(encodeUTF8(JSON.stringify(user))));
      return resp;
    });
  }

  update(params) {
    const i = params.avatar_url.indexOf(',');
    const data = {nickname: params.nickname, biography: params.biography, avatar: params.avatar_url.slice(i+1)};
    return this.api.axios.post('/me', data).then((resp) => {
      if (resp.error) {
        return resp;
      }
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(resp.data)));
      return resp;
    });
  }

  show(id) {
    return this.api.axios.get(`/users/${id}`);
  }

  remote() {
    return this.api.axios.get('/me').then((resp) => {
      if (resp.error) {
        return resp;
      }
      window.localStorage.setItem('user', this.base64.encode(JSON.stringify(resp.data)));
      return resp;
    });
  }

  local() {
    const source = window.localStorage.getItem('user');
    if (!source) {
      window.localStorage.clear();
      return {};
    }
    const user = JSON.parse(decodeUTF8(decodeBase64(source)));
    if (!user.private) {
      window.localStorage.clear();
      return {};
    }
    return user;
  }

  loggedIn() {
    const user = this.local();
    return user.user_id !== undefined || user.username !== undefined || user.nickname !== undefined;
  }

  isAdmin() {
    return this.local().role === 'admin';
  }

  topics(id) {
    return this.api.axios.get(`/users/${id}/topics`);
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
    return this.api.axios.get('/admin/users');
  }
}

export default User;
