import React from 'react';
import {useParams, useLocation, Link} from 'react-router-dom';
import PropTypes from 'prop-types';
import Widget from 'components/widget.js';
import Chart from './chart.js';

import style from './index.module.scss';

const View = ({children}) => {
  const {genre} = useParams();
  const location = useLocation();

  return (
    <div className='container'>
      <main className='column main'>
        <div className={style.nodes}>
          <Link to="/" className={`${style.node} ${location.pathname === '/' ? style.current : ''}`}> News </Link>
          <Link to="/genres/release" className={`${style.node} ${genre === 'release' ? style.current : ''}`}> Releases </Link>
          <Link to="/faucets" className={`${style.node} ${location.pathname === '/faucets' ? style.current : ''}`}> Faucets </Link>
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
