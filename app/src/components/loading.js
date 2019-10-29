import style from './loading.module.scss';
import React from 'react';

const Loading = (props) => {
  let classes = props.class || 'medium';
  classes = classes.split(' ').map((name) => {
    return style[[name]]
  }).join(' ');
  return (
    <div className={`${style.ring} ${classes}`}>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
    </div>
  )
};

export default Loading;
