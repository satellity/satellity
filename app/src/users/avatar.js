import style from './avatar.module.scss';
import React from 'react';
import PropTypes from 'prop-types';

const Avatar = ({user, classes}) => {
  return (
    <img src={user.avatar_url} alt={user.nickname} className={`${style.avatar} ${style[[classes]]}`} />
  );
};

Avatar.propTypes = {
  user: PropTypes.object,
  classes: PropTypes.string,
};

export default Avatar;
