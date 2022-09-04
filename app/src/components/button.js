import style from './button.module.scss';
import React from 'react';
import {Link} from 'react-router-dom';
import Loading from './loading.js';
import PropTypes from 'prop-types';

const Button = (props) => {
  const {classes, type, text, original, action, disabled} = props;
  const klass = classes.split(' ').map((name) => {
    return style[[name]];
  }).join(' ');

  if (original && type === 'link') {
    return (
      <a href={action} className={klass}>{text}</a>
    );
  }

  if (type === 'link') {
    return (
      <Link to={action} className={klass}>{text}</Link>
    );
  }

  if (type === 'button') {
    return (
      <button type={type} className={klass} disabled={disabled} onClick={click}>
        {disabled && <Loading class='small white' />}
        &nbsp;{text}
      </button>
    );
  }

  return (
    <button type={type} className={klass} disabled={disabled}>
      {disabled && <Loading class='small white' />}
      &nbsp;{text}
    </button>
  );
};

Button.propTypes = {
  classes: PropTypes.string,
  disabled: PropTypes.boolean,
  type: PropTypes.string,
  text: PropTypes.string,
  original: PropTypes.string,
  action: PropTypes.string,
};

export default Button;
