import React from 'react';
import TimeAgo from 'react-timeago';
import PropTypes from 'prop-types';

import style from './index.module.scss';

const Item = ({gist}) => {
  let author = gist.author;
  if (gist.author.trim().toLowerCase() !== gist.source.author.trim().toLowerCase()) {
    author = `${gist.author}, ${gist.source.author}`;
  }
  return (
    <div className={style.gist}>
      <a href={gist.link}>{gist.title}</a>
      <div className={style.meta}>
        {author} · {gist.source.host} · <TimeAgo date={gist.publish_at} />
      </div>
    </div>
  );
};

Item.propTypes = {
  gist: PropTypes.any,
};

export default Item;
