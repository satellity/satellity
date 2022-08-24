const short = require('short-uuid');

const translator = short(short.constants.flickrBase58, {
  consistentLength: false,
});

export const seoTitle = (title, id) => {
  return title.replace(/\W+/mgsi, ' ').trim().replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '') + '-' + translator.fromUUID(id);
};

export const titleToId = (title) => {
  const array = title.split('-');
  if (array.length < 1) {
    return '';
  }
  return translator.toUUID(array[array.length -1]);
};
