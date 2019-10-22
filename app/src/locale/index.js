import Polyglot from 'node-polyglot';

function Locale(lang) {
  var locale = 'en-US';
  //if (lang && lang.indexOf('zh') >= 0) {
  //  locale = 'zh-Hans';
  //}
  this.polyglot = new Polyglot({locale: locale});
  this.polyglot.extend(require(`./en-US.json`));
}

Locale.prototype = {
  t: function(key, options) {
    return this.polyglot.t(key, options);
  }
};

export default Locale;
