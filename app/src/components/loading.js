import style from './loading.module.scss';
import React from 'react';
import PropTypes from 'prop-types';

const Loading = ({classes}) => {
  classes = classes || 'medium';
  classes = classes.split(' ').map((name) => {
    return style[[name]];
  }).join(' ');
  return (
    <div className={`${style.ring} ${classes}`}>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
    </div>
  );
};

Loading.propTypes = {
  classes: PropTypes.string,
};

export default Loading;
