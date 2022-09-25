import React from 'react';
import PropTypes from 'prop-types';

import style from './widget.module.scss';

const Widget = ({children}) => {
  return (
    <div className={style.widget}>
      {children}
      <div className={style.slogan}>
        News, updates, in-depth analysis, about blockchain, Web3, DAO.
      </div>
      <div className={style.copyright}>
        © 2022 - Now
      </div>
    </div>
  );
};

Widget.propTypes = {
  children: PropTypes.any,
};

export default Widget;
