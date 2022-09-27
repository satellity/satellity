import React from 'react';
import Config from 'components/config.js';

import style from './index.module.scss';

const Index = () => {
  return (
    <h1 className={style.welcome}>
      This is the Dashboard for {Config.Name}.
    </h1>
  );
};

export default Index;
