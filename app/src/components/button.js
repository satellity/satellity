import style from './button.module.scss';
import React from 'react';
import {Link} from 'react-router-dom';
import Loading from './loading.js';
import PropTypes from 'prop-types';

const Button = ({classes, type, text, original, action}) => {
  classes = classes.split(' ').map((name) => {
    return style[[name]];
  }).join(' ');

  if (original && type === 'link') {
    return (
      <a href={action} className={classes}>{text}</a>
    );
  }

  if (type === 'link') {
    return (
      <Link to={action} className={classes}>{text}</Link>
    );
  }

  if (type === 'button') {
    return (
      <button type={type} className={classes} disabled={disabled} onClick={click}>
        {disabled && <Loading class='small white' />}
        &nbsp;{text}
      </button>
    );
  }

  return (
    <button type={type} className={classes} disabled={disabled}>
      {disabled && <Loading class='small white' />}
      &nbsp;{text}
    </button>
  );
};

Button.propTypes = {
  classes: PropTypes.string,
  type: PropTypes.string,
  text: PropTypes.string,
  original: PropTypes.string,
  action: PropTypes.string,
};

export default Button;
