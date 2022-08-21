const short = require('short-uuid');

const translator = short(short.constants.flickrBase58, {
  consistentLength: false,
});

export const seoTitle = (title, id) => {
  return title.replace(/\W+/mgsi, ' ').trim().replace(/\s+/mgsi, '-').replace(/[^\w-]/mgsi, '') + '-' + translator.fromUUID(id);
};
