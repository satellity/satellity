import {decode as decodeBase64} from '@stablelib/base64';
import {decode as decodeUTF8} from '@stablelib/utf8';

class Me {
  constructor() {
  }

  value() {
    const source = window.localStorage.getItem('user');
    if (!source) {
      window.localStorage.clear();
      return undefined;
    }
    const user = JSON.parse(decodeUTF8(decodeBase64(source)));
    if (!user.user_id || !user.session_id || !user.private) {
      window.localStorage.clear();
      return undefined;
    }
    return user;
  }

  signOut() {
    window.localStorage.clear();
  }
}

export default Me;
