import React from 'react';
import {useParams, Link} from 'react-router-dom';
import PropTypes from 'prop-types';
import Widget from 'components/widget.js';
import Chart from './chart.js';

import style from './index.module.scss';

const View = ({children}) => {
  const {genre} = useParams();
  console.log(genre);

  return (
    <div className='container'>
      <main className='column main'>
        <div className={style.nodes}>
          <Link to="/" className={`${style.node} ${!genre ? style.current : ''}`}> News </Link>
          <Link to="/genres/release" className={`${style.node} ${genre === 'release' ? style.current : ''}`}> Releases </Link>
        </div>
        {children}
      </main>
      <aside className='column aside'>
        <Widget>
          <Chart />
        </Widget>
      </aside>
    </div>
  );
};

View.propTypes = {
  children: PropTypes.any,
};

export default View;
