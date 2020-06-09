import Polyglot from 'node-polyglot';

function Locale(lang) {
  var locale = 'en-US';
  this.polyglot = new Polyglot({locale: locale});
  this.polyglot.extend(require(`./en-US.json`));
  if (languages[lang]) {
    this.polyglot.extend(require(`./${lang}.json`));
  }
}

Locale.prototype = {
  t: function(key, options) {
    return this.polyglot.t(key, options);
  }
};

const languages = {
  'zh': true,
  'ru': true
};

export default Locale;
