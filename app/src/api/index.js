import forge from 'node-forge';
import moment from 'moment';
import {v4 as uuidv4} from 'uuid';
import axios from 'axios';
import Noty from 'noty';
import {decode as decodeHex} from '@stablelib/hex';
import {encodeURLSafe} from '@stablelib/base64';
import {sign as signED25519} from '@stablelib/ed25519';
import {encode as encodeUTF8} from '@stablelib/utf8';
import Category from './category.js';
import Topic from './topic.js';
import Comment from './comment.js';
import Client from './client.js';
import User from './user.js';
import Me from './me.js';
import Verification from './verification.js';
import Gist from './gist.js';
import Source from './source.js';
import CoinGecko from './cgc.js';
import Chains from './chains.js';
import Ratio from './ratio.js';

Noty.overrideDefaults({
  type: 'error',
  layout: 'topCenter',
  killer: true,
  theme: 'nest',
  timeout: 1000,
  progressBar: false,
  animation: {
    open: 'noty_effects_open',
    close: 'noty_effects_close',
  },
});

const signToken = (method, uri, body) => {
  if (typeof body !== 'string') {
    body = '';
  }
  const me = new Me().value();
  if (!me) {
    return '';
  }

  const expire = moment.utc().add(30, 'minutes').unix();
  const md = forge.md.sha256.create();
  md.update(method + uri + body);

  const payload = {
    uid: me.user_id,
    sid: me.session_id,
    exp: expire,
    jti: uuidv4(),
    sig: md.digest().toHex(),
  };

  const header = encodeURLSafe(encodeUTF8(JSON.stringify({alg: 'EdDSA', typ: 'JWT'}))).replaceAll('=', '');
  const payloadStr = encodeURLSafe(encodeUTF8(JSON.stringify(payload))).replaceAll('=', '');
  const sig = encodeURLSafe(signED25519(decodeHex(me.private), encodeUTF8(`${header}.${payloadStr}`))).replaceAll('=', '');
  return `${header}.${payloadStr}.${sig}`;
};

axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.interceptors.request.use((config) => {
  config.url = '/api' + config.url;
  const {method, url, data} = config;
  const token = signToken(method, url, data);
  config.headers.common['Authorization'] = `Bearer ${token}`;
  return config;
}, (error) => {
  return Promise.reject(error);
});

axios.interceptors.response.use((response) => {
  if (!!response.status && (response.status >= 200 && response.status < 300)) {
    const data = response.data;
    if (!!data.error) {
      const error = data.error;
      new Noty({
        text: window.i18n.t(`errors.${error.code}`),
      }).show();
      if (error.code === 401) {
        window.localStorage.removeItem('user');
        window.location.href = '/';
      } else if (error.code === 404) {
        window.location.href = '/404';
      }
    }
    return data;
  }
  return response;
}, (error) => {
  let status; let data;
  // TODO: should clear error.request and error
  if (error.response) {
    status = error.response.status;
    data = error.response.data;
  } else if (error.request) {
    status = 500;
    data = 'Initialize request error';
  } else {
    status = 500;
    data = error.message;
  }
  new Noty({
    text: window.i18n.t(`errors.${status}`),
  }).show();
  return {error: {code: status, description: data}};
});

class API {
  constructor() {
    this.axios = axios;
    this.category = new Category(this);
    this.topic = new Topic(this);
    this.comment = new Comment(this);
    this.user = new User(this);
    this.me = new Me();
    this.verification = new Verification(this);
    this.client = new Client(this);
    this.gist = new Gist(this);
    this.source = new Source(this);
    this.cgc = new CoinGecko();
    this.chain = new Chains();
    this.ratio = new Ratio(this);
  }
}

export default API;
