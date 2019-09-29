import {TextEncoderLite,TextDecoderLite} from 'text-encoder-lite';
import {toByteArray,fromByteArray} from 'base64-js';

class Base64 {
  constructor() {
  }

  encode(str) {
    var bytes = new (typeof TextEncoder === "undefined" ? TextEncoderLite : TextEncoder)('utf-8').encode(str);
    return fromByteArray(bytes);
  }

  decode(str) {
    var bytes = toByteArray(str);
    return new (typeof TextDecoder === "undefined" ? TextDecoderLite : TextDecoder)('utf-8').decode(bytes);
  }
}

export default Base64;
